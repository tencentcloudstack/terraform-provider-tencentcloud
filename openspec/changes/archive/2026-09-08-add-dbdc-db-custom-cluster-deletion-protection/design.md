## Context

The `tencentcloud_dbdc_db_custom_cluster` resource currently supports creating DB Custom clusters with `cluster_name`, `container_network`, `api_server_network`, `cluster_description`, and `tags` parameters. The DBDC `CreateDBCustomCluster` API also accepts a `DeletionProtection` parameter (bool, default `true`) that controls whether cluster deletion protection is enabled, but the Terraform resource does not expose this parameter.

**Current state:**
- Resource file: `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_cluster.go`
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029` (vendored, already includes `DeletionProtection`)

**API behavior analysis:**

| API | DeletionProtection in Request | DeletionProtection in Response |
|-----|-------------------------------|-------------------------------|
| `CreateDBCustomCluster` | Yes (`DeletionProtection *bool`, default `true`) | No |
| `DescribeDBCustomClusterDetail` | N/A | Yes (`DeletionProtection *bool`) |
| `ModifyDBCustomClusterAttributes` | Yes (`DeletionProtection *bool`) | N/A |
| `DestroyDBCustomCluster` | No | N/A |

**Key observation:** `DeletionProtection` is available in Create, Read, and Update APIs, enabling full lifecycle management. This is unlike other parameters (e.g., `cluster_name`, `container_network`) which are only available in Create and are thus `ForceNew`.

## Goals / Non-Goals

**Goals:**
- Add `deletion_protection` (Optional, TypeBool) parameter to `tencentcloud_dbdc_db_custom_cluster` resource schema
- Pass `DeletionProtection` to `CreateDBCustomCluster` API request when specified by user
- Read `DeletionProtection` from `DescribeDBCustomClusterDetail` API response to support state refresh and import
- Update `DeletionProtection` via `ModifyDBCustomClusterAttributes` API when the user changes the parameter
- Maintain full backward compatibility — existing configurations continue to work unchanged

**Non-Goals:**
- Adding `deletion_protection` to any dbdc datasource (out of scope)
- Making other existing parameters mutable (they remain `ForceNew`)

## Decisions

### Decision 1: `deletion_protection` is Optional and updatable (not ForceNew)

**Rationale:** The `ModifyDBCustomClusterAttributes` API accepts `DeletionProtection`, so the parameter can be updated in-place without recreating the resource. Using `ForceNew: true` would unnecessarily destroy and recreate the cluster when the user only wants to toggle deletion protection. The parameter is `Optional` (not `Computed`) because the API has a default of `true`, but Terraform should only send the value when the user explicitly sets it.

### Decision 2: Update via `ModifyDBCustomClusterAttributes` in the existing Update function

**Rationale:** The existing `Update` function already handles `tags` changes via `ModifyDBCustomClusterTags`. We add a separate `if d.HasChange("deletion_protection")` block that calls `ModifyDBCustomClusterAttributes` with the `ClusterId` and `DeletionProtection` fields. This keeps the update logic modular and consistent with the existing pattern. The `ModifyDBCustomClusterAttributes` call is wrapped in `resource.Retry` with `tccommon.WriteRetryTimeout` for retry handling, consistent with the existing Create/Delete patterns.

### Decision 3: Use `d.GetOkExists` for reading the bool parameter in Create

**Rationale:** Since `DeletionProtection` is a bool with API default `true`, using `d.GetOk()` would not distinguish between "user set false" and "user did not set". However, since the API has a sensible default (`true`), we only need to pass the value when the user explicitly sets it. Using `d.GetOk()` is sufficient here because: if the user does not set the parameter, the API defaults to `true`; if the user sets it to `false`, `d.GetOk()` returns `ok=true` with `v=false`. If the user sets it to `true`, `d.GetOk()` returns `ok=true` with `v=true`. This correctly covers all cases.

### Decision 4: Read `DeletionProtection` from `DescribeDBCustomClusterDetail` response

**Rationale:** The `DescribeDBCustomClusterDetail` API response includes `DeletionProtection` as a field. The existing Read function calls `service.DescribeDBCustomClusterById` which returns the cluster detail. We add a nil-check for `respData.DeletionProtection` and set it to state, consistent with the existing pattern for other fields.

## Risks / Trade-offs

- **[Risk] API default vs Terraform default mismatch**: The API defaults `DeletionProtection` to `true` when not specified. If the user does not set `deletion_protection` in Terraform, the API will set it to `true`, and the Read function will populate `deletion_protection=true` in state. This is expected behavior and not a breaking change.
  - **Mitigation:** No special handling needed; the Read function correctly reflects the API's actual state.
