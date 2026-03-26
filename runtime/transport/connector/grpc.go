package connector

import (
	"time"

	"github.com/bubustack/core/contracts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	// DefaultMaxMessageSize is the fallback max receive/send message size for connector gRPC links.
	DefaultMaxMessageSize = 10 * 1024 * 1024
	// DefaultChannelBufferSize is the fallback buffered channel depth for connector stream loops.
	DefaultChannelBufferSize = 16
	// DefaultMessageTimeout is the fallback per-message timeout for connector runtime operations.
	DefaultMessageTimeout = 30 * time.Second
)

// ServerOptions returns standard gRPC server options derived from env overrides.
func ServerOptions(env Env, defaultRecv, defaultSend int) []grpc.ServerOption {
	env = ensureEnv(env)
	recv := parsePositiveInt(env, contracts.GRPCMaxRecvBytesEnv, defaultRecv)
	send := parsePositiveInt(env, contracts.GRPCMaxSendBytesEnv, defaultSend)
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(recv),
		grpc.MaxSendMsgSize(send),
	}
	if ka, ok := serverKeepaliveOption(env); ok {
		opts = append(opts, ka)
	}
	return opts
}

// ClientCallOptions returns dialing call options derived from env overrides.
func ClientCallOptions(env Env, defaultRecv, defaultSend int) []grpc.CallOption {
	env = ensureEnv(env)
	recv := parsePositiveInt(env, contracts.GRPCClientMaxRecvBytesEnv, defaultRecv)
	send := parsePositiveInt(env, contracts.GRPCClientMaxSendBytesEnv, defaultSend)
	return []grpc.CallOption{
		grpc.MaxCallRecvMsgSize(recv),
		grpc.MaxCallSendMsgSize(send),
	}
}

// ChannelBufferSize resolves the buffered channel size used by stream loops.
func ChannelBufferSize(env Env, def int) int {
	return parsePositiveInt(env, contracts.GRPCChannelBufferSizeEnv, def)
}

// MessageTimeout resolves the connector message timeout env.
func MessageTimeout(env Env, def time.Duration) time.Duration {
	if d := parsePositiveDuration(env, contracts.GRPCMessageTimeoutEnv); d > 0 {
		return d
	}
	return def
}

// ChannelSendTimeout resolves the optional send timeout env.
func ChannelSendTimeout(env Env) time.Duration {
	return parsePositiveDuration(env, contracts.GRPCChannelSendTimeoutEnv)
}

// HangTimeout resolves the passive hang timeout env.
func HangTimeout(env Env) time.Duration {
	return parsePositiveDuration(env, contracts.GRPCHangTimeoutEnv)
}

// DialTimeout resolves the shared GRPC dial timeout env with fallback.
func DialTimeout(env Env, def time.Duration) time.Duration {
	if d := parsePositiveDuration(env, contracts.GRPCDialTimeoutEnv); d > 0 {
		return d
	}
	return def
}

func serverKeepaliveOption(env Env) (grpc.ServerOption, bool) {
	var params keepalive.ServerParameters
	var have bool
	if d := parsePositiveDuration(env, contracts.GRPCKeepaliveTimeEnv); d > 0 {
		params.Time = d
		have = true
	}
	if d := parsePositiveDuration(env, contracts.GRPCKeepaliveTimeoutEnv); d > 0 {
		params.Timeout = d
		have = true
	}
	if !have {
		return nil, false
	}
	return grpc.KeepaliveParams(params), true
}
