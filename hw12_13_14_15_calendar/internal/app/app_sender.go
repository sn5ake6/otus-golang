package app

import (
	"context"
	"time"

	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/notification"
)

type SenderApp struct {
	Logger  Logger
	Storage Storage
}

func NewSenderApp(logger Logger, storage Storage) *SenderApp {
	return &SenderApp{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *SenderApp) SendNotification(
	ctx context.Context,
	notification notification.Notification,
) error {
	event, err := a.Storage.Get(notification.EventID)
	if err != nil {
		return err
	}

	if event.NotifiedAt.IsZero() {
		a.Logger.Info("send notification: " + notification.EventID.String())

		event.NotifiedAt = time.Now()

		err = a.Storage.Update(notification.EventID, event)
		if err != nil {
			return err
		}
	}

	return nil
}
