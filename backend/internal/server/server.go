// Package server implements the Builder pattern for constructing the HTTP
// server. ServerBuilder accumulates optional settings via fluent With* methods
// and produces a configured *http.Server via Build. This keeps main.go
// declarative and makes it easy to add future options (TLS, timeouts, etc.)
// without changing the call site.
package server

import (
	"net/http"
	"time"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/config"
)

const (
	defaultReadTimeout     = 15 * time.Second
	defaultWriteTimeout    = 15 * time.Second
	defaultIdleTimeout     = 60 * time.Second
	defaultShutdownTimeout = 10 * time.Second
)

// ServerBuilder assembles an *http.Server using the Builder pattern.
// Methods are chainable and return the same builder for fluency.
type ServerBuilder struct {
	cfg          *config.Config
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

// NewServerBuilder returns a ServerBuilder populated with defaults.
func NewServerBuilder(cfg *config.Config) *ServerBuilder {
	return &ServerBuilder{
		cfg:          cfg,
		readTimeout:  defaultReadTimeout,
		writeTimeout: defaultWriteTimeout,
		idleTimeout:  defaultIdleTimeout,
	}
}

// WithReadTimeout overrides the HTTP server read timeout.
func (b *ServerBuilder) WithReadTimeout(d time.Duration) *ServerBuilder {
	b.readTimeout = d
	return b
}

// WithWriteTimeout overrides the HTTP server write timeout.
func (b *ServerBuilder) WithWriteTimeout(d time.Duration) *ServerBuilder {
	b.writeTimeout = d
	return b
}

// WithIdleTimeout overrides the HTTP server idle (keep-alive) timeout.
func (b *ServerBuilder) WithIdleTimeout(d time.Duration) *ServerBuilder {
	b.idleTimeout = d
	return b
}

// Build constructs and returns the configured *http.Server.
// The Gin router is built here so that all routes and middleware are
// registered before the server starts accepting connections.
func (b *ServerBuilder) Build() *http.Server {
	router := newRouter(b.cfg.CORSAllowedOrigins)

	return &http.Server{
		Addr:         b.cfg.Addr(),
		Handler:      router,
		ReadTimeout:  b.readTimeout,
		WriteTimeout: b.writeTimeout,
		IdleTimeout:  b.idleTimeout,
	}
}

// ShutdownTimeout returns the grace period allowed for in-flight requests
// to complete before the server is forcibly closed.
func ShutdownTimeout() time.Duration {
	return defaultShutdownTimeout
}
