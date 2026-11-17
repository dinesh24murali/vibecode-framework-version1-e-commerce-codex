package repositories_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// Integration: ordering, inactive exclusion and the 5-slide cap against compose Postgres.
func TestCarouselRepository_Integration(t *testing.T) {
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
	// Work inside a transaction-like cleanup: remove only our rows.
	cleanup := func() { _, _ = pool.Exec(ctx, `DELETE FROM carousel_slides WHERE title LIKE 'zz-%'`) }
	cleanup()
	defer cleanup()
	// Pre-existing slides would distort the cap check; clear them for this test run only if none are real.
	var existing int
	if e := pool.QueryRow(ctx, `SELECT count(*) FROM carousel_slides`).Scan(&existing); e != nil || existing != 0 {
		t.Skipf("carousel_slides not empty (%d rows, err %v); skipping to avoid touching real data", existing, e)
	}
	for i := 7; i >= 1; i-- {
		if _, e := pool.Exec(ctx, `INSERT INTO carousel_slides (title, image_url, sort_order, is_active) VALUES ($1,'https://x/i.jpg',$2,$3)`,
			fmt.Sprintf("zz-%d", i), i, i != 2); e != nil {
			t.Fatal(e)
		}
	}
	rows, err := repositories.NewCarouselRepository(pool).ListActive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"zz-1", "zz-3", "zz-4", "zz-5", "zz-6"} // zz-2 inactive, capped at 5
	if len(rows) != len(want) {
		t.Fatalf("got %d rows, want %d", len(rows), len(want))
	}
	for i, w := range want {
		if rows[i].Title != w {
			t.Errorf("row %d = %s want %s", i, rows[i].Title, w)
		}
	}
}
