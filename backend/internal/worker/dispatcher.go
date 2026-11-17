// Package worker implements the Redis Streams background job consumer.
package worker

import (
	"context"
	"fmt"
	"sync"
)

// Job is a unit of background work read from the stream.
type Job struct {
	ID       string // stable job identifier, used as the idempotency key
	Type     string
	Payload  string
	Attempts int
}

// Handler processes one job. Implementations MUST be idempotent.
type Handler interface {
	Handle(ctx context.Context, job Job) error
}

// HandlerFunc adapts a function to Handler.
type HandlerFunc func(ctx context.Context, job Job) error

// Handle calls f.
func (f HandlerFunc) Handle(ctx context.Context, job Job) error { return f(ctx, job) }

// Dispatcher routes jobs to handlers by job type.
type Dispatcher struct {
	handlers map[string]Handler
	mu       sync.RWMutex
}

// NewDispatcher returns an empty Dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{handlers: map[string]Handler{}}
}

// Register binds a handler to a job type, replacing any existing binding.
func (d *Dispatcher) Register(jobType string, h Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[jobType] = h
}

// Dispatch runs the handler registered for job.Type.
func (d *Dispatcher) Dispatch(ctx context.Context, job Job) error {
	d.mu.RLock()
	h, ok := d.handlers[job.Type]
	d.mu.RUnlock()
	if !ok {
		return fmt.Errorf("worker: no handler registered for job type %q", job.Type)
	}
	return h.Handle(ctx, job)
}
