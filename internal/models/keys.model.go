package models

import (
	"time"

	"github.com/google/uuid"
)

type ClientKeys struct {
	ID         int       `json:"-" db:"id"`
	ClientId   int       `json:"client_id" db:"client_id"`
	UserId     uuid.UUID `json:"user_id" db:"user_id"`
	PrivateKey string    `json:"private_key" db:"private_key"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
