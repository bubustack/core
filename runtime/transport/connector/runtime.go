package connector

import "time"

// RuntimeTunables captures connector/channel runtime knobs shared across repos.
type RuntimeTunables struct {
	MessageTimeout     time.Duration
	ChannelSendTimeout time.Duration
	HangTimeout        time.Duration
	ChannelBufferSize  int
}

// RuntimeTunablesFromEnv resolves runtime knobs from env vars, falling back to
// the supplied defaults (or the package defaults when the provided values are
// zero). Env values always take precedence over defaults.
func RuntimeTunablesFromEnv(env Env, defaults RuntimeTunables) RuntimeTunables {
	env = ensureEnv(env)
	tunables := RuntimeTunables{
		MessageTimeout:    MessageTimeout(env, pickDuration(defaults.MessageTimeout, DefaultMessageTimeout)),
		ChannelBufferSize: ChannelBufferSize(env, pickInt(defaults.ChannelBufferSize, DefaultChannelBufferSize)),
	}

	if d := ChannelSendTimeout(env); d > 0 {
		tunables.ChannelSendTimeout = d
	} else {
		tunables.ChannelSendTimeout = defaults.ChannelSendTimeout
	}

	if d := HangTimeout(env); d > 0 {
		tunables.HangTimeout = d
	} else {
		tunables.HangTimeout = defaults.HangTimeout
	}

	return tunables
}

func pickDuration(candidate, fallback time.Duration) time.Duration {
	if candidate > 0 {
		return candidate
	}
	return fallback
}

func pickInt(candidate, fallback int) int {
	if candidate > 0 {
		return candidate
	}
	return fallback
}
