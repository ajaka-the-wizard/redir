package domain

import (
	"time"

	"github.com/google/uuid"
)

type LightUser struct {
	Id               uuid.UUID
	Email            string
	Admin            bool
	Paid             bool
	LastAccessedTime time.Time
	Expires          time.Time
}

type CreateUserDetails struct {
	FullName string `json:"full_name" binding:"required"`
	LoginUserDetails
}

type LoginUserDetails struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateProductDetails struct {
	ProductName string    `json:"product_name" binding:"required"`
	UserId      uuid.UUID `json:"-"`
}

type GoogleUser struct {
	ID            string `json:"sub"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}
