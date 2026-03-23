package domain

import (
	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/google/uuid"
)

type LightUser struct {
	Id    uuid.UUID
	Email string
	Admin bool
	Paid  bool
}

type PingResponseFormat struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
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
	UserId      uuid.UUID `json:"user_id"`
}

type LoginResponse struct {
	PingResponseFormat
	Errors []string `json:"errors"`
}

type CreateUserResponse struct {
	LoginResponse
}

type GetMeResponse struct {
	PingResponseFormat
	User models.User `json:"user" db:"user"`
}
