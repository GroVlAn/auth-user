package publisher

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Conf struct {
	Brokers []string
	Topic   string
}

type UserEventPublisher struct {
	producer *kafka.Writer
}

func New(cfg Conf) *UserEventPublisher {
	return &UserEventPublisher{
		producer: kafka.NewWriter(
			kafka.WriterConfig{
				Brokers: cfg.Brokers,
				Topic:   cfg.Topic,
			},
		),
	}
}

func (p *UserEventPublisher) PublishUserID(ctx context.Context, userID string) error {
	err := p.producer.WriteMessages(ctx, kafka.Message{
		Value: []byte(userID),
	})
	if err != nil {
		return fmt.Errorf("publishing user id: %w", err)
	}

	return nil
}

func (p *UserEventPublisher) Close() error {
	if err := p.producer.Close(); err != nil {
		return fmt.Errorf("closing kafka producer: %w", err)
	}

	return nil
}
