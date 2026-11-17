package services

import (
	"context"

	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/models"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// AttributeService serves the public filter-option reads.
type AttributeService struct {
	repo repositories.AttributeRepository
}

// NewAttributeService constructs an AttributeService.
func NewAttributeService(repo repositories.AttributeRepository) *AttributeService {
	return &AttributeService{repo: repo}
}

func toValue(v dbsqlc.ProductAttributeValue) models.AttributeValue {
	return models.AttributeValue{
		ID: v.ID, AttributeID: v.AttributeID, Value: v.Value, ValueSlug: v.Slug,
		SortOrder: v.SortOrder, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt,
	}
}

func toValues(rows []dbsqlc.ProductAttributeValue) []models.AttributeValue {
	out := make([]models.AttributeValue, 0, len(rows))
	for _, v := range rows {
		out = append(out, toValue(v))
	}
	return out
}

// ListFilterable returns every filterable attribute with its values.
func (s *AttributeService) ListFilterable(ctx context.Context) ([]models.AttributeWithValues, error) {
	attrs, err := s.repo.ListFilterable(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]models.AttributeWithValues, 0, len(attrs))
	for _, a := range attrs {
		vals, err := s.repo.ListValues(ctx, a.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, models.AttributeWithValues{
			ID: a.ID, Code: a.Code, Name: a.Name, IsFilterable: a.IsFilterable,
			SortOrder: a.SortOrder, Values: toValues(vals),
		})
	}
	return out, nil
}

// ListValues returns the values of one attribute, or ErrNotFound if it does not exist.
func (s *AttributeService) ListValues(ctx context.Context, attributeID uuid.UUID) ([]models.AttributeValue, error) {
	if _, err := s.repo.GetByID(ctx, attributeID); err != nil {
		return nil, mapNotFound(err)
	}
	vals, err := s.repo.ListValues(ctx, attributeID)
	if err != nil {
		return nil, err
	}
	return toValues(vals), nil
}
