package bootstrap

import (
	"github.com/go-logr/logr"

	stagemeta "github.com/bubustack/core/runtime/stage"
)

// ContractLogger emits structured bootstrap events that satisfy the Stage guard.
type ContractLogger struct {
	log       logr.Logger
	component string
	stage     *stagemeta.Metadata
}

// NewContractLogger creates a logger for the given component.
func NewContractLogger(log logr.Logger, component string) ContractLogger {
	return ContractLogger{log: log, component: component}
}

// WithComponent returns a copy targeting a nested component.
func (c ContractLogger) WithComponent(component string) ContractLogger {
	newComponent := component
	if newComponent == "" {
		newComponent = c.component
	} else if c.component != "" {
		newComponent = c.component + "/" + newComponent
	}
	return ContractLogger{
		log:       c.log,
		component: newComponent,
		stage:     c.stage,
	}
}

// WithStage attaches Stage metadata that will be appended to every event.
func (c ContractLogger) WithStage(meta stagemeta.Metadata) ContractLogger {
	return ContractLogger{
		log:       c.log,
		component: c.component,
		stage:     &meta,
	}
}

// Start records the beginning of a bootstrap event.
func (c ContractLogger) Start(event string, kv ...any) {
	c.emit(event, "start", nil, kv...)
}

// Success records a successful bootstrap event.
func (c ContractLogger) Success(event string, kv ...any) {
	c.emit(event, "success", nil, kv...)
}

// Failure records a failed bootstrap event.
func (c ContractLogger) Failure(event string, err error, kv ...any) {
	c.emit(event, "failure", err, kv...)
}

func (c ContractLogger) emit(event, result string, err error, kv ...any) {
	fields := append([]any{
		"component", c.component,
		"event", event,
		"result", result,
	}, kv...)
	if c.stage != nil {
		fields = append(fields, c.stage.Fields()...)
	}
	if err != nil {
		c.log.Error(err, "bootstrap", fields...)
		return
	}
	c.log.Info("bootstrap", fields...)
}
