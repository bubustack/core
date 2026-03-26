package featuretoggles

import "testing"

func TestApply(t *testing.T) {
	called := make(map[string]bool)
	features := Features{
		TelemetryEnabled:         true,
		TracePropagationEnabled:  true,
		VerboseLoggingEnabled:    true,
		StepOutputLoggingEnabled: false,
		MetricsEnabled:           true,
	}
	sink := Sink{
		EnableTelemetry: func(v bool) {
			called["telemetry"] = v
		},
		EnableTracePropagation: func(v bool) {
			called["propagation"] = v
		},
		EnableVerboseLogging: func(v bool) {
			called["verbose"] = v
		},
		EnableStepOutputLogging: func(v bool) {
			called["step"] = v
		},
		EnableMetrics: func(v bool) {
			called["metrics"] = v
		},
	}

	Apply(features, sink)

	if !called["telemetry"] || !called["propagation"] || !called["verbose"] || called["step"] || !called["metrics"] {
		t.Fatalf("unexpected sink calls: %+v", called)
	}
}

func TestApplySkipsNilCallbacksAndSupportsPartialSinks(t *testing.T) {
	called := false
	Apply(Features{MetricsEnabled: true}, Sink{
		EnableMetrics: func(v bool) {
			called = v
		},
	})
	if !called {
		t.Fatalf("expected partial sink callback to run")
	}
}
