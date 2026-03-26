---
name: "Bug report"
about: "Report a reproducible issue in core contracts, templating, or runtime helpers"
labels: ["kind/bug", "status/triage"]
---

## Area
- [ ] `contracts`
- [ ] `templating`
- [ ] `runtime/bootstrap`
- [ ] `runtime/env`
- [ ] `runtime/featuretoggles`
- [ ] `runtime/identity`
- [ ] `runtime/naming`
- [ ] `runtime/operatorconfig`
- [ ] `runtime/stage`
- [ ] `runtime/storage`
- [ ] `runtime/transport`
- [ ] Repo health / CI / release automation

## What happened?
Tell us what broke, what you expected, and what happened instead.

## Minimal reproduction
1. Input, config, or template snippet
2. Commands you ran (`make`, `go test`, or downstream reproduction)
3. Exact package/function involved

```go
// paste the smallest possible reproduction here
```

## Logs / errors
- Stack traces or returned errors
- Relevant template/config payloads
- Downstream impact in bobrapet, bobravoz-grpc, or bubu-sdk-go if applicable

## Additional context
Anything else we should know, including recent behavior changes or local environment details.
