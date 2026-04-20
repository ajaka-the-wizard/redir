package models

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	PublicKey string    `json:"public_key" db:"public_key"`
	InnerKey  string    `json:"-" db:"inner_key"`
	BatchId   uuid.UUID `json:"batch_id" db:"batch_id"`
	SeqId     int       `json:"seq_id" db:"seq_id"`
	UserId    uuid.UUID `json:"user_id" db:"user_id"`
	FileSize  int64     `json:"file_size" db:"file_size"`
	Status    string    `json:"status" db:"status"`
	FileType  string    `json:"file_type" db:"file_type"`
	Active    bool      `json:"active" db:"active"`
	Public    bool      `json:"public" db:"public"`
	FileName  string    `json:"file_name" db:"file_name"`
	MimeType  string    `json:"mime_type" db:"mime_type"`
	Hits      int       `json:"hits" db:"hits"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
