// Package config implements the Adapter pattern for loading application
// configuration from environment variables. The EnvAdapter reads raw env
// vars and produces a validated Config struct. Call Load at startup; the
// process exits with a clear message if any required var is absent.
package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Config holds all runtime configuration for the API server.
// All fields are required except Port, which defaults to "8080".
type Config struct {
	// DatabaseURL is the full Postgres connection string.
	DatabaseURL string

	// RedisURL is the Redis connection URL (redis://host:port).
	RedisURL string

	// JWTPrivateKey is the PEM-encoded RSA/EC private key used to sign JWTs.
	JWTPrivateKey string

	// JWTPublicKey is the PEM-encoded public key used to verify JWTs.
	JWTPublicKey string

	// Port is the TCP port the HTTP server listens on. Defaults to "8080".
	Port string

	// CORSAllowedOrigins is a comma-separated list of exact allowed origins.
	// Wildcards are not permitted per architecture constraints.
	CORSAllowedOrigins []string
}

// EnvAdapter is the Adapter that bridges raw OS environment variables to
// the domain Config struct.
type EnvAdapter struct{}

// NewEnvAdapter constructs an EnvAdapter.
func NewEnvAdapter() *EnvAdapter {
	return &EnvAdapter{}
}

// Load reads environment variables, validates that all required vars are
// present, and returns a populated Config. If any required variable is
// missing the process calls log.Fatal with a descriptive message so the
// operator has immediate, actionable feedback at startup.
func (a *EnvAdapter) Load() *Config {
	missing := []string{}

	get := func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			missing = append(missing, key)
		}
		return val
	}

	cfg := &Config{
		DatabaseURL:   get("DATABASE_URL"),
		RedisURL:      get("REDIS_URL"),
		JWTPrivateKey: get("JWT_PRIVATE_KEY"),
		JWTPublicKey:  get("JWT_PUBLIC_KEY"),
		Port:          os.Getenv("PORT"),
	}

	corsRaw := get("CORS_ALLOWED_ORIGINS")

	if len(missing) > 0 {
		log.Fatalf(
			"config: missing required environment variables: %s",
			strings.Join(missing, ", "),
		)
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	cfg.CORSAllowedOrigins = parseCORSOrigins(corsRaw)
	if len(cfg.CORSAllowedOrigins) == 0 {
		log.Fatal("config: CORS_ALLOWED_ORIGINS must contain at least one origin")
	}

	return cfg
}

// parseCORSOrigins splits a comma-separated origin list and trims whitespace.
func parseCORSOrigins(raw string) []string {
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}

// Addr returns the TCP listen address in ":port" form.
func (c *Config) Addr() string {
	return fmt.Sprintf(":%s", c.Port)
}
