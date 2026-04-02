// Package migrations provides the embedded filesystem containing all goose
// migration SQL files. Import this package to get the FS for use with
// github.com/pressly/goose/v3.
package migrations

import "embed"

// FS holds all *.sql migration files embedded at compile time.
// goose reads from this FS so migrations are bundled into the binary and
// run via a one-shot container at deploy time — not at server startup.
//
//go:embed *.sql
var FS embed.FS
