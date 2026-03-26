package connector

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestWaitForReadyRejectsNilConnection(t *testing.T) {
	if err := WaitForReady(context.Background(), nil); err == nil {
		t.Fatalf("expected nil connection error")
	}
}

func TestCallWithTimeoutReturnsDeadlineError(t *testing.T) {
	ctx := context.Background()
	err := CallWithTimeout(ctx, 10*time.Millisecond, "call", func(callCtx context.Context) error {
		<-callCtx.Done()
		return callCtx.Err()
	})
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("expected timeout error, got %v", err)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded classification, got %v", err)
	}
}

func TestRecvWithTimeoutReturnsDeadlineError(t *testing.T) {
	ctx := context.Background()
	_, err := RecvWithTimeout(ctx, 10*time.Millisecond, "recv", func(callCtx context.Context) (string, error) {
		<-callCtx.Done()
		return "", callCtx.Err()
	})
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("expected timeout error, got %v", err)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded classification, got %v", err)
	}
}

func TestCallWithTimeoutPropagatesParentCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := CallWithTimeout(ctx, time.Second, "call", func(callCtx context.Context) error {
		<-callCtx.Done()
		return callCtx.Err()
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}

func TestRecvWithTimeoutReturnsValueWithoutTimeout(t *testing.T) {
	value, err := RecvWithTimeout(context.Background(), time.Second, "recv", func(context.Context) (string, error) {
		return "ok", nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if value != "ok" {
		t.Fatalf("expected ok, got %q", value)
	}
}
