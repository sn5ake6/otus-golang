package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/notification"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/queue"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/version"
	"github.com/streadway/amqp"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config_sender.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		version.PrintVersion()
		return
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.NewSenderConfig()
	err := config.LoadConfig(configFile, &cfg)
	if err != nil {
		return err
	}

	logg, err := logger.New(cfg.Logger.Level)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := app.NewStorage(cfg.Storage)
	err = storage.Connect(ctx)
	if err != nil {
		return err
	}

	sender := app.NewSenderApp(logg, storage)

	go startSender(ctx, cfg, logg, sender)

	logg.Info("sender is running...")

	<-ctx.Done()

	return nil
}

func startSender(ctx context.Context, cfg config.SenderConfig, logg *logger.Logger, sender *app.SenderApp) {
	consumer := queue.NewConsumer(
		cfg.Queue.Dsn,
		cfg.Queue.Exchange,
		cfg.Queue.ExchangeType,
		cfg.Queue.Queue,
	)

	err := consumer.Handle(ctx, getWorker(logg, sender), 1)
	if err != nil {
		logg.Error(err.Error())
	}
}

func getWorker(logg *logger.Logger, sender *app.SenderApp) queue.Worker {
	return func(ctx context.Context, messages <-chan amqp.Delivery) {
		for {
			select {
			case <-ctx.Done():
				return
			case message := <-messages:
				if len(message.Body) == 0 {
					continue
				}

				var notification notification.Notification
				err := json.Unmarshal(message.Body, &notification)
				if err != nil {
					logg.Error(err.Error())
					nackMessage(message, logg)

					continue
				}

				err = sender.SendNotification(ctx, notification)
				if err != nil {
					logg.Error(err.Error())
					nackMessage(message, logg)

					continue
				}

				err = message.Ack(false)
				if err != nil {
					logg.Error(err.Error())

					return
				}
			}
		}
	}
}

func nackMessage(message amqp.Delivery, logg *logger.Logger) {
	err := message.Nack(false, false)
	if err != nil {
		logg.Error(err.Error())
	}
}
