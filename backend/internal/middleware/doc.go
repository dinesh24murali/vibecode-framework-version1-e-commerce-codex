// Package middleware contains Gin middleware functions used across the API.
// Middleware must be stateless where possible; any stateful middleware
// (e.g. rate limiters backed by Redis) should receive its dependencies via
// constructor injection rather than accessing global state.
package middleware
