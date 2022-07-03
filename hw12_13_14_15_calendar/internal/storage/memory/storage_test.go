package memorystorage

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	memorystorage "github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	storage := New()

	userID := uuid.New()

	now := time.Now()
	notifyAt := now.AddDate(0, 0, 1)
	beginAt := now.AddDate(0, 0, 2)
	endAt := now.AddDate(0, 0, 3)

	event := memorystorage.NewEvent(
		"some event",
		beginAt,
		endAt,
		"some description",
		userID,
		notifyAt,
	)

	err := storage.Create(event)
	require.NoError(t, err)

	t.Run("create cases", func(t *testing.T) {
		err = storage.Create(event)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, memorystorage.ErrEventAlreadyExists))
	})

	t.Run("update cases", func(t *testing.T) {
		event.Title = "changed title"
		event.Description = "changed description"

		onDayEventsBeforeUpdate, err := storage.SelectOnDay(event.BeginAt.Add(-1 * time.Hour))
		require.NoError(t, err)
		require.Len(t, onDayEventsBeforeUpdate, 1)

		onDayEvent := onDayEventsBeforeUpdate[0]
		require.NotEqual(t, event, onDayEvent)

		err = storage.Update(event.ID, event)
		require.NoError(t, err)

		updatedOnDayEvents, err := storage.SelectOnDay(event.BeginAt.Add(-1 * time.Hour))
		require.NoError(t, err)
		require.Len(t, updatedOnDayEvents, 1)

		updatedEvent := updatedOnDayEvents[0]
		require.Equal(t, event, updatedEvent)
	})

	t.Run("select cases", func(t *testing.T) {
		emptyOnDayEvents, err := storage.SelectOnDay(event.BeginAt.AddDate(0, 0, -2))
		require.NoError(t, err)
		require.Len(t, emptyOnDayEvents, 0)

		onDayEvents, err := storage.SelectOnDay(event.BeginAt.Add(-1 * time.Hour))
		require.NoError(t, err)
		require.Len(t, onDayEvents, 1)
		require.Equal(t, event, onDayEvents[0])

		emptyOnWeekEvents, err := storage.SelectOnWeek(event.BeginAt.AddDate(0, 0, -8))
		require.NoError(t, err)
		require.Len(t, emptyOnWeekEvents, 0)

		onWeekEvents, err := storage.SelectOnWeek(event.BeginAt.AddDate(0, 0, -1))
		require.NoError(t, err)
		require.Len(t, onWeekEvents, 1)
		require.Equal(t, event, onWeekEvents[0])

		emptyOnMonthEvents, err := storage.SelectOnWeek(event.BeginAt.AddDate(0, -1, -1))
		require.NoError(t, err)
		require.Len(t, emptyOnMonthEvents, 0)

		onMonthEvents, err := storage.SelectOnWeek(event.BeginAt.AddDate(0, 0, -1))
		require.NoError(t, err)
		require.Len(t, onMonthEvents, 1)
		require.Equal(t, event, onMonthEvents[0])
	})

	t.Run("delete cases", func(t *testing.T) {
		err = storage.Delete(event.ID)
		require.NoError(t, err)

		deletedOnDayEvents, err := storage.SelectOnDay(event.BeginAt.Add(-1 * time.Hour))
		require.NoError(t, err)
		require.Len(t, deletedOnDayEvents, 0)

		err = storage.Delete(uuid.New())
		require.True(t, errors.Is(err, memorystorage.ErrEventNotExists))
	})
}
