package templating

import (
	"strings"
	"testing"
)

func TestNormalizeTemplateRootsPreservesActionBoundariesInsideStrings(t *testing.T) {
	input := `{{ printf "%s" "}}" }}`
	got := normalizeTemplateRoots(input)
	if got != input {
		t.Fatalf("expected normalization to preserve action boundaries, got %q", got)
	}
}

func TestNormalizeTemplateRootsEscapesSubscriptKeys(t *testing.T) {
	input := `{{ inputs["a\"b"] }}`
	got := normalizeTemplateRoots(input)
	if !strings.Contains(got, `(index .inputs "a\"b")`) {
		t.Fatalf("expected escaped index conversion, got %q", got)
	}
}
