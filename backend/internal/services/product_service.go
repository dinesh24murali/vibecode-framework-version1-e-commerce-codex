package services

import (
	"context"
	"errors"
	"math"
	"strconv"

	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/models"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// SortOptions are the accepted sort_by values.
var SortOptions = []string{"price_asc", "price_desc", "newest", "popular"}

// ErrInvalidSort is returned for an unknown sort_by value.
var ErrInvalidSort = errors.New("invalid sort_by")

// ProductQuery is the validated input of the catalogue list endpoint.
type ProductQuery struct {
	Query         string
	Sort          string
	CategorySlugs []string
	ValueIDs      []uuid.UUID
	Page          int
	Limit         int
}

// ProductService serves the public product and category reads.
type ProductService struct {
	products   repositories.ProductRepository
	categories repositories.CategoryRepository
	attributes *AttributeService
}

// NewProductService constructs a ProductService.
func NewProductService(p repositories.ProductRepository, c repositories.CategoryRepository, a *AttributeService) *ProductService {
	return &ProductService{products: p, categories: c, attributes: a}
}

func price(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func round2(f float64) float64 { return math.Round(f*100) / 100 }

func toCategory(c dbsqlc.Category) models.Category {
	var img *string
	if c.ImageUrl.Valid {
		img = &c.ImageUrl.String
	}
	return models.Category{ID: c.ID, Name: c.Name, Slug: c.Slug, ImageURL: img, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}
}

// ListCategories returns all categories ordered by name.
func (s *ProductService) ListCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := s.categories.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]models.Category, 0, len(rows))
	for _, c := range rows {
		out = append(out, toCategory(c))
	}
	return out, nil
}

// GetCategory returns one category or ErrNotFound.
func (s *ProductService) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	c, err := s.categories.GetByID(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	m := toCategory(*c)
	return &m, nil
}

// ListProducts returns one page of active products matching the filters.
func (s *ProductService) ListProducts(ctx context.Context, q ProductQuery) (*models.ProductListResponse, error) {
	valid := false
	for _, o := range SortOptions {
		valid = valid || o == q.Sort
	}
	if !valid {
		return nil, ErrInvalidSort
	}
	f := repositories.CatalogFilter{Query: q.Query, Sort: q.Sort, CategorySlugs: q.CategorySlugs, ValueIDs: q.ValueIDs}
	total, err := s.products.Count(ctx, f)
	if err != nil {
		return nil, err
	}
	rows, err := s.products.List(ctx, f, int32(q.Limit), int32((q.Page-1)*q.Limit)) //nolint:gosec // bounded by handler
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ID)
	}
	rel, err := s.related(ctx, ids)
	if err != nil {
		return nil, err
	}

	data := summarise(rows, rel)

	filters, err := s.filters(ctx)
	if err != nil {
		return nil, err
	}
	return &models.ProductListResponse{
		Data: data,
		Meta: models.PaginationMeta{
			Total: total, Page: q.Page, Limit: q.Limit,
			TotalPages: int(math.Ceil(float64(total) / float64(q.Limit))),
		},
		Filters: *filters,
	}, nil
}

// GetProduct returns one active product or ErrNotFound.
func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	r, err := s.products.GetActive(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	rel, err := s.related(ctx, []uuid.UUID{id})
	if err != nil {
		return nil, err
	}
	disc := float64(r.DiscountPercent.Int32)
	p := price(r.PriceInr)
	out := &models.Product{
		ID: r.ID, Name: r.Title, Slug: r.Slug, Description: r.Description.String, Price: p,
		DiscountPercent: disc, DiscountedPrice: round2(p * (100 - disc) / 100),
		Stock: rel.stock[id], Attributes: rel.attrs[id], Images: rel.images[id],
		Categories: rel.categoryFor(r.CategoryID), IsActive: r.IsActive,
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
	if out.Attributes == nil {
		out.Attributes = []models.ProductAttribute{}
	}
	if out.Images == nil {
		out.Images = []models.ProductImage{}
	}
	return out, nil
}

type relations struct {
	attrs      map[uuid.UUID][]models.ProductAttribute
	images     map[uuid.UUID][]models.ProductImage
	stock      map[uuid.UUID]int32
	categories map[uuid.UUID]models.Category
}

func (r *relations) categoryFor(id uuid.UUID) []models.Category {
	if c, ok := r.categories[id]; ok {
		return []models.Category{c}
	}
	return []models.Category{}
}

func (s *ProductService) related(ctx context.Context, ids []uuid.UUID) (*relations, error) {
	rel := &relations{
		attrs: map[uuid.UUID][]models.ProductAttribute{}, images: map[uuid.UUID][]models.ProductImage{},
		stock: map[uuid.UUID]int32{}, categories: map[uuid.UUID]models.Category{},
	}
	attrs, err := s.products.Attributes(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, a := range attrs {
		rel.attrs[a.ProductID] = append(rel.attrs[a.ProductID], models.ProductAttribute{
			AttributeValueID: a.AttributeValueID, AttributeID: a.AttributeID, Code: a.Code,
			AttributeName: a.AttributeName, Value: a.Value, ValueSlug: a.ValueSlug, IsFilterable: a.IsFilterable,
		})
	}
	imgs, err := s.products.Images(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, i := range imgs {
		var alt *string
		if i.AltText.Valid {
			alt = &i.AltText.String
		}
		rel.images[i.ProductID] = append(rel.images[i.ProductID], models.ProductImage{
			ID: i.ID, URL: i.ImageUrl, AltText: alt, SortOrder: i.SortOrder,
		})
	}
	inv, err := s.products.Inventory(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, i := range inv {
		rel.stock[i.ProductID] = max(i.OnHand-i.Reserved, 0)
	}
	cats, err := s.categories.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, c := range cats {
		rel.categories[c.ID] = toCategory(c)
	}
	return rel, nil
}

func (s *ProductService) filters(ctx context.Context) (*models.ProductFilters, error) {
	slugs, err := s.products.ActiveCategorySlugs(ctx)
	if err != nil {
		return nil, err
	}
	attrs, err := s.attributes.ListFilterable(ctx)
	if err != nil {
		return nil, err
	}
	groups := make([]models.AttributeFilterGroup, 0, len(attrs))
	for _, a := range attrs {
		vals := make([]models.FilterValue, 0, len(a.Values))
		for _, v := range a.Values {
			vals = append(vals, models.FilterValue{ID: v.ID, ValueSlug: v.ValueSlug, Value: v.Value})
		}
		groups = append(groups, models.AttributeFilterGroup{
			AttributeID: a.ID, Code: a.Code, AttributeName: a.Name, IsFilterable: a.IsFilterable, Values: vals,
		})
	}
	if slugs == nil {
		slugs = []string{}
	}
	return &models.ProductFilters{Categories: slugs, Attributes: groups, SortOptions: SortOptions}, nil
}

func summarise(rows []dbsqlc.ListCatalogProductsRow, rel *relations) []models.ProductSummary {
	data := make([]models.ProductSummary, 0, len(rows))
	for _, r := range rows {
		disc := float64(r.DiscountPercent.Int32)
		p := price(r.PriceInr)
		sum := models.ProductSummary{
			ID: r.ID, Name: r.Title, Slug: r.Slug, Price: p, DiscountPercent: disc,
			DiscountedPrice: round2(p * (100 - disc) / 100), Attributes: rel.attrs[r.ID],
			Categories: rel.categoryFor(r.CategoryID), IsActive: r.IsActive,
		}
		if imgs := rel.images[r.ID]; len(imgs) > 0 {
			sum.CoverImage = &imgs[0]
		}
		if sum.Attributes == nil {
			sum.Attributes = []models.ProductAttribute{}
		}
		data = append(data, sum)
	}
	return data
}

// SummariesByIDs returns summaries of the active products among ids, keyed by id.
// Inactive or missing products are omitted.
func (s *ProductService) SummariesByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]models.ProductSummary, error) {
	out := make(map[uuid.UUID]models.ProductSummary, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	rows := make([]dbsqlc.GetActiveProductByIDRow, 0, len(ids))
	found := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		r, err := s.products.GetActive(ctx, id)
		if err != nil {
			if errors.Is(mapNotFound(err), ErrNotFound) {
				continue
			}
			return nil, err
		}
		rows = append(rows, *r)
		found = append(found, id)
	}
	rel, err := s.related(ctx, found)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		disc := float64(r.DiscountPercent.Int32)
		p := price(r.PriceInr)
		sum := models.ProductSummary{
			ID: r.ID, Name: r.Title, Slug: r.Slug, Price: p, DiscountPercent: disc,
			DiscountedPrice: round2(p * (100 - disc) / 100), Attributes: rel.attrs[r.ID],
			Categories: rel.categoryFor(r.CategoryID), IsActive: r.IsActive,
		}
		if imgs := rel.images[r.ID]; len(imgs) > 0 {
			sum.CoverImage = &imgs[0]
		}
		if sum.Attributes == nil {
			sum.Attributes = []models.ProductAttribute{}
		}
		out[r.ID] = sum
	}
	return out, nil
}

// Featured returns the n most recent active products (the schema has no featured flag).
func (s *ProductService) Featured(ctx context.Context, n int) ([]models.ProductSummary, error) {
	rows, err := s.products.List(ctx, repositories.CatalogFilter{Sort: "newest"}, int32(n), 0) //nolint:gosec // small constant
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ID)
	}
	rel, err := s.related(ctx, ids)
	if err != nil {
		return nil, err
	}
	return summarise(rows, rel), nil
}
