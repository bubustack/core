package templating

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type pathSegment struct {
	key   string
	index *int
}

type errMissingKey struct {
	path string
}

func (e *errMissingKey) Error() string {
	return fmt.Sprintf("missing key at %s", e.path)
}

//nolint:gocyclo // Path parsing stays more readable as a single small state machine.
func parseSimplePath(expr string) (string, []pathSegment, bool) {
	trimmed := strings.TrimSpace(expr)
	if trimmed == "" {
		return "", nil, false
	}
	if !strings.HasPrefix(trimmed, ".") {
		return "", nil, false
	}
	trimmed = strings.TrimPrefix(trimmed, ".")
	if trimmed == "" {
		return "", nil, false
	}

	readIdent := func(s string, start int) (string, int) {
		i := start
		for i < len(s) {
			ch := s[i]
			if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' {
				i++
				continue
			}
			break
		}
		if i == start {
			return "", start
		}
		return s[start:i], i
	}

	root, idx := readIdent(trimmed, 0)
	if root == "" {
		return "", nil, false
	}
	segments := []pathSegment{}
	for idx < len(trimmed) {
		switch trimmed[idx] {
		case '.':
			idx++
			ident, next := readIdent(trimmed, idx)
			if ident == "" {
				return "", nil, false
			}
			segments = append(segments, pathSegment{key: ident})
			idx = next
		case '[':
			end := strings.IndexByte(trimmed[idx:], ']')
			if end == -1 {
				return "", nil, false
			}
			end = idx + end
			raw := strings.TrimSpace(trimmed[idx+1 : end])
			if raw == "" {
				return "", nil, false
			}
			if isQuotedIndexToken(raw) {
				key, err := strconv.Unquote(raw)
				if err != nil {
					return "", nil, false
				}
				segments = append(segments, pathSegment{key: key})
			} else if num, err := strconv.Atoi(raw); err == nil {
				segments = append(segments, pathSegment{index: &num})
			} else {
				return "", nil, false
			}
			idx = end + 1
		default:
			return "", nil, false
		}
	}

	return root, segments, true
}

func evalPath(root string, segments []pathSegment, vars map[string]any) (any, error) {
	current, ok := vars[root]
	if !ok {
		return nil, &errMissingKey{path: root}
	}
	path := root
	for _, seg := range segments {
		if seg.index != nil {
			idx := *seg.index
			switch v := current.(type) {
			case []any:
				if idx < 0 || idx >= len(v) {
					return nil, &errMissingKey{path: path}
				}
				current = v[idx]
			default:
				return nil, fmt.Errorf("cannot index into %T at %s", current, path)
			}
			path = fmt.Sprintf("%s[%d]", path, idx)
			continue
		}
		key := seg.key
		switch v := current.(type) {
		case map[string]any:
			val, exists := v[key]
			if !exists {
				return nil, &errMissingKey{path: path + "." + key}
			}
			current = val
		case map[any]any:
			val, exists := v[key]
			if !exists {
				return nil, &errMissingKey{path: path + "." + key}
			}
			current = val
		default:
			return nil, fmt.Errorf("cannot access key %q on %T at %s", key, current, path)
		}
		path = path + "." + key
	}
	return current, nil
}

func (e *Evaluator) evaluateSimplePath(expr string, vars map[string]any) (any, bool, error) {
	root, segments, ok := parseSimplePath(expr)
	if !ok {
		return nil, false, nil
	}
	if selector, ok := resolveStorageSelector(root, segments, vars); ok {
		return selector, true, nil
	}
	val, err := evalPath(root, segments, vars)
	if err != nil {
		var missing *errMissingKey
		if errors.As(err, &missing) && root == RootSteps {
			return nil, true, &ErrEvaluationBlocked{Reason: missing.Error()}
		}
		return nil, true, err
	}
	return val, true, nil
}

func (e *Evaluator) evaluateIndexPath(expr string, vars map[string]any) (any, bool, error) {
	inner, suffix, ok := splitIndexExpression(expr)
	if !ok {
		return nil, false, nil
	}

	root, key, ok := parseIndexCall(inner)
	if !ok {
		return nil, false, nil
	}

	segments, ok := parseSuffixSegments(suffix)
	if !ok {
		return nil, false, nil
	}

	fullSegments := make([]pathSegment, 0, 1+len(segments))
	fullSegments = append(fullSegments, pathSegment{key: key})
	fullSegments = append(fullSegments, segments...)

	if selector, ok := resolveStorageSelector(root, fullSegments, vars); ok {
		return selector, true, nil
	}
	val, err := evalPath(root, fullSegments, vars)
	if err != nil {
		var missing *errMissingKey
		if errors.As(err, &missing) && root == RootSteps {
			return nil, true, &ErrEvaluationBlocked{Reason: missing.Error()}
		}
		return nil, true, err
	}
	return val, true, nil
}

func splitIndexExpression(expr string) (string, string, bool) {
	trimmed := strings.TrimSpace(expr)
	if trimmed == "" {
		return "", "", false
	}
	if strings.HasPrefix(trimmed, "index ") {
		return trimmed, "", true
	}
	if !strings.HasPrefix(trimmed, "(") {
		return "", "", false
	}
	end := findMatchingParen(trimmed)
	if end <= 0 {
		return "", "", false
	}
	inner := strings.TrimSpace(trimmed[1:end])
	suffix := strings.TrimSpace(trimmed[end+1:])
	if !strings.HasPrefix(inner, "index ") {
		return "", "", false
	}
	if suffix == "" || strings.HasPrefix(suffix, ".") || strings.HasPrefix(suffix, "[") {
		return inner, suffix, true
	}
	return "", "", false
}

func findMatchingParen(s string) int {
	depth := 0
	inSingle := false
	inDouble := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if inSingle {
			if ch == '\'' {
				inSingle = false
			} else if ch == '\\' && i+1 < len(s) {
				i++
			}
			continue
		}
		if inDouble {
			if ch == '"' {
				inDouble = false
			} else if ch == '\\' && i+1 < len(s) {
				i++
			}
			continue
		}
		switch ch {
		case '\'':
			inSingle = true
		case '"':
			inDouble = true
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func parseIndexCall(expr string) (string, string, bool) {
	trimmed := strings.TrimSpace(expr)
	if !strings.HasPrefix(trimmed, "index") {
		return "", "", false
	}
	rest := strings.TrimSpace(trimmed[len("index"):])
	rootToken, rest := nextToken(rest)
	if rootToken == "" {
		return "", "", false
	}
	key, rest, ok := parseQuotedToken(rest)
	if !ok {
		return "", "", false
	}
	if strings.TrimSpace(rest) != "" {
		return "", "", false
	}
	root := strings.TrimPrefix(rootToken, ".")
	if root == "" {
		return "", "", false
	}
	if _, ok := templateRootNames[root]; !ok {
		return "", "", false
	}
	return root, key, true
}

func nextToken(s string) (string, string) {
	trimmed := strings.TrimLeft(s, " \t")
	if trimmed == "" {
		return "", ""
	}
	i := 0
	for i < len(trimmed) {
		if trimmed[i] == ' ' || trimmed[i] == '\t' {
			break
		}
		i++
	}
	return trimmed[:i], trimmed[i:]
}

func parseQuotedToken(s string) (string, string, bool) {
	trimmed := strings.TrimLeft(s, " \t")
	if trimmed == "" {
		return "", "", false
	}
	val, next, ok := parseQuotedAt(trimmed, 0)
	if !ok {
		return "", "", false
	}
	return val, trimmed[next:], true
}

//nolint:gocyclo // Suffix parsing must accept both field and bracket notation.
func parseSuffixSegments(suffix string) ([]pathSegment, bool) {
	trimmed := strings.TrimSpace(suffix)
	if trimmed == "" {
		return nil, true
	}
	var segments []pathSegment
	i := 0
	for i < len(trimmed) {
		switch trimmed[i] {
		case '.':
			i++
			start := i
			for i < len(trimmed) && isPathIdentChar(trimmed[i]) {
				i++
			}
			if start == i {
				return nil, false
			}
			segments = append(segments, pathSegment{key: trimmed[start:i]})
		case '[':
			if i+1 >= len(trimmed) {
				return nil, false
			}
			if trimmed[i+1] == '"' || trimmed[i+1] == '\'' {
				key, next, ok := parseQuotedAt(trimmed, i+1)
				if !ok {
					return nil, false
				}
				if next >= len(trimmed) || trimmed[next] != ']' {
					return nil, false
				}
				segments = append(segments, pathSegment{key: key})
				i = next + 1
				continue
			}
			j := i + 1
			for j < len(trimmed) && trimmed[j] >= '0' && trimmed[j] <= '9' {
				j++
			}
			if j == i+1 {
				return nil, false
			}
			if j >= len(trimmed) || trimmed[j] != ']' {
				return nil, false
			}
			idx, err := strconv.Atoi(trimmed[i+1 : j])
			if err != nil {
				return nil, false
			}
			segments = append(segments, pathSegment{index: &idx})
			i = j + 1
		default:
			return nil, false
		}
	}
	return segments, true
}

func parseQuotedAt(s string, start int) (string, int, bool) {
	if start >= len(s) {
		return "", 0, false
	}
	quote := s[start]
	if quote != '"' && quote != '\'' {
		return "", 0, false
	}
	i := start + 1
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			i += 2
			continue
		}
		if s[i] == quote {
			break
		}
		i++
	}
	if i >= len(s) {
		return "", 0, false
	}
	raw := s[start : i+1]
	val, err := strconv.Unquote(raw)
	if err != nil {
		return "", 0, false
	}
	return val, i + 1, true
}

func isQuotedIndexToken(raw string) bool {
	return (strings.HasPrefix(raw, "\"") && strings.HasSuffix(raw, "\"")) ||
		(strings.HasPrefix(raw, "'") && strings.HasSuffix(raw, "'"))
}

func isPathIdentChar(ch byte) bool {
	if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '-' {
		return true
	}
	return false
}

func resolveStorageSelector(root string, segments []pathSegment, vars map[string]any) (any, bool) {
	if root != RootSteps {
		return nil, false
	}
	outputIdx := -1
	for i, seg := range segments {
		if seg.index != nil {
			continue
		}
		if seg.key == "output" || seg.key == "outputs" {
			outputIdx = i
			break
		}
	}
	if outputIdx == -1 {
		return nil, false
	}
	val, err := evalPath(root, segments[:outputIdx+1], vars)
	if err != nil {
		return nil, false
	}
	outputMap, ok := val.(map[string]any)
	if !ok {
		return nil, false
	}
	refRaw, ok := outputMap[StorageRefKey]
	if !ok {
		return nil, false
	}
	refPath, ok := refRaw.(string)
	if !ok || strings.TrimSpace(refPath) == "" {
		return nil, false
	}
	path := buildStoragePath(segments[outputIdx+1:])
	if path == "" {
		return nil, false
	}
	return map[string]any{
		StorageRefKey:  refPath,
		StoragePathKey: path,
	}, true
}

func buildStoragePath(segments []pathSegment) string {
	if len(segments) == 0 {
		return ""
	}
	var b strings.Builder
	for _, seg := range segments {
		if seg.index != nil {
			b.WriteString("[")
			b.WriteString(strconv.Itoa(*seg.index))
			b.WriteString("]")
			continue
		}
		if b.Len() > 0 {
			b.WriteByte('.')
		}
		b.WriteString(seg.key)
	}
	return b.String()
}
