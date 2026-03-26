package connector

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Env provides a minimal interface for resolving environment variables.
type Env interface {
	Lookup(key string) string
}

// EnvFunc adapts functions like os.Getenv to the Env interface.
type EnvFunc func(string) string

// Lookup returns the variable value for EnvFunc.
func (f EnvFunc) Lookup(key string) string {
	if f == nil {
		return ""
	}
	return f(key)
}

// OSEnv resolves variables from the host process environment.
var OSEnv Env = EnvFunc(os.Getenv)

func ensureEnv(env Env) Env {
	if env == nil {
		return OSEnv
	}
	return env
}

func trimEnv(env Env, key string) string {
	if env == nil {
		return ""
	}
	return strings.TrimSpace(env.Lookup(key))
}

func parsePositiveInt(env Env, key string, def int) int {
	val := trimEnv(env, key)
	if val == "" {
		return def
	}
	if n, err := strconv.Atoi(val); err == nil && n > 0 {
		return n
	}
	return def
}

func parsePositiveDuration(env Env, key string) time.Duration {
	val := trimEnv(env, key)
	if val == "" {
		return 0
	}
	if d, err := time.ParseDuration(val); err == nil && d > 0 {
		return d
	}
	return 0
}
