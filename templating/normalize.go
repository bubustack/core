package templating

import (
	"fmt"
	"strconv"
	"strings"
)

var templateRootNames = map[string]struct{}{
	RootInputs: {},
	RootSteps:  {},
	RootPacket: {},
}

func normalizeTemplateRoots(text string) string {
	if !strings.Contains(text, "{{") {
		return text
	}
	var out strings.Builder
	i := 0
	for {
		start := strings.Index(text[i:], "{{")
		if start < 0 {
			out.WriteString(text[i:])
			break
		}
		start += i
		out.WriteString(text[i:start])

		actionStart := start + 2
		hasTrimLeft := actionStart < len(text) && text[actionStart] == '-'
		if hasTrimLeft {
			out.WriteString("{{-")
			actionStart++
		} else {
			out.WriteString("{{")
		}

		end, ok := findActionEnd(text, actionStart)
		if !ok {
			out.WriteString(text[actionStart:])
			break
		}
		actionEnd := end
		hasTrimRight := actionEnd > actionStart && text[actionEnd-1] == '-'
		if hasTrimRight {
			actionEnd--
		}

		action := text[actionStart:actionEnd]
		trimmed := strings.TrimSpace(action)
		if strings.HasPrefix(trimmed, "/*") {
			out.WriteString(action)
		} else {
			normalized := normalizeActionRoots(action)
			normalized = normalizeSubscriptAccess(normalized)
			out.WriteString(normalized)
		}

		if hasTrimRight {
			out.WriteString("-}}")
		} else {
			out.WriteString("}}")
		}

		i = end + 2
	}
	return out.String()
}

//nolint:gocyclo // String-literal aware token scanning is clearer as one state machine.
func normalizeActionRoots(action string) string {
	var out strings.Builder
	inSingle := false
	inDouble := false
	inRaw := false
	var prev byte
	for i := 0; i < len(action); {
		ch := action[i]
		if inRaw {
			out.WriteByte(ch)
			if ch == '`' {
				inRaw = false
			}
			prev = ch
			i++
			continue
		}
		if inSingle {
			out.WriteByte(ch)
			if ch == '\\' && i+1 < len(action) {
				out.WriteByte(action[i+1])
				prev = action[i+1]
				i += 2
				continue
			}
			if ch == '\'' {
				inSingle = false
			}
			prev = ch
			i++
			continue
		}
		if inDouble {
			out.WriteByte(ch)
			if ch == '\\' && i+1 < len(action) {
				out.WriteByte(action[i+1])
				prev = action[i+1]
				i += 2
				continue
			}
			if ch == '"' {
				inDouble = false
			}
			prev = ch
			i++
			continue
		}

		switch ch {
		case '\'':
			inSingle = true
			out.WriteByte(ch)
			prev = ch
			i++
			continue
		case '"':
			inDouble = true
			out.WriteByte(ch)
			prev = ch
			i++
			continue
		case '`':
			inRaw = true
			out.WriteByte(ch)
			prev = ch
			i++
			continue
		}

		if isIdentStart(ch) {
			start := i
			for i < len(action) && isIdentChar(action[i]) {
				i++
			}
			token := action[start:i]
			if _, ok := templateRootNames[token]; ok {
				if prev != '.' && prev != '$' && !isIdentChar(prev) {
					out.WriteByte('.')
				}
			}
			out.WriteString(token)
			if len(token) > 0 {
				prev = token[len(token)-1]
			}
			continue
		}

		out.WriteByte(ch)
		prev = ch
		i++
	}
	return out.String()
}

func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isIdentChar(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9')
}

// normalizeSubscriptAccess converts bracket string-key access to Go template index calls.
// For example: .steps['fetch-feed'].output → (index .steps "fetch-feed").output
// Applied iteratively to handle nested brackets like obj['a']['b'].
func normalizeSubscriptAccess(s string) string {
	for {
		n := applyOneSubscriptConversion(s)
		if n == s {
			return s
		}
		s = n
	}
}

// applyOneSubscriptConversion replaces the leftmost ['key'] or ["key"] subscript
// with an equivalent (index receiver "key") Go template expression.
func applyOneSubscriptConversion(s string) string {
	inSingle := false
	inDouble := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if inSingle {
			if ch == '\'' {
				inSingle = false
			}
			continue
		}
		if inDouble {
			if ch == '"' {
				inDouble = false
			}
			continue
		}
		switch ch {
		case '\'':
			inSingle = true
		case '"':
			inDouble = true
		case '[':
			if i+1 >= len(s) {
				continue
			}
			key, next, ok := parseQuotedAt(s, i+1)
			if !ok || next >= len(s) || s[next] != ']' {
				continue
			}
			recvStart := subscriptReceiverStart(s, i)
			if recvStart == i {
				continue
			}
			recv := s[recvStart:i]
			return s[:recvStart] + fmt.Sprintf(`(index %s %s)`, recv, strconv.Quote(key)) + s[next+1:]
		}
	}
	return s
}

//nolint:gocyclo // Template action scanning keeps quote/comment handling in one place.
func findActionEnd(text string, start int) (int, bool) {
	inSingle := false
	inDouble := false
	inRaw := false
	inComment := false
	allowComment := true

	for i := start; i < len(text)-1; i++ {
		ch := text[i]
		next := text[i+1]

		if inComment {
			if ch == '*' && next == '/' {
				inComment = false
				i++
			}
			continue
		}
		if inRaw {
			if ch == '`' {
				inRaw = false
			}
			continue
		}
		if inSingle {
			if ch == '\\' && i+1 < len(text) {
				i++
				continue
			}
			if ch == '\'' {
				inSingle = false
			}
			continue
		}
		if inDouble {
			if ch == '\\' && i+1 < len(text) {
				i++
				continue
			}
			if ch == '"' {
				inDouble = false
			}
			continue
		}

		if allowComment && (ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r') {
			continue
		}
		if allowComment && ch == '/' && next == '*' {
			inComment = true
			i++
			continue
		}
		allowComment = false

		switch ch {
		case '\'':
			inSingle = true
		case '"':
			inDouble = true
		case '`':
			inRaw = true
		case '}':
			if next == '}' {
				return i, true
			}
		}
	}
	return -1, false
}

// subscriptReceiverStart scans backwards from pos to find the start of the
// expression that immediately precedes a subscript access.
//
//nolint:gocyclo // Backward scan must track nested parens while preserving parser behavior.
func subscriptReceiverStart(s string, pos int) int {
	i := pos - 1
	for i >= 0 {
		ch := s[i]
		// Dotted identifier path characters (including hyphen for step names)
		if ch == '.' || ch == '_' || ch == '-' ||
			(ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') {
			i--
			continue
		}
		// Closing paren: find the matching open paren
		if ch == ')' {
			depth := 1
			i--
			for i >= 0 && depth > 0 {
				switch s[i] {
				case ')':
					depth++
				case '(':
					depth--
				}
				i--
			}
			continue
		}
		break
	}
	start := i + 1
	for start < pos && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	return start
}
