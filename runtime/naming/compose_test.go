package naming

import (
	"strings"
	"testing"
)

func TestComposeDNS1123WithinLimit(t *testing.T) {
	name := ComposeDNS1123("story", "step")
	if name != "story-step" {
		t.Fatalf("expected story-step, got %s", name)
	}
	if len(name) > maxDNS1123Length {
		t.Fatalf("expected length <= %d, got %d", maxDNS1123Length, len(name))
	}
}

func TestComposeDNS1123TruncatesAndHashes(t *testing.T) {
	longPart := "this-is-a-very-long-name-component-that-will-force-truncation"
	name := ComposeDNS1123(longPart, longPart, longPart)
	if len(name) > maxDNS1123Length {
		t.Fatalf("expected length <= %d, got %d", maxDNS1123Length, len(name))
	}
	if len(name) < 9 {
		t.Fatalf("expected hashed suffix, got %s", name)
	}
}

func TestComposeDNS1123HashIsDeterministicAndDistinct(t *testing.T) {
	longA := "this-is-a-very-long-name-component-that-will-force-truncation-a"
	longB := "this-is-a-very-long-name-component-that-will-force-truncation-b"

	first := ComposeDNS1123(longA, longA, longA)
	second := ComposeDNS1123(longA, longA, longA)
	other := ComposeDNS1123(longB, longB, longB)

	if first != second {
		t.Fatalf("expected deterministic truncation, got %q and %q", first, second)
	}
	if first == other {
		t.Fatalf("expected distinct hashes for distinct long inputs, got %q", first)
	}
	suffix := first[strings.LastIndex(first, "-")+1:]
	if len(suffix) != 12 {
		t.Fatalf("expected 12-character hash suffix, got %q", suffix)
	}
}

func TestComposeDNS1123WithSuffixPreservesSuffix(t *testing.T) {
	longPart := "this-is-a-very-long-name-component-that-will-force-truncation"
	name := ComposeDNS1123WithSuffix(longPart, "stdio")
	if len(name) > maxDNS1123Length {
		t.Fatalf("expected length <= %d, got %d", maxDNS1123Length, len(name))
	}
	if name[len(name)-5:] != "stdio" {
		t.Fatalf("expected suffix stdio, got %s", name)
	}
}

func TestComposeDNS1123SanitizesInvalidParts(t *testing.T) {
	name := ComposeDNS1123("  Story Name  ", "%%%")
	if name != "story-name" {
		t.Fatalf("expected sanitized name, got %q", name)
	}
}

func TestComposeDNS1123WithSuffixHandlesEmptyBase(t *testing.T) {
	name := ComposeDNS1123WithSuffix("!!!", "stdio")
	if name != "resource-stdio" {
		t.Fatalf("expected resource fallback, got %q", name)
	}
}
