package identity

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"unicode"

	"github.com/bubustack/core/contracts"
	"github.com/bubustack/core/runtime/naming"
)

// StoryRunEngramRunnerServiceAccount returns the canonical ServiceAccount name
// used by Engram runner workloads spawned for the provided StoryRun.
func StoryRunEngramRunnerServiceAccount(storyRunName string) string {
	return naming.ComposeDNS1123WithSuffix(storyRunName, "engram-runner")
}

// LabelValueFromName returns a DNS-1123 label-safe value for the provided name.
// This mirrors the internal hashing/truncation used for StoryRun labels.
func LabelValueFromName(name string) string {
	return labelValueFromName(name)
}

// StoryRunSelectorLabels returns the set of base labels used to select objects
// associated with the provided StoryRun. The returned map is a fresh copy on
// every invocation so callers can safely mutate it with extra selectors.
func StoryRunSelectorLabels(storyRunName string) map[string]string {
	return map[string]string{
		contracts.StoryRunLabelKey: labelValueFromName(storyRunName),
	}
}

const (
	maxLabelValueLength = 63
	labelHashLength     = 10
)

func labelValueFromName(name string) string {
	name = sanitizeLabelValue(name)
	if name == "" {
		return "resource"
	}
	if len(name) <= maxLabelValueLength {
		return name
	}

	hash := sha256.Sum256([]byte(name))
	suffix := hex.EncodeToString(hash[:])[:labelHashLength]
	maxPrefix := maxLabelValueLength - 1 - labelHashLength
	prefix := name
	if len(prefix) > maxPrefix {
		prefix = prefix[:maxPrefix]
	}
	prefix = strings.Trim(prefix, "-_.")
	if prefix == "" {
		return suffix
	}
	return prefix + "-" + suffix
}

func sanitizeLabelValue(value string) string {
	var out strings.Builder
	lastDash := false

	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			out.WriteRune(r)
			lastDash = false
		case r == '-', r == '_', r == '.', unicode.IsSpace(r):
			if out.Len() == 0 || lastDash {
				continue
			}
			out.WriteByte('-')
			lastDash = true
		default:
			if out.Len() == 0 || lastDash {
				continue
			}
			out.WriteByte('-')
			lastDash = true
		}
	}

	return strings.Trim(out.String(), "-")
}

// SafeLabelValue returns a DNS-1123 compliant label value (max 63 chars).
// It preserves readability when possible and falls back to a hash suffix when needed.
func SafeLabelValue(value string) string {
	return labelValueFromName(value)
}
