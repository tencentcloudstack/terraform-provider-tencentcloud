## Why

The `docker_graph_path` field in `tencentcloud_kubernetes_cluster_attachment` currently defaults to `/var/lib/docker`, but TKE clusters have migrated from Docker to containerd as the container runtime. The correct default path for containerd is `/var/lib/containerd`. Continuing to ship `/var/lib/docker` as the default causes misconfigured nodes for users who do not explicitly set this field, while changing it outright would break existing Terraform state for users who rely on the current default.

## What Changes

- Remove the hardcoded `Default: "/var/lib/docker"` from both `docker_graph_path` field definitions in the `worker_config` and `worker_config_overrides` blocks of `resource_tc_kubernetes_cluster_attachment.go`.
- Add `Computed: true` to both `docker_graph_path` fields so that the actual value returned by the API is stored in state without forcing a diff for existing resources that already have `/var/lib/docker` in their state.
- Update the field `Description` to reflect that the default is determined by the API/platform.

## Capabilities

### New Capabilities
<!-- None – this is a bugfix/compatibility change only -->

### Modified Capabilities
<!-- No spec-level requirement change; this is an implementation-level fix within an existing resource. -->

## Impact

- **File**: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_attachment.go`
  - Two `docker_graph_path` schema entries updated (one in `worker_config`, one in `worker_config_overrides`).
- **Backward compatibility**: Removing `Default` + adding `Computed` ensures existing state values are preserved; no plan diff for users who have not explicitly set the field.
- **New users**: Will receive the platform default (currently `/var/lib/containerd`) without needing to set the field explicitly.
- **No API or SDK changes required.**
