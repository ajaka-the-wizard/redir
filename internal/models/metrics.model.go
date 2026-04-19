package models

import (
	"time"

	"github.com/google/uuid"
)

type Metrics struct {
	Id         int       `json:"id" db:"id"`
	MediaId    uuid.UUID `json:"media_id" db:"media_id"`
	Ip         string    `json:"ip" db:"ip"`
	Browser    string    `json:"browser" db:"browser"`
	Os         string    `json:"os" db:"os"`
	Country    string    `json:"country" db:"country"`
	Referrer   string    `json:"referrer" db:"referrer"`
	CapturedAt time.Time `json:"captured_at" db:"captured_at"`
}
