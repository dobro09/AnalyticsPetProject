package kafka

import (
	"analytics/internal/model"
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
    producer sarama.AsyncProducer
}

func NewProducer(brokers ...string) (*Producer, error) {
	config := sarama.NewConfig()

	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = false
	config.Producer.RequiredAcks = sarama.WaitForLocal

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for err := range producer.Errors() {
			log.Printf("failed to produce event: %v", err)
		}
	}()

	return &Producer{
		producer: producer,
	}, nil
}

func (p *Producer) Produce(ctx context.Context, event model.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: "user-events",
		Key:   sarama.StringEncoder(event.SessionID),
		Value: sarama.ByteEncoder(data),
	}

	select {
	case p.producer.Input() <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Producer) Close() error {
	return p.producer.Close()
}