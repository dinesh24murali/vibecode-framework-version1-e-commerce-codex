package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// ProductImageInput is one image to persist for a product.
type ProductImageInput struct {
	AltText   *string
	URL       string
	SortOrder int32
}

// ProductWrite is the full set of columns written on create/update.
// Stock and Images are only written when non-nil.
type ProductWrite struct { //nolint:govet // readability over packing
	Stock       *int32
	Images      *[]ProductImageInput
	CategoryID  uuid.UUID
	Slug        string
	Title       string
	Description string
	PriceInr    string
	Discount    int32 // 0 means no discount (stored as NULL)
	IsActive    bool
}

// AdminProductRepository is the write-side persistence interface for products.
type AdminProductRepository interface {
	List(ctx context.Context, query string, limit, offset int32) ([]dbsqlc.AdminListProductsRow, int64, error)
	Get(ctx context.Context, id uuid.UUID) (*dbsqlc.GetProductByIDRow, error)
	Create(ctx context.Context, w ProductWrite) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, w ProductWrite) error
	SetActive(ctx context.Context, id uuid.UUID, active bool) (bool, error)
	AssignValue(ctx context.Context, productID, valueID uuid.UUID) error
	RemoveValue(ctx context.Context, productID, valueID uuid.UUID) (bool, error)
	ValueExists(ctx context.Context, valueID uuid.UUID) (bool, error)
	CategoryExists(ctx context.Context, id uuid.UUID) (bool, error)
}

type pgAdminProductRepository struct {
	db *sql.DB
	q  *dbsqlc.Queries
}

// NewAdminProductRepository returns a Postgres-backed AdminProductRepository.
func NewAdminProductRepository(pool *pgxpool.Pool) AdminProductRepository {
	db := stdlib.OpenDBFromPool(pool)
	return &pgAdminProductRepository{db: db, q: dbsqlc.New(db)}
}

func (r *pgAdminProductRepository) List(ctx context.Context, query string, limit, offset int32) ([]dbsqlc.AdminListProductsRow, int64, error) {
	total, err := r.q.AdminCountProducts(ctx, nullQuery(query))
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.q.AdminListProducts(ctx, dbsqlc.AdminListProductsParams{Query: nullQuery(query), PageLimit: limit, PageOffset: offset})
	return rows, total, err
}

func (r *pgAdminProductRepository) Get(ctx context.Context, id uuid.UUID) (*dbsqlc.GetProductByIDRow, error) {
	p, err := r.q.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func discount(d int32) sql.NullInt32 { return sql.NullInt32{Int32: d, Valid: d > 0} }

func (r *pgAdminProductRepository) Create(ctx context.Context, w ProductWrite) (uuid.UUID, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}
	defer func() { _ = tx.Rollback() }()
	q := r.q.WithTx(tx)
	p, err := q.CreateProduct(ctx, dbsqlc.CreateProductParams{
		CategoryID: w.CategoryID, Sku: w.Slug, Slug: w.Slug, Title: w.Title,
		Description: sql.NullString{String: w.Description, Valid: true}, PriceInr: w.PriceInr,
		DiscountPercent: discount(w.Discount), CurrencyCode: "INR", IsActive: w.IsActive,
	})
	if err != nil {
		return uuid.Nil, err
	}
	if err := writeRelated(ctx, q, p.ID, w); err != nil {
		return uuid.Nil, err
	}
	return p.ID, tx.Commit()
}

func (r *pgAdminProductRepository) Update(ctx context.Context, id uuid.UUID, w ProductWrite) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	q := r.q.WithTx(tx)
	if _, err := q.AdminUpdateProduct(ctx, dbsqlc.AdminUpdateProductParams{
		ID: id, CategoryID: w.CategoryID, Slug: w.Slug, Title: w.Title,
		Description: sql.NullString{String: w.Description, Valid: true}, PriceInr: w.PriceInr,
		DiscountPercent: discount(w.Discount), IsActive: w.IsActive,
	}); err != nil {
		return err
	}
	if err := writeRelated(ctx, q, id, w); err != nil {
		return err
	}
	return tx.Commit()
}

func writeRelated(ctx context.Context, q *dbsqlc.Queries, id uuid.UUID, w ProductWrite) error {
	if w.Stock != nil {
		if err := q.UpsertInventoryOnHand(ctx, dbsqlc.UpsertInventoryOnHandParams{ProductID: id, OnHand: *w.Stock}); err != nil {
			return err
		}
	}
	if w.Images == nil {
		return nil
	}
	if err := q.DeleteImagesByProductID(ctx, id); err != nil {
		return err
	}
	for _, img := range *w.Images {
		alt := sql.NullString{}
		if img.AltText != nil {
			alt = sql.NullString{String: *img.AltText, Valid: true}
		}
		if err := q.InsertProductImage(ctx, dbsqlc.InsertProductImageParams{
			ProductID: id, ImageUrl: img.URL, AltText: alt, SortOrder: img.SortOrder,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *pgAdminProductRepository) SetActive(ctx context.Context, id uuid.UUID, active bool) (bool, error) {
	n, err := r.q.SetProductActive(ctx, dbsqlc.SetProductActiveParams{ID: id, IsActive: active})
	return n > 0, err
}

func (r *pgAdminProductRepository) AssignValue(ctx context.Context, productID, valueID uuid.UUID) error {
	return r.q.AssignProductAttributeValue(ctx, dbsqlc.AssignProductAttributeValueParams{ProductID: productID, AttributeValueID: valueID})
}

func (r *pgAdminProductRepository) RemoveValue(ctx context.Context, productID, valueID uuid.UUID) (bool, error) {
	n, err := r.q.RemoveProductAttributeValue(ctx, dbsqlc.RemoveProductAttributeValueParams{ProductID: productID, AttributeValueID: valueID})
	return n > 0, err
}

func (r *pgAdminProductRepository) ValueExists(ctx context.Context, valueID uuid.UUID) (bool, error) {
	_, err := r.q.GetAttributeValueByID(ctx, valueID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *pgAdminProductRepository) CategoryExists(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := r.q.GetCategoryByID(ctx, id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}
