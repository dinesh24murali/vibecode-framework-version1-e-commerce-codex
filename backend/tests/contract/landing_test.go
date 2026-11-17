package contract

import (
	"context"
	"database/sql"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/handlers"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

type fakeCarousel struct{ rows []dbsqlc.CarouselSlide }

func (f *fakeCarousel) ListActive(context.Context) ([]dbsqlc.CarouselSlide, error) {
	return f.rows, nil
}

func landingServer(slides []dbsqlc.CarouselSlide) *gin.Engine {
	r, fp, fc := catalogServer() // reuse the catalogue fakes (public, no auth)
	_ = r
	attrSvc := services.NewAttributeService(&fakeAttrs{values: map[uuid.UUID][]dbsqlc.ProductAttributeValue{}})
	productSvc := services.NewProductService(fp, fc, attrSvc)
	h := handlers.NewLandingHandler(services.NewLandingService(&fakeCarousel{rows: slides}, productSvc))
	e := gin.New()
	e.GET("/api/v1/landing", h.Get)
	return e
}

func TestLanding(t *testing.T) {
	slides := []dbsqlc.CarouselSlide{{
		ID: uuid.New(), Title: "Discover", ImageUrl: "https://cdn/s1.jpg", SortOrder: 1,
		Subtitle: sql.NullString{String: "New", Valid: true}, CtaLabel: sql.NullString{String: "Shop", Valid: true},
		CtaUrl: sql.NullString{String: "/products", Valid: true},
	}, {ID: uuid.New(), Title: "Second", ImageUrl: "https://cdn/s2.jpg", SortOrder: 2}}

	code, body := getJSON(landingServer(slides), "/api/v1/landing")
	if code != http.StatusOK {
		t.Fatalf("status %d %v", code, body)
	}
	requireKeys(t, body, "carousel_slides", "featured_products")
	got := arr(body["carousel_slides"])
	if len(got) != 2 {
		t.Fatalf("slides %v", got)
	}
	requireKeys(t, obj(got[0]), "id", "title", "subtitle", "image_url", "cta_label", "cta_url", "sort_order")
	if s2 := obj(got[1]); s2["subtitle"] != nil || s2["cta_url"] != nil {
		t.Fatalf("optional fields must be null: %v", s2)
	}
	feat := arr(body["featured_products"])
	if len(feat) != 1 {
		t.Fatalf("featured %v", feat)
	}
	requireKeys(t, obj(feat[0]), "id", "name", "slug", "price", "discount_percent", "discounted_price", "attributes", "cover_image", "categories", "is_active")
}

func TestLanding_NoSlidesStillOK(t *testing.T) {
	code, body := getJSON(landingServer(nil), "/api/v1/landing")
	if code != http.StatusOK || body["carousel_slides"] == nil || len(arr(body["carousel_slides"])) != 0 {
		t.Fatalf("empty slides must be [] not null: %d %v", code, body)
	}
}
