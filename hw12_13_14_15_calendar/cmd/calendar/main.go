package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.toml", "Path to configuration file")
}

func NewStorage(storageConfig StorageConf) app.Storage {
	var storage app.Storage

	switch storageConfig.Type {
	case "memory":
		storage = memorystorage.New()
	case "sql":
		storage = sqlstorage.New(storageConfig.Dsn)
	default:
		log.Fatal("Unknown storage type: " + storageConfig.Type)
	}

	return storage
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config, err := LoadConfig(configFile)
	if err != nil {
		return err
	}

	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := NewStorage(config.Storage)
	err = storage.Connect(ctx)
	if err != nil {
		return err
	}

	calendar := app.New(logg, storage)

	addr := net.JoinHostPort(config.HTTPServer.Host, config.HTTPServer.Port)
	server := internalhttp.NewServer(addr, logg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())

		return err
	}

	return nil
}
