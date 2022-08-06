package queue

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	dsn     string
	queue   string
}

func NewProducer(dsn, queue string) *Producer {
	return &Producer{
		dsn:   dsn,
		queue: queue,
	}
}

func (p *Producer) Publish(ctx context.Context, body []byte) error {
	if p.conn == nil || p.conn.IsClosed() {
		if err := p.Connect(); err != nil {
			return fmt.Errorf("connect: %w", err)
		}
	}

	queue, err := p.channel.QueueDeclarePassive(p.queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queueDeclarePassive: %w", err)
	}

	err = p.channel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

func (p *Producer) Connect() error {
	var err error

	p.conn, err = amqp.Dial(p.dsn)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	p.channel, err = p.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}

	return nil
}

func (p *Producer) Disconnect() error {
	if p.conn == nil {
		return nil
	}

	return p.conn.Close()
}
