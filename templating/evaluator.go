package templating

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"strings"
	"text/template"
	"text/template/parse"
)

// Evaluator resolves BubuStack template expressions against JSON-like values.
type Evaluator struct {
	cfg   Config
	cache *TemplateCache
	funcs template.FuncMap
}

const maxResolveDepth = 128

// New creates an Evaluator with the supplied configuration.
func New(cfg Config) (*Evaluator, error) {
	if cfg.EvaluationTimeout < 0 {
		return nil, fmt.Errorf("evaluation timeout must be >= 0")
	}
	if cfg.MaxOutputBytes < 0 {
		return nil, fmt.Errorf("max output bytes must be >= 0")
	}
	funcs := buildFuncMap(!cfg.Deterministic, !cfg.Deterministic)
	return &Evaluator{
		cfg:   cfg,
		cache: NewTemplateCache(DefaultCacheConfig()),
		funcs: funcs,
	}, nil
}

// Close releases the evaluator's background resources.
func (e *Evaluator) Close() {
	if e == nil || e.cache == nil {
		return
	}
	e.cache.Stop()
}

// ResolveValue resolves template expressions within any JSON-like value.
// This is a general-purpose helper for SDK callers that need to resolve templates
// outside of an object-only context.
func (e *Evaluator) ResolveValue(ctx context.Context, value any, vars map[string]any) (any, error) {
	return e.resolveValue(ctx, value, vars)
}

// ResolveTemplateString resolves a single template string with the provided vars.
func (e *Evaluator) ResolveTemplateString(ctx context.Context, value string, vars map[string]any) (any, error) {
	return e.resolveTemplateString(ctx, value, vars)
}

// ResolveWithInputs resolves template expressions inside a JSON-like map.
func (e *Evaluator) ResolveWithInputs(
	ctx context.Context,
	with map[string]any,
	vars map[string]any,
) (map[string]any, error) {
	if with == nil {
		return nil, nil
	}
	resolved, err := e.resolveValue(ctx, with, vars)
	if err != nil {
		return nil, err
	}
	out, ok := resolved.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("resolved with block must be an object, got %T", resolved)
	}
	return out, nil
}

// EvaluateCondition evaluates a boolean template.
func (e *Evaluator) EvaluateCondition(ctx context.Context, expr string, vars map[string]any) (bool, error) {
	if strings.TrimSpace(expr) == "" {
		return true, nil
	}
	val, err := e.resolveTemplateString(ctx, expr, vars)
	if err != nil {
		return false, err
	}
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		s := strings.TrimSpace(strings.ToLower(v))
		switch s {
		case "true", "1", "yes", "y":
			return true, nil
		case "false", "0", "no", "n", "":
			return false, nil
		default:
			return false, fmt.Errorf("condition did not resolve to boolean: %q", v)
		}
	default:
		return false, fmt.Errorf("condition did not resolve to boolean, got %T", val)
	}
}

func mergeTemplateVars(base map[string]any, overrides map[string]any) map[string]any {
	if len(overrides) == 0 {
		return base
	}
	merged := make(map[string]any, len(base)+len(overrides))
	maps.Copy(merged, base)
	maps.Copy(merged, overrides)
	return merged
}

func (e *Evaluator) resolveValue(ctx context.Context, value any, vars map[string]any) (any, error) {
	return e.resolveValueWithState(ctx, value, vars, newResolveState(), 0)
}

func (e *Evaluator) resolveValueWithState(
	ctx context.Context,
	value any,
	vars map[string]any,
	state *resolveState,
	depth int,
) (any, error) {
	if err := ctxErr(ctx); err != nil {
		return nil, err
	}
	if depth > state.maxDepth {
		return nil, fmt.Errorf("template value exceeded max resolve depth %d", state.maxDepth)
	}

	switch typed := value.(type) {
	case map[string]any:
		return e.resolveMapValue(ctx, typed, vars, state, depth)
	case []any:
		return e.resolveSliceValue(ctx, typed, vars, state, depth)
	case string:
		return e.resolveTemplateString(ctx, typed, vars)
	default:
		return value, nil
	}
}

func (e *Evaluator) resolveMapValue(
	ctx context.Context,
	value map[string]any,
	vars map[string]any,
	state *resolveState,
	depth int,
) (any, error) {
	if err := state.enter(value); err != nil {
		return nil, err
	}
	defer state.leave(value)

	if exprRaw, ok := value[TemplateExprKey]; ok {
		expr, ok := exprRaw.(string)
		if !ok {
			return nil, fmt.Errorf("%s must be a string", TemplateExprKey)
		}

		mergedVars, err := e.resolveTemplateVars(ctx, value, vars, state, depth)
		if err != nil {
			return nil, err
		}
		return e.resolveTemplateString(ctx, expr, mergedVars)
	}

	result := make(map[string]any, len(value))
	for key, entry := range value {
		resolved, err := e.resolveMapEntry(ctx, key, entry, vars, state, depth)
		if err != nil {
			return nil, err
		}
		result[key] = resolved
	}
	return result, nil
}

func (e *Evaluator) resolveTemplateVars(
	ctx context.Context,
	value map[string]any,
	vars map[string]any,
	state *resolveState,
	depth int,
) (map[string]any, error) {
	mergedVars := vars
	rawVars, ok := value[TemplateVarsKey]
	if !ok {
		return mergedVars, nil
	}
	resolvedVars, err := e.resolveValueWithState(ctx, rawVars, vars, state, depth+1)
	if err != nil {
		return nil, fmt.Errorf("template vars: %w", err)
	}
	if resolvedVars == nil {
		return mergedVars, nil
	}
	varsMap, ok := resolvedVars.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s must be an object", TemplateVarsKey)
	}
	return mergeTemplateVars(vars, varsMap), nil
}

func (e *Evaluator) resolveMapEntry(
	ctx context.Context,
	key string,
	value any,
	vars map[string]any,
	state *resolveState,
	depth int,
) (any, error) {
	resolved, err := e.resolveValueWithState(ctx, value, vars, state, depth+1)
	if err != nil {
		if key == TemplateVarsKey {
			return nil, fmt.Errorf("template vars: %w", err)
		}
		return nil, err
	}
	return resolved, nil
}

func (e *Evaluator) resolveSliceValue(
	ctx context.Context,
	value []any,
	vars map[string]any,
	state *resolveState,
	depth int,
) (any, error) {
	if err := state.enter(value); err != nil {
		return nil, err
	}
	defer state.leave(value)

	result := make([]any, len(value))
	for i, elem := range value {
		resolved, err := e.resolveValueWithState(ctx, elem, vars, state, depth+1)
		if err != nil {
			return nil, fmt.Errorf("index %d: %w", i, err)
		}
		result[i] = resolved
	}
	return result, nil
}

func (e *Evaluator) resolveTemplateString(ctx context.Context, value string, vars map[string]any) (any, error) {
	if !strings.Contains(value, "{{") {
		return value, nil
	}
	normalized := normalizeTemplateRoots(value)
	trimmed := strings.TrimSpace(normalized)
	if isSingleActionTemplate(trimmed) {
		expr := strings.TrimSpace(trimmed[2 : len(trimmed)-2])
		if expr != "" {
			if val, ok, err := e.evaluateSimplePath(expr, vars); ok || err != nil {
				if err != nil {
					return nil, err
				}
				return val, nil
			}
			if val, ok, err := e.evaluateIndexPath(expr, vars); ok || err != nil {
				if err != nil {
					return nil, err
				}
				return val, nil
			}
		}
	}

	// Before falling through to Go template engine, check if step vars
	// contain storage refs. Compound templates (with pipes/functions)
	// crash on storage ref maps — defer to SDK.
	if HasStorageRefsInStepVars(normalized, vars) {
		return nil, &ErrOffloadedDataUsage{
			Reason: "compound template references step output containing storage ref",
		}
	}

	rendered, err := e.renderTemplate(ctx, normalized, vars)
	if err != nil {
		return nil, err
	}
	if isSingleActionTemplate(trimmed) {
		if parsed, ok := parseJSONValue(rendered); ok {
			return parsed, nil
		}
	}
	return rendered, nil
}

func (e *Evaluator) renderTemplate(ctx context.Context, text string, vars map[string]any) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if e.cfg.EvaluationTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.cfg.EvaluationTimeout)
		defer cancel()
	}

	key := hashTemplate(text)
	tpl, err := e.cache.GetOrParse(ctx, key, func() (*template.Template, error) {
		return template.New("bubu").Funcs(e.funcs).Option("missingkey=zero").Parse(text)
	})
	if err != nil {
		return "", err
	}

	writer := newLimitedBuffer(ctx, e.cfg.MaxOutputBytes)
	if err := tpl.Execute(writer, vars); err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return "", fmt.Errorf("template evaluation timed out after %s", e.cfg.EvaluationTimeout)
		case errors.Is(err, context.Canceled):
			return "", fmt.Errorf("template evaluation canceled: %w", ctx.Err())
		default:
			return "", err
		}
	}
	if err := ctx.Err(); err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return "", fmt.Errorf("template evaluation timed out after %s", e.cfg.EvaluationTimeout)
		default:
			return "", fmt.Errorf("template evaluation canceled: %w", err)
		}
	}
	return writer.String(), nil
}

type limitedBuffer struct {
	buf bytes.Buffer
	max int
	ctx context.Context
}

func newLimitedBuffer(ctx context.Context, max int) *limitedBuffer {
	return &limitedBuffer{max: max, ctx: ctx}
}

func (l *limitedBuffer) Write(p []byte) (int, error) {
	if err := ctxErr(l.ctx); err != nil {
		return 0, err
	}
	if l.max > 0 && l.buf.Len()+len(p) > l.max {
		return 0, fmt.Errorf("template output exceeds max bytes (current=%d, writing=%d, max=%d)", l.buf.Len(), len(p), l.max)
	}
	return l.buf.Write(p)
}

func (l *limitedBuffer) String() string {
	return l.buf.String()
}

func isSingleActionTemplate(value string) bool {
	trimmed := strings.TrimSpace(value)
	if !strings.HasPrefix(trimmed, "{{") || !strings.HasSuffix(trimmed, "}}") {
		return false
	}
	tpl, err := template.New("check").Parse(trimmed)
	if err != nil {
		return false
	}
	if tpl.Tree == nil || tpl.Root == nil {
		return false
	}
	nodes := tpl.Root.Nodes
	if len(nodes) != 1 {
		return false
	}
	_, ok := nodes[0].(*parse.ActionNode)
	return ok
}

type resolveState struct {
	visiting map[containerVisit]struct{}
	maxDepth int
}

type containerVisit struct {
	kind reflect.Kind
	ptr  uintptr
}

func newResolveState() *resolveState {
	return &resolveState{
		visiting: make(map[containerVisit]struct{}),
		maxDepth: maxResolveDepth,
	}
}

func (s *resolveState) enter(value any) error {
	if s == nil {
		return nil
	}
	visit, ok := containerVisitFor(value)
	if !ok {
		return nil
	}
	if _, exists := s.visiting[visit]; exists {
		return fmt.Errorf("cyclic template value detected")
	}
	s.visiting[visit] = struct{}{}
	return nil
}

func (s *resolveState) leave(value any) {
	if s == nil {
		return
	}
	if visit, ok := containerVisitFor(value); ok {
		delete(s.visiting, visit)
	}
}

func containerVisitFor(value any) (containerVisit, bool) {
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return containerVisit{}, false
	}
	switch rv.Kind() {
	case reflect.Map, reflect.Slice:
		if rv.IsNil() {
			return containerVisit{}, false
		}
		ptr := rv.Pointer()
		if ptr == 0 {
			return containerVisit{}, false
		}
		return containerVisit{kind: rv.Kind(), ptr: ptr}, true
	default:
		return containerVisit{}, false
	}
}
