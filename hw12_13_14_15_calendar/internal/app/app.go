package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface { // TODO
}

type Storage interface {
	Connect(ctx context.Context) error
	Create(event storage.Event) error
	Update(id uuid.UUID, event storage.Event) error
	Delete(id uuid.UUID) error
	SelectOnDay(t time.Time) ([]storage.Event, error)
	SelectOnWeek(t time.Time) ([]storage.Event, error)
	SelectOnMonth(t time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
