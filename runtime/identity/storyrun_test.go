package identity

import (
	"strings"
	"testing"

	"github.com/bubustack/core/contracts"
)

func TestStoryRunEngramRunnerServiceAccount(t *testing.T) {
	got := StoryRunEngramRunnerServiceAccount("demo-run")
	if got != "demo-run-engram-runner" {
		t.Fatalf("expected demo-run-engram-runner, got %s", got)
	}
}

func TestStoryRunSelectorLabelsReturnsCopy(t *testing.T) {
	labelsA := StoryRunSelectorLabels("demo-run")
	if labelsA[contracts.StoryRunLabelKey] != "demo-run" {
		t.Fatalf("expected label to equal demo-run, got %s", labelsA[contracts.StoryRunLabelKey])
	}

	labelsB := StoryRunSelectorLabels("other")
	if labelsB[contracts.StoryRunLabelKey] != "other" {
		t.Fatalf("expected new map for other StoryRun, got %s", labelsB[contracts.StoryRunLabelKey])
	}

	labelsA[contracts.StoryRunLabelKey] = "mutated"
	if StoryRunSelectorLabels("demo-run")[contracts.StoryRunLabelKey] != "demo-run" {
		t.Fatalf("expected fresh map per call")
	}
}

func TestStoryRunSelectorLabelsTruncatesLongName(t *testing.T) {
	longName := strings.Repeat("a", 80)
	labelA := StoryRunSelectorLabels(longName)[contracts.StoryRunLabelKey]
	if len(labelA) > maxLabelValueLength {
		t.Fatalf("expected label length <= %d, got %d", maxLabelValueLength, len(labelA))
	}
	if labelA == longName {
		t.Fatalf("expected label to differ for long names")
	}

	labelB := StoryRunSelectorLabels(longName)[contracts.StoryRunLabelKey]
	if labelA != labelB {
		t.Fatalf("expected label to be deterministic, got %s and %s", labelA, labelB)
	}
}

func TestSafeLabelValueSanitizesInvalidShortNames(t *testing.T) {
	value := SafeLabelValue("  Hello_World!!!  ")
	if value != "hello-world" {
		t.Fatalf("expected sanitized value, got %q", value)
	}
}

func TestSafeLabelValueFallsBackWhenInputIsAllInvalid(t *testing.T) {
	value := SafeLabelValue("!!!")
	if value != "resource" {
		t.Fatalf("expected resource fallback, got %q", value)
	}
}
