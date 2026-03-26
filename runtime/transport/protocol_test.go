package transport

import "testing"

func TestValidateProtocolVersionRequiresExactMatch(t *testing.T) {
	if err := ValidateProtocolVersion(ProtocolVersion); err != nil {
		t.Fatalf("expected exact protocol version to pass, got %v", err)
	}
	if err := ValidateProtocolVersion("1.2.3"); err == nil {
		t.Fatalf("expected non-exact protocol version to fail")
	}
	if err := ValidateProtocolVersion(""); err == nil {
		t.Fatalf("expected missing protocol version to fail")
	}
}
