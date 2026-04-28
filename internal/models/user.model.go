package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `json:"id" db:"id" redis:"id"`
	FullName    string    `json:"full_name" db:"full_name" redis:"full_name"`
	Email       string    `json:"email" db:"email" redis:"email"`
	Provider    string    `json:"-" db:"provider" redis:"provider"`
	ProviderSub string    `json:"-" db:"provider_sub" redis:"provider_sub"`
	Verified    bool      `json:"verified" db:"verified" redis:"verified"`
	Paid        bool      `json:"paid" db:"paid" redis:"paid"`
	Admin       bool      `json:"admin" db:"admin" redis:"admin"`
	Active      bool      `json:"active" db:"active" redis:"active"`
	Password    string    `json:"-" db:"password" redis:"password"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" redis:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" redis:"created_at"`
}
