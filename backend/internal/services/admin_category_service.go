package services

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/models"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// AdminCategoryInput carries create/update fields; nil means "not provided".
type AdminCategoryInput struct {
	Name     *string
	Slug     *string
	ImageURL *string
}

// AdminCategoryService implements admin category management.
type AdminCategoryService struct {
	repo repositories.CategoryAdminRepository
}

// NewAdminCategoryService constructs an AdminCategoryService.
func NewAdminCategoryService(repo repositories.CategoryAdminRepository) *AdminCategoryService {
	return &AdminCategoryService{repo: repo}
}

// List returns one page of categories and the total.
func (s *AdminCategoryService) List(ctx context.Context, page, limit int) ([]models.Category, int64, error) {
	rows, total, err := s.repo.ListPaged(ctx, int32(limit), int32((page-1)*limit)) //nolint:gosec // bounded by handler
	if err != nil {
		return nil, 0, err
	}
	out := make([]models.Category, 0, len(rows))
	for _, c := range rows {
		out = append(out, toCategory(c))
	}
	return out, total, nil
}

// Get returns one category or ErrNotFound.
func (s *AdminCategoryService) Get(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	m := toCategory(*c)
	return &m, nil
}

// Create validates and creates a category; a duplicate slug yields ErrConflict.
func (s *AdminCategoryService) Create(ctx context.Context, in AdminCategoryInput) (*models.Category, error) {
	if in.Name == nil || in.Slug == nil {
		return nil, ErrValidation
	}
	name, slug := strings.TrimSpace(*in.Name), strings.TrimSpace(*in.Slug)
	if name == "" || !productSlugRe.MatchString(slug) {
		return nil, ErrValidation
	}
	c, err := s.repo.Create(ctx, slug, name, in.ImageURL)
	if err != nil {
		if repositories.IsUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	m := toCategory(*c)
	return &m, nil
}

// Update applies a partial update.
func (s *AdminCategoryService) Update(ctx context.Context, id uuid.UUID, in AdminCategoryInput) (*models.Category, error) {
	cur, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	name, slug := cur.Name, cur.Slug
	var img *string
	if cur.ImageUrl.Valid {
		img = &cur.ImageUrl.String
	}
	if in.Name != nil {
		if name = strings.TrimSpace(*in.Name); name == "" {
			return nil, ErrValidation
		}
	}
	if in.Slug != nil {
		if slug = strings.TrimSpace(*in.Slug); !productSlugRe.MatchString(slug) {
			return nil, ErrValidation
		}
	}
	if in.ImageURL != nil {
		img = in.ImageURL
	}
	c, err := s.repo.Update(ctx, id, slug, name, img)
	if err != nil {
		if repositories.IsUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, mapNotFound(err)
	}
	m := toCategory(*c)
	return &m, nil
}

// Delete removes a category; ErrConflict if any product still references it.
func (s *AdminCategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return mapNotFound(err)
	}
	n, err := s.repo.CountProducts(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrConflict
	}
	ok, err := s.repo.Delete(ctx, id)
	if err != nil {
		if repositories.IsForeignKeyViolation(err) {
			return ErrConflict
		}
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}
