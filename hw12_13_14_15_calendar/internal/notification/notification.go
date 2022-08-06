package notification

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	EventID uuid.UUID
	Title   string
	BeginAt time.Time
	UserID  uuid.UUID
}
