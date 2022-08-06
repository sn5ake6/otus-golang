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
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	BeginAt     time.Time `json:"beginAt" db:"begin_at"`
	EndAt       time.Time `json:"endAt" db:"end_at"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"userId" db:"user_id"`
	NotifyAt    time.Time `json:"notifyAt" db:"notify_at"`
	NotifiedAt  time.Time `json:"notifiedAt" db:"notified_at"`
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
