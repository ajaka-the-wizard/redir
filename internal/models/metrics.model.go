package models

import (
	"time"

	"github.com/google/uuid"
)

type Metrics struct {
	Id           uuid.UUID `json:"id" db:"id"`
	MediaId      uuid.UUID `json:"media_id" db:"media_id"`
	UploadId     int       `json:"upload_id" db:"upload_id"`
	ErrorMessage string    `json:"error_message" db:"error_message"`
	CapturedAt   time.Time `json:"captured_at" db:"captured_at"`
}
