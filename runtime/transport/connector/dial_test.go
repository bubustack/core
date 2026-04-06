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

func TestCallWithTimeoutReturnsDeadlineErrorWhenCallbackIgnoresContext(t *testing.T) {
	release := make(chan struct{})
	defer close(release)

	start := time.Now()
	err := CallWithTimeout(context.Background(), 10*time.Millisecond, "call", func(context.Context) error {
		<-release
		return nil
	})
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("expected timeout error, got %v", err)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded classification, got %v", err)
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Fatalf("expected timeout return before callback release, got %s", elapsed)
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

func TestRecvWithTimeoutReturnsDeadlineErrorWhenCallbackIgnoresContext(t *testing.T) {
	release := make(chan struct{})
	defer close(release)

	start := time.Now()
	_, err := RecvWithTimeout(context.Background(), 10*time.Millisecond, "recv", func(context.Context) (string, error) {
		<-release
		return "", nil
	})
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("expected timeout error, got %v", err)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded classification, got %v", err)
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Fatalf("expected timeout return before callback release, got %s", elapsed)
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

func TestCallWithTimeoutPropagatesParentCancellationWhenCallbackIgnoresContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	release := make(chan struct{})
	defer close(release)

	done := make(chan error, 1)
	go func() {
		done <- CallWithTimeout(ctx, time.Second, "call", func(context.Context) error {
			<-release
			return nil
		})
	}()

	time.Sleep(10 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context canceled, got %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected cancellation to return before callback release")
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
