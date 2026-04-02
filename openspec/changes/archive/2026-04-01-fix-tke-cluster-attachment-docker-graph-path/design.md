## Context

`tencentcloud_kubernetes_cluster_attachment` has two `docker_graph_path` fields — one inside `worker_config` and one inside `worker_config_overrides` — both carrying `Default: "/var/lib/docker"`. TKE nodes now use containerd instead of Docker as the container runtime, so `/var/lib/docker` is stale. The challenge is making the change without breaking existing Terraform state for users who never set the field explicitly (their state currently records `/var/lib/docker` as the value).

## Goals / Non-Goals

**Goals:**
- Stop shipping `/var/lib/docker` as the hardcoded default so new deployments receive the platform's actual default (`/var/lib/containerd`).
- Preserve backward compatibility: existing resources whose state already holds `/var/lib/docker` must not generate a spurious plan diff.
- Keep the schema change minimal and non-breaking.

**Non-Goals:**
- Migrating existing state files automatically.
- Changing any API call logic or read/update paths.
- Modifying the `worker_config_overrides` deprecated fields beyond `docker_graph_path`.

## Decisions

### Remove `Default`, add `Computed: true`

**Decision**: Replace `Default: "/var/lib/docker"` with `Computed: true` on both `docker_graph_path` fields.

**Rationale**:
- `Default` in Terraform SDK means the provider injects the value client-side when the user omits the field. This bypasses whatever the API returns and permanently stamps `/var/lib/docker` into state.
- `Computed: true` tells Terraform "the API may fill in this value"; if the user doesn't set it, the provider reads back the API value and stores that. Existing state values are honoured as-is, so no diff is generated for already-created resources.
- Alternative considered: keep `Default` and change its value to `/var/lib/containerd`. **Rejected** because this would create a plan diff for every existing resource that currently has `/var/lib/docker` in state, forcing a `ForceNew` recreation (since `docker_graph_path` is `ForceNew`), which is a destructive breaking change.

### Update Description

Update both descriptions from `"Docker graph path. Default is /var/lib/docker."` to reflect that the default is now platform-determined.

## Risks / Trade-offs

| Risk | Mitigation |
|---|---|
| Existing users who relied on the implicit default and never set the field: their state holds `/var/lib/docker`; after this change, Read will overwrite state with what the API returns (likely still `/var/lib/docker` for old nodes). | `Computed` only updates state on Read; `ForceNew` is not triggered by a state-only change, so no recreation. |
| Users on mixed clusters (some old, some new nodes) may see different values across instances. | No impact — the field is per-resource; each resource independently reads its own state. |
| The field is `ForceNew`, so any actual value change in config would still trigger recreation. | This is intended existing behavior; we are not changing it. |
