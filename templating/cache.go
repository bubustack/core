package templating

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"text/template"
	"time"
)

// CachedTemplate holds a parsed template and metadata.
type CachedTemplate struct {
	Template *template.Template
	CachedAt time.Time
	LastUsed time.Time
	HitCount int64
	Hash     string
}

// TemplateCache caches parsed templates.
type TemplateCache struct {
	cache           map[string]*CachedTemplate
	mu              sync.RWMutex
	maxSize         int
	ttl             time.Duration
	cleanupInterval time.Duration
	stopCh          chan struct{}
	stopOnce        sync.Once
}

// NewTemplateCache builds a cache for parsed templates.
func NewTemplateCache(cfg *CacheConfig) *TemplateCache {
	if cfg == nil {
		cfg = DefaultCacheConfig()
	}
	c := &TemplateCache{
		cache:           make(map[string]*CachedTemplate),
		maxSize:         cfg.MaxSize,
		ttl:             cfg.TTL,
		cleanupInterval: deriveCleanupInterval(cfg.TTL),
		stopCh:          make(chan struct{}),
	}
	if c.ttl > 0 {
		go c.cleanupRoutine()
	}
	return c
}

// Stop terminates the background cleanup loop.
func (c *TemplateCache) Stop() {
	if c == nil {
		return
	}
	c.stopOnce.Do(func() {
		close(c.stopCh)
	})
}

// GetOrParse returns a cached template or parses and stores a new one.
func (c *TemplateCache) GetOrParse(
	ctx context.Context,
	key string,
	parseFn func() (*template.Template, error),
) (*template.Template, error) {
	if c == nil {
		return parseFn()
	}
	if err := ctxErr(ctx); err != nil {
		return nil, err
	}

	if tpl, ok := c.lookup(key); ok {
		return tpl, nil
	}

	tpl, err := parseFn()
	if err != nil {
		return nil, err
	}
	if err := ctxErr(ctx); err != nil {
		return nil, err
	}

	now := time.Now()
	cached := &CachedTemplate{
		Template: tpl,
		CachedAt: now,
		LastUsed: now,
		Hash:     key,
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if existing, ok := c.cache[key]; ok && !c.isExpired(existing) {
		existing.HitCount++
		existing.LastUsed = now
		return existing.Template, nil
	}
	if len(c.cache) >= c.maxSize && c.maxSize > 0 {
		c.evictLRU()
	}
	c.cache[key] = cached
	return tpl, nil
}

func (c *TemplateCache) isExpired(entry *CachedTemplate) bool {
	if c.ttl <= 0 {
		return false
	}
	return time.Since(entry.CachedAt) > c.ttl
}

func (c *TemplateCache) evictLRU() {
	var oldestKey string
	var oldest time.Time
	first := true
	for key, entry := range c.cache {
		if first || entry.LastUsed.Before(oldest) {
			oldest = entry.LastUsed
			oldestKey = key
			first = false
		}
	}
	if oldestKey != "" {
		delete(c.cache, oldestKey)
	}
}

func (c *TemplateCache) cleanupRoutine() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCh:
			return
		}
	}
}

func (c *TemplateCache) cleanup() {
	if c == nil || c.ttl <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, entry := range c.cache {
		if c.isExpired(entry) {
			delete(c.cache, key)
		}
	}
}

func hashTemplate(text string) string {
	sum := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", sum)
}

func (c *TemplateCache) lookup(key string) (*template.Template, bool) {
	c.mu.RLock()
	cached, ok := c.cache[key]
	if !ok || c.isExpired(cached) {
		c.mu.RUnlock()
		return nil, false
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	cached, ok = c.cache[key]
	if !ok || c.isExpired(cached) {
		return nil, false
	}
	cached.HitCount++
	cached.LastUsed = time.Now()
	return cached.Template, true
}

func deriveCleanupInterval(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return 5 * time.Minute
	}
	interval := ttl / 2
	if interval < time.Second {
		return time.Second
	}
	if interval > 5*time.Minute {
		return 5 * time.Minute
	}
	return interval
}

func ctxErr(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	return ctx.Err()
}
