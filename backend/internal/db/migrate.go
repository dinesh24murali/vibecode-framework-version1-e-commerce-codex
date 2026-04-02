package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql
	"github.com/pressly/goose/v3"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/migrations"
)

func init() {
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(fmt.Sprintf("db: set goose dialect: %v", err))
	}
}

// openSQL opens a *sql.DB from a Postgres URL for use with goose.
// Callers must close it after the migration operation completes.
func openSQL(url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("db: open sql: %w", err)
	}
	return db, nil
}

// RunUp applies all pending migrations.
func RunUp(ctx context.Context, url string) error {
	db, err := openSQL(url)
	if err != nil {
		return err
	}
	defer db.Close()
	return goose.UpContext(ctx, db, ".")
}

// RunDown rolls back the most recent migration.
func RunDown(ctx context.Context, url string) error {
	db, err := openSQL(url)
	if err != nil {
		return err
	}
	defer db.Close()
	return goose.DownContext(ctx, db, ".")
}

// PrintStatus prints the current migration state to stdout.
func PrintStatus(ctx context.Context, url string) error {
	db, err := openSQL(url)
	if err != nil {
		return err
	}
	defer db.Close()
	return goose.StatusContext(ctx, db, ".")
}
