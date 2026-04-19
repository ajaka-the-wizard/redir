package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `json:"id" db:"id"`
	FullName    string    `json:"full_name" db:"full_name"`
	Email       string    `json:"email" db:"email"`
	Provider    string    `json:"-" db:"provider"`
	ProviderSub string    `json:"-" db:"provider_sub"`
	Verified    bool      `json:"verified" db:"verified"`
	Paid        bool      `json:"paid" db:"paid"`
	Admin       bool      `json:"admin" db:"admin"`
	Active      bool      `json:"active" db:"active"`
	Password    string    `json:"-" db:"password"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
