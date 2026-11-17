## Context

The `tencentcloud_dcx` resource (in `tencentcloud/services/dc/resource_tc_dcx.go`) manages TencentCloud Direct Connect dedicated tunnels. The resource currently defines a `bandwidth` schema field (`TypeInt`, `Optional`, `Computed`) but marks it `ForceNew: true`. This means changing the bandwidth in a Terraform configuration destroys and recreates the dedicated tunnel, which is disruptive.

**Current state of the `bandwidth` field:**
- Schema: `Optional`, `Computed`, `ForceNew: true`, `TypeInt`
- Create: passes `Bandwidth` to `CreateDirectConnectTunnel` request
- Read: reads `Bandwidth` from `DirectConnectTunnel.Bandwidth` (Describe response)
- Update: only handles `name` changes via `ModifyDirectConnectTunnelAttribute`; `bandwidth` is ignored (ForceNew handles it)

**API behavior analysis:**

| API | `Bandwidth` in Request | `Bandwidth` in Response |
|-----|------------------------|------------------------|
| `CreateDirectConnectTunnel` | Yes (`*int64`) | N/A |
| `DescribeDirectConnectTunnels` | N/A | Yes (`DirectConnectTunnel.Bandwidth *int64`) |
| `ModifyDirectConnectTunnelAttribute` | Yes (`*int64`, "专用通道带宽值，单位为M") | N/A |
| `DeleteDirectConnectTunnel` | No | N/A |

The `ModifyDirectConnectTunnelAttribute` API request struct (`ModifyDirectConnectTunnelAttributeRequest` in `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410`) already includes the `Bandwidth *int64` field. No SDK update is required.

## Goals / Non-Goals

**Goals:**
- Make the `bandwidth` parameter updatable in place by removing `ForceNew: true` from the schema field.
- In the Update function, when `bandwidth` changes, call `ModifyDirectConnectTunnelAttribute` with the `Bandwidth` request parameter.
- Maintain full backward compatibility — existing configurations and state continue to work unchanged.

**Non-Goals:**
- Changing the Create or Read logic for `bandwidth` (already correct).
- Adding new parameters beyond `bandwidth` to the resource.
- Modifying the `tencentcloud_dcx_instances` datasource.
- Adding new immutable-args handling (this resource has other updatable/non-updatable fields already handled by ForceNew individually; this change only affects `bandwidth`).

## Decisions

### Decision 1: Remove `ForceNew: true` from `bandwidth` schema field

**Rationale:** The `ModifyDirectConnectTunnelAttribute` API now accepts `Bandwidth`, so bandwidth can be updated in place. Removing `ForceNew` lets Terraform route bandwidth changes to the Update function instead of destroying and recreating the resource. The field remains `Optional`, `Computed`, and `TypeInt`.

**Alternatives considered:** Keep `ForceNew` and do nothing — rejected because it forces disruptive recreation when an in-place update is supported by the API.

### Decision 2: Handle `bandwidth` change in the Update function via `ModifyDirectConnectTunnelAttribute`

**Rationale:** The existing Update function already calls `ModifyDirectConnectTunnelAttribute` for `name` changes. We add a parallel handling block for `bandwidth`: when `d.HasChange("bandwidth")`, build a `ModifyDirectConnectTunnelAttributeRequest`, set `DirectConnectTunnelId` from `d.Id()`, set `Bandwidth` from the schema value, and call the API wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError(e)` error wrapping — consistent with the existing `name` update path and provider conventions.

**Alternatives considered:** Combine `name` and `bandwidth` into a single `ModifyDirectConnectTunnelAttribute` call — rejected to keep changes minimal and consistent with the existing per-field pattern (the existing code issues a separate API call for `name`). Each changed field triggers its own modify call, matching the current resource style.

### Decision 3: No changes to Read or Create

**Rationale:** The Read function already reads `DirectConnectTunnel.Bandwidth` and sets it in state. The Create function already passes `Bandwidth` to `CreateDirectConnectTunnel`. No modifications are needed.

## Risks / Trade-offs

- **[Risk] State drift if API update partially fails** → Mitigation: The Update call is wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError(e)`, and the function calls Read at the end to refresh state, consistent with the existing pattern.
- **[Risk] Users relying on ForceNew behavior may be surprised by in-place update** → Mitigation: This is strictly an improvement (less disruption). Existing configurations continue to work; the only behavioral change is that bandwidth updates no longer recreate the resource.
- **[Trade-off] Separate API calls for `name` and `bandwidth` when both change** → Accepted for consistency with the existing per-field update pattern in this resource.
