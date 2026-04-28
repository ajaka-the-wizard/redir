package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          int       `json:"-" db:"id" redis:"id"`
	ProductId   int       `json:"product_id" db:"product_id" redis:"product_id"`
	ProductName string    `json:"product_name" db:"product_name" redis:"product_name"`
	UserId      uuid.UUID `json:"user_id" db:"user_id" redis:"user_id"`
	Public      bool      `json:"public" db:"public" redis:"public"`
	PrivateKey  string    `json:"private_key" db:"private_key" redis:"private_key"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" redis:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" redis:"updated_at"`
}
