package env

import (
	"strconv"
	"testing"

	corev1 "k8s.io/api/core/v1"

	"github.com/bubustack/core/contracts"
)

func TestBuildBaseEnvNormalizesInvalidConfigValues(t *testing.T) {
	envs := BuildBaseEnv(Metadata{}, Config{
		DefaultGRPCPort:                       -1,
		DefaultGRPCHeartbeatIntervalSeconds:   0,
		DefaultMaxInlineSize:                  0,
		DefaultStorageTimeoutSeconds:          0,
		DefaultGracefulShutdownTimeoutSeconds: 0,
		DefaultMaxRecvMsgBytes:                0,
		DefaultMaxSendMsgBytes:                0,
		DefaultDialTimeoutSeconds:             0,
		DefaultChannelBufferSize:              0,
		DefaultReconnectMaxRetries:            0,
		DefaultReconnectBaseBackoffMillis:     0,
		DefaultReconnectMaxBackoffSeconds:     0,
		DefaultMessageTimeoutSeconds:          0,
	}, false)

	assertEnvValue(t, envs, contracts.GRPCPortEnv, strconv.Itoa(contracts.DefaultGRPCPort))
	assertEnvValue(t, envs, contracts.MaxInlineSizeEnv, strconv.Itoa(contracts.DefaultMaxInlineSize))
	assertEnvValue(t, envs, contracts.StorageTimeoutEnv, "300s")
	assertEnvValue(t, envs, contracts.GRPCHeartbeatIntervalEnv, "10s")
	assertEnvValue(t, envs, contracts.GRPCDialTimeoutEnv, "10s")
}

func TestBuildBaseEnvUsesExplicitDownwardAPIversion(t *testing.T) {
	envs := BuildBaseEnv(Metadata{}, Config{}, false)
	for _, env := range envs {
		if env.ValueFrom == nil || env.ValueFrom.FieldRef == nil {
			continue
		}
		if env.ValueFrom.FieldRef.APIVersion != "v1" {
			t.Fatalf("expected APIVersion v1 for %s, got %q", env.Name, env.ValueFrom.FieldRef.APIVersion)
		}
	}
}

func assertEnvValue(t *testing.T, envs []corev1.EnvVar, name, expected string) {
	t.Helper()
	for _, env := range envs {
		if env.Name == name {
			if env.Value != expected {
				t.Fatalf("expected %s=%s, got %s", name, expected, env.Value)
			}
			return
		}
	}
	t.Fatalf("expected env %s to be present", name)
}
