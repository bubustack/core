package templating

import "time"

// Config controls template evaluation behavior.
type Config struct {
	EvaluationTimeout time.Duration
	MaxOutputBytes    int
	// Deterministic disables non-deterministic helpers (now/rand).
	Deterministic bool
}

// CacheConfig controls template caching.
type CacheConfig struct {
	MaxSize int
	TTL     time.Duration
}

// DefaultCacheConfig returns the production-oriented cache defaults.
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		MaxSize: 1000,
		TTL:     30 * time.Minute,
	}
}

// ExpressionScope defines allowed template contexts for validation.
type ExpressionScope struct {
	Name         string
	AllowNow     bool
	AllowRandom  bool
	AllowedRoots map[string]struct{}
	AllowRange   bool
}

// NewExpressionScope builds a scope with provided roots.
func NewExpressionScope(
	name string,
	allowNow bool,
	allowRandom bool,
	allowRange bool,
	roots ...string,
) ExpressionScope {
	allowed := make(map[string]struct{}, len(roots))
	for _, root := range roots {
		allowed[root] = struct{}{}
	}
	return ExpressionScope{
		Name:         name,
		AllowNow:     allowNow,
		AllowRandom:  allowRandom,
		AllowedRoots: allowed,
		AllowRange:   allowRange,
	}
}
