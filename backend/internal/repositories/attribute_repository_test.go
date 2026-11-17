package repositories_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// Integration: runs against the compose Postgres; skipped when DATABASE_URL is unset.
func TestAttributeRepository_Integration(t *testing.T) {
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
	if _, e := pool.Exec(ctx, `DELETE FROM product_attributes WHERE code = 'zz_test_attr'`); e != nil {
		t.Fatal(e)
	}
	var id string
	if e := pool.QueryRow(ctx, `INSERT INTO product_attributes (code, name, is_filterable, sort_order) VALUES ('zz_test_attr','Test',true,99) RETURNING id`).Scan(&id); e != nil {
		t.Fatal(e)
	}
	defer func() { _, _ = pool.Exec(ctx, `DELETE FROM product_attributes WHERE code = 'zz_test_attr'`) }()
	if _, e := pool.Exec(ctx, `INSERT INTO product_attribute_values (attribute_id, value, slug, sort_order) VALUES ($1,'B','b',2),($1,'A','a',1)`, id); e != nil {
		t.Fatal(e)
	}

	repo := repositories.NewAttributeRepository(pool)
	attrs, err := repo.ListFilterable(ctx)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, a := range attrs {
		if a.Code == "zz_test_attr" {
			found = true
			vals, err := repo.ListValues(ctx, a.ID)
			if err != nil || len(vals) != 2 || vals[0].Slug != "a" {
				t.Fatalf("values %v err %v (want sorted by sort_order)", vals, err)
			}
		}
	}
	if !found {
		t.Fatal("seeded attribute not returned by ListFilterable")
	}
}
