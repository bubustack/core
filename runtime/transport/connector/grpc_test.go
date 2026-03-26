package connector

import (
	"testing"
	"time"

	"github.com/bubustack/core/contracts"
)

func TestDialTimeoutUsesEnvOverride(t *testing.T) {
	env := EnvFunc(func(key string) string {
		if key == contracts.GRPCDialTimeoutEnv {
			return "7s"
		}
		return ""
	})
	if got := DialTimeout(env, time.Second); got != 7*time.Second {
		t.Fatalf("expected 7s, got %s", got)
	}
}
