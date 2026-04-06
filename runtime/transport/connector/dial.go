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
// The callback should honor the provided context so any background work can unwind
// promptly after timeout or cancellation.
func CallWithTimeout(ctx context.Context, timeout time.Duration, opName string, fn func(context.Context) error) error {
	callCtx, cancel := deriveCallContext(ctx, timeout)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- fn(callCtx)
	}()

	select {
	case err := <-errCh:
		return classifyCallError(callCtx, timeout, opName, err)
	case <-callCtx.Done():
		return classifyCallError(callCtx, timeout, opName, callCtx.Err())
	}
}

// RecvWithTimeout executes fn with a derived context that carries the timeout.
// The callback should honor the provided context so any background work can unwind
// promptly after timeout or cancellation.
func RecvWithTimeout[T any](
	ctx context.Context,
	timeout time.Duration,
	opName string,
	fn func(context.Context) (T, error),
) (T, error) {
	var zero T
	callCtx, cancel := deriveCallContext(ctx, timeout)
	defer cancel()

	type result struct {
		val T
		err error
	}
	resCh := make(chan result, 1)
	go func() {
		val, err := fn(callCtx)
		resCh <- result{val: val, err: err}
	}()

	select {
	case res := <-resCh:
		if err := classifyCallError(callCtx, timeout, opName, res.err); err != nil {
			return zero, err
		}
		return res.val, nil
	case <-callCtx.Done():
		return zero, classifyCallError(callCtx, timeout, opName, callCtx.Err())
	}
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

func classifyCallError(ctx context.Context, timeout time.Duration, opName string, err error) error {
	if timeout > 0 && errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("%s timed out after %s: %w", opName, timeout, context.DeadlineExceeded)
	}
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}
