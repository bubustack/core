package bootstrap_test

import (
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"

	"github.com/bubustack/core/runtime/bootstrap"
	stagemeta "github.com/bubustack/core/runtime/stage"
)

func TestContractLoggerEmitsFields(t *testing.T) {
	records := &capturedRecords{}
	logger := logr.New(&captureSink{records: records})
	meta := stagemeta.StoryRunMetadata("story-run", "ns").WithStep("step-a")
	contract := bootstrap.NewContractLogger(logger, "manager").WithStage(meta)

	contract.Start("init")
	contract.Success("init", "detail", "ok")
	contract.Failure("init", errors.New("boom"))

	if len(records.entries) != 3 {
		t.Fatalf("expected 3 log entries, got %d", len(records.entries))
	}

	start := records.entries[0]
	if start.msg != "bootstrap" {
		t.Fatalf("unexpected start message: %#v", start)
	}
	assertLogField(t, start.kv, "component", "manager")
	assertLogField(t, start.kv, "event", "init")
	assertLogField(t, start.kv, "result", "start")
	assertLogField(t, start.kv, "storyRun", "story-run")
	assertLogField(t, start.kv, "namespace", "ns")
	assertLogField(t, start.kv, "step", "step-a")

	success := records.entries[1]
	assertLogField(t, success.kv, "result", "success")
	assertLogField(t, success.kv, "detail", "ok")

	failure := records.entries[2]
	if failure.err == nil || failure.err.Error() != "boom" {
		t.Fatalf("expected failure log to carry error, got %#v", failure)
	}
	assertLogField(t, failure.kv, "result", "failure")
}

func TestContractLoggerNestedComponent(t *testing.T) {
	logger := testr.New(t)
	parent := bootstrap.NewContractLogger(logger, "bootstrap")
	child := parent.WithComponent("controller/Story")

	child.Start("register")
}

type capturedRecords struct {
	entries []capturedEntry
}

type capturedEntry struct {
	msg string
	err error
	kv  []any
}

type captureSink struct {
	records *capturedRecords
	values  []any
}

func (s *captureSink) Init(logr.RuntimeInfo) {}

func (s *captureSink) Enabled(int) bool {
	return true
}

func (s *captureSink) Info(_ int, msg string, keysAndValues ...any) {
	s.records.entries = append(s.records.entries, capturedEntry{
		msg: msg,
		kv:  append(append([]any{}, s.values...), keysAndValues...),
	})
}

func (s *captureSink) Error(err error, msg string, keysAndValues ...any) {
	s.records.entries = append(s.records.entries, capturedEntry{
		msg: msg,
		err: err,
		kv:  append(append([]any{}, s.values...), keysAndValues...),
	})
}

func (s *captureSink) WithValues(keysAndValues ...any) logr.LogSink {
	clone := *s
	clone.values = append(append([]any{}, s.values...), keysAndValues...)
	return &clone
}

func (s *captureSink) WithName(string) logr.LogSink {
	clone := *s
	return &clone
}

func assertLogField(t *testing.T, kv []any, key string, expected any) {
	t.Helper()
	for i := 0; i+1 < len(kv); i += 2 {
		if kv[i] == key {
			if kv[i+1] != expected {
				t.Fatalf("expected %s=%v, got %v", key, expected, kv[i+1])
			}
			return
		}
	}
	t.Fatalf("expected field %s in %#v", key, kv)
}
