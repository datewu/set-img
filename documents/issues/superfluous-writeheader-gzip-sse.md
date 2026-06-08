# Superfluous WriteHeader from gzipResponseWriter

## Symptom

The server logs the following error at runtime:

```json
{
  "level": "ERROR",
  "time": "2026-06-08T14:16:09Z",
  "message": "http: superfluous response.WriteHeader call from github.com/datewu/gtea/handler.(*gzipResponseWriter).WriteHeader (middleware.go:247)\n",
  "properties": [null]
}
```

## Root Cause

The `/my/logs` endpoint uses Server-Sent Events (SSE) via `sse.SSE()`, but was registered under a route group that applies `GzipMiddleware`:

```go
my := r.Group("/my", h.auth)
my.Use(handler.GzipMiddleware)      // ← gzip applied to ALL /my/* routes
my.Get("/logs", h.logs)             // ← h.logs calls sse.SSE() internally
```

The `gzipResponseWriter` in `gtea/handler/middleware.go` does **not** implement `http.Flusher`. This triggers the following chain:

1. `GzipMiddleware` wraps `w` in `gzipResponseWriter`
2. `sse.SSE()` calls `w.WriteHeader(200)` → `gzipResponseWriter.WriteHeader(200)` → `w.ResponseWriter.WriteHeader(200)` — **first WriteHeader, OK**
3. `sse.SSE()` checks `f, ok := w.(http.Flusher)` → **`ok` is `false`** (gzipResponseWriter lacks `Flush()`)
4. `sse.SSE()` falls into the error path: `http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)`
5. `http.Error` calls `w.WriteHeader(500)` → `gzipResponseWriter.WriteHeader(500)` → `w.ResponseWriter.WriteHeader(500)` — **SUPERFLUOUS** (200 already sent in step 2)

Go's `net/http` package detects that `WriteHeader` was called on a response that already had its headers written and logs the warning.

## Why Gzip and SSE Are Incompatible

Even if `gzipResponseWriter` implemented `http.Flusher`, combining gzip with SSE is still problematic:

- **Gzip buffers data** for compression efficiency — it waits to accumulate enough bytes before emitting a compressed block
- **SSE requires real-time flushing** — each event must be sent to the client immediately
- Flushing the gzip writer after every small SSE message defeats the purpose of compression and can even increase payload size due to gzip overhead per block

## Fix Applied

### Application (this repo) — `cmd/api/routes.go`

Moved the `/logs` SSE route into a separate route group **without** `GzipMiddleware`:

```go
// Before (broken):
my := r.Group("/my", h.auth)
my.Use(handler.GzipMiddleware)
my.Get("/logs", h.logs)

// After (fixed):
my := r.Group("/my", h.auth)
my.Use(handler.GzipMiddleware)
// ... other gzip-safe routes ...

// SSE route without gzip — SSE and gzip are incompatible
myNoGzip := r.Group("/my", h.auth)
myNoGzip.Get("/logs", h.logs)
```

### gtea Library — `handler/middleware.go`

The `gzipResponseWriter` should also be hardened in the `gtea` repo with two improvements:

#### 1. Guard against superfluous WriteHeader calls

```go
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
	wroteHeader bool
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.Writer.Write(b)
}
```

#### 2. Implement `http.Flusher` for downstream compatibility

```go
func (w *gzipResponseWriter) Flush() {
	if flusher, ok := w.Writer.(interface{ Flush() }); ok {
		flusher.Flush()
	}
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
```

These changes ensure:
- `WriteHeader` is called at most once on the underlying writer (eliminates superfluous warnings)
- Downstream code that relies on `http.Flusher` (SSE, streaming, etc.) works correctly through the gzip wrapper
- `Write()` auto-calls `WriteHeader(200)` if not yet called (consistent with `http.ResponseWriter` semantics)

## Secondary Bug: gzPool Initialization

The `gzPool` in `GzipMiddleware` has a wasteful initialization:

```go
var gzPool = sync.Pool{
	New: func() any {
		w := gzip.NewWriter(io.Discard)                // writer #1 → io.Discard
		gzip.NewWriterLevel(w, gzip.BestCompression)    // writer #2 → writer #1 (discarded!)
		return w                                         // returns writer #1
	},
}
```

The second `gzip.NewWriterLevel` call creates a writer that compresses into the first writer (which goes to `io.Discard`), but its return value is thrown away. The fix:

```go
var gzPool = sync.Pool{
	New: func() any {
		w, _ := gzip.NewWriterLevel(io.Discard, gzip.BestCompression)
		return w
	},
}
```

This is harmless in practice because `gz.Reset(w)` is always called before use, but it wastes a gzip writer allocation on every pool miss.

## Files Involved

| File | Change |
|---|---|
| `cmd/api/routes.go` | Moved `/logs` to a non-gzip route group |
| `gtea/handler/middleware.go` | (External) Added `wroteHeader` guard, `Flush()`, fixed `gzPool` |
