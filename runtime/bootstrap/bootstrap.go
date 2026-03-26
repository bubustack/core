package bootstrap

import (
	"fmt"

	"github.com/go-logr/logr"
)

// Runner coordinates startup registrations (controllers, webhooks, health checks, etc.)
// and emits consistent structured logs for both success and failure paths.
type Runner struct {
	Log logr.Logger
}

// Entry describes a single component registration.
type Entry struct {
	Kind           string
	Name           string
	Register       func() error
	ErrMessage     string
	SuccessMessage string
	Fields         []any
}

// Register executes every entry sequentially. The first failure is logged and
// returned so callers decide whether to exit or retry; successes emit a log entry
// so operators know which components finished wiring.
func (r Runner) Register(entries ...Entry) error {
	for _, entry := range entries {
		if entry.Register == nil {
			continue
		}

		if err := entry.Register(); err != nil {
			r.logError(entry, err)
			return fmt.Errorf("bootstrap: register %s %s: %w", entry.Kind, entry.Name, err)
		}
		r.logSuccess(entry)
	}
	return nil
}

func (r Runner) logError(entry Entry, err error) {
	msg := entry.ErrMessage
	if msg == "" {
		msg = "unable to register component"
	}
	fields := append([]any{}, entry.Fields...)
	if entry.Kind != "" {
		fields = append(fields, "kind", entry.Kind)
	}
	if entry.Name != "" {
		fields = append(fields, "name", entry.Name)
	}
	r.Log.Error(err, msg, fields...)
}

func (r Runner) logSuccess(entry Entry) {
	msg := entry.SuccessMessage
	if msg == "" {
		msg = "component registered"
	}
	fields := append([]any{}, entry.Fields...)
	if entry.Kind != "" {
		fields = append(fields, "kind", entry.Kind)
	}
	if entry.Name != "" {
		fields = append(fields, "name", entry.Name)
	}
	r.Log.Info(msg, fields...)
}
