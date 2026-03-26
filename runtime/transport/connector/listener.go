package connector

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// ListenLocalEndpoint normalizes tcp/unix endpoints and returns a listener.
func ListenLocalEndpoint(endpoint string) (net.Listener, error) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint must be provided")
	}
	if path, ok := strings.CutPrefix(endpoint, "unix://"); ok {
		if path == "" {
			return nil, fmt.Errorf("unix endpoint missing path")
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, fmt.Errorf("create unix socket dir: %w", err)
		}
		if info, err := os.Lstat(path); err == nil {
			if info.Mode()&os.ModeSocket == 0 {
				return nil, fmt.Errorf("unix endpoint path exists and is not a socket: %s", path)
			}
			if err := os.Remove(path); err != nil {
				return nil, fmt.Errorf("remove unix socket %s: %w", path, err)
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("stat unix socket %s: %w", path, err)
		}
		lis, err := net.Listen("unix", path)
		if err != nil {
			return nil, fmt.Errorf("listen on unix socket %s: %w", path, err)
		}
		return lis, nil
	}
	if !strings.Contains(endpoint, ":") {
		endpoint = fmt.Sprintf("127.0.0.1:%s", endpoint)
	}
	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		return nil, fmt.Errorf("listen on %s: %w", endpoint, err)
	}
	return lis, nil
}
