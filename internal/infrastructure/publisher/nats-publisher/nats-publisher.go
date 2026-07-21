package nats_publisher

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

type Conf struct {
	NatsURL string
	Subj    string
}

type NatsPublisher struct {
	nc *nats.Conn
	Conf
}

func New(cfg Conf) (*NatsPublisher, error) {
	conn, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		return nil, fmt.Errorf("connection to nats: %w", err)
	}

	return &NatsPublisher{
		nc:   conn,
		Conf: cfg,
	}, nil
}

func (p *NatsPublisher) PublishUserID(ctx context.Context, userID string) error {
	if err := p.nc.Publish(p.Conf.Subj, []byte(userID)); err != nil {
		return fmt.Errorf("publishing user id: %w", err)
	}

	return nil
}

func (p *NatsPublisher) Close() {
	p.nc.Close()
}
