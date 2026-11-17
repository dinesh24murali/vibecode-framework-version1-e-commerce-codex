package contract

import (
	"context"
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/handlers"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

type fakeProducts struct {
	product dbsqlc.GetActiveProductByIDRow
	lastF   repositories.CatalogFilter
	total   int64
}

func (f *fakeProducts) Count(_ context.Context, fl repositories.CatalogFilter) (int64, error) {
	f.lastF = fl
	return f.total, nil
}

func (f *fakeProducts) List(context.Context, repositories.CatalogFilter, int32, int32) ([]dbsqlc.ListCatalogProductsRow, error) {
	p := f.product
	return []dbsqlc.ListCatalogProductsRow{{
		ID: p.ID, CategoryID: p.CategoryID, Slug: p.Slug, Title: p.Title, PriceInr: p.PriceInr,
		DiscountPercent: p.DiscountPercent, IsActive: true, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}}, nil
}

func (f *fakeProducts) GetActive(_ context.Context, id uuid.UUID) (*dbsqlc.GetActiveProductByIDRow, error) {
	if id == f.product.ID {
		return &f.product, nil
	}
	return nil, sql.ErrNoRows
}

func (f *fakeProducts) Images(context.Context, []uuid.UUID) ([]dbsqlc.ProductImage, error) {
	return []dbsqlc.ProductImage{
		{ID: uuid.New(), ProductID: f.product.ID, ImageUrl: "https://cdn/x.jpg", SortOrder: 1},
	}, nil
}

func (f *fakeProducts) Attributes(context.Context, []uuid.UUID) ([]dbsqlc.ListAttributesByProductsRow, error) {
	return []dbsqlc.ListAttributesByProductsRow{{
		ProductID: f.product.ID, AttributeValueID: uuid.New(), AttributeID: uuid.New(), Code: "level",
		AttributeName: "Level", Value: "Beginner", ValueSlug: "beginner", IsFilterable: true,
	}}, nil
}

func (f *fakeProducts) Inventory(context.Context, []uuid.UUID) ([]dbsqlc.InventoryItem, error) {
	return []dbsqlc.InventoryItem{{ProductID: f.product.ID, OnHand: 10, Reserved: 3}}, nil
}

func (f *fakeProducts) ActiveCategorySlugs(context.Context) ([]string, error) {
	return []string{"journals"}, nil
}

type fakeCategories struct{ rows []dbsqlc.Category }

func (f *fakeCategories) List(context.Context) ([]dbsqlc.Category, error) { return f.rows, nil }
func (f *fakeCategories) GetByID(_ context.Context, id uuid.UUID) (*dbsqlc.Category, error) {
	for i := range f.rows {
		if f.rows[i].ID == id {
			return &f.rows[i], nil
		}
	}
	return nil, sql.ErrNoRows
}

func catalogServer() (*gin.Engine, *fakeProducts, *fakeCategories) {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	cat := dbsqlc.Category{ID: uuid.New(), Slug: "journals", Name: "Journals", CreatedAt: now, UpdatedAt: now}
	fp := &fakeProducts{total: 45, product: dbsqlc.GetActiveProductByIDRow{
		ID: uuid.New(), CategoryID: cat.ID, Slug: "mandala-journal", Title: "Mandala Journal",
		Description: sql.NullString{String: "Nice", Valid: true}, PriceInr: "499.00",
		DiscountPercent: sql.NullInt32{Int32: 10, Valid: true}, IsActive: true, CreatedAt: now, UpdatedAt: now,
	}}
	fc := &fakeCategories{rows: []dbsqlc.Category{cat}}
	attrSvc := services.NewAttributeService(&fakeAttrs{values: map[uuid.UUID][]dbsqlc.ProductAttributeValue{}})
	svc := services.NewProductService(fp, fc, attrSvc)
	ph, ch := handlers.NewProductHandler(svc), handlers.NewCategoryHandler(svc)
	r := gin.New() // public: no auth middleware
	r.GET("/api/v1/products", ph.List)
	r.GET("/api/v1/products/:id", ph.Get)
	r.GET("/api/v1/categories", ch.List)
	r.GET("/api/v1/categories/:id", ch.Get)
	return r, fp, fc
}

func TestProducts_List(t *testing.T) {
	r, fp, _ := catalogServer()
	attr1, attr2 := uuid.NewString(), uuid.NewString()
	code, body := getJSON(r, "/api/v1/products?q=mandala&category[]=journals&category[]=planners&attr_value_id[]="+attr1+"&attr_value_id[]="+attr2+"&sort_by=price_asc&limit=20&page=2")
	if code != http.StatusOK {
		t.Fatalf("status %d %v", code, body)
	}
	requireKeys(t, body, "data", "meta", "filters")
	meta := obj(body["meta"])
	requireKeys(t, meta, "total", "page", "limit", "total_pages")
	if meta["total"] != float64(45) || meta["page"] != float64(2) || meta["total_pages"] != float64(3) {
		t.Fatalf("meta %v", meta)
	}
	f := fp.lastF
	if f.Query != "mandala" || f.Sort != "price_asc" || len(f.CategorySlugs) != 2 || len(f.ValueIDs) != 2 {
		t.Fatalf("filters not passed through: %+v", f)
	}

	item := obj(arr(body["data"])[0])
	requireKeys(t, item, "id", "name", "slug", "price", "discount_percent", "discounted_price", "attributes", "cover_image", "categories", "is_active")
	if item["discounted_price"] != 449.1 || obj(item["cover_image"])["url"] != "https://cdn/x.jpg" {
		t.Fatalf("item %v", item)
	}
	filters := obj(body["filters"])
	requireKeys(t, filters, "categories", "attributes", "sort_options")
	if len(arr(filters["sort_options"])) != 4 {
		t.Fatalf("sort options %v", filters["sort_options"])
	}
}

func TestProducts_ListValidation(t *testing.T) {
	r, _, _ := catalogServer()
	for _, q := range []string{"sort_by=bogus", "page=0", "limit=101", "limit=abc", "attr_value_id[]=nope"} {
		if code, body := getJSON(r, "/api/v1/products?"+q); code != http.StatusBadRequest || body["message"] == nil {
			t.Errorf("%s: status %d body %v", q, code, body)
		}
	}
}

func TestProducts_Get(t *testing.T) {
	r, fp, _ := catalogServer()
	code, body := getJSON(r, "/api/v1/products/"+fp.product.ID.String())
	if code != http.StatusOK {
		t.Fatalf("status %d %v", code, body)
	}
	requireKeys(t, body, "id", "name", "slug", "description", "price", "discount_percent", "discounted_price",
		"stock", "attributes", "images", "categories", "is_active", "created_at", "updated_at")
	if body["stock"] != float64(7) {
		t.Fatalf("stock must be on_hand - reserved: %v", body["stock"])
	}
	if code, _ = getJSON(r, "/api/v1/products/"+uuid.NewString()); code != http.StatusNotFound {
		t.Fatalf("unknown product: %d", code)
	}
	if code, _ = getJSON(r, "/api/v1/products/not-a-uuid"); code != http.StatusBadRequest {
		t.Fatalf("bad id: %d", code)
	}
}

func TestCategories(t *testing.T) {
	r, _, fc := catalogServer()
	code, body := getJSON(r, "/api/v1/categories")
	if code != http.StatusOK || len(arr(body["data"])) != 1 {
		t.Fatalf("list: %d %v", code, body)
	}
	requireKeys(t, obj(arr(body["data"])[0]), "id", "name", "slug", "image_url", "created_at", "updated_at")

	code, body = getJSON(r, "/api/v1/categories/"+fc.rows[0].ID.String())
	if code != http.StatusOK || body["slug"] != "journals" {
		t.Fatalf("get: %d %v", code, body)
	}
	if code, _ = getJSON(r, "/api/v1/categories/"+uuid.NewString()); code != http.StatusNotFound {
		t.Fatalf("unknown: %d", code)
	}
	if code, _ = getJSON(r, "/api/v1/categories/x"); code != http.StatusBadRequest {
		t.Fatalf("bad id: %d", code)
	}
}
