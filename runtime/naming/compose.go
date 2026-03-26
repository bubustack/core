package naming

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"

	"k8s.io/apimachinery/pkg/util/validation"
)

const maxDNS1123Length = validation.DNS1123LabelMaxLength

const defaultDNS1123Name = "resource"

// ComposeDNS1123 builds a DNS-1123 compliant name from the provided parts.
// It truncates and appends a deterministic hash suffix when needed.
func ComposeDNS1123(parts ...string) string {
	base := sanitizeAndJoin(parts...)
	if base == "" {
		return defaultDNS1123Name
	}
	if len(base) <= maxDNS1123Length {
		return base
	}

	suffix := hashParts(parts...)

	prefixLen := maxDNS1123Length - len(suffix) - 1
	if prefixLen < 1 {
		prefixLen = maxDNS1123Length - len(suffix)
	}
	if prefixLen < 1 {
		if len(suffix) > maxDNS1123Length {
			return suffix[:maxDNS1123Length]
		}
		return suffix
	}

	prefix := base[:prefixLen]
	prefix = strings.TrimSuffix(prefix, "-")
	if len(prefix) == 0 {
		prefix = base[:prefixLen]
		prefix = strings.Trim(prefix, "-")
		if len(prefix) == 0 {
			prefix = defaultDNS1123Name
		}
	}

	result := fmt.Sprintf("%s-%s", prefix, suffix)
	if len(result) > maxDNS1123Length {
		result = result[:maxDNS1123Length]
		result = strings.TrimSuffix(result, "-")
		if len(result) == 0 {
			if len(suffix) > maxDNS1123Length {
				return suffix[:maxDNS1123Length]
			}
			return suffix
		}
	}
	return result
}

// ComposeDNS1123WithSuffix preserves a suffix while keeping the result DNS-1123 compliant.
// When truncation is required, it inserts a short hash before the suffix.
func ComposeDNS1123WithSuffix(base, suffix string) string {
	base = sanitizeDNS1123Segment(base)
	suffix = sanitizeDNS1123Segment(suffix)
	if suffix == "" {
		return ComposeDNS1123(base)
	}
	if base == "" {
		base = defaultDNS1123Name
	}
	name := base + "-" + suffix
	if len(name) <= maxDNS1123Length {
		return name
	}

	digest := hashParts(base, suffix)

	maxBase := maxDNS1123Length - len(suffix) - 1 - len(digest) - 1
	if maxBase < 1 {
		return ComposeDNS1123(suffix, digest)
	}
	prefix := base
	if len(prefix) > maxBase {
		prefix = prefix[:maxBase]
	}
	prefix = strings.TrimSuffix(prefix, "-")
	if prefix == "" {
		prefix = defaultDNS1123Name
	}
	return fmt.Sprintf("%s-%s-%s", prefix, digest, suffix)
}

func sanitizeAndJoin(parts ...string) string {
	sanitized := make([]string, 0, len(parts))
	for _, part := range parts {
		normalized := sanitizeDNS1123Segment(part)
		if normalized == "" {
			continue
		}
		sanitized = append(sanitized, normalized)
	}
	return strings.Join(sanitized, "-")
}

func sanitizeDNS1123Segment(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return ""
	}

	var out strings.Builder
	lastDash := false
	for _, r := range value {
		switch {
		case unicode.IsLower(r), unicode.IsDigit(r):
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

func hashParts(parts ...string) string {
	hasher := sha256.New()
	for _, part := range parts {
		_, _ = hasher.Write([]byte(part))
		_, _ = hasher.Write([]byte{0})
	}
	sum := hasher.Sum(nil)
	return hex.EncodeToString(sum[:])[:12]
}
