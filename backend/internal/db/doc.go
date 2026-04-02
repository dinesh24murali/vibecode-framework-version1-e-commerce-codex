// Package db manages the PostgreSQL connection pool (pgxpool) and migration
// tooling (goose) for the API. Use Open to obtain a pool and RunUp/RunDown
// for migration operations.
package db
