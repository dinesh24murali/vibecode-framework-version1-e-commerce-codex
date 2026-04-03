.PHONY: run build test lint migrate-up migrate-down migrate-status sqlc-gen sqlc-notgen

# Load .env from repo root if it exists and DATABASE_URL is not already set.
ifneq (,$(wildcard ../.env))
  include ../.env
  export
endif

# run starts the API server using go run. Requires all env vars to be set.
run:
	go run ./cmd/api

# build compiles the API binary into bin/api.
build:
	go build -o bin/api ./cmd/api

# test runs the full test suite.
test:
	go test ./...

# lint runs golangci-lint against all packages.
lint:
	golangci-lint run

# migrate-up applies all pending goose migrations.
# Requires DATABASE_URL to be set in the environment.
migrate-up:
	go run ./cmd/migrate up

# migrate-down rolls back the most recent goose migration.
# Requires DATABASE_URL to be set in the environment.
migrate-down:
	go run ./cmd/migrate down

# migrate-status prints the current migration state.
# Requires DATABASE_URL to be set in the environment.
migrate-status:
	go run ./cmd/migrate status

# sqlc-gen runs sqlc code generation.
# Requires sqlc to be installed: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
sqlc-gen:
	sqlc generate
