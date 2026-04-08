# Contributing to core

`core` is shared API surface for the rest of the BubuStack ecosystem. Changes
here affect bobrapet, bobravoz-grpc, and bubu-sdk-go, so contributions need to
be explicit, well-tested, and easy to review.

## Reporting bugs

- Check [existing issues](https://github.com/bubustack/core/issues?q=is%3Aissue) first.
- Include the following details when filing a bug:
  - The affected package or helper (`contracts`, `templating`, `runtime/*`, or release/CI/docs).
  - A minimal reproduction: template snippet, input payload, config example, or direct Go call.
  - The `core` version or commit you tested, plus any downstream repo/version involved.
  - Returned errors, stack traces, or logs showing the incorrect behaviour.
  - Downstream impact in bobrapet, bobravoz-grpc, or bubu-sdk-go if applicable.
- Apply the relevant `kind/*`, `area/*`, and `priority/*` labels when triaging.

## Requesting enhancements

- Use the [feature request template](https://github.com/bubustack/core/issues/new?template=feature_request.md).
- Describe the shared problem the change solves across the BubuStack ecosystem.
- Sketch the desired API, config shape, or contract change when possible.
- If the change requires coordinated updates in downstream repos, say that up front so releases can be planned.

## Pull requests

1. Fork the repo, branch from `main`, and keep the PR focused.
2. Preserve package boundaries: `contracts` stay canonical, `templating` stays shared and consumer-agnostic, and runtime helpers stay dependency-light.
3. Run the quality gates before requesting review:
   ```bash
   make fmt
   make vet
   make test
   make lint
   make test-coverage # optional but recommended for larger changes
   make tidy          # if dependencies changed
   ```
4. Update docs, templates, and examples when behaviour changes. Call out any required downstream follow-up in bobrapet, bobravoz-grpc, or bubu-sdk-go.
5. Use the PR template to record commands run, test evidence, and linked issues.

## Development workflow

### Prerequisites

- Go 1.26+ (matching `go.mod`)
- `make` and bash

### Setup

1. Fork the repository and clone your fork.
2. `cd core`
3. `make help` to inspect the available targets.
4. `make test` before opening your first PR.

### Running tests

```bash
make test             # race-enabled unit suite
make test-coverage    # optional coverage profile
```

### Linting

```bash
make lint
```

### Commit style & Code of Conduct

- Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) so release automation can generate accurate changelog entries (for example: `feat: add shared transport env helper`, `fix: guard nil template path lookup`).
- Participation in this project is governed by the [Contributor Covenant Code of Conduct](./CODE_OF_CONDUCT.md). Report unacceptable behaviour to community@bubustack.io or via the org Discussions moderation channel.
