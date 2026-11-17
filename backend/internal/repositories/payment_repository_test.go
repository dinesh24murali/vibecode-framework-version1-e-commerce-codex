package repositories_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// Integration: capture marks the order paid and deducts stock exactly once; replays and
// late failures are no-ops; failure moves the order to payment_failed.
func TestPaymentRepository_Integration(t *testing.T) {
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

	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	sfx := uuid.NewString()[:8]
	var userID, catID, productID, orderA, orderB uuid.UUID
	must(pool.QueryRow(ctx, `INSERT INTO users (email, full_name) VALUES ($1, 'zz') RETURNING id`, "zzp-"+sfx+"@example.com").Scan(&userID))
	must(pool.QueryRow(ctx, `INSERT INTO categories (slug, name) VALUES ($1, 'zz') RETURNING id`, "zzp-"+sfx).Scan(&catID))
	must(pool.QueryRow(ctx, `INSERT INTO products (category_id, sku, slug, title, price_inr) VALUES ($1, $2, $2, 'zz', 100) RETURNING id`, catID, "zzp-"+sfx).Scan(&productID))
	_, err = pool.Exec(ctx, `INSERT INTO inventory_items (product_id, on_hand, reserved) VALUES ($1, 10, 4)`, productID)
	must(err)
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM payments WHERE order_id IN ($1, $2)`, orderA, orderB)
		_, _ = pool.Exec(ctx, `DELETE FROM orders WHERE user_id = $1`, userID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
		_, _ = pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE id = $1`, catID)
	}()
	mkOrder := func(qty int) uuid.UUID {
		var id uuid.UUID
		must(pool.QueryRow(ctx, `INSERT INTO orders (user_id, order_number, subtotal_inr, total_inr) VALUES ($1, $2, 100, 118) RETURNING id`, userID, "ZZP-"+uuid.NewString()[:10]).Scan(&id))
		_, e := pool.Exec(ctx, `INSERT INTO order_items (order_id, product_id, product_title, product_sku, quantity, unit_price_inr, line_total_inr) VALUES ($1,$2,'zz','zz',$3,100,100)`, id, productID, qty)
		must(e)
		return id
	}
	orderA, orderB = mkOrder(2), mkOrder(1)
	repo := repositories.NewPaymentRepository(pool)
	rzpA, rzpB := "order_zzA"+sfx, "order_zzB"+sfx
	_, err = repo.Create(ctx, orderA, rzpA, "118.00")
	must(err)
	_, err = repo.Create(ctx, orderB, rzpB, "118.00")
	must(err)

	if _, err = repo.Apply(ctx, "order_missing", "pay_x", repositories.PaymentCaptured, nil); !errors.Is(err, repositories.ErrPaymentNotFound) {
		t.Fatalf("want not found, got %v", err)
	}

	res, err := repo.Apply(ctx, rzpA, "pay_A", repositories.PaymentCaptured, []byte(`{"event":"payment.captured"}`))
	must(err)
	if !res.Applied || res.OrderID != orderA {
		t.Fatalf("first capture: %+v", res)
	}
	if res, err = repo.Apply(ctx, rzpA, "pay_A", repositories.PaymentCaptured, nil); err != nil || res.Applied {
		t.Fatalf("replay must be a no-op: %+v %v", res, err)
	}
	if res, err = repo.Apply(ctx, rzpA, "pay_A", repositories.PaymentFailed, nil); err != nil || res.Applied {
		t.Fatalf("late failure must be a no-op: %+v %v", res, err)
	}
	var onHand, reserved int
	var status, payStatus, payID string
	must(pool.QueryRow(ctx, `SELECT on_hand, reserved FROM inventory_items WHERE product_id = $1`, productID).Scan(&onHand, &reserved))
	must(pool.QueryRow(ctx, `SELECT status FROM orders WHERE id = $1`, orderA).Scan(&status))
	must(pool.QueryRow(ctx, `SELECT status, provider_payment_id FROM payments WHERE provider_order_id = $1`, rzpA).Scan(&payStatus, &payID))
	if onHand != 8 || reserved != 2 || status != "paid" || payStatus != "captured" || payID != "pay_A" {
		t.Fatalf("on_hand=%d reserved=%d order=%s payment=%s/%s", onHand, reserved, status, payStatus, payID)
	}

	if res, err = repo.Apply(ctx, rzpB, "pay_B", repositories.PaymentFailed, nil); err != nil || !res.Applied {
		t.Fatalf("failure: %+v %v", res, err)
	}
	must(pool.QueryRow(ctx, `SELECT status FROM orders WHERE id = $1`, orderB).Scan(&status))
	must(pool.QueryRow(ctx, `SELECT on_hand, reserved FROM inventory_items WHERE product_id = $1`, productID).Scan(&onHand, &reserved))
	if status != "payment_failed" || onHand != 8 || reserved != 2 {
		t.Fatalf("after failure: order=%s on_hand=%d reserved=%d", status, onHand, reserved)
	}
}
