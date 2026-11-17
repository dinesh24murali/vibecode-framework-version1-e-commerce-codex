package repositories_test

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// Integration: stock is reserved atomically, so concurrent orders cannot oversell,
// and a failed order leaves no trace (transaction rolled back).
func TestCheckoutRepository_Integration(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	suffix := uuid.NewString()[:8]
	var userID, catID, productID uuid.UUID
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(pool.QueryRow(ctx, `INSERT INTO users (email, full_name) VALUES ($1, 'zz') RETURNING id`, "zz-"+suffix+"@example.com").Scan(&userID))
	must(pool.QueryRow(ctx, `INSERT INTO categories (slug, name) VALUES ($1, 'zz') RETURNING id`, "zz-"+suffix).Scan(&catID))
	must(pool.QueryRow(ctx, `INSERT INTO products (category_id, sku, slug, title, price_inr) VALUES ($1, $2, $2, 'zz', 100) RETURNING id`, catID, "zz-"+suffix).Scan(&productID))
	_, err = pool.Exec(ctx, `INSERT INTO inventory_items (product_id, on_hand, reserved) VALUES ($1, 3, 0)`, productID)
	must(err)
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM orders WHERE user_id = $1`, userID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
		_, _ = pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE id = $1`, catID)
	}()

	repo := repositories.NewCheckoutRepository(pool)
	draft := func(qty int32) *repositories.OrderDraft {
		d := &repositories.OrderDraft{
			UserID: userID, OrderNumber: "ZZ-" + uuid.NewString()[:10], Subtotal: 100 * float64(qty), Taxable: 100 * float64(qty),
			GST: 18 * float64(qty), Total: 118 * float64(qty),
			Items: []repositories.DraftItem{{ProductID: productID, Title: "zz", SKU: "zz", Quantity: qty, UnitPrice: 100, LineTotal: 100 * float64(qty)}},
		}
		d.Address.RecipientName, d.Address.Line1, d.Address.City = "Ann", "1 St", "Chennai"
		d.Address.State, d.Address.PostalCode, d.Address.CountryCode = "TN", "600001", "IN"
		return d
	}

	// 5 concurrent orders of 1 against 3 units: exactly 3 succeed.
	var wg sync.WaitGroup
	results := make(chan error, 5)
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, e := repo.Place(ctx, draft(1))
			results <- e
		}()
	}
	wg.Wait()
	close(results)
	ok, short := 0, 0
	for e := range results {
		switch {
		case e == nil:
			ok++
		case errors.Is(e, repositories.ErrInsufficientStock):
			short++
		default:
			t.Fatalf("unexpected error: %v", e)
		}
	}
	if ok != 3 || short != 2 {
		t.Fatalf("ok=%d short=%d, want 3 and 2", ok, short)
	}
	var reserved, orders int
	must(pool.QueryRow(ctx, `SELECT reserved FROM inventory_items WHERE product_id = $1`, productID).Scan(&reserved))
	must(pool.QueryRow(ctx, `SELECT count(*) FROM orders WHERE user_id = $1`, userID).Scan(&orders))
	if reserved != 3 || orders != 3 {
		t.Fatalf("reserved=%d orders=%d, want 3 and 3", reserved, orders)
	}

	// Rollback: an order that cannot be fully reserved leaves stock and orders untouched.
	if _, err = repo.Place(ctx, draft(1)); !errors.Is(err, repositories.ErrInsufficientStock) {
		t.Fatalf("want insufficient stock, got %v", err)
	}
	must(pool.QueryRow(ctx, `SELECT count(*) FROM orders WHERE user_id = $1`, userID).Scan(&orders))
	if orders != 3 {
		t.Fatalf("orders = %d after failed order", orders)
	}
}
