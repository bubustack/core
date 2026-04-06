package transport

import (
	"fmt"
	"strings"
)

const (
	// ProtocolVersion defines the transport contract version shared across hub, connector, and SDK.
	ProtocolVersion = "1.0.0"
	// ProtocolMetadataKey is the gRPC metadata key used to advertise the protocol version.
	ProtocolMetadataKey = "bubu-transport-protocol"
	// StartupCapabilitiesMetadataKey declares whether an initial
	// connector.capabilities snapshot is expected during startup. It is carried
	// on connector.ready and mirrored into the connector -> hub stream metadata
	// for passive hub-side observation.
	StartupCapabilitiesMetadataKey = "startup.capabilities"
	// StartupCapabilitiesRequired means the connector will send an initial
	// connector.capabilities update as part of startup.
	StartupCapabilitiesRequired = "required"
	// StartupCapabilitiesNone means the connector has no initial capability
	// snapshot to send during startup.
	StartupCapabilitiesNone = "none"
)

// ValidateProtocolVersion checks that the provided version matches this runtime exactly.
func ValidateProtocolVersion(version string) error {
	version = strings.TrimSpace(version)
	if version == "" {
		return fmt.Errorf("missing transport protocol version")
	}
	if version != ProtocolVersion {
		return fmt.Errorf("unsupported transport protocol version %q (expected %s)", version, ProtocolVersion)
	}
	return nil
}

// NormalizeStartupCapabilitiesMode validates the startup capability declaration
// and returns the normalized mode.
func NormalizeStartupCapabilitiesMode(mode string) (string, error) {
	mode = strings.ToLower(strings.TrimSpace(mode))
	switch mode {
	case StartupCapabilitiesRequired, StartupCapabilitiesNone:
		return mode, nil
	case "":
		return "", fmt.Errorf("missing %s metadata", StartupCapabilitiesMetadataKey)
	default:
		return "", fmt.Errorf(
			"invalid %s metadata %q (expected %s or %s)",
			StartupCapabilitiesMetadataKey, mode, StartupCapabilitiesRequired, StartupCapabilitiesNone,
		)
	}
}
