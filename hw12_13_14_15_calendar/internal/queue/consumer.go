package queue

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	backoff "github.com/cenkalti/backoff/v3"
	"github.com/streadway/amqp"
)

type Consumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	done         chan error
	dsn          string
	exchange     string
	exchangeType string
	queue        string
}

func NewConsumer(dsn, exchange, exchangeType, queue string) *Consumer {
	return &Consumer{
		dsn:          dsn,
		exchange:     exchange,
		exchangeType: exchangeType,
		queue:        queue,
		done:         make(chan error),
	}
}

type Worker func(context.Context, <-chan amqp.Delivery)

func (c *Consumer) Handle(ctx context.Context, fn Worker, threads int) error {
	var err error
	if err = c.connect(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	msgs, err := c.announceQueue()
	if err != nil {
		return fmt.Errorf("announceQueue: %w", err)
	}

	for {
		for i := 0; i < threads; i++ {
			go fn(ctx, msgs)
		}

		if <-c.done != nil {
			msgs, err = c.reConnect(ctx)
			if err != nil {
				return fmt.Errorf("reConnect: %w", err)
			}
		}
		fmt.Println("Reconnected... possibly")
	}
}

func (c *Consumer) connect() error {
	var err error

	c.conn, err = amqp.Dial(c.dsn)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}

	go func() {
		log.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		// Понимаем, что канал сообщений закрыт, надо пересоздать соединение.
		c.done <- errors.New("channel Closed")
	}()

	if err = c.channel.ExchangeDeclare(
		c.exchange,
		c.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("exchange declare: %w", err)
	}

	return nil
}

// Задекларировать очередь, которую будем слушать.
func (c *Consumer) announceQueue() (<-chan amqp.Delivery, error) {
	queue, err := c.channel.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queueDeclare: %w", err)
	}

	// Число сообщений, которые можно подтвердить за раз.
	err = c.channel.Qos(50, 0, false)
	if err != nil {
		return nil, fmt.Errorf("qos: %w", err)
	}

	// Создаём биндинг (правило маршрутизации).
	if err = c.channel.QueueBind(
		queue.Name,
		"",
		c.exchange,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("queueBind: %w", err)
	}

	msgs, err := c.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	return msgs, nil
}

func (c *Consumer) reConnect(ctx context.Context) (<-chan amqp.Delivery, error) {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, ctx)
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return nil, fmt.Errorf("stop reconnecting")
		}

		for range time.After(d) {
			if err := c.connect(); err != nil {
				log.Printf("could not connect in reconnect call: %+v", err)
				continue
			}
			msgs, err := c.announceQueue()
			if err != nil {
				fmt.Printf("couldn't connect: %+v", err)
				continue
			}

			return msgs, nil
		}
	}
}
