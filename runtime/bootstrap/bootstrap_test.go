package bootstrap_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/go-logr/logr/testr"

	"github.com/bubustack/core/runtime/bootstrap"
)

func TestRunnerRegisterSuccess(t *testing.T) {
	logger := testr.New(t)
	var registered []string

	runner := bootstrap.Runner{Log: logger}
	err := runner.Register(
		bootstrap.Entry{
			Kind: "controller",
			Name: "Story",
			Register: func() error {
				registered = append(registered, "Story")
				return nil
			},
		},
		bootstrap.Entry{
			Kind: "webhook",
			Name: "Impulse",
			Register: func() error {
				registered = append(registered, "Impulse")
				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(registered) != 2 {
		t.Fatalf("expected 2 registrations, got %d", len(registered))
	}
	if registered[0] != "Story" || registered[1] != "Impulse" {
		t.Fatalf("registrations executed out of order: %v", registered)
	}
}

func TestRunnerRegisterStopsOnError(t *testing.T) {
	logger := testr.New(t)
	calls := 0
	expectedErr := errors.New("boom")

	runner := bootstrap.Runner{Log: logger}
	err := runner.Register(
		bootstrap.Entry{
			Kind: "controller",
			Name: "Story",
			Register: func() error {
				calls++
				return expectedErr
			},
		},
		bootstrap.Entry{
			Kind: "controller",
			Name: "StepRun",
			Register: func() error {
				calls++
				return nil
			},
		},
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) && !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Fatalf("error did not wrap original: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected runner to stop after first failure, call count=%d", calls)
	}
}
