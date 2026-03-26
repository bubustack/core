package connector

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// WaitForReady blocks until the connection reports connectivity.Ready or the context expires.
func WaitForReady(ctx context.Context, conn *grpc.ClientConn) error {
	if conn == nil {
		return fmt.Errorf("connection is nil")
	}
	for {
		state := conn.GetState()
		if state == connectivity.Ready {
			return nil
		}
		conn.Connect()
		if !conn.WaitForStateChange(ctx, state) {
			if err := ctx.Err(); err != nil {
				return err
			}
			return fmt.Errorf("connection stuck in %s state", state)
		}
	}
}

// CallWithTimeout executes fn with a derived context that carries the timeout.
// The callback must honor the provided context for cancellation to take effect.
func CallWithTimeout(ctx context.Context, timeout time.Duration, opName string, fn func(context.Context) error) error {
	callCtx, cancel := deriveCallContext(ctx, timeout)
	defer cancel()

	err := fn(callCtx)
	if err != nil {
		if timeout > 0 && errors.Is(callCtx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("%s timed out after %s: %w", opName, timeout, context.DeadlineExceeded)
		}
		return err
	}
	if err := callCtx.Err(); err != nil {
		if timeout > 0 && errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%s timed out after %s: %w", opName, timeout, context.DeadlineExceeded)
		}
		return err
	}
	return nil
}

// RecvWithTimeout executes fn with a derived context that carries the timeout.
// The callback must honor the provided context for cancellation to take effect.
func RecvWithTimeout[T any](
	ctx context.Context,
	timeout time.Duration,
	opName string,
	fn func(context.Context) (T, error),
) (T, error) {
	var zero T
	callCtx, cancel := deriveCallContext(ctx, timeout)
	defer cancel()

	val, err := fn(callCtx)
	if err != nil {
		if timeout > 0 && errors.Is(callCtx.Err(), context.DeadlineExceeded) {
			return zero, fmt.Errorf("%s timed out after %s: %w", opName, timeout, context.DeadlineExceeded)
		}
		return zero, err
	}
	if err := callCtx.Err(); err != nil {
		if timeout > 0 && errors.Is(err, context.DeadlineExceeded) {
			return zero, fmt.Errorf("%s timed out after %s: %w", opName, timeout, context.DeadlineExceeded)
		}
		return zero, err
	}
	return val, nil
}

func deriveCallContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	if timeout <= 0 {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, timeout)
}
