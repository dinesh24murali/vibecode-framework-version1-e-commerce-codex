package models

import (
	"time"

	"github.com/google/uuid"
)

// The structs below are JSON response DTOs whose field order follows the API
// spec, so the fieldalignment check is suppressed on them.

// ProductAttribute is an attribute value assigned to a product.
type ProductAttribute struct { //nolint:govet // DTO field order follows the API spec
	AttributeValueID uuid.UUID `json:"attribute_value_id"`
	AttributeID      uuid.UUID `json:"attribute_id"`
	Code             string    `json:"code"`
	AttributeName    string    `json:"attribute_name"`
	Value            string    `json:"value"`
	ValueSlug        string    `json:"value_slug"`
	IsFilterable     bool      `json:"is_filterable"`
}

// ProductImage is one product image.
type ProductImage struct { //nolint:govet // DTO field order follows the API spec
	ID        uuid.UUID `json:"id"`
	URL       string    `json:"url"`
	AltText   *string   `json:"alt_text"`
	SortOrder int32     `json:"sort_order"`
}

// ProductSummary is a product as listed in the catalogue.
type ProductSummary struct { //nolint:govet // DTO field order follows the API spec
	ID              uuid.UUID          `json:"id"`
	Name            string             `json:"name"`
	Slug            string             `json:"slug"`
	Price           float64            `json:"price"`
	DiscountPercent float64            `json:"discount_percent"`
	DiscountedPrice float64            `json:"discounted_price"`
	Attributes      []ProductAttribute `json:"attributes"`
	CoverImage      *ProductImage      `json:"cover_image"`
	Categories      []Category         `json:"categories"`
	IsActive        bool               `json:"is_active"`
}

// Product is the full product detail.
type Product struct { //nolint:govet // DTO field order follows the API spec
	ID              uuid.UUID          `json:"id"`
	Name            string             `json:"name"`
	Slug            string             `json:"slug"`
	Description     string             `json:"description"`
	Price           float64            `json:"price"`
	DiscountPercent float64            `json:"discount_percent"`
	DiscountedPrice float64            `json:"discounted_price"`
	Stock           int32              `json:"stock"`
	Attributes      []ProductAttribute `json:"attributes"`
	Images          []ProductImage     `json:"images"`
	Categories      []Category         `json:"categories"`
	IsActive        bool               `json:"is_active"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

// FilterValue is one selectable value inside an AttributeFilterGroup.
type FilterValue struct { //nolint:govet // DTO field order follows the API spec
	ID        uuid.UUID `json:"id"`
	ValueSlug string    `json:"value_slug"`
	Value     string    `json:"value"`
}

// AttributeFilterGroup is one filterable attribute dimension.
type AttributeFilterGroup struct { //nolint:govet // DTO field order follows the API spec
	AttributeID   uuid.UUID     `json:"attribute_id"`
	Code          string        `json:"code"`
	AttributeName string        `json:"attribute_name"`
	IsFilterable  bool          `json:"is_filterable"`
	Values        []FilterValue `json:"values"`
}

// ProductFilters describes the filter options available for the catalogue.
type ProductFilters struct { //nolint:govet // DTO field order follows the API spec
	Categories  []string               `json:"categories"`
	Attributes  []AttributeFilterGroup `json:"attributes"`
	SortOptions []string               `json:"sort_options"`
}

// PaginationMeta is the pagination block of every list response.
type PaginationMeta struct { //nolint:govet // DTO field order follows the API spec
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

// ProductListResponse is the GET /products response body.
type ProductListResponse struct { //nolint:govet // DTO field order follows the API spec
	Data    []ProductSummary `json:"data"`
	Meta    PaginationMeta   `json:"meta"`
	Filters ProductFilters   `json:"filters"`
}
