package connector

import (
	"testing"
	"time"

	"github.com/bubustack/core/contracts"
	coretransport "github.com/bubustack/core/runtime/transport"
	transportpb "github.com/bubustack/tractatus/gen/go/proto/transport/v1"
)

func TestLoadConfigFromEnvRejectsMissingProtocolVersion(t *testing.T) {
	env := testConnectorEnv(t, &transportpb.BindingInfo{
		Driver:   "grpc",
		Endpoint: "hub:7443",
	})

	if _, err := LoadConfigFromEnv(mapEnv(env)); err == nil {
		t.Fatalf("expected missing protocol version error")
	}
}

func TestLoadConfigFromEnvRejectsArbitraryGRPCPortEndpoint(t *testing.T) {
	env := testConnectorEnv(t, &transportpb.BindingInfo{
		Driver:          "grpc",
		Endpoint:        "hub:7443",
		ProtocolVersion: coretransport.ProtocolVersion,
	})
	env[contracts.GRPCPortEnv] = "127.0.0.1:50051"

	if _, err := LoadConfigFromEnv(mapEnv(env)); err == nil {
		t.Fatalf("expected invalid GRPC_PORT error")
	}
}

func TestLoadConfigFromEnvUsesSeparateHubDialTimeout(t *testing.T) {
	env := testConnectorEnv(t, &transportpb.BindingInfo{
		Driver:          "grpc",
		Endpoint:        "hub:7443",
		ProtocolVersion: coretransport.ProtocolVersion,
	})
	env[contracts.GRPCDialTimeoutEnv] = "5s"
	env[contracts.GRPCHubDialTimeoutEnv] = "17s"

	cfg, err := LoadConfigFromEnv(mapEnv(env))
	if err != nil {
		t.Fatalf("LoadConfigFromEnv failed: %v", err)
	}
	if cfg.LocalDialTimeout != 5*time.Second {
		t.Fatalf("expected local dial timeout 5s, got %s", cfg.LocalDialTimeout)
	}
	if cfg.HubDialTimeout != 17*time.Second {
		t.Fatalf("expected hub dial timeout 17s, got %s", cfg.HubDialTimeout)
	}
}

func TestLoadConfigFromEnvDoesNotFallbackHubDialTimeoutToSharedDialTimeout(t *testing.T) {
	env := testConnectorEnv(t, &transportpb.BindingInfo{
		Driver:          "grpc",
		Endpoint:        "hub:7443",
		ProtocolVersion: coretransport.ProtocolVersion,
	})
	env[contracts.GRPCDialTimeoutEnv] = "5s"

	cfg, err := LoadConfigFromEnv(mapEnv(env))
	if err != nil {
		t.Fatalf("LoadConfigFromEnv failed: %v", err)
	}
	if cfg.LocalDialTimeout != 5*time.Second {
		t.Fatalf("expected local dial timeout 5s, got %s", cfg.LocalDialTimeout)
	}
	if cfg.HubDialTimeout != DefaultHubDialTimeout {
		t.Fatalf("expected hub dial timeout default %s, got %s", DefaultHubDialTimeout, cfg.HubDialTimeout)
	}
}

func TestLoadConfigFromEnvRejectsMissingPodNamespaceEnv(t *testing.T) {
	env := testConnectorEnv(t, &transportpb.BindingInfo{
		Driver:          "grpc",
		Endpoint:        "hub:7443",
		ProtocolVersion: coretransport.ProtocolVersion,
	})
	delete(env, contracts.PodNamespaceEnv)
	env["POD_NAMESPACE"] = "legacy"

	if _, err := LoadConfigFromEnv(mapEnv(env)); err == nil {
		t.Fatalf("expected missing BUBU pod namespace error")
	}
}

func TestLoadConfigFromEnvRejectsLegacyPlaintextSecurityMode(t *testing.T) {
	env := testConnectorEnv(t, &transportpb.BindingInfo{
		Driver:          "grpc",
		Endpoint:        "hub:7443",
		ProtocolVersion: coretransport.ProtocolVersion,
	})
	env[contracts.TransportSecurityModeEnv] = "plaintext"

	if _, err := LoadConfigFromEnv(mapEnv(env)); err == nil {
		t.Fatalf("expected legacy plaintext security mode to be rejected")
	}
}

func TestLoadConfigFromEnvRejectsInvalidGeneration(t *testing.T) {
	env := testConnectorEnv(t, &transportpb.BindingInfo{
		Driver:          "grpc",
		Endpoint:        "hub:7443",
		ProtocolVersion: coretransport.ProtocolVersion,
	})
	env[contracts.ConnectorGenerationEnv] = "oops"

	if _, err := LoadConfigFromEnv(mapEnv(env)); err == nil {
		t.Fatalf("expected invalid generation error")
	}
}

func TestLoadConfigFromEnvRejectsMissingGeneration(t *testing.T) {
	env := testConnectorEnv(t, &transportpb.BindingInfo{
		Driver:          "grpc",
		Endpoint:        "hub:7443",
		ProtocolVersion: coretransport.ProtocolVersion,
	})
	delete(env, contracts.ConnectorGenerationEnv)

	if _, err := LoadConfigFromEnv(mapEnv(env)); err == nil {
		t.Fatalf("expected missing generation error")
	}
}

type mapEnv map[string]string

func (m mapEnv) Lookup(key string) string {
	return m[key]
}

func testConnectorEnv(t *testing.T, info *transportpb.BindingInfo) map[string]string {
	t.Helper()

	payload, err := coretransport.EncodeBindingEnvelope(coretransport.BindingReference{
		Name:      "binding",
		Namespace: "transport",
	}, info)
	if err != nil {
		t.Fatalf("encode binding envelope: %v", err)
	}

	return map[string]string{
		contracts.HubEndpointEnv:         "hub:7443",
		contracts.StoryRunIDEnv:          "storyrun-1",
		contracts.PodNamespaceEnv:        "default",
		contracts.StepNameEnv:            "step-a",
		contracts.TransportBindingEnv:    payload,
		contracts.ConnectorGenerationEnv: "1",
		contracts.GRPCPortEnv:            "50051",
		contracts.EngramNameEnv:          "engram-a",
	}
}
