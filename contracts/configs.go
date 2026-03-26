/*
Copyright 2025 BubuStack.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package contracts

// ConfigMap keys for operator configuration.
// These constants define the keys used in the operator ConfigMap to configure
// controller behavior, resource limits, timeouts, and feature flags.
const (
	// Controller timing configuration
	// Note: Per-controller max-concurrent-reconciles are under storyrun.*, steprun.*, etc.
	KeyRequeueBaseDelay                  = "controller.requeue-base-delay"
	KeyRequeueMaxDelay                   = "controller.requeue-max-delay"
	KeyCleanupInterval                   = "controller.cleanup-interval"
	KeyReconcileTimeout                  = "controller.reconcile-timeout"
	KeyMaxStoryWithBlockSizeBytes        = "controller.max-story-with-block-size-bytes"
	KeyControllerMaxConcurrentReconciles = "controller.max-concurrent-reconciles"

	// Image configuration
	KeyDefaultEngramImage  = "images.default-engram"
	KeyDefaultImpulseImage = "images.default-impulse"
	KeyImagePullPolicy     = "images.pull-policy"

	// Resource limit configuration
	KeyDefaultCPURequest    = "resources.default.cpu-request"
	KeyDefaultCPULimit      = "resources.default.cpu-limit"
	KeyDefaultMemoryRequest = "resources.default.memory-request"
	KeyDefaultMemoryLimit   = "resources.default.memory-limit"
	KeyEngramCPURequest     = "resources.engram.cpu-request"
	KeyEngramCPULimit       = "resources.engram.cpu-limit"
	KeyEngramMemoryRequest  = "resources.engram.memory-request"
	KeyEngramMemoryLimit    = "resources.engram.memory-limit"

	// Retry and timeout configuration
	KeyMaxRetries             = "retry.max-retries"
	KeyDefaultStepTimeout     = "timeout.default-step"
	KeyApprovalDefaultTimeout = "timeout.approval-default"
	KeyExternalDataTimeout    = "timeout.external-data-default"
	KeyConditionalTimeout     = "timeout.conditional-default"

	// Loop configuration
	KeyMaxLoopIterations    = "loop.max-iterations"
	KeyDefaultLoopBatchSize = "loop.default-batch-size"
	KeyMaxLoopBatchSize     = "loop.max-batch-size"
	KeyMaxLoopConcurrency   = "loop.max-concurrency"
	KeyMaxConcurrencyLimit  = "loop.max-concurrency-limit"

	// Security configuration
	KeyRunAsNonRoot                 = "security.run-as-non-root"
	KeyReadOnlyRootFilesystem       = "security.read-only-root-filesystem"
	KeyAllowPrivilegeEscalation     = "security.allow-privilege-escalation"
	KeyDropCapabilities             = "security.drop-capabilities"
	KeyRunAsUser                    = "security.run-as-user"
	KeyAutomountServiceAccountToken = "security.automount-service-account-token"
	KeyServiceAccountName           = "security.service-account-name"

	// Transport configuration
	KeyTransportEnabled              = "transport.enabled"
	KeyTransportHeartbeatInterval    = "controller.transport.heartbeat-interval"
	KeyTransportHeartbeatTimeout     = "controller.transport.heartbeat-timeout"
	KeyTransportSecurityMode         = "transport.security-mode"
	KeyTransportGRPCEnableDownstream = "controller.transport.grpc.enable-downstream-targets"
	KeyTransportGRPCDefaultTLSSecret = "controller.transport.grpc.default-tls-secret"

	// Job configuration
	KeyJobBackoffLimit                 = "job.backoff-limit"
	KeyJobTTLSecondsAfterFinished      = "job.ttl-seconds-after-finished"
	KeyRealtimeTTLSecondsAfterFinished = "realtime.ttl-seconds-after-finished"
	KeyJobRestartPolicy                = "job.restart-policy"
	KeyStoryRunRetentionSeconds        = "storyrun.retention-seconds"

	// Templating configuration
	KeyTemplatingEvaluationTimeout = "templating.evaluation-timeout"
	KeyTemplatingMaxOutputBytes    = "templating.max-output-bytes"
	KeyTemplatingDeterministic     = "templating.deterministic"
	KeyTemplatingOffloadedPolicy   = "templating.offloaded-data-policy"
	KeyTemplatingMaterializeEngram = "templating.materialize-engram"

	// Reference configuration
	KeyReferencesCrossNamespacePolicy = "references.cross-namespace-policy"

	// Telemetry configuration
	KeyTelemetryEnabled = "telemetry.enabled"
	KeyTracePropagation = "telemetry.trace-propagation"

	// Debug configuration
	KeyVerboseLogging    = "debug.enable-verbose-logging"
	KeyStepOutputLogging = "debug.enable-step-output-logging"
	KeyEnableMetrics     = "debug.enable-metrics"

	// Engram defaults
	KeyEngramDefaultInlineSize = "engram.default-max-inline-size"

	// StoryRun controller configuration
	KeyStoryRunMaxConcurrentReconciles     = "storyrun.max-concurrent-reconciles"
	KeyStoryRunRateLimiterBaseDelay        = "storyrun.rate-limiter.base-delay"
	KeyStoryRunRateLimiterMaxDelay         = "storyrun.rate-limiter.max-delay"
	KeyStoryRunMaxInlineInputsSize         = "storyrun.max-inline-inputs-size"
	KeyStoryRunBindingMaxMutations         = "storyrun.binding.max-mutations-per-reconcile"
	KeyStoryRunBindingThrottleRequeueDelay = "storyrun.binding.throttle-requeue-delay"
	KeyStoryRunGlobalConcurrency           = "storyrun.global-concurrency"
	KeyStoryRunQueuePrefix                 = "storyrun.queue."
	KeyStoryRunQueueConcurrencySuffix      = "concurrency"
	KeyStoryRunQueueDefaultPrioritySuffix  = "default-priority"
	KeyStoryRunQueuePriorityAgingSuffix    = "priority-aging-seconds"

	// StepRun controller configuration
	KeyStepRunMaxConcurrentReconciles = "steprun.max-concurrent-reconciles"
	KeyStepRunRateLimiterBaseDelay    = "steprun.rate-limiter.base-delay"
	KeyStepRunRateLimiterMaxDelay     = "steprun.rate-limiter.max-delay"

	// Story controller configuration
	KeyStoryMaxConcurrentReconciles     = "story.max-concurrent-reconciles"
	KeyStoryRateLimiterBaseDelay        = "story.rate-limiter.base-delay"
	KeyStoryRateLimiterMaxDelay         = "story.rate-limiter.max-delay"
	KeyStoryBindingMaxMutations         = "story.binding.max-mutations-per-reconcile"
	KeyStoryBindingThrottleRequeueDelay = "story.binding.throttle-requeue-delay"

	// Engram controller configuration
	KeyEngramMaxConcurrentReconciles           = "engram.max-concurrent-reconciles"
	KeyEngramRateLimiterBaseDelay              = "engram.rate-limiter.base-delay"
	KeyEngramRateLimiterMaxDelay               = "engram.rate-limiter.max-delay"
	KeyEngramDefaultGRPCPort                   = "engram.default-grpc-port"
	KeyEngramDefaultHeartbeatIntervalSeconds   = "engram.default-grpc-heartbeat-interval-seconds"
	KeyEngramDefaultStorageTimeoutSeconds      = "engram.default-storage-timeout-seconds"
	KeyEngramDefaultGracefulShutdownSeconds    = "engram.default-graceful-shutdown-timeout-seconds"
	KeyEngramDefaultTerminationGraceSeconds    = "engram.default-termination-grace-period-seconds"
	KeyEngramDefaultMaxRecvMsgBytes            = "engram.default-max-recv-msg-bytes"
	KeyEngramDefaultMaxSendMsgBytes            = "engram.default-max-send-msg-bytes"
	KeyEngramDefaultDialTimeoutSeconds         = "engram.default-dial-timeout-seconds"
	KeyEngramDefaultChannelBufferSize          = "engram.default-channel-buffer-size"
	KeyEngramDefaultReconnectMaxRetries        = "engram.default-reconnect-max-retries"
	KeyEngramDefaultReconnectBaseBackoffMillis = "engram.default-reconnect-base-backoff-millis"
	KeyEngramDefaultReconnectMaxBackoffSeconds = "engram.default-reconnect-max-backoff-seconds"
	KeyEngramDefaultHangTimeoutSeconds         = "engram.default-hang-timeout-seconds"
	KeyEngramDefaultMessageTimeoutSeconds      = "engram.default-message-timeout-seconds"

	// Impulse controller configuration
	KeyImpulseMaxConcurrentReconciles = "impulse.max-concurrent-reconciles"
	KeyImpulseRateLimiterBaseDelay    = "impulse.rate-limiter.base-delay"
	KeyImpulseRateLimiterMaxDelay     = "impulse.rate-limiter.max-delay"

	// Template controller configuration
	KeyTemplateMaxConcurrentReconciles = "template.max-concurrent-reconciles"
	KeyTemplateRateLimiterBaseDelay    = "template.rate-limiter.base-delay"
	KeyTemplateRateLimiterMaxDelay     = "template.rate-limiter.max-delay"

	// Storage defaults
	KeyStorageProvider        = "controller.storage.provider"
	KeyStorageS3Bucket        = "controller.storage.s3.bucket"
	KeyStorageS3Region        = "controller.storage.s3.region"
	KeyStorageS3Endpoint      = "controller.storage.s3.endpoint"
	KeyStorageS3PathStyle     = "controller.storage.s3.use-path-style"
	KeyStorageS3AuthSecret    = "controller.storage.s3.auth-secret-name"
	KeyStorageFilePath        = "controller.storage.file.path"
	KeyStorageFileVolumeClaim = "controller.storage.file.volume-claim-name"
)

// Image pull policy values
const (
	PullPolicyAlways       = "Always"
	PullPolicyNever        = "Never"
	PullPolicyIfNotPresent = "IfNotPresent"
)

// Restart policy values
const (
	RestartPolicyAlways    = "Always"
	RestartPolicyNever     = "Never"
	RestartPolicyOnFailure = "OnFailure"
)

// Storage provider values
const (
	StorageProviderS3   = "s3"
	StorageProviderFile = "file"
)

// Default configuration values
// These constants document the default values used by DefaultControllerConfig()
// and allow tests/docs to reference them without magic numbers.
const (
	// Controller concurrency defaults
	DefaultControllerMaxConcurrentReconciles = 10
	DefaultStoryRunMaxConcurrentReconciles   = 8
	DefaultStepRunMaxConcurrentReconciles    = 15
	DefaultStoryMaxConcurrentReconciles      = 5
	DefaultEngramMaxConcurrentReconciles     = 5
	DefaultImpulseMaxConcurrentReconciles    = 5
	DefaultTemplateMaxConcurrentReconciles   = 2
	DefaultTransportMaxConcurrentReconciles  = 2

	// Loop configuration defaults
	DefaultMaxLoopIterations   = 10000
	DefaultLoopBatchSize       = 100
	DefaultMaxLoopBatchSize    = 1000
	DefaultMaxLoopConcurrency  = 10
	DefaultMaxConcurrencyLimit = 50

	// Size limits
	DefaultMaxStoryWithBlockSizeBytes = 64 * 1024        // 64 KiB
	DefaultMaxInlineInputsSize        = 1 * 1024         // 1 KiB
	DefaultMaxInlineSize              = 1 * 1024         // 1 KiB
	DefaultMaxRecvMsgBytes            = 10 * 1024 * 1024 // 10 MiB
	DefaultMaxSendMsgBytes            = 10 * 1024 * 1024 // 10 MiB

	// gRPC defaults
	DefaultGRPCPort                       = 50051
	DefaultGRPCHeartbeatIntervalSeconds   = 10
	DefaultStorageTimeoutSeconds          = 300 // 5 minutes
	DefaultGracefulShutdownTimeoutSeconds = 20
	DefaultTerminationGracePeriodSeconds  = 30
	DefaultDialTimeoutSeconds             = 10
	DefaultChannelBufferSize              = 16
	DefaultReconnectMaxRetries            = 10
	DefaultReconnectBaseBackoffMillis     = 500
	DefaultReconnectMaxBackoffSeconds     = 30
	DefaultHangTimeoutSeconds             = 0
	DefaultMessageTimeoutSeconds          = 30

	// Binding controller defaults
	DefaultBindingMaxMutationsPerReconcile = 8
)
