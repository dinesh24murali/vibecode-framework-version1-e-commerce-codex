package repositories_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// Integration test of the catalogue filter SQL against the compose Postgres.
// Seeds its own rows under the "zz-" prefix and removes them afterwards.
func TestProductRepository_Filters_Integration(t *testing.T) {
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

	cleanup := func() {
		_, _ = pool.Exec(ctx, `DELETE FROM products WHERE sku LIKE 'zz-%'`)
		_, _ = pool.Exec(ctx, `DELETE FROM product_attributes WHERE code LIKE 'zz_%'`)
		_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE slug LIKE 'zz-%'`)
	}
	cleanup()
	defer cleanup()

	must := func(sql string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, sql, args...); e != nil {
			t.Fatalf("%s: %v", sql, e)
		}
	}
	scan := func(sql string, args ...any) uuid.UUID {
		t.Helper()
		var id uuid.UUID
		if e := pool.QueryRow(ctx, sql, args...).Scan(&id); e != nil {
			t.Fatalf("%s: %v", sql, e)
		}
		return id
	}

	catA := scan(`INSERT INTO categories (slug, name) VALUES ('zz-cat-a','A') RETURNING id`)
	catB := scan(`INSERT INTO categories (slug, name) VALUES ('zz-cat-b','B') RETURNING id`)
	level := scan(`INSERT INTO product_attributes (code, name, is_filterable) VALUES ('zz_level','Level',true) RETURNING id`)
	theme := scan(`INSERT INTO product_attributes (code, name, is_filterable) VALUES ('zz_theme','Theme',true) RETURNING id`)
	beginner := scan(`INSERT INTO product_attribute_values (attribute_id, value, slug) VALUES ($1,'Beginner','beginner') RETURNING id`, level)
	advanced := scan(`INSERT INTO product_attribute_values (attribute_id, value, slug) VALUES ($1,'Advanced','advanced') RETURNING id`, level)
	mandala := scan(`INSERT INTO product_attribute_values (attribute_id, value, slug) VALUES ($1,'Mandala','mandala') RETURNING id`, theme)

	mk := func(sku, title string, cat uuid.UUID, price string, disc any, active bool) uuid.UUID {
		return scan(`INSERT INTO products (category_id, sku, slug, title, price_inr, discount_percent, is_active)
			VALUES ($1,$2,$2,$3,$4,$5,$6) RETURNING id`, cat, sku, title, price, disc, active)
	}
	p1 := mk("zz-1", "Zebra mandala journal", catA, "100", nil, true)
	p2 := mk("zz-2", "Zebra sketch pad", catB, "200", 50, true)
	p3 := mk("zz-3", "Zebra hidden", catA, "50", nil, false) // inactive: never visible
	for _, a := range [][2]uuid.UUID{{p1, beginner}, {p1, mandala}, {p2, advanced}, {p3, beginner}} {
		must(`INSERT INTO product_attribute_assignments (product_id, attribute_value_id) VALUES ($1,$2)`, a[0], a[1])
	}

	repo := repositories.NewProductRepository(pool)
	check := func(name string, f repositories.CatalogFilter, want ...uuid.UUID) {
		t.Helper()
		f.Sort = "newest"
		total, e := repo.Count(ctx, f)
		if e != nil {
			t.Fatalf("%s count: %v", name, e)
		}
		rows, e := repo.List(ctx, f, 50, 0)
		if e != nil {
			t.Fatalf("%s list: %v", name, e)
		}
		if int(total) != len(want) || len(rows) != len(want) {
			t.Fatalf("%s: total=%d rows=%d want %d", name, total, len(rows), len(want))
		}
		got := map[uuid.UUID]bool{}
		for _, r := range rows {
			got[r.ID] = true
		}
		for _, w := range want {
			if !got[w] {
				t.Errorf("%s: missing %s", name, w)
			}
		}
	}

	search := repositories.CatalogFilter{Query: "zebra"}
	check("search excludes inactive", search, p1, p2)
	check("category", repositories.CatalogFilter{Query: "zebra", CategorySlugs: []string{"zz-cat-b"}}, p2)
	check("two categories (OR)", repositories.CatalogFilter{Query: "zebra", CategorySlugs: []string{"zz-cat-a", "zz-cat-b"}}, p1, p2)
	check("one value", repositories.CatalogFilter{Query: "zebra", ValueIDs: []uuid.UUID{beginner}}, p1)
	check("same dimension is OR", repositories.CatalogFilter{Query: "zebra", ValueIDs: []uuid.UUID{beginner, advanced}}, p1, p2)
	check("different dimensions are AND", repositories.CatalogFilter{Query: "zebra", ValueIDs: []uuid.UUID{advanced, mandala}})
	check("AND match", repositories.CatalogFilter{Query: "zebra", ValueIDs: []uuid.UUID{beginner, mandala}}, p1)
	check("no search term matches", repositories.CatalogFilter{Query: "nomatchxyz"})

	// Sorting by effective (discounted) price: p2 costs 100 after 50% off, p1 costs 100 - tie broken by newest.
	rows, err := repo.List(ctx, repositories.CatalogFilter{Query: "zebra", Sort: "price_desc"}, 50, 0)
	if err != nil || len(rows) != 2 {
		t.Fatalf("sort: %v %d", err, len(rows))
	}

	if _, err := repo.GetActive(ctx, p3); err == nil {
		t.Error("inactive product must not be returned by GetActive")
	}
	if _, err := repo.GetActive(ctx, p1); err != nil {
		t.Errorf("active product: %v", err)
	}
}
