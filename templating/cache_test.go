package templating

import (
	"context"
	"errors"
	"testing"
	"text/template"
	"time"
)

func TestTemplateCacheStopIsIdempotent(t *testing.T) {
	cache := NewTemplateCache(&CacheConfig{MaxSize: 1, TTL: time.Second})
	cache.Stop()
	cache.Stop()
}

func TestTemplateCacheGetOrParseHonorsCanceledContextBeforeParse(t *testing.T) {
	cache := NewTemplateCache(DefaultCacheConfig())
	t.Cleanup(cache.Stop)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	called := false
	_, err := cache.GetOrParse(ctx, "key", func() (*template.Template, error) {
		called = true
		return template.New("x").Parse("ok")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
	if called {
		t.Fatalf("expected parseFn not to run for canceled context")
	}
}

func TestTemplateCacheGetOrParseHonorsCanceledContextAfterParse(t *testing.T) {
	cache := NewTemplateCache(DefaultCacheConfig())
	t.Cleanup(cache.Stop)

	ctx, cancel := context.WithCancel(context.Background())
	_, err := cache.GetOrParse(ctx, "key", func() (*template.Template, error) {
		cancel()
		return template.New("x").Parse("ok")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}
