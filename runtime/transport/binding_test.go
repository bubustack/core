package transport

import (
	"strings"
	"testing"

	transportpb "github.com/bubustack/tractatus/gen/go/proto/transport/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestEncodeAndDecodeBindingEnvelope(t *testing.T) {
	ref := BindingReference{Name: "binding-one", Namespace: "default"}
	info := &transportpb.BindingInfo{
		TransportRef: "demo",
		Driver:       "grpc",
		Endpoint:     "bubravoz:443",
		AudioCodecs:  []string{"opus"},
	}

	payload, err := EncodeBindingEnvelope(ref, info)
	if err != nil {
		t.Fatalf("EncodeBindingEnvelope() error = %v", err)
	}
	decoded, err := DecodeBindingInfo(payload)
	if err != nil {
		t.Fatalf("DecodeBindingInfo() error = %v", err)
	}

	if decoded.TransportRef != info.TransportRef || decoded.Driver != info.Driver || decoded.Endpoint != info.Endpoint {
		t.Fatalf("decoded info mismatch: %+v", decoded)
	}
}

func TestSanitizeBindingAnnotationValue(t *testing.T) {
	ref := BindingReference{Name: "binding-two", Namespace: "ns"}
	info := &transportpb.BindingInfo{TransportRef: "demo"}
	payload, err := EncodeBindingEnvelope(ref, info)
	if err != nil {
		t.Fatalf("EncodeBindingEnvelope() error = %v", err)
	}
	got := SanitizeBindingAnnotationValue(payload)
	want := "ns/binding-two"
	if got != want {
		t.Fatalf("SanitizeBindingAnnotationValue() = %q, want %q", got, want)
	}
	if SanitizeBindingAnnotationValue("") != "" {
		t.Fatalf("SanitizeBindingAnnotationValue(empty) should be empty")
	}
}

func TestParseBindingPayload_JSONEnvelope(t *testing.T) {
	ref := BindingReference{Name: "binding-three", Namespace: "transport"}
	info := &transportpb.BindingInfo{
		TransportRef: "demo",
		Driver:       "grpc",
	}
	payload, err := EncodeBindingEnvelope(ref, info)
	if err != nil {
		t.Fatalf("EncodeBindingEnvelope() error = %v", err)
	}
	parsed, err := ParseBindingPayload(payload)
	if err != nil {
		t.Fatalf("ParseBindingPayload() error = %v", err)
	}
	if parsed.Info == nil || parsed.Info.Driver != "grpc" {
		t.Fatalf("ParseBindingPayload() info mismatch: %+v", parsed.Info)
	}
	if parsed.Reference.Name != ref.Name || parsed.Reference.Namespace != ref.Namespace {
		t.Fatalf("ParseBindingPayload() ref mismatch: %+v", parsed.Reference)
	}
	if parsed.Raw == "" {
		t.Fatalf("ParseBindingPayload() raw empty")
	}
}

func TestParseBindingPayloadRejectsRawBindingInfoJSON(t *testing.T) {
	info := &transportpb.BindingInfo{
		TransportRef: "legacy",
		Driver:       "hub",
	}
	raw, err := protojson.Marshal(info)
	if err != nil {
		t.Fatalf("protojson.Marshal() error = %v", err)
	}
	if _, err := ParseBindingPayload(string(raw)); err == nil {
		t.Fatalf("expected raw binding info payload to be rejected")
	}
}

func TestParseBindingPayloadRejectsLegacyBase64Payload(t *testing.T) {
	if _, err := ParseBindingPayload("Zm9v"); err == nil {
		t.Fatalf("expected base64 payload to be rejected")
	}
}

func TestParseBindingPayloadRejectsEnvelopeWithoutBindingName(t *testing.T) {
	info := &transportpb.BindingInfo{
		TransportRef: "demo",
		Driver:       "grpc",
	}
	raw, err := protojson.Marshal(info)
	if err != nil {
		t.Fatalf("protojson.Marshal() error = %v", err)
	}
	value := `{"binding":` + string(raw) + `}`
	if _, err := ParseBindingPayload(value); err == nil {
		t.Fatalf("expected missing name to be rejected")
	}
}

func TestParseBindingPayloadRejectsEnvelopeWithoutBindingPayload(t *testing.T) {
	value := `{"name":"binding"}`
	if _, err := ParseBindingPayload(value); err == nil {
		t.Fatalf("expected missing binding payload to be rejected")
	}
}

func TestParseBindingPayloadRejectsOversizedInput(t *testing.T) {
	value := strings.Repeat("a", maxBindingPayloadBytes+1)
	if _, err := ParseBindingPayload(value); err == nil {
		t.Fatalf("expected oversize binding payload error")
	}
}
