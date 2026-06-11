# Table Template Container Evaluation Error

## Symptom

When loading the deployments or statefulsets dashboard, the frontend fails to load and instead returns a `500 Internal Server Error` with the following message:

```json
{
  "error": {
    "detail": "template: table.html:98:186: executing \"content\" at \u003c.containers\u003e: can't evaluate field containers in type front.container",
    "error": "the server encountered a problem and could not process your request"
  }
}
```

This prevents users from accessing the workloads dashboard and viewing container logs.

## Root Cause

In `front/table.html` at line 98, the template was written as:

```html
<button class="btn btn-secondary btn-sm btn-icon-only btn-logs" onclick="toggleLogRow('log-{{ $resourceId }}-0', '{{ $resNs }}', '{{ $kind }}', '{{ $name }}', '{{ (index .Containers 0).Name }}')" title="View Container Logs">
```

This code is executed inside a range block iterating over the workload's containers:

```html
{{ range $index, $element := .Containers }}
  ...
{{ end }}
```

Inside a `range` block, the dot `.` represents the current element being iterated, which is of type `front.Container`. Because `front.Container` does not have a `Containers` field, Go's template engine failed to evaluate `(index .Containers 0)` and panicked.

## Fix

Since we are already iterating over the `.Containers` list, the current container object is available directly as `.Name`.

We replaced the evaluation expression with `{{ .Name }}`:

```html
<button class="btn btn-secondary btn-sm btn-icon-only btn-logs" onclick="toggleLogRow('log-{{ $resourceId }}-0', '{{ $resNs }}', '{{ $kind }}', '{{ $name }}', '{{ .Name }}')" title="View Container Logs">
```

## Router Group Cleanup & Robustness

To ensure that real-time Server-Sent Events (SSE) `/my/logs` are never accidentally wrapped by the `GzipMiddleware` (which would buffer logs and break streaming), the route registrations in `cmd/api/routes.go` were refactored to make middleware application explicit and independent of line ordering:

```go
func myRoutes(app *gtea.App, r *router.RoutesGroup) {
	h := &myHandler{app: app}
	my := r.Group("/my", h.auth)
	my.Get("/logs", h.logs) // Registered on a group WITHOUT GzipMiddleware

	myGzip := my.Group("", handler.GzipMiddleware) // Explicit sub-group with GzipMiddleware
	myGzip.Delete("/logout", h.logout)
	myGzip.Get("/profile", h.profile)
	myGzip.Get("/deploys", h.deploys)
	myGzip.Get("/sts", h.sts)
	myGzip.Put("/update/resource", h.updateResouce)
	myGzip.Get("/pods", h.listPods)
}
```
