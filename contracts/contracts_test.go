package contracts

import "testing"

func TestTransportConnectorEnvConstantsRemainCanonical(t *testing.T) {
	if GRPCLocalEndpointEnv != PrefixEnv+"GRPC_LOCAL_ENDPOINT" {
		t.Fatalf("unexpected local endpoint env: %q", GRPCLocalEndpointEnv)
	}
	if GRPCHubDialTimeoutEnv != PrefixEnv+"GRPC_HUB_DIAL_TIMEOUT" {
		t.Fatalf("unexpected hub dial timeout env: %q", GRPCHubDialTimeoutEnv)
	}
	if GRPCCAFileEnv != PrefixEnv+"GRPC_CA_FILE" {
		t.Fatalf("unexpected grpc ca env: %q", GRPCCAFileEnv)
	}
	if GRPCHubServerNameEnv != PrefixEnv+"GRPC_HUB_SERVER_NAME" {
		t.Fatalf("unexpected hub server name env: %q", GRPCHubServerNameEnv)
	}
}
