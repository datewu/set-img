# Go Dependency Upgrade (2026-11-06)

## 1. Context

Upgrade every Go module to its latest stable release, with
`k8s.io/client-go` as the priority.

## 2. Result

`go get -u ./...` followed by `go mod tidy` upgraded the module graph:

| Module | Before | After |
| --- | --- | --- |
| `k8s.io/client-go` | v0.36.2 | v0.37.2 |
| `k8s.io/api` | v0.36.2 | v0.37.2 |
| `k8s.io/apimachinery` | v0.36.2 | v0.37.2 |
| `golang.org/x/net` | v0.56.0 | v0.62.0 |
| `golang.org/x/sys` | v0.46.0 | v0.49.0 |
| `golang.org/x/text` | v0.38.0 | v0.43.0 |
| `golang.org/x/oauth2` | v0.36.0 | v0.37.0 |
| `github.com/spf13/pflag` | v1.0.10 | v1.1.0 |
| `google.golang.org/protobuf` | v1.36.12-0.2026… | v1.36.12 |
| `sigs.k8s.io/structured-merge-diff/v6` | v6.4.0 | v6.4.2 |
| `k8s.io/utils` | 2026-06 pseudo | 2026-10 pseudo |

`github.com/datewu/gtea` (v0.6.0) and `github.com/fsnotify/fsnotify` (v1.10.1)
were already at their latest versions.

`k8s.io/client-go` v0.37.2 is the latest **stable** release as of 2026-11-06.
`go list -m k8s.io/client-go@latest` resolves there; v0.38.0 exists only as a
beta (`v0.38.0-beta.0`) and was deliberately not adopted.

## 3. kube-openapi must stay on the k8s-pinned version

Letting `-u` float `k8s.io/kube-openapi` to the newest pseudo-version breaks the
build:

```
cannot use typeSchema.Types (variable of type
[]"sigs.k8s.io/structured-merge-diff/v7/schema".TypeDef) as
[]"sigs.k8s.io/structured-merge-diff/v6/schema".TypeDef value in struct literal
```

The newer `kube-openapi` switched to `structured-merge-diff/v7`, while
`k8s.io/apimachinery` v0.37.2 still imports `/v6` in
`pkg/util/managedfields/internal/typeconverter.go:51`.

Fix: pin `k8s.io/kube-openapi` to the revision both k8s v0.37.2 module files
require (`v0.0.0-20260721132016-d427ff9ee9ad`). `v7` then drops out of the module
graph.

The initial `-u` run also bumped the `go` directive to `1.27.0` because that
newer kube-openapi required it. `go mod tidy -go=1.26.3` restored `go 1.26.3`
after the pin, so `deploy/Dockerfile` keeps building with
`golang:1.26.3-alpine3.22` — no Docker change was needed. Every module in the
final graph requires Go ≤ 1.26.3.

## 4. Files changed

| File | Change |
| --- | --- |
| `go.mod` | dependency upgrade, kube-openapi pin, `go` directive unchanged |
| `go.sum` | checksums for the new graph |

## 5. Verification

- `go mod verify` — all modules verified.
- `go mod tidy` is idempotent (no further diff).
- `go build ./...` and `go vet ./...` clean.
- `go test -race -vet=off ./...` — all packages pass.

This repository has no `CHANGELOG.md`, so no changelog entry was added; the
documented convention is not to create the file where the structure is absent.
