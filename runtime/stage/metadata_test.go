package stage_test

import (
	"testing"

	"github.com/go-logr/logr/testr"

	"github.com/bubustack/core/runtime/stage"
)

func TestFieldsIncludeMetadata(t *testing.T) {
	meta := stage.StoryRunMetadata("storyrun-a", "bobrapet-system").WithStep("step-1")
	fields := meta.Fields("extra", "value")

	if len(fields) != 8 {
		t.Fatalf("expected 8 fields, got %d", len(fields))
	}
	expected := map[string]string{
		"storyRun":  "storyrun-a",
		"namespace": "bobrapet-system",
		"step":      "step-1",
		"extra":     "value",
	}
	for i := 0; i < len(fields); i += 2 {
		key := fields[i].(string)
		value, ok := fields[i+1].(string)
		if !ok {
			t.Fatalf("expected string value for key %s", key)
		}
		if expected[key] != value {
			t.Fatalf("unexpected field %s=%s", key, value)
		}
	}
}

func TestInfoAndErrorDoNotPanic(t *testing.T) {
	logger := testr.New(t)
	meta := stage.StoryRunMetadata("storyrun-a", "default")
	meta.Info(logger, "test info")
	meta.Error(logger, nil, "test error")
}

func TestFieldsNormalizesOddAndNonStringKeyValues(t *testing.T) {
	meta := stage.StoryRunMetadata("storyrun-a", "default")
	fields := meta.Fields(123, "value", "lonely")

	expected := []any{"storyRun", "storyrun-a", "namespace", "default", "123", "value", "lonely", "<missing>"}
	if len(fields) != len(expected) {
		t.Fatalf("unexpected fields length: %#v", fields)
	}
	for i := range expected {
		if fields[i] != expected[i] {
			t.Fatalf("unexpected field at %d: got %#v want %#v", i, fields[i], expected[i])
		}
	}
}
