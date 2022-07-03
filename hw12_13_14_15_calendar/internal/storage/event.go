package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEventAlreadyExists = errors.New("event already exists")
	ErrEventNotExists     = errors.New("event not exists")
)

type Event struct {
	ID          uuid.UUID
	Title       string
	BeginAt     time.Time
	EndAt       time.Time
	Description string
	UserID      uuid.UUID
	NotifyAt    time.Time
}

func NewEvent(
	title string,
	beginAt time.Time,
	endAt time.Time,
	description string,
	userID uuid.UUID,
	notifyAt time.Time,
) Event {
	return Event{
		ID:          uuid.New(),
		Title:       title,
		BeginAt:     beginAt,
		EndAt:       endAt,
		Description: description,
		UserID:      userID,
		NotifyAt:    notifyAt,
	}
}
