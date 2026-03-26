package transport

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/bubustack/core/contracts"
	transportpb "github.com/bubustack/tractatus/gen/go/proto/transport/v1"
)

var envVarNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// AppendTransportMetadataEnv adds convenience env vars derived from BindingInfo so
// runtimes can avoid re-parsing the transport envelope.
func AppendTransportMetadataEnv(envVars []corev1.EnvVar, info *transportpb.BindingInfo) []corev1.EnvVar {
	if info == nil {
		return envVars
	}
	if driver := strings.TrimSpace(info.GetDriver()); driver != "" {
		envVars = SetOrReplaceEnvVar(envVars, corev1.EnvVar{Name: contracts.TransportDriverEnv, Value: driver})
	}
	if endpoint := strings.TrimSpace(info.GetEndpoint()); endpoint != "" {
		envVars = SetOrReplaceEnvVar(envVars, corev1.EnvVar{Name: contracts.TransportEndpointEnv, Value: endpoint})
	}
	if len(info.AudioCodecs) > 0 {
		envVars = SetOrReplaceEnvVar(envVars, corev1.EnvVar{
			Name:  contracts.TransportAudioCodecsEnv,
			Value: strings.Join(info.AudioCodecs, ","),
		})
	}
	if len(info.VideoCodecs) > 0 {
		envVars = SetOrReplaceEnvVar(envVars, corev1.EnvVar{
			Name:  contracts.TransportVideoCodecsEnv,
			Value: strings.Join(info.VideoCodecs, ","),
		})
	}
	if len(info.BinaryTypes) > 0 {
		envVars = SetOrReplaceEnvVar(envVars, corev1.EnvVar{
			Name:  contracts.TransportBinaryTypesEnv,
			Value: strings.Join(info.BinaryTypes, ","),
		})
	}
	return envVars
}

// AppendTransportHeartbeatEnv surfaces the heartbeat interval so runtimes can tune
// their own watchdogs without parsing controller config.
func AppendTransportHeartbeatEnv(envVars []corev1.EnvVar, interval time.Duration) []corev1.EnvVar {
	if interval <= 0 {
		return envVars
	}
	value := interval.String()
	for i := range envVars {
		if envVars[i].Name == contracts.TransportHeartbeatIntervalEnv {
			envVars[i].Value = value
			return envVars
		}
	}
	return append(envVars, corev1.EnvVar{Name: contracts.TransportHeartbeatIntervalEnv, Value: value})
}

// AppendBindingEnvOverrides injects provider-specific env overrides encoded in the
// binding payload, validating env var names to prevent invalid workloads.
func AppendBindingEnvOverrides(envVars []corev1.EnvVar, info *transportpb.BindingInfo) []corev1.EnvVar {
	overrides := envOverrides(info)
	if len(overrides) == 0 {
		return envVars
	}
	keys := make([]string, 0, len(overrides))
	for key := range overrides {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := overrides[key]
		name := strings.TrimSpace(key)
		if name == "" || !envVarNamePattern.MatchString(name) {
			continue
		}
		if !strings.HasPrefix(name, contracts.PrefixEnv) && !hasEnvVar(envVars, name) {
			continue
		}
		envVars = SetOrReplaceEnvVar(envVars, corev1.EnvVar{Name: name, Value: value})
	}
	return envVars
}

// SetOrReplaceEnvVar replaces the env var with the same name if it exists or appends
// it otherwise, returning the updated slice.
func SetOrReplaceEnvVar(envVars []corev1.EnvVar, env corev1.EnvVar) []corev1.EnvVar {
	for i := range envVars {
		if envVars[i].Name == env.Name {
			envVars[i].Value = env.Value
			envVars[i].ValueFrom = env.ValueFrom
			return envVars
		}
	}
	return append(envVars, env)
}

func envOverrides(info *transportpb.BindingInfo) map[string]string {
	if info == nil || len(info.Payload) == 0 {
		return nil
	}
	var payload struct {
		Env map[string]any `json:"env"`
	}
	if err := json.Unmarshal(info.Payload, &payload); err != nil || len(payload.Env) == 0 {
		return nil
	}
	out := make(map[string]string, len(payload.Env))
	for key, raw := range payload.Env {
		name := strings.TrimSpace(key)
		if name == "" {
			continue
		}
		out[name] = fmt.Sprint(raw)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// BindingEnvOverrides returns the sanitized env override map encoded inside a
// BindingInfo payload so non-controller workflows (connectors, SDKs) can apply
// the shared parsing logic without duplicating it.
func BindingEnvOverrides(info *transportpb.BindingInfo) map[string]string {
	return envOverrides(info)
}

func hasEnvVar(envVars []corev1.EnvVar, name string) bool {
	for _, env := range envVars {
		if env.Name == name {
			return true
		}
	}
	return false
}
