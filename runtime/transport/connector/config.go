package connector

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/bubustack/core/contracts"
	coretransport "github.com/bubustack/core/runtime/transport"
)

const (
	defaultLocalHost = "127.0.0.1"
	defaultLocalPort = "50051"
	// DefaultHeartbeatInterval is the connector heartbeat interval used when the
	// binding payload does not provide an override.
	DefaultHeartbeatInterval = 30 * time.Second
	// DefaultLocalDialTimeout is the timeout for connector-to-local-runtime dials.
	DefaultLocalDialTimeout = 5 * time.Second
	// DefaultHubDialTimeout is the timeout for connector-to-hub dials.
	DefaultHubDialTimeout = 10 * time.Second
)

// Config captures the ambient connector runtime settings.
type Config struct {
	LocalEndpoint            string
	LocalServerName          string
	HubEndpoint              string
	StoryRunName             string
	Namespace                string
	StepID                   string
	LocalDialTimeout         time.Duration
	HubDialTimeout           time.Duration
	Binding                  coretransport.BindingPayload
	BindingHeartbeatInterval time.Duration
	Generation               int32
}

// LoadConfigFromEnv builds a Config from shared connector env vars.
func LoadConfigFromEnv(env Env) (*Config, error) {
	env = ensureEnv(env)
	if err := validateTransportSecurityMode(env); err != nil {
		return nil, err
	}
	cfg, err := loadBaseConfigFromEnv(env)
	if err != nil {
		return nil, err
	}
	payload, err := loadBindingPayload(env, cfg.Namespace)
	if err != nil {
		return nil, err
	}
	cfg.Binding = payload

	cfg.Generation, err = loadConnectorGeneration(env)
	if err != nil {
		return nil, err
	}
	cfg.LocalServerName = resolveLocalServerName(env, cfg.Namespace)

	return cfg, nil
}

func loadBaseConfigFromEnv(env Env) (*Config, error) {
	localEndpoint, err := resolveLocalEndpoint(env)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		LocalEndpoint:            localEndpoint,
		HubEndpoint:              trimEnv(env, contracts.HubEndpointEnv),
		StoryRunName:             trimEnv(env, contracts.StoryRunIDEnv),
		Namespace:                resolveNamespace(env),
		StepID:                   trimEnv(env, contracts.StepNameEnv),
		LocalDialTimeout:         durationWithDefault(env, DefaultLocalDialTimeout, contracts.GRPCDialTimeoutEnv),
		HubDialTimeout:           durationWithDefault(env, DefaultHubDialTimeout, contracts.GRPCHubDialTimeoutEnv),
		BindingHeartbeatInterval: durationWithDefault(env, DefaultHeartbeatInterval, contracts.TransportHeartbeatIntervalEnv),
	}
	if err := validateRequiredConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func validateRequiredConfig(cfg *Config) error {
	switch {
	case cfg.HubEndpoint == "":
		return fmt.Errorf("%s must be set", contracts.HubEndpointEnv)
	case cfg.StoryRunName == "":
		return fmt.Errorf("%s must be set", contracts.StoryRunIDEnv)
	case cfg.Namespace == "":
		return fmt.Errorf("%s must be set", contracts.PodNamespaceEnv)
	case cfg.StepID == "":
		return fmt.Errorf("%s must be set", contracts.StepNameEnv)
	default:
		return nil
	}
}

func loadBindingPayload(env Env, namespace string) (coretransport.BindingPayload, error) {
	bindingValue := trimEnv(env, contracts.TransportBindingEnv)
	if bindingValue == "" {
		return coretransport.BindingPayload{}, fmt.Errorf(
			"%s must be set",
			contracts.TransportBindingEnv,
		)
	}
	payload, err := coretransport.ParseBindingPayload(bindingValue)
	if err != nil {
		return coretransport.BindingPayload{}, fmt.Errorf(
			"parse %s: %w",
			contracts.TransportBindingEnv,
			err,
		)
	}
	if payload.Info == nil {
		return coretransport.BindingPayload{}, fmt.Errorf(
			"binding info missing in %s",
			contracts.TransportBindingEnv,
		)
	}
	if err := coretransport.ValidateProtocolVersion(payload.Info.GetProtocolVersion()); err != nil {
		return coretransport.BindingPayload{}, fmt.Errorf(
			"invalid transport protocol version in binding: %w",
			err,
		)
	}
	if strings.TrimSpace(payload.Reference.Namespace) == "" {
		payload.Reference.Namespace = namespace
	}
	return payload, nil
}

func loadConnectorGeneration(env Env) (int32, error) {
	genStr := trimEnv(env, contracts.ConnectorGenerationEnv)
	if genStr == "" {
		return 0, fmt.Errorf("%s must be set", contracts.ConnectorGenerationEnv)
	}
	n, err := strconv.ParseInt(genStr, 10, 32)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("invalid %s %q", contracts.ConnectorGenerationEnv, genStr)
	}
	return int32(n), nil
}

func resolveLocalServerName(env Env, namespace string) string {
	engramName := trimEnv(env, contracts.EngramNameEnv)
	if engramName == "" || namespace == "" {
		return ""
	}
	return fmt.Sprintf("%s.%s.svc.cluster.local", engramName, namespace)
}

func resolveLocalEndpoint(env Env) (string, error) {
	if endpoint := trimEnv(env, contracts.GRPCLocalEndpointEnv); endpoint != "" {
		return normalizeLocalEndpoint(endpoint)
	}
	port := trimEnv(env, contracts.GRPCPortEnv)
	if port == "" {
		port = defaultLocalPort
	}
	normalizedPort, err := normalizePortValue(port, contracts.GRPCPortEnv)
	if err != nil {
		return "", fmt.Errorf("%w; use %s for full endpoints", err, contracts.GRPCLocalEndpointEnv)
	}
	return net.JoinHostPort(defaultLocalHost, normalizedPort), nil
}

func resolveNamespace(env Env) string {
	return trimEnv(env, contracts.PodNamespaceEnv)
}

func validateTransportSecurityMode(env Env) error {
	mode := strings.ToLower(trimEnv(env, contracts.TransportSecurityModeEnv))
	if mode == "" || mode == contracts.TransportSecurityModeTLS {
		return nil
	}
	return fmt.Errorf(
		"invalid %s %q: only %q is supported",
		contracts.TransportSecurityModeEnv,
		mode,
		contracts.TransportSecurityModeTLS,
	)
}

func durationWithDefault(env Env, def time.Duration, keys ...string) time.Duration {
	for _, key := range keys {
		if d := parsePositiveDuration(env, key); d > 0 {
			return d
		}
	}
	return def
}

func normalizeLocalEndpoint(endpoint string) (string, error) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return "", fmt.Errorf("%s must not be empty", contracts.GRPCLocalEndpointEnv)
	}
	if path, ok := strings.CutPrefix(endpoint, "unix://"); ok {
		if strings.TrimSpace(path) == "" {
			return "", fmt.Errorf("%s unix endpoint missing path", contracts.GRPCLocalEndpointEnv)
		}
		return endpoint, nil
	}
	if !strings.Contains(endpoint, ":") {
		port, err := normalizePortValue(endpoint, contracts.GRPCLocalEndpointEnv)
		if err != nil {
			return "", err
		}
		return net.JoinHostPort(defaultLocalHost, port), nil
	}
	host, port, err := net.SplitHostPort(endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid %s %q: %w", contracts.GRPCLocalEndpointEnv, endpoint, err)
	}
	if host == "" {
		host = defaultLocalHost
	}
	port, err = normalizePortValue(port, contracts.GRPCLocalEndpointEnv)
	if err != nil {
		return "", err
	}
	return net.JoinHostPort(host, port), nil
}

func normalizePortValue(port string, key string) (string, error) {
	port = strings.TrimSpace(port)
	if port == "" {
		return "", fmt.Errorf("%s must be set", key)
	}
	n, err := strconv.Atoi(port)
	if err != nil || n < 1 || n > 65535 {
		return "", fmt.Errorf("invalid %s %q", key, port)
	}
	return strconv.Itoa(n), nil
}
