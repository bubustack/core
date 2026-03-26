package stage

import (
	"fmt"

	"github.com/go-logr/logr"
)

// Metadata captures the StoryRun-centric identifiers every Stage participant
// (controllers, cleanup helpers, hub routing) propagates through logs/metrics.
type Metadata struct {
	StoryRun  string
	Namespace string
	Step      string
}

// StoryRunMetadata creates metadata for a StoryRun scoped operation.
func StoryRunMetadata(storyRun, namespace string) Metadata {
	return Metadata{
		StoryRun:  storyRun,
		Namespace: namespace,
	}
}

// WithStep returns a copy of the metadata that also tracks the current step ID.
func (m Metadata) WithStep(step string) Metadata {
	m.Step = step
	return m
}

// Fields returns the canonical logging fields for Stage metadata followed by
// any caller-provided key/value pairs.
func (m Metadata) Fields(kv ...any) []any {
	fields := []any{
		"storyRun", m.StoryRun,
		"namespace", m.Namespace,
	}
	if m.Step != "" {
		fields = append(fields, "step", m.Step)
	}
	return append(fields, normalizeKeyValues(kv...)...)
}

// Info logs using the supplied logger with the metadata fields appended.
func (m Metadata) Info(log logr.Logger, msg string, kv ...any) {
	log.Info(msg, m.Fields(kv...)...)
}

// Error logs using the supplied logger with the metadata fields appended.
func (m Metadata) Error(log logr.Logger, err error, msg string, kv ...any) {
	log.Error(err, msg, m.Fields(kv...)...)
}

func normalizeKeyValues(kv ...any) []any {
	if len(kv) == 0 {
		return nil
	}
	out := make([]any, 0, len(kv)+len(kv)%2)
	for i := 0; i < len(kv); i += 2 {
		value := any("<missing>")
		if i+1 < len(kv) {
			value = kv[i+1]
		}
		out = append(out, fmt.Sprint(kv[i]), value)
	}
	return out
}
