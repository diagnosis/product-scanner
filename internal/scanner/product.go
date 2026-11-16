package scanner

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID `json:"id"`
	Category  string    `json:"category"`
	Name      string    `json:"name"`
	Weight    *string   `json:"weight"`
	Calories  *int      `json:"calories"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
