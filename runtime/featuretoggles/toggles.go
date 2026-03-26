package featuretoggles

// Features captures the runtime feature toggles that need to remain consistent
// across controllers, webhooks, and SDK clients.
type Features struct {
	TelemetryEnabled         bool
	TracePropagationEnabled  bool
	VerboseLoggingEnabled    bool
	StepOutputLoggingEnabled bool
	MetricsEnabled           bool
}

// Sink defines callbacks that apply feature toggles inside a specific process.
// Callers may provide only the functions they support.
type Sink struct {
	EnableTelemetry         func(bool)
	EnableTracePropagation  func(bool)
	EnableVerboseLogging    func(bool)
	EnableStepOutputLogging func(bool)
	EnableMetrics           func(bool)
}

// Apply updates the provided Sink with the supplied feature toggle values. Any
// nil callbacks are skipped, which allows each consumer (controllers, SDK, etc.)
// to provide only the toggles it needs.
func Apply(features Features, sink Sink) {
	if sink.EnableTelemetry != nil {
		sink.EnableTelemetry(features.TelemetryEnabled)
	}
	if sink.EnableTracePropagation != nil {
		sink.EnableTracePropagation(features.TracePropagationEnabled)
	}
	if sink.EnableVerboseLogging != nil {
		sink.EnableVerboseLogging(features.VerboseLoggingEnabled)
	}
	if sink.EnableStepOutputLogging != nil {
		sink.EnableStepOutputLogging(features.StepOutputLoggingEnabled)
	}
	if sink.EnableMetrics != nil {
		sink.EnableMetrics(features.MetricsEnabled)
	}
}
