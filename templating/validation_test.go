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
