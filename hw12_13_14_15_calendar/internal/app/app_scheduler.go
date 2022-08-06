package app

import (
	"context"
	"time"

	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type SchedulerApp struct {
	Logger  Logger
	Storage Storage
}

func NewSchedulerApp(logger Logger, storage Storage) *SchedulerApp {
	return &SchedulerApp{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *SchedulerApp) GetForNotify(ctx context.Context) ([]storage.Event, error) {
	return a.Storage.GetForNotify(time.Now())
}

func (a *SchedulerApp) DeleteOld(ctx context.Context) error {
	return a.Storage.DeleteOld(time.Now().AddDate(-1, 0, 0))
}
