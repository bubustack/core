package templating

import (
	"strings"
	"testing"
)

func TestValidateTemplateStringRejectsUnknownRoot(t *testing.T) {
	scope := NewExpressionScope("step", false, false, false, RootInputs)
	err := ValidateTemplateString(`{{ .mystery.value }}`, scope)
	if err == nil || !strings.Contains(err.Error(), "unknown context 'mystery'") {
		t.Fatalf("expected unknown root validation error, got %v", err)
	}
}

func TestValidateTemplateStringRejectsIndexBypass(t *testing.T) {
	scope := NewExpressionScope("step", false, false, false, RootInputs)
	err := ValidateTemplateString(`{{ index .steps "build" }}`, scope)
	if err == nil || !strings.Contains(err.Error(), "context 'steps' is not allowed") {
		t.Fatalf("expected steps validation error, got %v", err)
	}
}

func TestValidateTemplateStringRejectsGetBypassViaDotRoot(t *testing.T) {
	scope := NewExpressionScope("step", false, false, false, RootInputs)
	err := ValidateTemplateString(`{{ get . "steps" }}`, scope)
	if err == nil || !strings.Contains(err.Error(), "context 'steps' is not allowed") {
		t.Fatalf("expected steps validation error, got %v", err)
	}
}

func TestValidateJSONTemplatesRejectsOversizedInput(t *testing.T) {
	scope := NewExpressionScope("step", false, false, false, RootInputs)
	raw := `{"value":"` + strings.Repeat("a", maxTemplateValidationJSONBytes) + `"}`
	err := ValidateJSONTemplates([]byte(raw), scope)
	if err == nil || !strings.Contains(err.Error(), "exceeds max validation size") {
		t.Fatalf("expected oversize validation error, got %v", err)
	}
}

func TestValidateTemplateStringAllowsIndexedStepsChainWhenStepsAllowed(t *testing.T) {
	scope := NewExpressionScope("story-output", false, false, false, RootInputs, RootSteps)
	err := ValidateTemplateString(`{{ (index .steps "fetch-feed").output.body }}`, scope)
	if err != nil {
		t.Fatalf("expected indexed step chain to validate when steps are allowed, got %v", err)
	}
}

func TestValidateJSONTemplatesAllowsIndexedStepsChainWhenStepsAllowed(t *testing.T) {
	scope := NewExpressionScope("batch-runtime", false, false, false, RootInputs, RootSteps)
	raw := []byte(`{"userPrompt":"{{ (index .steps \"fetch-feed\").output.body }}"}`)
	err := ValidateJSONTemplates(raw, scope)
	if err != nil {
		t.Fatalf("expected indexed step chain in JSON template to validate when steps are allowed, got %v", err)
	}
}
