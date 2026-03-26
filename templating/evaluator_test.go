package templating

import (
	"context"
	"testing"
	"time"
)

func TestNewRejectsNegativeLimits(t *testing.T) {
	if _, err := New(Config{EvaluationTimeout: -time.Second}); err == nil {
		t.Fatalf("expected negative timeout to fail")
	}
	if _, err := New(Config{MaxOutputBytes: -1}); err == nil {
		t.Fatalf("expected negative max output bytes to fail")
	}
}

func TestResolveValueDetectsCycles(t *testing.T) {
	evaluator, err := New(Config{})
	if err != nil {
		t.Fatalf("new evaluator: %v", err)
	}
	t.Cleanup(evaluator.Close)

	cyclic := map[string]any{}
	cyclic["self"] = cyclic

	if _, err := evaluator.ResolveValue(context.Background(), cyclic, nil); err == nil {
		t.Fatalf("expected cycle detection error")
	}
}

func TestResolveTemplateStringReturnsStorageSelectorForOffloadedOutput(t *testing.T) {
	evaluator, err := New(Config{})
	if err != nil {
		t.Fatalf("new evaluator: %v", err)
	}
	t.Cleanup(evaluator.Close)

	value, err := evaluator.ResolveTemplateString(
		context.Background(),
		`{{ (index .steps "fetch").output.body }}`,
		map[string]any{
			RootSteps: map[string]any{
				"fetch": map[string]any{
					"output": map[string]any{
						StorageRefKey: "s3://bucket/object",
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("resolve template string: %v", err)
	}

	selector, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("expected selector map, got %T", value)
	}
	if selector[StorageRefKey] != "s3://bucket/object" {
		t.Fatalf("unexpected storage ref: %#v", selector)
	}
	if selector[StoragePathKey] != "body" {
		t.Fatalf("unexpected storage path: %#v", selector)
	}
}

func TestResolveTemplateStringSupportsEscapedQuotedKeys(t *testing.T) {
	evaluator, err := New(Config{})
	if err != nil {
		t.Fatalf("new evaluator: %v", err)
	}
	t.Cleanup(evaluator.Close)

	value, err := evaluator.ResolveTemplateString(context.Background(), `{{ .inputs["a\"b"] }}`, map[string]any{
		RootInputs: map[string]any{
			`a"b`: "value",
		},
	})
	if err != nil {
		t.Fatalf("resolve template string: %v", err)
	}
	if value != "value" {
		t.Fatalf("expected value, got %#v", value)
	}
}

func TestResolveTemplateStringEnforcesOutputLimit(t *testing.T) {
	evaluator, err := New(Config{MaxOutputBytes: 4})
	if err != nil {
		t.Fatalf("new evaluator: %v", err)
	}
	t.Cleanup(evaluator.Close)

	if _, err := evaluator.ResolveTemplateString(context.Background(), `{{ "abcdef" }}`, nil); err == nil {
		t.Fatalf("expected output limit error")
	}
}
