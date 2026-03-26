package templating

import (
	"sort"
	"strings"
	"text/template"
	"text/template/parse"
)

// ExtractStepReferences parses a Go template string and returns the names of
// steps referenced via `(index .steps "name")` or `.steps.name` patterns.
// Returns nil if no step references are found.
func ExtractStepReferences(value string) []string {
	refs, _ := ExtractStepReferencesWithError(value)
	return refs
}

// ExtractStepReferencesWithError returns step references or a parse error.
func ExtractStepReferencesWithError(value string) ([]string, error) {
	if !strings.Contains(value, "{{") {
		return nil, nil
	}
	funcs := buildFuncMap(true, true) // permissive for parsing
	tpl, err := template.New("refs").Funcs(funcs).Option("missingkey=zero").Parse(value)
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	walkStepRefs(tpl.Root, seen)
	if len(seen) == 0 {
		return nil, nil
	}
	result := make([]string, 0, len(seen))
	for name := range seen {
		result = append(result, name)
	}
	sort.Strings(result)
	return result, nil
}

func walkStepRefs(node parse.Node, seen map[string]struct{}) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *parse.ListNode:
		for _, child := range n.Nodes {
			walkStepRefs(child, seen)
		}
	case *parse.ActionNode:
		walkStepRefsPipe(n.Pipe, seen)
	case *parse.IfNode:
		walkStepRefsPipe(n.Pipe, seen)
		walkStepRefs(n.List, seen)
		if n.ElseList != nil {
			walkStepRefs(n.ElseList, seen)
		}
	case *parse.RangeNode:
		walkStepRefsPipe(n.Pipe, seen)
		walkStepRefs(n.List, seen)
		if n.ElseList != nil {
			walkStepRefs(n.ElseList, seen)
		}
	case *parse.WithNode:
		walkStepRefsPipe(n.Pipe, seen)
		walkStepRefs(n.List, seen)
		if n.ElseList != nil {
			walkStepRefs(n.ElseList, seen)
		}
	}
}

func walkStepRefsPipe(pipe *parse.PipeNode, seen map[string]struct{}) {
	if pipe == nil {
		return
	}
	for _, cmd := range pipe.Cmds {
		walkStepRefsCmd(cmd, seen)
	}
}

//nolint:gocyclo // Template AST matching is clearer when the node cases stay together.
func walkStepRefsCmd(cmd *parse.CommandNode, seen map[string]struct{}) {
	if cmd == nil {
		return
	}
	// Check for `index .steps "name"` pattern
	if len(cmd.Args) >= 3 {
		if ident, ok := cmd.Args[0].(*parse.IdentifierNode); ok && ident.Ident == "index" {
			if field, ok := cmd.Args[1].(*parse.FieldNode); ok {
				if len(field.Ident) >= 1 && field.Ident[0] == RootSteps {
					if str, ok := cmd.Args[2].(*parse.StringNode); ok {
						seen[str.Text] = struct{}{}
					}
				}
			}
		}
	}
	// Check for `.steps.name` field access and chained expressions
	for _, arg := range cmd.Args {
		switch n := arg.(type) {
		case *parse.FieldNode:
			if len(n.Ident) >= 2 && n.Ident[0] == RootSteps {
				seen[n.Ident[1]] = struct{}{}
			}
		case *parse.ChainNode:
			if field, ok := n.Node.(*parse.FieldNode); ok {
				if len(field.Ident) >= 1 && field.Ident[0] == RootSteps && len(n.Field) >= 1 {
					seen[n.Field[0]] = struct{}{}
				}
			}
			// Handle `(index .steps "name").field` — ChainNode wrapping a PipeNode
			if pipe, ok := n.Node.(*parse.PipeNode); ok {
				walkStepRefsPipe(pipe, seen)
			}
		case *parse.PipeNode:
			walkStepRefsPipe(n, seen)
		}
	}
}
