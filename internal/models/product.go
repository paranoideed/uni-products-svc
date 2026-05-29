package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Price     float32    `json:"price"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
