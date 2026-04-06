package templating

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"
)

const maxTemplateValidationJSONBytes = 1 << 20

// ValidateTemplateString validates a single template string against the scope.
func ValidateTemplateString(value string, scope ExpressionScope) error {
	if value == "" || !strings.Contains(value, "{{") {
		return nil
	}
	value = normalizeTemplateRoots(value)
	funcs := buildFuncMap(scope.AllowNow, scope.AllowRandom)
	tpl, err := template.New("validate").Funcs(funcs).Option("missingkey=error").Parse(value)
	if err != nil {
		return err
	}
	var errs []string
	walkTemplateNode(tpl.Root, scope, &errs)
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// ValidateJSONTemplates walks a JSON blob and validates any template strings.
func ValidateJSONTemplates(raw []byte, scope ExpressionScope) error {
	if len(raw) == 0 {
		return nil
	}
	if len(raw) > maxTemplateValidationJSONBytes {
		return fmt.Errorf("template JSON exceeds max validation size of %d bytes", maxTemplateValidationJSONBytes)
	}
	var node any
	if err := json.Unmarshal(raw, &node); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	var errs []string
	walkJSONTemplates(node, "$", scope, &errs)
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func walkJSONTemplates(node any, path string, scope ExpressionScope, errs *[]string) {
	switch typed := node.(type) {
	case map[string]any:
		if rawExpr, ok := typed[TemplateExprKey]; ok {
			expr, ok := rawExpr.(string)
			if !ok {
				*errs = append(*errs, fmt.Sprintf("%s.%s must be a string", path, TemplateExprKey))
			} else if err := ValidateTemplateString(expr, scope); err != nil {
				*errs = append(*errs, fmt.Sprintf("%s.%s: %v", path, TemplateExprKey, err))
			}
		}
		for key, value := range typed {
			if key == TemplateExprKey {
				continue
			}
			walkJSONTemplates(value, path+"."+key, scope, errs)
		}
	case []any:
		for i, value := range typed {
			walkJSONTemplates(value, fmt.Sprintf("%s[%d]", path, i), scope, errs)
		}
	case string:
		if err := ValidateTemplateString(typed, scope); err != nil {
			*errs = append(*errs, fmt.Sprintf("%s: %v", path, err))
		}
	}
}

func walkTemplateNode(node parse.Node, scope ExpressionScope, errs *[]string) {
	switch n := node.(type) {
	case *parse.ListNode:
		for _, child := range n.Nodes {
			walkTemplateNode(child, scope, errs)
		}
	case *parse.ActionNode:
		walkPipe(n.Pipe, scope, errs)
	case *parse.IfNode:
		walkPipe(n.Pipe, scope, errs)
		walkTemplateNode(n.List, scope, errs)
		if n.ElseList != nil {
			walkTemplateNode(n.ElseList, scope, errs)
		}
	case *parse.RangeNode:
		if !scope.AllowRange {
			*errs = append(*errs, fmt.Sprintf("range is not allowed for %s templates", scopeName(scope)))
		}
		walkPipe(n.Pipe, scope, errs)
		walkTemplateNode(n.List, scope, errs)
		if n.ElseList != nil {
			walkTemplateNode(n.ElseList, scope, errs)
		}
	case *parse.WithNode:
		walkPipe(n.Pipe, scope, errs)
		walkTemplateNode(n.List, scope, errs)
		if n.ElseList != nil {
			walkTemplateNode(n.ElseList, scope, errs)
		}
	case *parse.TemplateNode:
		// template inclusion not supported
		*errs = append(*errs, fmt.Sprintf("template inclusion is not supported (%s)", n.Name))
	}
}

func walkPipe(pipe *parse.PipeNode, scope ExpressionScope, errs *[]string) {
	if pipe == nil {
		return
	}
	for _, cmd := range pipe.Cmds {
		walkCommand(cmd, scope, errs)
	}
}

func walkCommand(cmd *parse.CommandNode, scope ExpressionScope, errs *[]string) {
	if cmd == nil {
		return
	}
	checkCommandRootUsage(cmd, scope, errs)
	for _, arg := range cmd.Args {
		walkArg(arg, scope, errs)
	}
}

func walkArg(arg parse.Node, scope ExpressionScope, errs *[]string) {
	switch n := arg.(type) {
	case *parse.FieldNode:
		checkRootUsage(n.Ident, scope, errs)
	case *parse.ChainNode:
		checkChainUsage(n, scope, errs)
	case *parse.VariableNode:
		// allow variables
	case *parse.PipeNode:
		walkPipe(n, scope, errs)
	case *parse.CommandNode:
		walkCommand(n, scope, errs)
	case *parse.BoolNode, *parse.DotNode, *parse.NumberNode, *parse.StringNode, *parse.NilNode:
		return
	case *parse.IdentifierNode:
		// function names are validated by parser (must exist in func map)
		return
	}
}

func checkChainUsage(node *parse.ChainNode, scope ExpressionScope, errs *[]string) {
	if node == nil {
		return
	}
	switch node.Node.(type) {
	case *parse.VariableNode:
		// skip variables
		return
	}

	root, _, ok := rootNameFromNode(node.Node)
	if !ok {
		return
	}
	var ident []string
	if root != "" {
		ident = append(ident, root)
	}
	ident = append(ident, node.Field...)
	if len(ident) == 0 {
		return
	}
	checkRootUsage(ident, scope, errs)
}

func checkRootUsage(ident []string, scope ExpressionScope, errs *[]string) {
	if len(ident) == 0 {
		return
	}
	root := ident[0]
	switch root {
	case RootInputs, RootSteps, RootPacket:
		if _, ok := scope.AllowedRoots[root]; !ok {
			*errs = append(*errs, fmt.Sprintf("context '%s' is not allowed for %s templates", root, scopeName(scope)))
		}
	default:
		*errs = append(*errs, fmt.Sprintf("unknown context '%s' is not allowed for %s templates", root, scopeName(scope)))
	}
}

func scopeName(scope ExpressionScope) string {
	if strings.TrimSpace(scope.Name) == "" {
		return "this"
	}
	return scope.Name
}

func checkCommandRootUsage(cmd *parse.CommandNode, scope ExpressionScope, errs *[]string) {
	if cmd == nil || len(cmd.Args) == 0 {
		return
	}
	ident, ok := cmd.Args[0].(*parse.IdentifierNode)
	if !ok {
		return
	}

	switch ident.Ident {
	case fnIndex, fnGet:
		if root, ok := rootFromAccessorArgs(cmd.Args[1:]); ok {
			checkRootUsage([]string{root}, scope, errs)
		}
	case fnDig:
		if root, ok := rootFromDigArgs(cmd.Args[1:]); ok {
			checkRootUsage([]string{root}, scope, errs)
		}
	}
}

func rootFromAccessorArgs(args []parse.Node) (string, bool) {
	if len(args) == 0 {
		return "", false
	}
	root, isDot, ok := rootNameFromNode(args[0])
	if !ok {
		return "", false
	}
	if root != "" {
		return root, true
	}
	if isDot && len(args) > 1 {
		if str, ok := args[1].(*parse.StringNode); ok {
			return str.Text, true
		}
	}
	return "", false
}

func rootFromDigArgs(args []parse.Node) (string, bool) {
	if len(args) == 0 {
		return "", false
	}
	root, isDot, ok := rootNameFromNode(args[len(args)-1])
	if ok && root != "" {
		return root, true
	}
	if isDot && len(args) > 0 {
		if str, ok := args[0].(*parse.StringNode); ok {
			return str.Text, true
		}
	}
	return "", false
}

func rootNameFromNode(node parse.Node) (root string, isDot bool, ok bool) {
	switch typed := node.(type) {
	case *parse.FieldNode:
		if len(typed.Ident) == 0 {
			return "", false, false
		}
		return typed.Ident[0], false, true
	case *parse.ChainNode:
		switch base := typed.Node.(type) {
		case *parse.FieldNode:
			if len(base.Ident) == 0 {
				return "", false, false
			}
			return base.Ident[0], false, true
		case *parse.DotNode:
			if len(typed.Field) == 0 {
				return "", true, true
			}
			return typed.Field[0], false, true
		case *parse.PipeNode:
			return rootNameFromPipe(base)
		case *parse.CommandNode:
			return rootNameFromCommand(base)
		}
	case *parse.PipeNode:
		return rootNameFromPipe(typed)
	case *parse.CommandNode:
		return rootNameFromCommand(typed)
	case *parse.DotNode:
		return "", true, true
	}
	return "", false, false
}

func rootNameFromPipe(pipe *parse.PipeNode) (root string, isDot bool, ok bool) {
	if pipe == nil || len(pipe.Cmds) == 0 {
		return "", false, false
	}
	return rootNameFromCommand(pipe.Cmds[0])
}

func rootNameFromCommand(cmd *parse.CommandNode) (root string, isDot bool, ok bool) {
	if cmd == nil || len(cmd.Args) == 0 {
		return "", false, false
	}

	if ident, ok := cmd.Args[0].(*parse.IdentifierNode); ok {
		switch ident.Ident {
		case fnIndex, fnGet:
			if root, ok := rootFromAccessorArgs(cmd.Args[1:]); ok {
				return root, false, true
			}
			return "", false, false
		case fnDig:
			if root, ok := rootFromDigArgs(cmd.Args[1:]); ok {
				return root, false, true
			}
			return "", false, false
		}
	}

	return rootNameFromNode(cmd.Args[0])
}
