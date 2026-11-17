// Package handlers holds background job handlers. Phase 2 tasks replace the stubs.
package handlers

import (
	"context"
	"log/slog"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/worker"
)

// JobTypes pre-registered as stubs for Phase 2.
var JobTypes = []string{"order_confirmation_email", "payment_reconciliation", "stock_release", "audit_fan_out"}

// RegisterAll registers every handler (currently logging stubs) on d.
func RegisterAll(d *worker.Dispatcher) {
	for _, t := range JobTypes {
		d.Register(t, worker.HandlerFunc(func(_ context.Context, j worker.Job) error {
			slog.Info("stub handler", "id", j.ID, "type", j.Type)
			return nil
		}))
	}
}
