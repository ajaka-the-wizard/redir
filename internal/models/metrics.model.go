package models

import (
	"time"

	"github.com/google/uuid"
)

type Metrics struct {
	Id         int       `json:"id" db:"id" redis:"id"`
	MediaId    uuid.UUID `json:"media_id" db:"media_id" redis:"media_id"`
	Ip         string    `json:"ip" db:"ip" redis:"ip"`
	Browser    string    `json:"browser" db:"browser" redis:"browser"`
	Os         string    `json:"os" db:"os" redis:"os"`
	Country    string    `json:"country" db:"country" redis:"country"`
	Referrer   string    `json:"referrer" db:"referrer" redis:"referrer"`
	CapturedAt time.Time `json:"captured_at" db:"captured_at" redis:"captured_at"`
}
