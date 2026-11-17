package services

import (
	"context"
	"database/sql"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/models"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

var productSlugRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// AdminImageInput is one image of a product write request.
type AdminImageInput struct {
	AltText   *string
	URL       string
	SortOrder int32
}

// AdminProductInput carries create/update fields; nil means "not provided".
type AdminProductInput struct {
	Name            *string
	Slug            *string
	Description     *string
	Price           *float64
	DiscountPercent *float64
	Stock           *int32
	IsActive        *bool
	CategoryIDs     *[]uuid.UUID
	Images          *[]AdminImageInput
}

// AdminProductService implements admin product management.
// Persistence note: a product has exactly one category column, so only the first
// of category_ids is stored; sku mirrors the slug; stock is inventory on_hand.
type AdminProductService struct {
	repo repositories.AdminProductRepository
	pub  *ProductService
}

// NewAdminProductService constructs an AdminProductService.
func NewAdminProductService(repo repositories.AdminProductRepository, pub *ProductService) *AdminProductService {
	return &AdminProductService{repo: repo, pub: pub}
}

// build assembles the spec Product (admin view: stock is on_hand, inactive allowed).
func (s *AdminProductService) build(ctx context.Context, rows []productCore) ([]models.Product, error) {
	ids := make([]uuid.UUID, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.id)
	}
	rel, err := s.pub.related(ctx, ids)
	if err != nil {
		return nil, err
	}
	inv, err := s.pub.products.Inventory(ctx, ids)
	if err != nil {
		return nil, err
	}
	onHand := map[uuid.UUID]int32{}
	for _, i := range inv {
		onHand[i.ProductID] = i.OnHand
	}
	out := make([]models.Product, 0, len(rows))
	for _, r := range rows {
		disc := float64(r.discount.Int32)
		p := price(r.priceInr)
		m := models.Product{
			ID: r.id, Name: r.title, Slug: r.slug, Description: r.description.String, Price: p,
			DiscountPercent: disc, DiscountedPrice: round2(p * (100 - disc) / 100),
			Stock: onHand[r.id], Attributes: rel.attrs[r.id], Images: rel.images[r.id],
			Categories: rel.categoryFor(r.categoryID), IsActive: r.active, CreatedAt: r.created, UpdatedAt: r.updated,
		}
		if m.Attributes == nil {
			m.Attributes = []models.ProductAttribute{}
		}
		if m.Images == nil {
			m.Images = []models.ProductImage{}
		}
		out = append(out, m)
	}
	return out, nil
}

type productCore struct { //nolint:govet // readability over packing
	created, updated time.Time
	description      sql.NullString
	discount         sql.NullInt32
	priceInr         string
	title, slug      string
	id, categoryID   uuid.UUID
	active           bool
}

// List returns one page of products (including inactive) matching the search term.
func (s *AdminProductService) List(ctx context.Context, q string, page, limit int) ([]models.Product, int64, error) {
	rows, total, err := s.repo.List(ctx, strings.TrimSpace(q), int32(limit), int32((page-1)*limit)) //nolint:gosec // bounded by handler
	if err != nil {
		return nil, 0, err
	}
	cores := make([]productCore, 0, len(rows))
	for _, r := range rows {
		cores = append(cores, productCore{
			created: r.CreatedAt, updated: r.UpdatedAt, description: r.Description, discount: r.DiscountPercent,
			priceInr: r.PriceInr, title: r.Title, slug: r.Slug, id: r.ID, categoryID: r.CategoryID, active: r.IsActive,
		})
	}
	out, err := s.build(ctx, cores)
	return out, total, err
}

// Get returns one product (active or not) or ErrNotFound.
func (s *AdminProductService) Get(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	r, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	out, err := s.build(ctx, []productCore{{
		created: r.CreatedAt, updated: r.UpdatedAt, description: r.Description, discount: r.DiscountPercent,
		priceInr: r.PriceInr, title: r.Title, slug: r.Slug, id: r.ID, categoryID: r.CategoryID, active: r.IsActive,
	}})
	if err != nil {
		return nil, err
	}
	return &out[0], nil
}

func (s *AdminProductService) validate(ctx context.Context, in AdminProductInput, w *repositories.ProductWrite) error {
	if in.Name != nil {
		if w.Title = strings.TrimSpace(*in.Name); w.Title == "" {
			return ErrValidation
		}
	}
	if in.Slug != nil {
		if w.Slug = strings.TrimSpace(*in.Slug); !productSlugRe.MatchString(w.Slug) {
			return ErrValidation
		}
	}
	if in.Description != nil {
		w.Description = *in.Description
	}
	if in.Price != nil {
		if *in.Price < 0 || math.IsNaN(*in.Price) {
			return ErrValidation
		}
		w.PriceInr = strconv.FormatFloat(*in.Price, 'f', 2, 64)
	}
	if in.DiscountPercent != nil {
		if *in.DiscountPercent < 0 || *in.DiscountPercent > 100 {
			return ErrValidation
		}
		w.Discount = int32(math.Round(*in.DiscountPercent))
	}
	if in.Stock != nil {
		if *in.Stock < 0 {
			return ErrValidation
		}
		w.Stock = in.Stock
	}
	if in.IsActive != nil {
		w.IsActive = *in.IsActive
	}
	if in.CategoryIDs != nil {
		if len(*in.CategoryIDs) == 0 {
			return ErrValidation
		}
		ok, err := s.repo.CategoryExists(ctx, (*in.CategoryIDs)[0])
		if err != nil {
			return err
		}
		if !ok {
			return ErrValidation
		}
		w.CategoryID = (*in.CategoryIDs)[0]
	}
	if in.Images != nil {
		imgs := make([]repositories.ProductImageInput, 0, len(*in.Images))
		for _, i := range *in.Images {
			if strings.TrimSpace(i.URL) == "" {
				return ErrValidation
			}
			imgs = append(imgs, repositories.ProductImageInput{URL: i.URL, AltText: i.AltText, SortOrder: i.SortOrder})
		}
		w.Images = &imgs
	}
	return nil
}

// Create validates and creates a product; a duplicate slug yields ErrConflict.
func (s *AdminProductService) Create(ctx context.Context, in AdminProductInput) (*models.Product, error) {
	if in.Name == nil || in.Slug == nil || in.Price == nil || in.CategoryIDs == nil {
		return nil, ErrValidation
	}
	w := repositories.ProductWrite{IsActive: true}
	if err := s.validate(ctx, in, &w); err != nil {
		return nil, err
	}
	if w.Stock == nil {
		zero := int32(0)
		w.Stock = &zero
	}
	id, err := s.repo.Create(ctx, w)
	if err != nil {
		if repositories.IsUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	return s.Get(ctx, id)
}

// Update applies a partial update.
func (s *AdminProductService) Update(ctx context.Context, id uuid.UUID, in AdminProductInput) (*models.Product, error) {
	cur, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	w := repositories.ProductWrite{
		CategoryID: cur.CategoryID, Slug: cur.Slug, Title: cur.Title, Description: cur.Description.String,
		PriceInr: cur.PriceInr, Discount: cur.DiscountPercent.Int32, IsActive: cur.IsActive,
	}
	if err := s.validate(ctx, in, &w); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, id, w); err != nil {
		if repositories.IsUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	return s.Get(ctx, id)
}

// Delete soft-deletes a product (is_active = false); it disappears from public endpoints.
func (s *AdminProductService) Delete(ctx context.Context, id uuid.UUID) error {
	ok, err := s.repo.SetActive(ctx, id, false)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

// ListAttributes returns the attribute values assigned to a product.
func (s *AdminProductService) ListAttributes(ctx context.Context, id uuid.UUID) ([]models.ProductAttribute, error) {
	if _, err := s.repo.Get(ctx, id); err != nil {
		return nil, mapNotFound(err)
	}
	rel, err := s.pub.related(ctx, []uuid.UUID{id})
	if err != nil {
		return nil, err
	}
	if rel.attrs[id] == nil {
		return []models.ProductAttribute{}, nil
	}
	return rel.attrs[id], nil
}

// AssignAttribute links an existing attribute value to a product (idempotent).
func (s *AdminProductService) AssignAttribute(ctx context.Context, id, valueID uuid.UUID) (*models.ProductAttribute, error) {
	if _, err := s.repo.Get(ctx, id); err != nil {
		return nil, mapNotFound(err)
	}
	ok, err := s.repo.ValueExists(ctx, valueID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotFound
	}
	if err = s.repo.AssignValue(ctx, id, valueID); err != nil {
		return nil, err
	}
	attrs, err := s.ListAttributes(ctx, id)
	if err != nil {
		return nil, err
	}
	for i := range attrs {
		if attrs[i].AttributeValueID == valueID {
			return &attrs[i], nil
		}
	}
	return nil, ErrNotFound
}

// RemoveAttribute unlinks one attribute value from a product.
func (s *AdminProductService) RemoveAttribute(ctx context.Context, id, valueID uuid.UUID) error {
	ok, err := s.repo.RemoveValue(ctx, id, valueID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}
