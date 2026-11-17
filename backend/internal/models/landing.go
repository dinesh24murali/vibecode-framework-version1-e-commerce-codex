package models

import "github.com/google/uuid"

// CarouselSlide is one hero carousel slide.
type CarouselSlide struct { //nolint:govet // DTO field order follows the API spec
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Subtitle  *string   `json:"subtitle"`
	ImageURL  string    `json:"image_url"`
	CTALabel  *string   `json:"cta_label"`
	CTAURL    *string   `json:"cta_url"`
	SortOrder int32     `json:"sort_order"`
}

// LandingPageResponse is the GET /landing response body.
type LandingPageResponse struct { //nolint:govet // DTO field order follows the API spec
	CarouselSlides   []CarouselSlide  `json:"carousel_slides"`
	FeaturedProducts []ProductSummary `json:"featured_products"`
}
