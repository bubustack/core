package connector

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestListenLocalEndpoint_TCP(t *testing.T) {
	lis, err := ListenLocalEndpoint("127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen tcp: %v", err)
	}
	cleanupListener(t, lis)
	if lis.Addr().Network() != "tcp" {
		t.Fatalf("expected tcp listener, got %s", lis.Addr().Network())
	}
}

func TestListenLocalEndpoint_UNIX(t *testing.T) {
	dir, err := os.MkdirTemp("/tmp", "listener")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})
	socketPath := filepath.Join(dir, "connector.sock")
	lis, err := ListenLocalEndpoint("unix://" + socketPath)
	if err != nil {
		t.Fatalf("listen unix: %v", err)
	}
	cleanupListener(t, lis)
	if lis.Addr().Network() != "unix" {
		t.Fatalf("expected unix listener, got %s", lis.Addr().Network())
	}
	if _, err := net.Dial("unix", socketPath); err != nil {
		t.Fatalf("dial unix socket: %v", err)
	}
}

func TestListenLocalEndpointRejectsRegularFileAtUnixPath(t *testing.T) {
	dir := t.TempDir()
	socketPath := filepath.Join(dir, "connector.sock")
	if err := os.WriteFile(socketPath, []byte("not-a-socket"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	if _, err := ListenLocalEndpoint("unix://" + socketPath); err == nil {
		t.Fatalf("expected regular-file unix path to be rejected")
	}
}

func TestListenLocalEndpointBarePortUsesLoopback(t *testing.T) {
	lis, err := ListenLocalEndpoint("0")
	if err != nil {
		t.Fatalf("listen tcp: %v", err)
	}
	cleanupListener(t, lis)

	host, _, err := net.SplitHostPort(lis.Addr().String())
	if err != nil {
		t.Fatalf("split host port: %v", err)
	}
	if host != "127.0.0.1" {
		t.Fatalf("expected loopback host, got %q", host)
	}
}

func cleanupListener(t *testing.T, lis net.Listener) {
	t.Helper()
	t.Cleanup(func() {
		if err := lis.Close(); err != nil {
			t.Errorf("close listener: %v", err)
		}
	})
}
