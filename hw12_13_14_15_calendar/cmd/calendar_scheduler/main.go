package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/notification"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/queue"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/version"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config_scheduler.toml", "Path to configuration file")
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
	cfg := config.NewSchedulerConfig()
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

	scheduler := app.NewSchedulerApp(logg, storage)

	go startNotify(ctx, cfg, logg, scheduler)
	go startDelete(ctx, cfg, logg, scheduler)

	logg.Info("scheduler is running...")

	<-ctx.Done()

	return nil
}

func startNotify(ctx context.Context, cfg config.SchedulerConfig, logg *logger.Logger, scheduler *app.SchedulerApp) {
	logg.Info("notify started")
	defer logg.Info("notify finished")

	producer := queue.NewProducer(
		cfg.Queue.Dsn,
		cfg.Queue.Queue,
	)

	notifyInterval, err := time.ParseDuration(cfg.Intervals.Notify)
	if err != nil {
		logg.Error(err.Error())
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(notifyInterval):
			events, err := scheduler.GetForNotify(ctx)
			if err != nil {
				logg.Error(err.Error())
			}

			for _, event := range events {
				message := notification.Notification{
					EventID: event.ID,
					Title:   event.Title,
					BeginAt: event.BeginAt,
					UserID:  event.UserID,
				}

				body, err := json.Marshal(message)
				if err != nil {
					logg.Error(err.Error())

					continue
				}

				err = producer.Publish(ctx, body)
				if err != nil {
					logg.Error(err.Error())
				}

				logg.Info("published to queue: " + message.EventID.String())
			}
		}
	}
}

func startDelete(ctx context.Context, cfg config.SchedulerConfig, logg *logger.Logger, scheduler *app.SchedulerApp) {
	logg.Info("delete started")
	defer logg.Info("delete finished")

	deleteInterval, err := time.ParseDuration(cfg.Intervals.Delete)
	if err != nil {
		logg.Error(err.Error())
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(deleteInterval):
			err := scheduler.DeleteOld(ctx)
			if err != nil {
				logg.Error(err.Error())

				continue
			}

			logg.Info("old events deleted")
		}
	}
}
