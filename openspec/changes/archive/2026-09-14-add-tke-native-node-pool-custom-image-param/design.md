## Context

The `tencentcloud_kubernetes_native_node_pool` resource manages TKE native node pools via the `tencentcloud-sdk-go/tencentcloud/tke/v20220501` SDK package. The resource schema defines a `native` block (`TypeList`, `MaxItems=1`) whose `Elem` schema maps to the `CreateNativeNodePoolParam` / `UpdateNativeNodePoolParam` / `NativeNodePoolInfo` SDK structs.

The SDK already exposes a `CustomImage *string` field on all three structs:
- `CreateNativeNodePoolParam.CustomImage` (used by `CreateNodePool`)
- `UpdateNativeNodePoolParam.CustomImage` (used by `ModifyNodePool`)
- `NativeNodePoolInfo.CustomImage` (returned by `DescribeNodePools` / `DescribeCluster`)

The Terraform resource does not currently expose this field. This change adds it.

**Current State:**
- The `native` block schema in `resource_tc_kubernetes_native_node_pool.go` already includes sibling fields such as `machine_type`, `key_ids`, `data_disks`, etc.
- The Create handler builds a `CreateNativeNodePoolParam` from the `native` schema map; the Update handler builds an `UpdateNativeNodePoolParam`; the Read handler flattens `NativeNodePoolInfo` back into the `native` schema map.

**Constraints:**
- Must maintain backward compatibility (the new field is `Optional`).
- Must follow existing patterns for reading values from the `nativeMap` (using `helper.InterfacesHeadMap` + map key lookup) and flattening into the `nativeMap` on Read.
- The field is mutable (no `ForceNew`) because `ModifyNodePool` accepts `Native.CustomImage`.

## Goals / Non-Goals

**Goals:**
- Add `custom_image` (String, Optional) to the `native` block schema.
- Wire it into Create (`CreateNativeNodePoolParam.CustomImage`), Update (`UpdateNativeNodePoolParam.CustomImage`), and Read (`NativeNodePoolInfo.CustomImage`).
- Keep the change minimal and consistent with the existing sibling fields.

**Non-Goals:**
- Not refactoring the existing `native` block handling.
- Not adding validation beyond what the API enforces.
- Not changing the resource ID format or provider registration.

## Decisions

### Decision 1: Field placement inside the `native` block

**Choice:** Add `custom_image` as a field inside the `native` block schema (alongside `machine_type`, `key_ids`, etc.).

**Rationale:**
- The cloud API path is `request.Native.CustomImage`, so the field belongs to the `Native` sub-structure.
- All existing `Native.*` parameters are modeled inside the `native` block, so this is consistent.

**Alternatives Considered:**
- Top-level schema field: Rejected because it does not match the `Native.CustomImage` API path and would break the established nesting pattern.

### Decision 2: Field is mutable (not ForceNew)

**Choice:** `custom_image` is `Optional` without `ForceNew`.

**Rationale:**
- `ModifyNodePool` accepts `Native.CustomImage` via `UpdateNativeNodePoolParam.CustomImage`, so the value can be updated in place.
- The existing Update handler already rebuilds the `UpdateNativeNodePoolParam` from the full `nativeMap` whenever `native` has changed, so the new field is picked up automatically once added to the map-reading code.

**Alternatives Considered:**
- `ForceNew: true`: Rejected because the API supports updating the value, and forcing recreation would be disruptive.

### Decision 3: Read handler flattens from `NativeNodePoolInfo.CustomImage`

**Choice:** In the Read handler, set `nativeMap["custom_image"]` from `respData.Native.CustomImage` when it is not nil.

**Rationale:**
- `NativeNodePoolInfo.CustomImage` is returned by the Describe API (marked "may return null").
- Following the existing pattern, nil-check before setting to avoid overwriting state with null.
- This keeps Terraform state in sync with the cloud and enables drift detection.

## Risks / Trade-offs

**Risk:** The API may return `CustomImage` as null for node pools created without a custom image.
- **Mitigation:** Nil-check before `d.Set`, consistent with all other `native` fields; existing configurations without `custom_image` produce no plan diff.

**Trade-off:** None significant — this is a straightforward additive change following established patterns.
