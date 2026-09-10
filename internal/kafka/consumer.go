package kafka

import (
	"analytics/internal/model"
	"analytics/internal/usecase"
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	group sarama.ConsumerGroup
}

func NewConsumer(brokers []string, groupID string) (*Consumer, error) {
	config := sarama.NewConfig()

	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
        sarama.NewBalanceStrategyRoundRobin(),
    }
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup(
		brokers,
		groupID,
		config,
	)
	if err != nil {
		log.Printf("ошибка сосздания конс группы: %v", err)
		return nil, err
	}
	go func() {
    	for err := range group.Errors() {
     	   log.Printf("consumer group error: %v", err)
    	}
	}()

	return &Consumer{
		group: group,
	}, nil
}

func (c *Consumer) Consume(ctx context.Context, topics []string, handler *ConsumerHandler) error {
	for {
		if err := c.group.Consume(ctx, topics, handler); err != nil {
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}


func (c *Consumer) Close() error {
	return c.group.Close()
}

type ConsumerHandler struct {
	usecase usecase.AnalyticsUsecase
}

func NewConsumerHandler(uc usecase.AnalyticsUsecase) *ConsumerHandler {
	return &ConsumerHandler{
		usecase: uc,
	}
}

func (h *ConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession,
										claim sarama.ConsumerGroupClaim) error{
	for message := range claim.Messages() {
		var event model.Event

		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("failed to unmarshal event, offset=%d: %v", message.Offset, err)
			continue
		}

		if err := h.usecase.Analyse(session.Context(), event); err != nil {
			log.Printf("failed to process event, offset=%d: %v", message.Offset, err)
			continue
		}
		log.Printf("event received: type=%s session=%s", event.EventType, event.SessionID)
		session.MarkMessage(message, "")
	}
	return nil
}