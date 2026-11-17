package models

import (
	"time"

	"github.com/google/uuid"
)

// AttributeValue is one allowed value of an attribute (API: AttributeDefinitionValue).
// Field order is chosen for the fieldalignment linter, not for readability.
type AttributeValue struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Value       string    `json:"value"`
	ValueSlug   string    `json:"value_slug"`
	ID          uuid.UUID `json:"id"`
	AttributeID uuid.UUID `json:"attribute_id"`
	SortOrder   int32     `json:"sort_order"`
}

// AttributeWithValues is a filterable attribute definition and its values
// (API: AttributeWithValues).
type AttributeWithValues struct {
	Code         string           `json:"code"`
	Name         string           `json:"name"`
	Values       []AttributeValue `json:"values"`
	ID           uuid.UUID        `json:"id"`
	SortOrder    int32            `json:"sort_order"`
	IsFilterable bool             `json:"is_filterable"`
}
