package sqlstorage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	dsn string
	db  *sqlx.DB
	ctx context.Context
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Open("pgx", s.dsn)
	if err != nil {
		return err
	}

	s.db = db
	s.ctx = ctx

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Create(event storage.Event) error {
	sql := `
		INSERT INTO events
		  (id, title, begin_at, end_at, description, user_id, notify_at)
		VALUES
		  ($1, $2, $3, $4, $5, $6, $7)
		;
	`

	_, err := s.db.ExecContext(
		s.ctx,
		sql,
		event.ID.String(),
		event.Title,
		event.BeginAt,
		event.EndAt,
		event.Description,
		event.UserID,
		event.NotifyAt,
	)

	return err
}

func (s *Storage) Update(id uuid.UUID, event storage.Event) error {
	sql := `
		UPDATE
		  events
		SET
		  title = $2,
		  start_at = $3,
		  end_at = $4,
		  description = $5,
		  user_id = $6,
		  notify_at = $7
		WHERE
		  id = $1
		;
	`

	_, err := s.db.ExecContext(
		s.ctx,
		sql,
		id,
		event.Title,
		event.BeginAt,
		event.EndAt,
		event.Description,
		event.UserID,
		event.NotifyAt,
	)

	return err
}

func (s *Storage) Delete(id uuid.UUID) error {
	sql := `
		DELETE
		FROM
		  events
		WHERE
		  id = $1
		;
	`
	_, err := s.db.ExecContext(s.ctx, sql, id)

	return err
}

func (s *Storage) Get(id uuid.UUID) (storage.Event, error) {
	sql := `
		SELECT
		 id,
		 title,
		 begin_at,
		 end_at,
		 description,
		 user_id,
		 notify_at
		FROM
		  events
		WHERE
		  id = $1
		;
	`
	row := s.db.QueryRowxContext(s.ctx, sql, id)

	var event storage.Event
	if err := row.StructScan(&event); err != nil {
		return storage.Event{}, err
	}

	return event, nil
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
	sql := `
		SELECT
		 id,
		 title,
		 begin_at,
		 end_at,
		 description,
		 user_id,
		 notify_at
		FROM
		  events
		WHERE
		  begin_at BETWEEN $1 AND $2
		;
	`
	rows, err := s.db.QueryxContext(s.ctx, sql, beginAt, endAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]storage.Event, 0)
	for rows.Next() {
		var event storage.Event
		err := rows.StructScan(&event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
