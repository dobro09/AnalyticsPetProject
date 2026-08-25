package kafka

import (
	"analytics/internal/model"
	"context"
	"encoding/json"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
}

func NewProducer(brokers string) (*Producer, error) {
	client, err := kgo.NewClient(kgo.SeedBrokers(brokers))
	if err != nil {
		return nil, err
	}
	return &Producer{
		client: client,
	}, nil
}

func (p *Producer) Produce(ctx context.Context, event model.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	record := &kgo.Record{
		Topic: "user-events",
		Key: []byte(event.SessionID),
		Value: data,
	}

	produceCtx := context.WithoutCancel(ctx)

	p.client.Produce(produceCtx, record, func(r *kgo.Record, err error) {
		if err != nil {
			log.Printf("failed to produce event: %v", err)
			return 
		}
		log.Printf("produced event to kafka: key=%s", record.Key)
	})
	return nil
}

func (p *Producer) Close() {
    p.client.Close()
}