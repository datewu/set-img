# Replica Scaling from the Dashboard Table (2026-11-06)

## 1. Context

The deployments and statefulsets table showed the replica count as read-only
text. Operators had to leave the page and run `kubectl scale` to change it.
The Replicas cell is now a number input that scales the workload in place.

## 2. Backend

- `internal/k8s/scale.go` — `ScalePath{Ns, Kind, Name, Replicas}` with
  `Scale(labels...)`:
  - validates a non-negative `int32`, a non-empty namespace/name, and
    `kind ∈ {deploy, sts}`;
  - fetches the workload to enforce the caller's ownership label through
    `checkLabels` (same check as image/env updates);
  - reads `GetScale`, sets `Spec.Replicas`, and writes `UpdateScale`.
  - `UpdateScale` is used instead of a full object update so a scale never
    rewrites the pod template, and no scale-to-zero/restore dance is involved.
- `cmd/api/my_handlers.go` — `myHandler.scaleResource` parses the form,
  rejects malformed/negative counts and invalid resource identity with `400`,
  then calls `ScalePath.Scale("image-user=<user>")`.
- `cmd/api/routes.go` — `PUT /my/scale` registered in the existing authenticated
  `/my` gzip group.

## 3. Frontend

- `front/table.html`:
  - the `.cell-replicas` cell now contains
    `<input type="number" name="replicas" min="0" step="1">` plus hidden
    `ns`/`kind`/`name` fields;
  - the cell fires `hx-put="/my/scale"` on `change` with `hx-include="this"`
    and `hx-swap="none"` (htmx includes every input inside the cell, including
    the hidden identity fields);
  - `data-current` keeps the last accepted value;
  - the live search reads `.replicas-input`'s value instead of the removed
    `.replicas-count` text.
- `front/layout.html`: the `htmx:afterRequest` hook handles `/my/scale`
  separately — success updates `data-current`, flashes the input green, and
  shows a toast; failure restores `data-current`, flashes red, and shows the
  server error.
- `front/static/page.css`: `.replicas-count` replaced by `.cell-replicas`
  width and `.replicas-input` styling.

## 4. Behaviour notes

- Scaling to `0` is allowed (that is a valid Kubernetes scale and useful for
  parking a workload).
- A resource not labeled `image-user=<current user>` is rejected server-side
  even if the form is forged. The response is the generic `500` produced by the
  existing label-check path (the image/env update endpoint behaves the same);
  a typed `403` would be a separate refinement.
- The image-update path still scales a workload to zero and restores the old
  replica count after 3 s (deployments) or 5 s (statefulsets). A scale issued
  inside that window can be overwritten by the restore. This is pre-existing
  behaviour and was left untouched.

## 5. Files changed

| File | Change |
| --- | --- |
| `internal/k8s/scale.go` | new `ScalePath` / `Scale` implementation |
| `internal/k8s/scale_test.go` | invalid-parameter coverage |
| `cmd/api/my_handlers.go` | `scaleResource` handler |
| `cmd/api/my_handlers_test.go` | handler validation coverage |
| `cmd/api/routes.go` | `PUT /my/scale` route |
| `front/table.html` | replicas number input, filter reads the input |
| `front/layout.html` | scale toast, success/error flash, value rollback |
| `front/static/page.css` | replicas input styling |
| `front/table_test.go` | asserts the rendered scaling markup |

## 6. Verification

- `go build ./...`, `go vet ./...` clean.
- `go test -race -vet=off ./...` — all packages pass, including the new handler
  and k8s tests and the table render assertions.
- Rendered `/my/deploys?ns=default` against the live microk8s cluster and
  confirmed the row contains the number input with the workload's current
  count and the hidden identity fields.
- End-to-end against scratch resources in the `default` namespace (created for
  the test, then deleted):
  - deployment scaled `1 → 3` via `PUT /my/scale` and Kubernetes reported `3`;
  - statefulset scaled `1 → 2` and Kubernetes reported `2`;
  - a deployment labeled `image-user=someoneelse` was rejected and stayed at
    replicas `1`;
  - `replicas=-1` returned `400` and changed nothing.

Remaining pre-existing `gofmt -l` drift (not introduced here, not changed):
`cmd/api/handlers.go`, `cmd/api/routes_test.go`, `cmd/main.go`.

This repository has no `CHANGELOG.md`, so no changelog entry was added; the
documented convention is not to create the file where the structure is absent.
