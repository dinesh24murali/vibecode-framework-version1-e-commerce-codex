.PHONY: run build test lint

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
