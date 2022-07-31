package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[uuid.UUID]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[uuid.UUID]storage.Event),
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) Create(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; ok {
		return storage.ErrEventAlreadyExists
	}

	s.events[event.ID] = event

	return nil
}

func (s *Storage) Update(id uuid.UUID, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return storage.ErrEventNotExists
	}

	s.events[id] = event

	return nil
}

func (s *Storage) Delete(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return storage.ErrEventNotExists
	}

	delete(s.events, id)

	return nil
}

func (s *Storage) Get(id uuid.UUID) (storage.Event, error) {
	if event, ok := s.events[id]; ok {
		return event, nil
	}

	return storage.Event{}, storage.ErrEventNotExists
}

func (s *Storage) SelectOnDay(t time.Time) ([]storage.Event, error) {
	return s.SelectBeetween(t, t.AddDate(0, 0, 1))
}

func (s *Storage) SelectOnWeek(t time.Time) ([]storage.Event, error) {
	return s.SelectBeetween(t, t.AddDate(0, 0, 7))
}

func (s *Storage) SelectOnMonth(t time.Time) ([]storage.Event, error) {
	return s.SelectBeetween(t, t.AddDate(0, 1, 0))
}

func (s *Storage) SelectBeetween(beginAt time.Time, endAt time.Time) ([]storage.Event, error) {
	events := make([]storage.Event, 0, len(s.events))
	for _, event := range s.events {
		if event.BeginAt.After(beginAt) && event.BeginAt.Before(endAt) {
			events = append(events, event)
		}
	}

	return events, nil
}
