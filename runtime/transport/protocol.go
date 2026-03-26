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
