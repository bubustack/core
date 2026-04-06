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

func TestNormalizeStartupCapabilitiesMode(t *testing.T) {
	got, err := NormalizeStartupCapabilitiesMode(StartupCapabilitiesRequired)
	if err != nil || got != StartupCapabilitiesRequired {
		t.Fatalf("expected %q to pass, got mode=%q err=%v", StartupCapabilitiesRequired, got, err)
	}
	if got, err := NormalizeStartupCapabilitiesMode(" NONE "); err != nil || got != StartupCapabilitiesNone {
		t.Fatalf("expected %q to normalize to %q, got mode=%q err=%v", " NONE ", StartupCapabilitiesNone, got, err)
	}
	if _, err := NormalizeStartupCapabilitiesMode(""); err == nil {
		t.Fatalf("expected missing startup capabilities mode to fail")
	}
	if _, err := NormalizeStartupCapabilitiesMode("legacy"); err == nil {
		t.Fatalf("expected invalid startup capabilities mode to fail")
	}
}
