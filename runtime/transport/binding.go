package transport

import (
	"encoding/json"
	"fmt"
	"strings"

	transportpb "github.com/bubustack/tractatus/gen/go/proto/transport/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

const maxBindingPayloadBytes = 1 << 20

// BindingReference identifies a TransportBinding (namespace/name).
type BindingReference struct {
	Name      string
	Namespace string
}

type bindingEnvelope struct {
	Name      string          `json:"name"`
	Namespace string          `json:"namespace,omitempty"`
	Binding   json.RawMessage `json:"binding,omitempty"`
}

// BindingPayload represents the parsed contents of a serialized binding env value.
type BindingPayload struct {
	Reference BindingReference
	Info      *transportpb.BindingInfo
	Raw       string
}

// EncodeBindingEnvelope serializes binding metadata into the shared JSON payload
// consumed by controllers, connectors, and SDK helpers.
func EncodeBindingEnvelope(ref BindingReference, info *transportpb.BindingInfo) (string, error) {
	name := strings.TrimSpace(ref.Name)
	if name == "" {
		return "", fmt.Errorf("binding name is required")
	}
	if info == nil {
		return "", fmt.Errorf("binding info is nil")
	}
	payload, err := protojson.Marshal(info)
	if err != nil {
		return "", fmt.Errorf("encode binding info: %w", err)
	}
	env := bindingEnvelope{
		Name:      name,
		Namespace: strings.TrimSpace(ref.Namespace),
		Binding:   payload,
	}
	data, err := json.Marshal(env)
	if err != nil {
		return "", fmt.Errorf("encode binding envelope: %w", err)
	}
	return string(data), nil
}

// DecodeBindingInfo extracts the BindingInfo proto from a serialized binding envelope.
func DecodeBindingInfo(value string) (*transportpb.BindingInfo, error) {
	payload, err := ParseBindingPayload(value)
	if err != nil {
		return nil, err
	}
	if payload.Info == nil {
		return nil, fmt.Errorf("binding info missing")
	}
	return payload.Info, nil
}

// SanitizeBindingAnnotationValue converts the serialized binding envelope into the canonical
// "namespace/name" annotation value (or just name when namespace is empty). Returns "" on failure.
func SanitizeBindingAnnotationValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	var env bindingEnvelope
	if err := json.Unmarshal([]byte(value), &env); err != nil {
		return ""
	}
	name := strings.TrimSpace(env.Name)
	if name == "" {
		return ""
	}
	namespace := strings.TrimSpace(env.Namespace)
	if namespace == "" {
		return name
	}
	return fmt.Sprintf("%s/%s", namespace, name)
}

func decodeBindingInfoJSON(data []byte) (*transportpb.BindingInfo, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("binding info payload empty")
	}
	var info transportpb.BindingInfo
	if err := protojson.Unmarshal(data, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// ParseBindingPayload returns the parsed binding reference, binding info, and raw payload.
func ParseBindingPayload(value string) (BindingPayload, error) {
	payload := BindingPayload{}
	value = strings.TrimSpace(value)
	if value == "" {
		return payload, fmt.Errorf("binding payload empty")
	}
	if len(value) > maxBindingPayloadBytes {
		return payload, fmt.Errorf("binding payload exceeds max size of %d bytes", maxBindingPayloadBytes)
	}
	payload.Raw = value

	if !strings.HasPrefix(value, "{") {
		return BindingPayload{}, fmt.Errorf("binding payload must be a JSON envelope")
	}
	envPayload, err := parseBindingJSON(value)
	if err != nil {
		return BindingPayload{}, err
	}
	return envPayload, nil
}

func parseBindingJSON(value string) (BindingPayload, error) {
	var env bindingEnvelope
	if err := json.Unmarshal([]byte(value), &env); err != nil {
		return BindingPayload{}, fmt.Errorf("unmarshal binding envelope: %w", err)
	}
	if strings.TrimSpace(env.Name) == "" {
		return BindingPayload{}, fmt.Errorf("binding envelope name is required")
	}
	if len(env.Binding) == 0 {
		return BindingPayload{}, fmt.Errorf("binding envelope binding payload is required")
	}

	info, err := decodeBindingInfoJSON(env.Binding)
	if err != nil {
		return BindingPayload{}, err
	}

	return BindingPayload{
		Reference: BindingReference{
			Name:      strings.TrimSpace(env.Name),
			Namespace: strings.TrimSpace(env.Namespace),
		},
		Info: info,
		Raw:  value,
	}, nil
}
