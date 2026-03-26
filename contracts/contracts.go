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

// PrefixEnv prefixes all BubuStack-owned environment variables.
const PrefixEnv = "BUBU_"

const (
	// StoryNameEnv stores the Story name for the active workload.
	StoryNameEnv = PrefixEnv + "STORY_NAME"
	// StoryVersionEnv stores the Story version for the active workload.
	StoryVersionEnv = PrefixEnv + "STORY_VERSION"
	// StoryRunIDEnv stores the StoryRun identifier for the active workload.
	StoryRunIDEnv = PrefixEnv + "STORYRUN_ID"
	// StepNameEnv stores the step name for the active workload.
	StepNameEnv = PrefixEnv + "STEP_NAME"
	// StepRunNameEnv stores the StepRun resource name for the active workload.
	StepRunNameEnv = PrefixEnv + "STEPRUN_NAME"
	// StepRunNamespaceEnv stores the StepRun namespace for the active workload.
	StepRunNamespaceEnv = PrefixEnv + "STEPRUN_NAMESPACE"
	// EngramNameEnv stores the Engram name bound to the active workload.
	EngramNameEnv = PrefixEnv + "ENGRAM_NAME"
	// EngramVersionEnv stores the Engram version bound to the active workload.
	EngramVersionEnv = PrefixEnv + "ENGRAM_VERSION"
	// ExecutionModeEnv stores the top-level execution mode selected for the workload.
	ExecutionModeEnv = PrefixEnv + "EXECUTION_MODE"
	// WorkloadModeEnv stores the specific runtime mode used by the current pod or process.
	WorkloadModeEnv = PrefixEnv + "WORKLOAD_MODE"
)

const (
	// StepConfigEnv stores static step configuration resolved from `step.config` or `step.with`.
	StepConfigEnv = PrefixEnv + "STEP_CONFIG"
	// TriggerDataEnv stores the resolved trigger payload for the active workload.
	TriggerDataEnv = PrefixEnv + "TRIGGER_DATA"
	// TemplateContextEnv stores the template evaluation context passed to the SDK.
	TemplateContextEnv = PrefixEnv + "TEMPLATE_CONTEXT"
	// SkipInputTemplatingEnv disables SDK-side template resolution when set.
	SkipInputTemplatingEnv = PrefixEnv + "SKIP_INPUT_TEMPLATING"
	// TransportsEnv stores serialized transport configuration for the active workload.
	TransportsEnv = PrefixEnv + "TRANSPORTS"
	// ConfigPrefixEnv prefixes per-key configuration environment variables.
	ConfigPrefixEnv = PrefixEnv + "CONFIG_"
	// StartedAtEnv stores the workload start timestamp.
	StartedAtEnv = PrefixEnv + "STARTED_AT"
	// ImpulseNameEnv stores the triggering Impulse name when execution originates from an Impulse.
	ImpulseNameEnv = PrefixEnv + "IMPULSE_NAME"
	// ImpulseNamespaceEnv stores the triggering Impulse namespace when execution originates from an Impulse.
	ImpulseNamespaceEnv = PrefixEnv + "IMPULSE_NAMESPACE"
	// TargetStoryNameEnv stores the Story name targeted by a trigger or bridge action.
	TargetStoryNameEnv = PrefixEnv + "TARGET_STORY_NAME"
	// TargetStoryNamespaceEnv stores the namespace of the targeted Story.
	TargetStoryNamespaceEnv = PrefixEnv + "TARGET_STORY_NAMESPACE"
	// TriggerTokenEnv stores the normalized trigger token for the current execution.
	TriggerTokenEnv = PrefixEnv + "TRIGGER_TOKEN"
	// PodNameEnv stores the pod name for downward-API consumers.
	PodNameEnv = PrefixEnv + "POD_NAME"
	// PodNamespaceEnv stores the pod namespace for downward-API consumers.
	PodNamespaceEnv = PrefixEnv + "POD_NAMESPACE"
	// ServiceAccountNameEnv stores the Kubernetes service account name used by the workload.
	ServiceAccountNameEnv = PrefixEnv + "SERVICE_ACCOUNT_NAME"
	// SecretPrefixEnv prefixes environment variables sourced from secrets.
	SecretPrefixEnv = PrefixEnv + "SECRET_"
	// TransportBindingEnv stores the serialized transport binding envelope for the workload.
	TransportBindingEnv = PrefixEnv + "TRANSPORT_BINDING"
	// TransportDriverEnv stores the selected transport driver identifier.
	TransportDriverEnv = PrefixEnv + "TRANSPORT_DRIVER"
	// TransportEndpointEnv stores the resolved transport endpoint for the workload.
	TransportEndpointEnv = PrefixEnv + "TRANSPORT_ENDPOINT"
	// TransportAudioCodecsEnv stores the allowed transport audio codec list.
	TransportAudioCodecsEnv = PrefixEnv + "TRANSPORT_AUDIO_CODECS"
	// TransportVideoCodecsEnv stores the allowed transport video codec list.
	TransportVideoCodecsEnv = PrefixEnv + "TRANSPORT_VIDEO_CODECS"
	// TransportBinaryTypesEnv stores the allowed binary payload type list.
	TransportBinaryTypesEnv = PrefixEnv + "TRANSPORT_BINARY_TYPES"
	// TransportHeartbeatIntervalEnv stores the transport heartbeat interval.
	TransportHeartbeatIntervalEnv = PrefixEnv + "TRANSPORT_HEARTBEAT_INTERVAL"
	// TransportSecurityModeEnv stores the transport security mode. The supported value is "tls".
	TransportSecurityModeEnv = PrefixEnv + "TRANSPORT_SECURITY_MODE"
	// DebugEnv enables shared debug behavior when set.
	DebugEnv = PrefixEnv + "DEBUG"
)

const (
	// TransportReadyAnnotation marks a transport-backed workload as ready.
	TransportReadyAnnotation = "transport.bobravoz.bubustack.io/ready"
	// TransportReadyMessageAnnotation stores human-readable readiness details for a transport workload.
	TransportReadyMessageAnnotation = "transport.bobravoz.bubustack.io/message"
	// TransportBindingAnnotation stores the transport binding reference attached to a resource.
	TransportBindingAnnotation = "transport.bubustack.io/binding"
	// StoryAnnotation stores the owning Story name on an annotated resource.
	StoryAnnotation = "bubustack.io/story"
	// StoryRunAnnotation stores the owning StoryRun name on an annotated resource.
	StoryRunAnnotation = "bubustack.io/storyrun"
	// StoryRunTriggerAnnotation marks that a StoryRun trigger has been recorded.
	StoryRunTriggerAnnotation = "storyrun.bubustack.io/trigger-recorded"
	// StoryRunTriggerTokenAnnotation stores the trigger token recorded on a StoryRun.
	StoryRunTriggerTokenAnnotation = "storyrun.bubustack.io/trigger-token"
	// StepAnnotation stores the owning step name on an annotated resource.
	StepAnnotation = "bubustack.io/step"
	// ResolvedInputsAnnotation stores resolved step inputs materialized by the controller.
	ResolvedInputsAnnotation = "bubustack.io/resolved-inputs"
	// EngramTLSSecretAnnotation stores the TLS secret associated with an Engram.
	EngramTLSSecretAnnotation = "engram.bubustack.io/tls-secret"
	// TLSSecretAnnotation stores a generic TLS secret reference on a workload resource.
	TLSSecretAnnotation = "bubustack.io/tls-secret"
	// StepRunTriggerAnnotation marks that a StepRun trigger has been recorded.
	StepRunTriggerAnnotation = "steprun.bubustack.io/trigger-recorded"
	// MaterializePurposeAnnotation stores why a materialization workload was created.
	MaterializePurposeAnnotation = "bubustack.io/materialize-purpose"
	// MaterializeTargetAnnotation stores the target step for a materialization workload.
	MaterializeTargetAnnotation = "bubustack.io/materialize-target-step"
	// MaterializeModeAnnotation stores the materialization mode requested by the controller.
	MaterializeModeAnnotation = "bubustack.io/materialize-mode"
	// MaterializeTransportAnnotation stores the transport selected for a materialization workload.
	MaterializeTransportAnnotation = "bubustack.io/materialize-transport"
	// CorrelationIDAnnotation stores the cross-component correlation identifier.
	CorrelationIDAnnotation = "bubustack.io/correlation-id"
)

// ControllerResolveAnnotation forces the controller to resolve step templates
// server-side instead of delegating evaluation to the SDK. Use the literal
// value `"true"` to enable controller-side resolution.
const ControllerResolveAnnotation = "bubustack.io/controller-resolve"

const (
	// StoryRunTriggerTokenStory identifies Story-triggered StoryRuns.
	StoryRunTriggerTokenStory = "story"
	// StoryRunTriggerTokenImpulse identifies Impulse-triggered StoryRuns.
	StoryRunTriggerTokenImpulse = "impulse"
	// StoryRunTriggerTokenImpulseSuccess identifies success callbacks from Impulse execution.
	StoryRunTriggerTokenImpulseSuccess = "impulse-success"
	// StoryRunTriggerTokenImpulseFailed identifies failure callbacks from Impulse execution.
	StoryRunTriggerTokenImpulseFailed = "impulse-failed"
	// StepRunTriggerTokenEngram identifies Engram-triggered StepRuns.
	StepRunTriggerTokenEngram = "engram"
)

const (
	// MaxInlineSizeEnv stores the maximum size for inline payload materialization.
	MaxInlineSizeEnv = PrefixEnv + "MAX_INLINE_SIZE"
	// MediaInlineSizeEnv stores the maximum size for inline media payloads.
	MediaInlineSizeEnv = PrefixEnv + "MEDIA_INLINE_SIZE"
	// MaxRecursionDepthEnv stores the maximum shared recursion depth for evaluation helpers.
	MaxRecursionDepthEnv = PrefixEnv + "MAX_RECURSION_DEPTH"
	// StepTimeoutEnv stores the step execution timeout.
	StepTimeoutEnv = PrefixEnv + "STEP_TIMEOUT"
	// StorageTimeoutEnv stores the storage operation timeout.
	StorageTimeoutEnv = PrefixEnv + "STORAGE_TIMEOUT"
)

const (
	// HybridBridgeEnv stores the hybrid bridge mode or endpoint selection.
	HybridBridgeEnv = PrefixEnv + "HYBRID_BRIDGE"
	// HybridBridgeTimeoutEnv stores the timeout applied to hybrid bridge operations.
	HybridBridgeTimeoutEnv = PrefixEnv + "HYBRID_BRIDGE_TIMEOUT"
	// ConnectorGenerationEnv stores the connector generation expected by the workload.
	ConnectorGenerationEnv = PrefixEnv + "CONNECTOR_GENERATION"
	// GRPCLocalEndpointEnv stores the local gRPC listen endpoint for an Engram or connector.
	GRPCLocalEndpointEnv = PrefixEnv + "GRPC_LOCAL_ENDPOINT"
	// GRPCPortEnv stores the gRPC listen port.
	GRPCPortEnv = PrefixEnv + "GRPC_PORT"
)

const (
	// HubEndpointEnv stores the explicitly configured hub endpoint.
	HubEndpointEnv = PrefixEnv + "HUB_ENDPOINT"
	// HubServiceNameEnv stores the Kubernetes Service name used for hub discovery.
	HubServiceNameEnv = PrefixEnv + "HUB_SERVICE_NAME"
	// HubServiceNamespaceEnv stores the Kubernetes Service namespace used for hub discovery.
	HubServiceNamespaceEnv = PrefixEnv + "HUB_SERVICE_NAMESPACE"
	// HubClusterDomainEnv stores the cluster DNS suffix used for hub discovery.
	HubClusterDomainEnv = PrefixEnv + "HUB_CLUSTER_DOMAIN"
	// HubPortEnv stores the discovered or configured hub port.
	HubPortEnv = PrefixEnv + "HUB_PORT"
)

const (
	// HubBufferMaxMessagesEnv stores the maximum buffered message count per hub buffer.
	HubBufferMaxMessagesEnv = PrefixEnv + "HUB_BUFFER_MAX_MESSAGES"
	// HubBufferMaxBytesEnv stores the maximum buffered byte size per hub buffer.
	HubBufferMaxBytesEnv = PrefixEnv + "HUB_BUFFER_MAX_BYTES"
	// HubBufferEvictionTTLEnv stores the TTL for buffered hub messages before eviction.
	HubBufferEvictionTTLEnv = PrefixEnv + "HUB_BUFFER_EVICTION_TTL"
	// HubBufferEvictionIntervalEnv stores how often hub buffer eviction runs.
	HubBufferEvictionIntervalEnv = PrefixEnv + "HUB_BUFFER_EVICTION_INTERVAL"
	// HubBufferFlushIntervalEnv stores how often buffered hub messages are flushed.
	HubBufferFlushIntervalEnv = PrefixEnv + "HUB_BUFFER_FLUSH_INTERVAL"
	// HubPerMessageTimeoutEnv stores the timeout applied to individual hub message handling.
	HubPerMessageTimeoutEnv = PrefixEnv + "HUB_PER_MESSAGE_TIMEOUT"
	// HubJoinCacheTTLEnv stores the TTL for cached hub join state.
	HubJoinCacheTTLEnv = PrefixEnv + "HUB_JOIN_CACHE_TTL"
	// HubJoinCacheMaxEntriesEnv stores the maximum join-cache entry count.
	HubJoinCacheMaxEntriesEnv = PrefixEnv + "HUB_JOIN_CACHE_MAX_ENTRIES"
	// HubMaxActiveStreamsEnv stores the maximum number of concurrent active streams.
	HubMaxActiveStreamsEnv = PrefixEnv + "HUB_MAX_ACTIVE_STREAMS"
	// HubMaxBuffersEnv stores the maximum number of in-memory hub buffers.
	HubMaxBuffersEnv = PrefixEnv + "HUB_MAX_BUFFERS"
	// HubMaxDownstreamsEnv stores the maximum number of downstream consumers attached to a buffer.
	HubMaxDownstreamsEnv = PrefixEnv + "HUB_MAX_DOWNSTREAMS"
)

const (
	// GRPCMaxRecvBytesEnv stores the server-side maximum received gRPC message size.
	GRPCMaxRecvBytesEnv = PrefixEnv + "GRPC_MAX_RECV_BYTES"
	// GRPCMaxSendBytesEnv stores the server-side maximum sent gRPC message size.
	GRPCMaxSendBytesEnv = PrefixEnv + "GRPC_MAX_SEND_BYTES"
	// GRPCClientMaxRecvBytesEnv stores the client-side maximum received gRPC message size.
	GRPCClientMaxRecvBytesEnv = PrefixEnv + "GRPC_CLIENT_MAX_RECV_BYTES"
	// GRPCClientMaxSendBytesEnv stores the client-side maximum sent gRPC message size.
	GRPCClientMaxSendBytesEnv = PrefixEnv + "GRPC_CLIENT_MAX_SEND_BYTES"
	// GRPCDialTimeoutEnv stores the default connector dial timeout.
	GRPCDialTimeoutEnv = PrefixEnv + "GRPC_DIAL_TIMEOUT"
	// GRPCHubDialTimeoutEnv stores the timeout for dialing the transport hub.
	GRPCHubDialTimeoutEnv = PrefixEnv + "GRPC_HUB_DIAL_TIMEOUT"
	// GRPCStreamTimeoutEnv stores the timeout applied to long-lived gRPC streams.
	GRPCStreamTimeoutEnv = PrefixEnv + "GRPC_STREAM_TIMEOUT"
	// GRPCChannelBufferSizeEnv stores the connector channel buffer size.
	GRPCChannelBufferSizeEnv = PrefixEnv + "GRPC_CHANNEL_BUFFER_SIZE"
	// GRPCChannelSendTimeoutEnv stores the timeout applied to connector channel sends.
	GRPCChannelSendTimeoutEnv = PrefixEnv + "GRPC_CHANNEL_SEND_TIMEOUT"
	// GRPCMessageTimeoutEnv stores the timeout applied to single message operations.
	GRPCMessageTimeoutEnv = PrefixEnv + "GRPC_MESSAGE_TIMEOUT"
	// GRPCHangTimeoutEnv stores the timeout used when detecting hung gRPC operations.
	GRPCHangTimeoutEnv = PrefixEnv + "GRPC_HANG_TIMEOUT"
	// GRPCGracefulShutdownTimeout stores the timeout for graceful gRPC shutdown.
	GRPCGracefulShutdownTimeout = PrefixEnv + "GRPC_GRACEFUL_SHUTDOWN_TIMEOUT"
	// GRPCHeartbeatIntervalEnv stores the heartbeat interval for gRPC keepalive logic.
	GRPCHeartbeatIntervalEnv = PrefixEnv + "GRPC_HEARTBEAT_INTERVAL"
	// GRPCKeepaliveTimeEnv stores the gRPC keepalive probe interval.
	GRPCKeepaliveTimeEnv = PrefixEnv + "GRPC_KEEPALIVE_TIME"
	// GRPCKeepaliveTimeoutEnv stores the gRPC keepalive acknowledgment timeout.
	GRPCKeepaliveTimeoutEnv = PrefixEnv + "GRPC_KEEPALIVE_TIMEOUT"
)

const (
	// GRPCReconnectMaxRetriesEnv stores the maximum number of connector reconnect attempts.
	GRPCReconnectMaxRetriesEnv = PrefixEnv + "GRPC_RECONNECT_MAX_RETRIES"
	// GRPCReconnectBaseBackoffEnv stores the initial reconnect backoff duration.
	GRPCReconnectBaseBackoffEnv = PrefixEnv + "GRPC_RECONNECT_BASE_BACKOFF"
	// GRPCReconnectMaxBackoffEnv stores the maximum reconnect backoff duration.
	GRPCReconnectMaxBackoffEnv = PrefixEnv + "GRPC_RECONNECT_MAX_BACKOFF"
)

const (
	// GRPCClientTLSSecretNameEnv stores the secret name used for client TLS material.
	GRPCClientTLSSecretNameEnv = PrefixEnv + "GRPC_CLIENT_TLS_SECRET_NAME"
	// GRPCRequireTLSEnv enables or requires TLS for gRPC connections.
	GRPCRequireTLSEnv = PrefixEnv + "GRPC_REQUIRE_TLS"
	// GRPCTLSCertFileEnv stores the server TLS certificate path.
	GRPCTLSCertFileEnv = PrefixEnv + "GRPC_TLS_CERT_FILE"
	// GRPCTLSKeyFileEnv stores the server TLS private-key path.
	GRPCTLSKeyFileEnv = PrefixEnv + "GRPC_TLS_KEY_FILE"
	// GRPCCAFileEnv stores the CA bundle path used to validate peers.
	GRPCCAFileEnv = PrefixEnv + "GRPC_CA_FILE"
	// GRPCClientCertFileEnv stores the client certificate path for mTLS.
	GRPCClientCertFileEnv = PrefixEnv + "GRPC_CLIENT_CERT_FILE"
	// GRPCClientKeyFileEnv stores the client private-key path for mTLS.
	GRPCClientKeyFileEnv = PrefixEnv + "GRPC_CLIENT_KEY_FILE"
	// GRPCHubServerNameEnv stores the expected TLS server name for hub connections.
	GRPCHubServerNameEnv = PrefixEnv + "GRPC_HUB_SERVER_NAME"
	// HubTLSCertFileEnv stores the hub server TLS certificate path.
	HubTLSCertFileEnv = PrefixEnv + "HUB_TLS_CERT_FILE"
	// HubTLSKeyFileEnv stores the hub server TLS private-key path.
	HubTLSKeyFileEnv = PrefixEnv + "HUB_TLS_KEY_FILE"
	// HubCAFileEnv stores the CA bundle path used by hub TLS.
	HubCAFileEnv = PrefixEnv + "HUB_CA_FILE"
)

const (
	// StorageProviderEnv stores the selected shared storage backend.
	StorageProviderEnv = PrefixEnv + "STORAGE_PROVIDER"
	// StoragePathEnv stores the local or mounted storage path.
	StoragePathEnv = PrefixEnv + "STORAGE_PATH"
	// StorageS3BucketEnv stores the S3 bucket name used for offloaded data.
	StorageS3BucketEnv = PrefixEnv + "STORAGE_S3_BUCKET"
	// StorageS3RegionEnv stores the S3 region used for offloaded data.
	StorageS3RegionEnv = PrefixEnv + "STORAGE_S3_REGION"
	// StorageS3EndpointEnv stores the custom S3 endpoint used for offloaded data.
	StorageS3EndpointEnv = PrefixEnv + "STORAGE_S3_ENDPOINT"
)

const (
	// StorageS3MaxRetriesEnv stores the maximum retry count for S3 requests.
	StorageS3MaxRetriesEnv = PrefixEnv + "S3_MAX_RETRIES"
	// StorageS3MaxBackoffEnv stores the maximum retry backoff for S3 requests.
	StorageS3MaxBackoffEnv = PrefixEnv + "S3_MAX_BACKOFF"
	// StorageS3ForcePathStyleEnv enables path-style S3 addressing when set.
	StorageS3ForcePathStyleEnv = PrefixEnv + "S3_FORCE_PATH_STYLE"
	// StorageS3AccessKeyIDEnv stores the S3 access key identifier.
	StorageS3AccessKeyIDEnv = PrefixEnv + "S3_ACCESS_KEY_ID"
	// StorageS3SecretAccessKeyEnv stores the S3 secret access key.
	StorageS3SecretAccessKeyEnv = PrefixEnv + "S3_SECRET_ACCESS_KEY"
	// StorageS3SessionTokenEnv stores the optional S3 session token.
	StorageS3SessionTokenEnv = PrefixEnv + "S3_SESSION_TOKEN"
	// StorageS3TimeoutEnv stores the timeout applied to S3 requests.
	StorageS3TimeoutEnv = PrefixEnv + "S3_TIMEOUT"
	// StorageS3TLSEnv enables TLS for S3 endpoints when set.
	StorageS3TLSEnv = PrefixEnv + "S3_TLS"
	// StorageS3MaxPartSizeEnv stores the multipart-upload part size for S3.
	StorageS3MaxPartSizeEnv = PrefixEnv + "S3_MAX_PART_SIZE"
	// StorageS3ConcurrencyEnv stores the multipart-upload concurrency for S3.
	StorageS3ConcurrencyEnv = PrefixEnv + "S3_CONCURRENCY"
	// StorageS3SSEEnv stores the server-side encryption mode for S3 objects.
	StorageS3SSEEnv = PrefixEnv + "S3_SSE"
	// StorageS3SSEKMSKeyEnv stores the KMS key identifier used for S3 SSE-KMS.
	StorageS3SSEKMSKeyEnv = PrefixEnv + "S3_SSE_KMS_KEY"
	// StorageS3OutputPrefixEnv stores the output prefix for offloaded S3 objects.
	StorageS3OutputPrefixEnv = PrefixEnv + "STORAGE_OUTPUT_PREFIX"
	// StorageS3InputPrefixEnv stores the input prefix for offloaded S3 objects.
	StorageS3InputPrefixEnv = PrefixEnv + "STORAGE_INPUT_PREFIX"
)

const (
	// S3IntegrationEnv enables S3 integration tests or integration-only code paths when set.
	S3IntegrationEnv = PrefixEnv + "S3_INTEGRATION_ENABLED"
	// SDKMetricsEnabledEnv enables shared SDK metrics instrumentation when set.
	SDKMetricsEnabledEnv = PrefixEnv + "SDK_METRICS_ENABLED"
	// SDKTracingEnabledEnv enables shared SDK tracing instrumentation when set.
	SDKTracingEnabledEnv = PrefixEnv + "SDK_TRACING_ENABLED"
	// TracePropagationEnv enables cross-component trace propagation behavior when set.
	TracePropagationEnv = PrefixEnv + "TRACE_PROPAGATION"
	// K8sUserAgentEnv stores the Kubernetes client user-agent string.
	K8sUserAgentEnv = PrefixEnv + "K8S_USER_AGENT"
	// K8sTimeoutEnv stores the default Kubernetes client timeout.
	K8sTimeoutEnv = PrefixEnv + "K8S_TIMEOUT"
	// K8sOperationTimeoutEnv stores the timeout for individual Kubernetes operations.
	K8sOperationTimeoutEnv = PrefixEnv + "K8S_OPERATION_TIMEOUT"
	// K8sPatchMaxRetriesEnv stores the maximum retry count for Kubernetes patch operations.
	K8sPatchMaxRetriesEnv = PrefixEnv + "K8S_PATCH_MAX_RETRIES"
)

const (
	// TransportSecurityModeTLS enforces TLS between engrams, connectors, and the hub.
	TransportSecurityModeTLS = "tls"
)

const (
	// StoryLabelKey labels resources with their owning Story.
	StoryLabelKey = "bubustack.io/story"
	// StoryNameLabelKey labels resources with a Story name field used for selectors.
	StoryNameLabelKey = "bubustack.io/story-name"
	// StepLabelKey labels resources with their owning step.
	StepLabelKey = "bubustack.io/step"
	// StepRunLabelKey labels resources with their owning StepRun.
	StepRunLabelKey = "bubustack.io/steprun"
	// StoryRunLabelKey labels resources with their owning StoryRun.
	StoryRunLabelKey = "bubustack.io/storyrun"
	// ParentStoryRunLabel labels resources with their parent StoryRun.
	ParentStoryRunLabel = "bubustack.io/parent-storyrun"
	// ParentStepLabel labels resources with their parent step.
	ParentStepLabel = "bubustack.io/parent-step"
	// LoopIndexLabelKey labels loop-generated resources with their iteration index.
	LoopIndexLabelKey = "bubustack.io/loop-index"
	// EngramLabelKey labels resources with their owning Engram.
	EngramLabelKey = "bubustack.io/engram"
	// MaterializeLabelKey labels resources created for template materialization.
	MaterializeLabelKey = "bubustack.io/materialize"
	// QueueLabelKey labels resources with the queue assigned for reconciliation.
	QueueLabelKey = "bubustack.io/queue"
	// QueuePriorityLabelKey labels resources with their queue priority.
	QueuePriorityLabelKey = "bubustack.io/queue-priority"
	// CorrelationIDLabelKey labels resources with a cross-component correlation identifier.
	CorrelationIDLabelKey = "bubustack.io/correlation-id"
)
