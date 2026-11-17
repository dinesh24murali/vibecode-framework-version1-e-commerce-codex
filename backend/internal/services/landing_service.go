package services

import (
	"context"
	"database/sql"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/models"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// FeaturedCount is how many featured products the landing page shows.
const FeaturedCount = 8

// LandingService assembles the landing page payload.
type LandingService struct {
	carousel repositories.CarouselRepository
	products *ProductService
}

// NewLandingService constructs a LandingService.
func NewLandingService(c repositories.CarouselRepository, p *ProductService) *LandingService {
	return &LandingService{carousel: c, products: p}
}

func strPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

// Landing returns the carousel slides and featured products.
func (s *LandingService) Landing(ctx context.Context) (*models.LandingPageResponse, error) {
	rows, err := s.carousel.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	slides := make([]models.CarouselSlide, 0, len(rows))
	for _, r := range rows {
		slides = append(slides, models.CarouselSlide{
			ID: r.ID, Title: r.Title, Subtitle: strPtr(r.Subtitle), ImageURL: r.ImageUrl,
			CTALabel: strPtr(r.CtaLabel), CTAURL: strPtr(r.CtaUrl), SortOrder: r.SortOrder,
		})
	}
	featured, err := s.products.Featured(ctx, FeaturedCount)
	if err != nil {
		return nil, err
	}
	return &models.LandingPageResponse{CarouselSlides: slides, FeaturedProducts: featured}, nil
}
