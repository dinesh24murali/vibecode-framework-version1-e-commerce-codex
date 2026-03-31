// Package repositories contains data-access logic for the e-commerce API.
// Each repository exposes an interface consumed by the service layer.
// Concrete implementations are backed by PostgreSQL via sqlc-generated code
// and are wired in TASK-004. The Bridge pattern is used here: the service
// layer depends on the repository interface, not the concrete type, so
// the concrete implementation can be swapped without touching service code.
package repositories
