# R23_URL Environment Variable Refinement

## Changes

### Hardcoded URL extracted to env var
The `set_cdn` function in `cmd/api/handlers.go` previously hardcoded `http://r2-s3/auto-origin`. The base URL is now read from the `R23_URL` environment variable, with `/auto-origin` appended at runtime.

### Default value in .envrc
`R23_URL` is exported in `.envrc` with the value `http://r2-s3` so existing development workflows are unaffected.

## Files modified
- `cmd/api/handlers.go` — replaced hardcoded URL with `os.Getenv("R23_URL") + "/auto-origin"`, added `"os"` import
- `.envrc` — added `R23_URL` export
