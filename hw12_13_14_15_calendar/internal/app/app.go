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

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	return a.Storage.Create(event)
}

func (a *App) UpdateEvent(ctx context.Context, id uuid.UUID, event storage.Event) error {
	return a.Storage.Update(id, event)
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return a.Storage.Delete(id)
}

func (a *App) GetEvent(ctx context.Context, id uuid.UUID) (storage.Event, error) {
	return a.Storage.Get(id)
}

func (a *App) SelectOnDayEvents(ctx context.Context, t time.Time) ([]storage.Event, error) {
	return a.Storage.SelectOnDay(t)
}

func (a *App) SelectOnWeekEvents(ctx context.Context, t time.Time) ([]storage.Event, error) {
	return a.Storage.SelectOnWeek(t)
}

func (a *App) SelectOnMonthEvents(ctx context.Context, t time.Time) ([]storage.Event, error) {
	return a.Storage.SelectOnMonth(t)
}
