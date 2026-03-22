package models

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	Id        uuid.UUID `json:"id" db:"id"`
	BucketKey string    `json:"bucket_key" db:"bucket_key"`
	UserId    uuid.UUID `json:"user_id" db:"user_id"`
	Bucket    string    `json:"bucket" db:"bucket"`
	FileSize  int64     `json:"file_size" db:"file_size"`
	Status    string    `json:"status" db:"status"`
	FileType  string    `json:"file_type" db:"file_type"`
	Active    bool      `json:"active" db:"active"`
	FileName  string    `json:"file_name" db:"file_name"`
	MimeType  string    `json:"mime_type" db:"mime_type"`
	Hits      int       `json:"hits" db:"hits"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
