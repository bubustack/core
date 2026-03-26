package connector

import (
	"testing"
	"time"

	"github.com/bubustack/core/contracts"
)

func TestRuntimeTunablesFromEnv_Defaults(t *testing.T) {
	env := EnvFunc(func(string) string { return "" })
	tunables := RuntimeTunablesFromEnv(env, RuntimeTunables{
		MessageTimeout:     2 * time.Second,
		ChannelSendTimeout: 500 * time.Millisecond,
		HangTimeout:        5 * time.Second,
		ChannelBufferSize:  32,
	})

	if tunables.MessageTimeout != 2*time.Second {
		t.Fatalf("expected message timeout fallback, got %s", tunables.MessageTimeout)
	}
	if tunables.ChannelSendTimeout != 500*time.Millisecond {
		t.Fatalf("expected channel send timeout fallback, got %s", tunables.ChannelSendTimeout)
	}
	if tunables.HangTimeout != 5*time.Second {
		t.Fatalf("expected hang timeout fallback, got %s", tunables.HangTimeout)
	}
	if tunables.ChannelBufferSize != 32 {
		t.Fatalf("expected channel buffer size fallback 32, got %d", tunables.ChannelBufferSize)
	}
}

func TestRuntimeTunablesFromEnv_EnvOverrides(t *testing.T) {
	env := EnvFunc(func(key string) string {
		switch key {
		case contracts.GRPCMessageTimeoutEnv:
			return "750ms"
		case contracts.GRPCChannelSendTimeoutEnv:
			return "3s"
		case contracts.GRPCHangTimeoutEnv:
			return "9s"
		case contracts.GRPCChannelBufferSizeEnv:
			return "64"
		default:
			return ""
		}
	})

	tunables := RuntimeTunablesFromEnv(env, RuntimeTunables{})
	if tunables.MessageTimeout != 750*time.Millisecond {
		t.Fatalf("expected env message timeout, got %s", tunables.MessageTimeout)
	}
	if tunables.ChannelSendTimeout != 3*time.Second {
		t.Fatalf("expected env channel send timeout, got %s", tunables.ChannelSendTimeout)
	}
	if tunables.HangTimeout != 9*time.Second {
		t.Fatalf("expected env hang timeout, got %s", tunables.HangTimeout)
	}
	if tunables.ChannelBufferSize != 64 {
		t.Fatalf("expected env channel buffer size 64, got %d", tunables.ChannelBufferSize)
	}
}
