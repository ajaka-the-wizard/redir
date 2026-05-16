package models

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	PublicKey string    `json:"public_key" db:"public_key" redis:"public_key"`
	InnerKey  string    `json:"-" db:"inner_key" redis:"inner_key"`
	BatchId   uuid.UUID `json:"batch_id" db:"batch_id" redis:"batch_id"`
	SeqId     int       `json:"seq_id" db:"seq_id" redis:"seq_id"`
	UserId    uuid.UUID `json:"user_id" db:"user_id" redis:"user_id"`
	FileSize  int64     `json:"file_size" db:"file_size" redis:"file_size"`
	Status    string    `json:"status" db:"status" redis:"status"`
	FileType  string    `json:"file_type" db:"file_type" redis:"file_type"`
	Active    bool      `json:"active" db:"active" redis:"active"`
	Public    bool      `json:"public" db:"public" redis:"public"`
	FileName  string    `json:"file_name" db:"file_name" redis:"file_name"`
	MimeType  string    `json:"mime_type" db:"mime_type" redis:"mime_type"`
	Hits      int       `json:"hits" db:"hits" redis:"hits"`
	CreatedAt time.Time `json:"created_at" db:"created_at" redis:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" redis:"updated_at"`
}
