package contract

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/handlers"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

type fakeContactRepo struct {
	rows []dbsqlc.CreateContactQueryParams
}

func (f *fakeContactRepo) Create(_ context.Context, a dbsqlc.CreateContactQueryParams) (dbsqlc.ContactQuery, error) {
	f.rows = append(f.rows, a)
	return dbsqlc.ContactQuery{Name: a.Name, Email: a.Email, Phone: a.Phone, Subject: a.Subject, Message: a.Message}, nil
}

type fakeNotifier struct{ sent int }

func (f *fakeNotifier) Notify(context.Context, dbsqlc.ContactQuery) error { f.sent++; return nil }

func contactServer() (*gin.Engine, *fakeContactRepo, *fakeNotifier) {
	repo, n := &fakeContactRepo{}, &fakeNotifier{}
	h := handlers.NewContactHandler(services.NewContactService(repo, n))
	e := gin.New()
	e.POST("/api/v1/contact", h.Submit)
	return e, repo, n
}

func postContact(e *gin.Engine, body any) (int, map[string]any) {
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/contact", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	e.ServeHTTP(w, req)
	var out map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &out)
	return w.Code, out
}

func TestContact_StoresAndNotifies(t *testing.T) {
	e, repo, n := contactServer()
	code, body := postContact(e, map[string]any{
		"name": "Bob Kumar", "email": "bob@example.com", "phone": "+919876500000",
		"subject": "Bulk order", "message": "50 journals please",
	})
	// The spec (04_api_spec.yaml) defines 200 + MessageResponse for this endpoint.
	if code != http.StatusOK {
		t.Fatalf("status %d %v", code, body)
	}
	requireKeys(t, body, "message")
	if len(repo.rows) != 1 || repo.rows[0].Email != "bob@example.com" || !repo.rows[0].Phone.Valid {
		t.Fatalf("not stored: %+v", repo.rows)
	}
	if n.sent != 1 {
		t.Fatalf("notification not sent")
	}
}

func TestContact_PhoneOptional(t *testing.T) {
	e, repo, _ := contactServer()
	code, _ := postContact(e, map[string]any{"name": "A", "email": "a@b.co", "phone": nil, "subject": "S", "message": "M"})
	if code != http.StatusOK || len(repo.rows) != 1 || repo.rows[0].Phone.Valid {
		t.Fatalf("code %d rows %+v", code, repo.rows)
	}
}

func TestContact_Validation(t *testing.T) {
	e, repo, _ := contactServer()
	if code, body := postContact(e, map[string]any{"name": "A", "subject": "S"}); code != http.StatusBadRequest {
		t.Fatalf("missing fields want 400, got %d %v", code, body)
	} else {
		requireKeys(t, body, "code", "message")
	}
	if code, _ := postContact(e, map[string]any{"name": "A", "email": "not-an-email", "subject": "S", "message": "M"}); code != http.StatusUnprocessableEntity {
		t.Fatalf("bad email want 422, got %d", code)
	}
	if code, _ := postContact(e, map[string]any{"name": "A", "email": "a@b.co", "phone": "abc", "subject": "S", "message": "M"}); code != http.StatusUnprocessableEntity {
		t.Fatalf("bad phone want 422, got %d", code)
	}
	if len(repo.rows) != 0 {
		t.Fatalf("invalid enquiries must not be stored")
	}
}
