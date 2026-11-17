package services

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// Errors returned by AttributeAdminService (in addition to ErrNotFound).
var (
	ErrConflict   = errors.New("conflict")
	ErrValidation = errors.New("validation failed")
)

var (
	codeRe    = regexp.MustCompile(`^[a-z0-9]+(_[a-z0-9]+)*$`)
	slugBadRe = regexp.MustCompile(`[^a-z0-9]+`)
)

// AttributeAdminService implements admin CRUD for attribute definitions and values.
type AttributeAdminService struct {
	repo repositories.AttributeAdminRepository
}

// NewAttributeAdminService constructs an AttributeAdminService.
func NewAttributeAdminService(repo repositories.AttributeAdminRepository) *AttributeAdminService {
	return &AttributeAdminService{repo: repo}
}

// AttributeDefInput carries create/update fields; nil means "not provided".
type AttributeDefInput struct {
	Name         *string
	IsFilterable *bool
	SortOrder    *int32
	Code         string
}

// AttributeValueInput carries create/update fields; nil means "not provided".
type AttributeValueInput struct {
	Value     *string
	ValueSlug *string
	SortOrder *int32
}

// Slugify turns a label into a URL-safe slug.
func Slugify(s string) string {
	return strings.Trim(slugBadRe.ReplaceAllString(strings.ToLower(strings.TrimSpace(s)), "-"), "-")
}

// ListDefinitions returns one page of attribute definitions and the total count.
func (s *AttributeAdminService) ListDefinitions(ctx context.Context, page, limit int32) ([]dbsqlc.ProductAttribute, int64, error) {
	return s.repo.ListPaged(ctx, limit, (page-1)*limit)
}

// GetDefinition returns one definition or ErrNotFound.
func (s *AttributeAdminService) GetDefinition(ctx context.Context, id uuid.UUID) (*dbsqlc.ProductAttribute, error) {
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	return a, nil
}

// CreateDefinition creates a definition; a duplicate code yields ErrConflict.
func (s *AttributeAdminService) CreateDefinition(ctx context.Context, in AttributeDefInput) (*dbsqlc.ProductAttribute, error) {
	name := ""
	if in.Name != nil {
		name = strings.TrimSpace(*in.Name)
	}
	code := strings.TrimSpace(in.Code)
	if name == "" || !codeRe.MatchString(code) {
		return nil, ErrValidation
	}
	filterable := true
	if in.IsFilterable != nil {
		filterable = *in.IsFilterable
	}
	var sort int32
	if in.SortOrder != nil {
		sort = *in.SortOrder
	}
	a, err := s.repo.Create(ctx, dbsqlc.CreateAttributeParams{Code: code, Name: name, IsFilterable: filterable, SortOrder: sort})
	if err != nil {
		if repositories.IsUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	return a, nil
}

// UpdateDefinition applies a partial update (code is immutable).
func (s *AttributeAdminService) UpdateDefinition(ctx context.Context, id uuid.UUID, in AttributeDefInput) (*dbsqlc.ProductAttribute, error) {
	cur, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	arg := dbsqlc.UpdateAttributeParams{ID: id, Name: cur.Name, IsFilterable: cur.IsFilterable, SortOrder: cur.SortOrder}
	if in.Name != nil {
		if arg.Name = strings.TrimSpace(*in.Name); arg.Name == "" {
			return nil, ErrValidation
		}
	}
	if in.IsFilterable != nil {
		arg.IsFilterable = *in.IsFilterable
	}
	if in.SortOrder != nil {
		arg.SortOrder = *in.SortOrder
	}
	a, err := s.repo.Update(ctx, arg)
	if err != nil {
		return nil, mapNotFound(err)
	}
	return a, nil
}

// DeleteDefinition removes a definition (values and assignments cascade, per spec).
func (s *AttributeAdminService) DeleteDefinition(ctx context.Context, id uuid.UUID) error {
	ok, err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

// ListValues returns all values of a definition.
func (s *AttributeAdminService) ListValues(ctx context.Context, attrID uuid.UUID) ([]dbsqlc.ProductAttributeValue, error) {
	if _, err := s.repo.GetByID(ctx, attrID); err != nil {
		return nil, mapNotFound(err)
	}
	return s.repo.ListValues(ctx, attrID)
}

// CreateValue adds a value; a duplicate slug yields ErrConflict.
func (s *AttributeAdminService) CreateValue(ctx context.Context, attrID uuid.UUID, in AttributeValueInput) (*dbsqlc.ProductAttributeValue, error) {
	if _, err := s.repo.GetByID(ctx, attrID); err != nil {
		return nil, mapNotFound(err)
	}
	if in.Value == nil || strings.TrimSpace(*in.Value) == "" {
		return nil, ErrValidation
	}
	label := strings.TrimSpace(*in.Value)
	slug := Slugify(label)
	if in.ValueSlug != nil && *in.ValueSlug != "" {
		slug = Slugify(*in.ValueSlug)
	}
	if slug == "" {
		return nil, ErrValidation
	}
	var sort int32
	if in.SortOrder != nil {
		sort = *in.SortOrder
	}
	v, err := s.repo.CreateValue(ctx, dbsqlc.CreateAttributeValueParams{AttributeID: attrID, Value: label, Slug: slug, SortOrder: sort})
	if err != nil {
		if repositories.IsUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	return v, nil
}

// UpdateValue applies a partial update to a value of the given definition.
func (s *AttributeAdminService) UpdateValue(ctx context.Context, attrID, valueID uuid.UUID, in AttributeValueInput) (*dbsqlc.ProductAttributeValue, error) {
	cur, err := s.findValue(ctx, attrID, valueID)
	if err != nil {
		return nil, err
	}
	arg := dbsqlc.UpdateAttributeValueParams{ID: valueID, AttributeID: attrID, Value: cur.Value, Slug: cur.Slug, SortOrder: cur.SortOrder}
	if in.Value != nil {
		if arg.Value = strings.TrimSpace(*in.Value); arg.Value == "" {
			return nil, ErrValidation
		}
	}
	if in.ValueSlug != nil {
		if arg.Slug = Slugify(*in.ValueSlug); arg.Slug == "" {
			return nil, ErrValidation
		}
	}
	if in.SortOrder != nil {
		arg.SortOrder = *in.SortOrder
	}
	v, err := s.repo.UpdateValue(ctx, arg)
	if err != nil {
		if repositories.IsUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, mapNotFound(err)
	}
	return v, nil
}

// DeleteValue removes a value unless products still use it (ErrConflict).
func (s *AttributeAdminService) DeleteValue(ctx context.Context, attrID, valueID uuid.UUID) error {
	if _, err := s.findValue(ctx, attrID, valueID); err != nil {
		return err
	}
	n, err := s.repo.CountAssignments(ctx, valueID)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrConflict
	}
	ok, err := s.repo.DeleteValue(ctx, valueID, attrID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

func (s *AttributeAdminService) findValue(ctx context.Context, attrID, valueID uuid.UUID) (*dbsqlc.ProductAttributeValue, error) {
	vals, err := s.ListValues(ctx, attrID)
	if err != nil {
		return nil, err
	}
	for i := range vals {
		if vals[i].ID == valueID {
			return &vals[i], nil
		}
	}
	return nil, ErrNotFound
}
