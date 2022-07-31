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
	internalgrpc "github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/server/grpc"
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

	startGRPCServer(ctx, config, logg, calendar)

	startHTTPServer(ctx, config, logg, calendar)

	logg.Info("calendar is running...")

	<-ctx.Done()

	return nil
}

func startGRPCServer(ctx context.Context, config *Config, logg *logger.Logger, calendar *app.App) {
	grpcAddr := net.JoinHostPort(config.GRPCServer.Host, config.GRPCServer.Port)
	grpcServer := internalgrpc.NewServer(grpcAddr, logg, calendar)

	go func() {
		if err := grpcServer.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
		}
	}()

	go func() {
		<-ctx.Done()
		grpcServer.Stop()
	}()
}

func startHTTPServer(ctx context.Context, config *Config, logg *logger.Logger, calendar *app.App) {
	httpAddr := net.JoinHostPort(config.HTTPServer.Host, config.HTTPServer.Port)
	httpServer := internalhttp.NewServer(httpAddr, logg, calendar)

	go func() {
		if err := httpServer.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
		}
	}()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()
}
