package transport

import (
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/bubustack/core/contracts"
	transportpb "github.com/bubustack/tractatus/gen/go/proto/transport/v1"
)

func TestAppendTransportMetadataEnvIncludesBinary(t *testing.T) {
	info := &transportpb.BindingInfo{
		Driver:      "livekit",
		Endpoint:    "unix:///tmp/connector.sock",
		AudioCodecs: []string{"pcm16"},
		VideoCodecs: []string{"h264"},
		BinaryTypes: []string{"application/json", "application/octet-stream"},
	}
	envs := AppendTransportMetadataEnv(nil, info)
	assertEnvValue(t, envs, contracts.TransportDriverEnv, "livekit")
	assertEnvValue(t, envs, contracts.TransportEndpointEnv, "unix:///tmp/connector.sock")
	assertEnvValue(t, envs, contracts.TransportAudioCodecsEnv, "pcm16")
	assertEnvValue(t, envs, contracts.TransportVideoCodecsEnv, "h264")
	assertEnvValue(t, envs, contracts.TransportBinaryTypesEnv, "application/json,application/octet-stream")
}

func TestAppendTransportHeartbeatEnv(t *testing.T) {
	envs := AppendTransportHeartbeatEnv(nil, 45*time.Second)
	assertEnvValue(t, envs, contracts.TransportHeartbeatIntervalEnv, "45s")

	envs = AppendTransportHeartbeatEnv([]corev1.EnvVar{
		{Name: contracts.TransportHeartbeatIntervalEnv, Value: "old"},
	}, 15*time.Second)
	assertEnvValue(t, envs, contracts.TransportHeartbeatIntervalEnv, "15s")

	envs = AppendTransportHeartbeatEnv(envs, 0)
	assertEnvValue(t, envs, contracts.TransportHeartbeatIntervalEnv, "15s")
}

func TestAppendBindingEnvOverrides(t *testing.T) {
	info := &transportpb.BindingInfo{
		Payload: []byte(`{"env":{"FOO":"bar","_VALID":"1","1INVALID":"noop","BUBU_DEBUG":"true"}}`),
	}
	envs := []corev1.EnvVar{{Name: "FOO", Value: "original"}}
	envs = AppendBindingEnvOverrides(envs, info)

	assertEnvValue(t, envs, "FOO", "bar")
	assertEnvMissing(t, envs, "_VALID")
	assertEnvValue(t, envs, "BUBU_DEBUG", "true")
	assertEnvMissing(t, envs, "1INVALID")
}

func TestAppendBindingEnvOverridesNoInfo(t *testing.T) {
	envs := []corev1.EnvVar{{Name: "FOO", Value: "bar"}}
	result := AppendBindingEnvOverrides(envs, nil)
	if len(result) != len(envs) {
		t.Fatalf("expected env slice unchanged when info is nil")
	}
}

func TestBindingEnvOverrides(t *testing.T) {
	info := &transportpb.BindingInfo{
		Payload: []byte(`{"env":{"FOO":"bar","BAR":"baz"}}`),
	}
	overrides := BindingEnvOverrides(info)
	if len(overrides) != 2 || overrides["FOO"] != "bar" || overrides["BAR"] != "baz" {
		t.Fatalf("unexpected overrides map: %#v", overrides)
	}
}

func TestBindingEnvOverridesNil(t *testing.T) {
	if overrides := BindingEnvOverrides(nil); overrides != nil {
		t.Fatalf("expected nil overrides for nil info, got %#v", overrides)
	}
}

func TestAppendTransportMetadataEnvReplacesExistingValues(t *testing.T) {
	info := &transportpb.BindingInfo{
		Driver:   "updated",
		Endpoint: "127.0.0.1:50051",
	}
	envs := []corev1.EnvVar{
		{Name: contracts.TransportDriverEnv, Value: "old"},
		{Name: contracts.TransportEndpointEnv, Value: "old"},
	}

	envs = AppendTransportMetadataEnv(envs, info)
	assertEnvValue(t, envs, contracts.TransportDriverEnv, "updated")
	assertEnvValue(t, envs, contracts.TransportEndpointEnv, "127.0.0.1:50051")
}

func TestAppendBindingEnvOverridesSortsInjectedBubuKeys(t *testing.T) {
	info := &transportpb.BindingInfo{
		Payload: []byte(`{"env":{"BUBU_Z":"z","BUBU_A":"a"}}`),
	}

	envs := AppendBindingEnvOverrides(nil, info)
	if len(envs) != 2 {
		t.Fatalf("expected two env vars, got %#v", envs)
	}
	if envs[0].Name != "BUBU_A" || envs[1].Name != "BUBU_Z" {
		t.Fatalf("expected deterministic ordering, got %#v", envs)
	}
}

func assertEnvValue(t *testing.T, envs []corev1.EnvVar, name, expected string) {
	t.Helper()
	for _, env := range envs {
		if env.Name == name {
			if env.Value != expected {
				t.Fatalf("expected env %s=%s, got %s", name, expected, env.Value)
			}
			return
		}
	}
	t.Fatalf("expected env %s to be present", name)
}

func assertEnvMissing(t *testing.T, envs []corev1.EnvVar, name string) {
	t.Helper()
	for _, env := range envs {
		if env.Name == name {
			t.Fatalf("expected env %s to be absent", name)
			return
		}
	}
}
