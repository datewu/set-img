# Manual CDN Update

## Overview

When a GitHub Actions deployment updates an image, the `set_cdn` goroutine is automatically triggered after a 30-second delay. If that goroutine fails, users can now manually trigger a CDN origin update from the web UI.

Only logged-in users with matching ingress resources can update CDN for their sites.

## User Interface

The **Update CDN** card appears on the index page for logged-in users. It contains a `<select>` dropdown populated with the user's site subdomains and a submit button.

### How It Works

1. User logs in via GitHub OAuth
2. The index page queries Kubernetes for ingresses labeled `ingress-user=<username>` in the `wu` namespace
3. Each ingress host matching `*.deoops.com` has its subdomain extracted and shown in the dropdown
4. User selects a site and clicks **Update CDN**
5. A background goroutine calls `set_cdn("<subdomain>.deoops.com")` to refresh the CDN origin

If no matching ingresses are found, a message is shown instead of the form.

## Security

- Only logged-in users see the form
- The server validates that the selected subdomain belongs to an ingress labeled `ingress-user=<username>`
- Users cannot update CDN for sites they do not own
- The dropdown only shows sites the user is authorized for (no free-text input)

## Data Flow

### Frontend → Backend

The HTMX form sends a `POST /my/update/cdn` request with:

| Field | Description |
|---|---|
| `subdomain` | The subdomain part of the site (e.g. `set-img`) |

### Backend → CDN Service

The `set_cdn` function:

1. Reads `R23_URL` from the environment
2. Sends `POST ${R23_URL}/auto-origin` with JSON body `{"origin": "<subdomain>.deoops.com"}`
3. Prints the response to stdout

## Files Involved

| File | Role |
|---|---|
| `cmd/api/handlers.go` | `index` handler fetches user's ingress hosts and populates `IndexView.Sites` |
| `cmd/api/my_handlers.go` | `updateCDN` handler validates subdomain ownership and triggers `set_cdn` |
| `cmd/api/routes.go` | Registers `POST /my/update/cdn` route under authenticated `myGzip` group |
| `internal/k8s/ingress.go` | `ListIngressHostsByLabel` queries ingresses and extracts subdomains |
| `front/index.tpl.go` | `IndexView` struct with `Sites []string` field |
| `front/index.html` | CDN form template with `<select>` dropdown, conditional rendering |
| `deploy/deployment.yaml` | ClusterRole grants `get`/`list` on `networking.k8s.io` ingresses |

## Prerequisites

- Ingress resources must be labeled with `ingress-user=<github-username>` in the `wu` namespace
- The `R23_URL` environment variable must be set to the CDN service endpoint
- The ClusterRole must include permissions for `networking.k8s.io` ingresses (`get`, `list`)
