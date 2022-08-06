package app

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
)

type Storage interface {
	Connect(ctx context.Context) error
	Create(event storage.Event) error
	Update(id uuid.UUID, event storage.Event) error
	Delete(id uuid.UUID) error
	Get(id uuid.UUID) (storage.Event, error)
	SelectOnDay(t time.Time) ([]storage.Event, error)
	SelectOnWeek(t time.Time) ([]storage.Event, error)
	SelectOnMonth(t time.Time) ([]storage.Event, error)
	GetForNotify(t time.Time) ([]storage.Event, error)
	DeleteOld(t time.Time) error
}

func NewStorage(storageConfig config.StorageConf) Storage {
	var storage Storage

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
