package contract

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/handlers"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

type fakeAttrs struct {
	values map[uuid.UUID][]dbsqlc.ProductAttributeValue
	attrs  []dbsqlc.ProductAttribute
}

func (f *fakeAttrs) ListFilterable(context.Context) ([]dbsqlc.ProductAttribute, error) {
	var out []dbsqlc.ProductAttribute
	for _, a := range f.attrs {
		if a.IsFilterable {
			out = append(out, a)
		}
	}
	return out, nil
}

func (f *fakeAttrs) GetByID(_ context.Context, id uuid.UUID) (*dbsqlc.ProductAttribute, error) {
	for i := range f.attrs {
		if f.attrs[i].ID == id {
			return &f.attrs[i], nil
		}
	}
	return nil, sql.ErrNoRows
}

func (f *fakeAttrs) ListValues(_ context.Context, id uuid.UUID) ([]dbsqlc.ProductAttributeValue, error) {
	return f.values[id], nil
}

func attrServer() (*gin.Engine, uuid.UUID, uuid.UUID) {
	gin.SetMode(gin.TestMode)
	level, hidden := uuid.New(), uuid.New()
	now := time.Now()
	f := &fakeAttrs{
		attrs: []dbsqlc.ProductAttribute{
			{ID: level, Code: "level", Name: "Skill Level", IsFilterable: true, SortOrder: 1, CreatedAt: now, UpdatedAt: now},
			{ID: hidden, Code: "isbn", Name: "ISBN", IsFilterable: false, CreatedAt: now, UpdatedAt: now},
		},
		values: map[uuid.UUID][]dbsqlc.ProductAttributeValue{
			level: {{ID: uuid.New(), AttributeID: level, Value: "Beginner", Slug: "beginner", SortOrder: 1, CreatedAt: now, UpdatedAt: now}},
		},
	}
	h := handlers.NewAttributeHandler(services.NewAttributeService(f))
	r := gin.New() // no auth middleware: both endpoints are public
	r.GET("/api/v1/attributes", h.List)
	r.GET("/api/v1/attributes/:attrDefId/values", h.ListValues)
	return r, level, hidden
}

func getJSON(r *gin.Engine, path string) (int, map[string]any) {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	var out map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &out)
	return w.Code, out
}

func requireKeys(t *testing.T, m map[string]any, keys ...string) {
	t.Helper()
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			t.Errorf("missing field %q in %v", k, m)
		}
	}
}

func TestAttributes_List(t *testing.T) {
	r, _, _ := attrServer()
	code, body := getJSON(r, "/api/v1/attributes")
	if code != http.StatusOK {
		t.Fatalf("status %d", code)
	}
	data := arr(body["data"])
	if len(data) != 1 {
		t.Fatalf("only filterable attributes expected, got %v", data)
	}
	a := obj(data[0])
	requireKeys(t, a, "id", "code", "name", "is_filterable", "sort_order", "values")
	vals := arr(a["values"])
	if len(vals) != 1 {
		t.Fatalf("values %v", vals)
	}
	requireKeys(t, obj(vals[0]), "id", "attribute_id", "value", "value_slug", "sort_order", "created_at", "updated_at")
}

func TestAttributes_ListValues(t *testing.T) {
	r, level, hidden := attrServer()
	code, body := getJSON(r, "/api/v1/attributes/"+level.String()+"/values")
	if code != http.StatusOK || len(arr(body["data"])) != 1 {
		t.Fatalf("status %d body %v", code, body)
	}
	requireKeys(t, obj(arr(body["data"])[0]), "id", "attribute_id", "value", "value_slug", "sort_order")

	code, body = getJSON(r, "/api/v1/attributes/"+hidden.String()+"/values")
	if code != http.StatusOK || len(arr(body["data"])) != 0 {
		t.Fatalf("existing attribute without values must return empty list: %d %v", code, body)
	}
	if code, body = getJSON(r, "/api/v1/attributes/"+uuid.NewString()+"/values"); code != http.StatusNotFound {
		t.Fatalf("unknown id: %d %v", code, body)
	}
	if code, _ = getJSON(r, "/api/v1/attributes/not-a-uuid/values"); code != http.StatusNotFound {
		t.Fatalf("bad id: %d", code)
	}
}
