# 🌐 core — Shared Primitives for BubuStack

[![Go Reference](https://pkg.go.dev/badge/github.com/bubustack/core.svg)](https://pkg.go.dev/github.com/bubustack/core)
[![Go Report Card](https://goreportcard.com/badge/github.com/bubustack/core)](https://goreportcard.com/report/github.com/bubustack/core)

`core` is the shared foundation layer for the BubuStack ecosystem. It
centralizes contracts, templating behavior, and runtime helpers that must stay
identical across **bobrapet** (operator), **bobravoz-grpc** (transport hub),
and **bubu-sdk-go** (runtime SDK).

## 🔗 Quick Links

- API reference: https://pkg.go.dev/github.com/bubustack/core
- Issues: https://github.com/bubustack/core/issues
- Support: [SUPPORT.md](./SUPPORT.md)
- Security policy: [SECURITY.md](./SECURITY.md) (use GitHub Security Advisories for vulnerability reports)

## 🌟 Key Features

- **Canonical shared contracts** for env vars, annotations, labels, config keys, and index fields.
- **Templating engine** with normalization, validation, path resolution, storage-ref detection, and shared evaluation rules.
- **Runtime helpers** for env injection, naming, identity, stage metadata, operator config, and feature toggles.
- **Transport primitives** for binding envelopes, protocol validation, connector config, dial/listen behavior, and TLS setup.
- **Shared API surface discipline** so downstream repos do not re-implement behavior that should remain centralized.

## 🏗️ Architecture

| Package | Description |
|---------|-------------|
| `contracts` | Canonical annotations, labels, env vars, config keys, and index constants shared across repos. |
| `templating` | Shared template normalization, validation, evaluation, path lookup, storage-ref detection, and caching. |
| `runtime/bootstrap` | Shared startup registration and structured bootstrap logging helpers. |
| `runtime/env` | Deterministic env builders for gRPC tuning, debug flags, and timestamps. |
| `runtime/featuretoggles` | Shared telemetry, trace propagation, logging, and metrics toggle wiring. |
| `runtime/identity` | StoryRun-safe label and service-account naming helpers. |
| `runtime/naming` | DNS-safe resource name composition helpers. |
| `runtime/operatorconfig` | Shared ConfigMap-backed operator config manager for controller-runtime consumers. |
| `runtime/stage` | Structured StoryRun/step metadata helpers for logs and related shared metadata flow. |
| `runtime/storage` | Storage provider env, secret, timeout, and volume wiring helpers. |
| `runtime/transport` | Binding envelopes, protocol checks, and shared transport env helpers. |
| `runtime/transport/connector` | Connector config, dial/listen helpers, TLS setup, and runtime tunables. |

> If a behavior must stay identical between operator, hub, and SDK, it belongs
> here with tests.

## 🚀 Quick Start

### 1. Add the module

```bash
go get github.com/bubustack/core@latest
```

### 2. Import the helpers you need

```go
import (
	"github.com/bubustack/core/contracts"
	coreenv "github.com/bubustack/core/runtime/env"
	operatorconfig "github.com/bubustack/core/runtime/operatorconfig"
	coretransport "github.com/bubustack/core/runtime/transport"
)
```

### 3. Use the shared primitives

- Use `contracts` for stable keys such as `contracts.StoryRunLabelKey`.
- Use `coreenv.BuildGRPCTuningEnv` and `coreenv.AppendStartedAtEnv` when wiring workloads.
- Use `coretransport.EncodeBindingEnvelope` and `coretransport.ParseBindingPayload` for transport binding serialization.
- Use `operatorconfig.NewManager` when you need shared ConfigMap-backed runtime settings.

## 🛠️ Local Development

### Prerequisites

- Go 1.25.1 or newer (matching `go.mod`)
- `make`

### Setup

1. Clone your fork of `core`.
2. Run `make help` to inspect available targets.
3. Run `make test` before opening a change.

### Common commands

```bash
make build
make test
make test-coverage
make lint
```

Release automation is managed by Release Please. Conventional Commits feed
`CHANGELOG.md` and release PR generation.

## 📢 Support, Security, and Changelog

- See [SUPPORT.md](./SUPPORT.md) for help channels and support expectations.
- See [SECURITY.md](./SECURITY.md) for vulnerability reporting. Security issues must go through GitHub Security Advisories, not public issues or email.
- See [CHANGELOG.md](./CHANGELOG.md) for release history and upcoming shared-library changes.

## 🤝 Community

- Code of Conduct: [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md)
- Discord: https://discord.gg/dysrB7D8H6

## 📄 License

Copyright 2025 BubuStack.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
