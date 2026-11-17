package models

import (
	"time"

	"github.com/google/uuid"
)

// Category is a product category.
type Category struct { //nolint:govet // DTO field order follows the API spec
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	ImageURL  *string   `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
