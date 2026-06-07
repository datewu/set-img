# Container Environment Variables Editing

## Overview

Users can view, add, update, and delete environment variables for any container within Deployments and StatefulSets directly from the web dashboard. Values are masked by default for security.

## User Interface

Each container row in the dashboard table has an **🔧 Env** button in the Actions column. Clicking it expands an inline editor panel below that container's row.

### Editor Layout

- Each environment variable is displayed as a **key-value input pair**:
  - **Key** input (text, monospace, indigo highlight)
  - `=` separator
  - **Value** input (`type="password"` by default, monospace)
  - **👁 Toggle** button to reveal/mask the value
  - **× Delete** button to remove the variable
- An **＋ Add Variable** button at the bottom-left adds a new empty key-value pair
- **Save Env** submits the changes
- **Close** collapses the editor panel

### Interactions

| Action | How |
|---|---|
| Edit a variable | Modify the key or value input directly |
| Add a variable | Click **＋ Add Variable**, fill in key and value |
| Delete a variable | Click the **×** button on that row |
| Show/hide value | Click the **👁** toggle icon |
| Save changes | Click **Save Env** |

## Data Flow

### Frontend → Backend

The HTMX form sends a `PUT /my/update/resource` request with:

| Field | Description |
|---|---|
| `ns` | Namespace |
| `kind` | `deploy` or `sts` |
| `name` | Resource name |
| `cname` | Container name |
| `image` | Current container image |
| `env_key` | Array of environment variable keys (ordered) |
| `env_val` | Array of environment variable values (ordered) |

The backend zips `env_key[]` and `env_val[]` by index to reconstruct the full environment variable list. Empty keys are skipped.

### Backend → Kubernetes

The update follows the same scale-to-zero pipeline used for image updates:

1. Fetch the current Deployment/StatefulSet spec
2. Locate the target container by name
3. Replace the container's `Env` slice with the parsed variables
4. If the image also changed, update it in the same operation
5. Scale replicas to 0, apply the update, then scale back to original count (in a background goroutine)

### Kubernetes `valueFrom` References

Environment variables sourced from ConfigMaps, Secrets, field references, or resource field references are displayed with a special format:

```
KEY=valueFrom(secretKeyRef:secret-name:key-name)
KEY=valueFrom(configMapKeyRef:configmap-name:key-name)
KEY=valueFrom(fieldRef:status.podIP)
KEY=valueFrom(resourceFieldRef:container-name:limits.cpu)
```

When saving:
- If the value string matches a `valueFrom(...)` pattern, it is parsed back into the corresponding Kubernetes `EnvVarSource` struct
- If the user changes a `valueFrom(...)` value to a plain string, it becomes a static `Value` field
- If the original `valueFrom` reference is left unchanged, the original struct is preserved exactly

## Files Involved

| File | Role |
|---|---|
| `internal/k8s/k8s.go` | `EnvVar` struct, `ContainerPath` fields, `ParseEnvStr`, `ParseEnvVarValue` helpers |
| `internal/k8s/deploy_crud.go` | Applies env updates to Deployments |
| `internal/k8s/sts_crud.go` | Applies env updates to StatefulSets |
| `cmd/api/my_handlers.go` | Parses `env_key`/`env_val` form arrays in `updateResouce` handler |
| `front/table.tpl.go` | `EnvKeyVal` struct, `mapEnv` helper, populates `Container.Env` |
| `front/table.html` | Env editor HTML template, inline JS for add/remove/toggle |
| `front/static/page.css` | Styling for `.env-pair`, `.env-key-input`, `.env-val-input`, visibility toggle, animations |
| `front/layout.html` | HTMX toast feedback for env form submissions |

## Validation Rules

- The `ContainerPath.valid()` method allows an empty `Img` field if `UpdateEnv` is true
- Empty key names are silently skipped during parsing
- At least one of image update or env update must be present for a valid request
