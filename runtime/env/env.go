package env

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/bubustack/core/contracts"
)

// Metadata captures the contextual identifiers exported to workloads via env vars.
type Metadata struct {
	StoryName        string
	StoryRunName     string
	StepName         string
	StepRunName      string
	StepRunNamespace string
	WorkloadMode     string
}

// Config represents the gRPC/runtime defaults injected into realtime workloads.
type Config struct {
	DefaultGRPCPort                       int
	DefaultGRPCHeartbeatIntervalSeconds   int
	DefaultMaxInlineSize                  int
	DefaultStorageTimeoutSeconds          int
	DefaultGracefulShutdownTimeoutSeconds int
	DefaultMaxRecvMsgBytes                int
	DefaultMaxSendMsgBytes                int
	DefaultDialTimeoutSeconds             int
	DefaultChannelBufferSize              int
	DefaultReconnectMaxRetries            int
	DefaultReconnectBaseBackoffMillis     int
	DefaultReconnectMaxBackoffSeconds     int
	DefaultHangTimeoutSeconds             int
	DefaultMessageTimeoutSeconds          int
}

// BuildBaseEnv returns the common env vars required by realtime workloads.
func BuildBaseEnv(meta Metadata, cfg Config, tracePropagationEnabled bool) []corev1.EnvVar {
	cfg = normalizeConfig(cfg)
	mode := meta.WorkloadMode
	if mode == "" {
		mode = "deployment"
	}
	tracePropagation := "false"
	if tracePropagationEnabled {
		tracePropagation = "true"
	}
	envVars := []corev1.EnvVar{
		{Name: contracts.WorkloadModeEnv, Value: mode},
		{Name: contracts.ExecutionModeEnv, Value: "realtime"},
		{Name: contracts.GRPCPortEnv, Value: fmt.Sprintf("%d", cfg.DefaultGRPCPort)},
		{Name: contracts.MaxRecursionDepthEnv, Value: "64"},
		downwardEnvVar(contracts.PodNameEnv, "metadata.name"),
		downwardEnvVar(contracts.PodNamespaceEnv, "metadata.namespace"),
		downwardEnvVar(contracts.ServiceAccountNameEnv, "spec.serviceAccountName"),
		{Name: contracts.MaxInlineSizeEnv, Value: fmt.Sprintf("%d", cfg.DefaultMaxInlineSize)},
		{Name: contracts.StorageTimeoutEnv, Value: fmt.Sprintf("%ds", cfg.DefaultStorageTimeoutSeconds)},
		{Name: contracts.TracePropagationEnv, Value: tracePropagation},
		{Name: contracts.GRPCGracefulShutdownTimeout, Value: fmt.Sprintf("%ds", cfg.DefaultGracefulShutdownTimeoutSeconds)},
	}
	envVars = append(envVars, BuildGRPCTuningEnv(cfg)...)
	envVars = append(envVars, corev1.EnvVar{
		Name:  contracts.GRPCHeartbeatIntervalEnv,
		Value: fmt.Sprintf("%ds", cfg.DefaultGRPCHeartbeatIntervalSeconds),
	})

	if meta.StoryName != "" {
		envVars = append(envVars, corev1.EnvVar{Name: contracts.StoryNameEnv, Value: meta.StoryName})
	}
	if meta.StoryRunName != "" {
		envVars = append(envVars, corev1.EnvVar{Name: contracts.StoryRunIDEnv, Value: meta.StoryRunName})
	}
	if meta.StepName != "" {
		envVars = append(envVars, corev1.EnvVar{Name: contracts.StepNameEnv, Value: meta.StepName})
	}
	if meta.StepRunName != "" {
		envVars = append(envVars, corev1.EnvVar{Name: contracts.StepRunNameEnv, Value: meta.StepRunName})
	}
	if meta.StepRunNamespace != "" {
		envVars = append(envVars, corev1.EnvVar{Name: contracts.StepRunNamespaceEnv, Value: meta.StepRunNamespace})
	}

	return envVars
}

// BuildGRPCTuningEnv exports the gRPC tuning knobs shared across runtimes.
func BuildGRPCTuningEnv(cfg Config) []corev1.EnvVar {
	cfg = normalizeConfig(cfg)
	return []corev1.EnvVar{
		{Name: contracts.GRPCMaxRecvBytesEnv, Value: fmt.Sprintf("%d", cfg.DefaultMaxRecvMsgBytes)},
		{Name: contracts.GRPCMaxSendBytesEnv, Value: fmt.Sprintf("%d", cfg.DefaultMaxSendMsgBytes)},
		{Name: contracts.GRPCClientMaxRecvBytesEnv, Value: fmt.Sprintf("%d", cfg.DefaultMaxRecvMsgBytes)},
		{Name: contracts.GRPCClientMaxSendBytesEnv, Value: fmt.Sprintf("%d", cfg.DefaultMaxSendMsgBytes)},
		{Name: contracts.GRPCDialTimeoutEnv, Value: fmt.Sprintf("%ds", cfg.DefaultDialTimeoutSeconds)},
		{Name: contracts.GRPCChannelBufferSizeEnv, Value: fmt.Sprintf("%d", cfg.DefaultChannelBufferSize)},
		{Name: contracts.GRPCReconnectMaxRetriesEnv, Value: fmt.Sprintf("%d", cfg.DefaultReconnectMaxRetries)},
		{Name: contracts.GRPCReconnectBaseBackoffEnv, Value: fmt.Sprintf("%dms", cfg.DefaultReconnectBaseBackoffMillis)},
		{Name: contracts.GRPCReconnectMaxBackoffEnv, Value: fmt.Sprintf("%ds", cfg.DefaultReconnectMaxBackoffSeconds)},
		{Name: contracts.GRPCHangTimeoutEnv, Value: fmt.Sprintf("%ds", cfg.DefaultHangTimeoutSeconds)},
		{Name: contracts.GRPCMessageTimeoutEnv, Value: fmt.Sprintf("%ds", cfg.DefaultMessageTimeoutSeconds)},
		{Name: contracts.GRPCChannelSendTimeoutEnv, Value: fmt.Sprintf("%ds", cfg.DefaultMessageTimeoutSeconds)},
	}
}

// AppendStartedAtEnv appends BUBU_STARTED_AT when the timestamp is set.
func AppendStartedAtEnv(envVars *[]corev1.EnvVar, startedAt metav1.Time) {
	if envVars == nil || startedAt.IsZero() {
		return
	}
	*envVars = append(*envVars, corev1.EnvVar{
		Name:  contracts.StartedAtEnv,
		Value: startedAt.Format(time.RFC3339Nano),
	})
}

func downwardEnvVar(name, fieldPath string) corev1.EnvVar {
	return corev1.EnvVar{
		Name: name,
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				// APIVersion must be set explicitly to match what the API server
				// returns. Without it, the field defaults to "" in the desired
				// spec while the live object has "v1", causing reflect.DeepEqual
				// to report a diff on every reconcile → Deployment rollout storm.
				APIVersion: "v1",
				FieldPath:  fieldPath,
			},
		},
	}
}

func normalizeConfig(cfg Config) Config {
	cfg.DefaultGRPCPort = normalizePort(cfg.DefaultGRPCPort, contracts.DefaultGRPCPort)
	cfg.DefaultGRPCHeartbeatIntervalSeconds = normalizePositive(
		cfg.DefaultGRPCHeartbeatIntervalSeconds,
		contracts.DefaultGRPCHeartbeatIntervalSeconds,
	)
	cfg.DefaultMaxInlineSize = normalizePositive(
		cfg.DefaultMaxInlineSize,
		contracts.DefaultMaxInlineSize,
	)
	cfg.DefaultStorageTimeoutSeconds = normalizePositive(
		cfg.DefaultStorageTimeoutSeconds,
		contracts.DefaultStorageTimeoutSeconds,
	)
	cfg.DefaultGracefulShutdownTimeoutSeconds = normalizePositive(
		cfg.DefaultGracefulShutdownTimeoutSeconds,
		contracts.DefaultGracefulShutdownTimeoutSeconds,
	)
	cfg.DefaultMaxRecvMsgBytes = normalizePositive(
		cfg.DefaultMaxRecvMsgBytes,
		contracts.DefaultMaxRecvMsgBytes,
	)
	cfg.DefaultMaxSendMsgBytes = normalizePositive(
		cfg.DefaultMaxSendMsgBytes,
		contracts.DefaultMaxSendMsgBytes,
	)
	cfg.DefaultDialTimeoutSeconds = normalizePositive(
		cfg.DefaultDialTimeoutSeconds,
		contracts.DefaultDialTimeoutSeconds,
	)
	cfg.DefaultChannelBufferSize = normalizePositive(
		cfg.DefaultChannelBufferSize,
		contracts.DefaultChannelBufferSize,
	)
	cfg.DefaultReconnectMaxRetries = normalizePositive(
		cfg.DefaultReconnectMaxRetries,
		contracts.DefaultReconnectMaxRetries,
	)
	cfg.DefaultReconnectBaseBackoffMillis = normalizePositive(
		cfg.DefaultReconnectBaseBackoffMillis,
		contracts.DefaultReconnectBaseBackoffMillis,
	)
	cfg.DefaultReconnectMaxBackoffSeconds = normalizePositive(
		cfg.DefaultReconnectMaxBackoffSeconds,
		contracts.DefaultReconnectMaxBackoffSeconds,
	)
	cfg.DefaultHangTimeoutSeconds = normalizeNonNegative(
		cfg.DefaultHangTimeoutSeconds,
		contracts.DefaultHangTimeoutSeconds,
	)
	cfg.DefaultMessageTimeoutSeconds = normalizePositive(
		cfg.DefaultMessageTimeoutSeconds,
		contracts.DefaultMessageTimeoutSeconds,
	)
	return cfg
}

func normalizePort(value, fallback int) int {
	if value <= 0 || value > 65535 {
		return fallback
	}
	return value
}

func normalizePositive(value, fallback int) int {
	if value <= 0 {
		return fallback
	}
	return value
}

func normalizeNonNegative(value, fallback int) int {
	if value < 0 {
		return fallback
	}
	return value
}
