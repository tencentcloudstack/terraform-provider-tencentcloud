## Context

The `tencentcloud_instance` resource is the core resource for managing CVM instances. It currently supports several `Placement` parameters through the schema:
- `availability_zone` → maps to `Placement.Zone`
- `project_id` → maps to `Placement.ProjectId`
- `dedicated_cluster_id` and `cdh_instance_type` → map to `Placement.HostIds` (CDH support)

However, it does not support the newer dedicated resource pack placement parameters (`DedicatedResourcePackTenancy` and `DedicatedResourcePackIds`) that were added to the TencentCloud CVM API. These parameters allow users to create instances using pre-purchased resource pool packs.

**Current State:**
- The `Placement` struct in the SDK already contains these fields (SDK is up-to-date)
- The `resource_tc_instance.go` creates a `Placement` object in the Create function (line ~604)
- The existing pattern uses flat schema fields (not nested structures) for placement parameters

**Constraints:**
- Must maintain backward compatibility (no breaking changes)
- Must follow existing patterns in `resource_tc_instance.go`
- Fields must be `ForceNew: true` since changing placement requires recreating the instance
- Need to validate that both parameters are specified together (API requirement)

## Goals / Non-Goals

**Goals:**
- Add support for `DedicatedResourcePackTenancy` and `DedicatedResourcePackIds` placement parameters
- Follow the existing flat schema pattern (not nested structures, as requested by user)
- Ensure proper validation (both parameters must be specified together)
- Maintain backward compatibility

**Non-Goals:**
- Not adding support for reading/updating these parameters (Read function will not populate them from API response)
- Not supporting the `RackId` output field (output-only, not useful for Terraform state)
- Not refactoring the existing placement parameter handling

## Decisions

### Decision 1: Flat Schema Fields vs Nested Structure

**Choice:** Use flat schema fields at the root level

**Rationale:**
- Consistent with existing placement parameters (`availability_zone`, `project_id`, etc.)
- Explicitly requested by user: "对于资源来讲,直接参数输入,不要作为结构体输入"
- Simpler user experience (no need to create nested blocks)

**Alternatives Considered:**
- Nested `placement` block: Would be more API-aligned but inconsistent with existing pattern and user requirements

### Decision 2: Field Naming Convention

**Choice:** Use snake_case field names: `dedicated_resource_pack_tenancy` and `dedicated_resource_pack_ids`

**Rationale:**
- Matches Terraform/Go naming conventions
- Consistent with other fields in the resource (`dedicated_cluster_id`)
- Clear and self-documenting

### Decision 3: Validation Strategy

**Choice:** Use `RequiredWith` validation on both fields

**Rationale:**
- API requires both parameters to be specified together when using resource pool packs
- Provides clear error message at plan time (before API call)
- Follows existing pattern (e.g., `force_replace_placement_group_id` uses `RequiredWith`)

**Implementation:**
```go
"dedicated_resource_pack_tenancy": {
    Type:         schema.TypeString,
    Optional:     true,
    ForceNew:     true,
    RequiredWith: []string{"dedicated_resource_pack_ids"},
    Description:  "...",
},
"dedicated_resource_pack_ids": {
    Type:         schema.TypeList,
    Optional:     true,
    ForceNew:     true,
    RequiredWith: []string{"dedicated_resource_pack_tenancy"},
    Elem:         &schema.Schema{Type: schema.TypeString},
    Description:  "...",
},
```

### Decision 4: Read Function Handling

**Choice:** DO populate these fields in the Read function from `instance.Placement`

**Rationale:**
- The SDK's `Placement` in `DescribeInstances` response includes `DedicatedResourcePackTenancy` and `DedicatedResourcePackIds` fields
- These fields are not marked as "input-only" in the SDK (unlike `HostId` and `RackId` which are marked as "仅用于出参")
- Populating these fields enables proper drift detection and state management
- Consistent with other placement parameters like `availability_zone` and `project_id` which are populated from `instance.Placement`
- Even though these are ForceNew parameters, reading them ensures Terraform state remains accurate

**Implementation:**
```go
// In resourceTencentCloudInstanceRead()
if instance.Placement != nil {
    if instance.Placement.DedicatedResourcePackTenancy != nil {
        _ = d.Set("dedicated_resource_pack_tenancy", instance.Placement.DedicatedResourcePackTenancy)
    }
    if len(instance.Placement.DedicatedResourcePackIds) > 0 {
        _ = d.Set("dedicated_resource_pack_ids", helper.StringsInterfaces(instance.Placement.DedicatedResourcePackIds))
    }
}
```

### Decision 5: Update Function Handling

**Choice:** No special handling needed in Update function (fields are ForceNew)

**Rationale:**
- ForceNew ensures Terraform triggers resource replacement if these values change
- No update API exists for placement parameters

## Risks / Trade-offs

**Risk:** API validation errors if invalid resource pack IDs are provided
- **Mitigation:** Clear documentation with examples; API will return validation errors during apply

**Risk:** Users may not understand the relationship between resource pool packs and these parameters
- **Mitigation:** Add documentation notes explaining these fields work with `tencentcloud_cvm_resource_pool_packs` resource

**Trade-off:** Not populating these fields in Read function means Terraform won't detect out-of-band changes
- **Acceptable:** These are immutable after creation, and out-of-band modifications are not possible through the API

**Trade-off:** Flat schema increases the number of root-level fields
- **Acceptable:** Matches user requirement and existing pattern; keeps UX simple

## Migration Plan

Not applicable - this is a backward-compatible addition. No migration needed.

## Open Questions

None - design is straightforward and follows established patterns.

---

## Disaster Recover Group IDs Priority

### Decision 6: `disaster_recover_group_ids` Takes Priority Over `placement_group_id`

**Choice:** When `disaster_recover_group_ids` is set, ignore `placement_group_id` in Create, skip `placement_group_id` readback, and reject `placement_group_id` changes in Update.

**Rationale:**
- The `ModifyInstancesDisasterRecoverGroup` API now supports `DisasterRecoverGroupIds` for batch setting (up to 3 group IDs)
- When both fields are configured, using both would lead to conflicting behaviors
- Priority model ensures deterministic behavior: `disaster_recover_group_ids` always wins

**Implementation:**
- **Create**: `if disaster_recover_group_ids is set → use it; else → fallback to placement_group_id`
- **Read**: `if disaster_recover_group_ids in state → skip placement_group_id readback` (avoids plan diffs)
- **Update**: `if disaster_recover_group_ids is set → reject placement_group_id changes`

### Decision 7: Remove `ConflictsWith`, Use Priority Model

**Choice:** Remove `ConflictsWith: ["placement_group_id"]` from `disaster_recover_group_ids`, replace with runtime priority checks.

**Rationale:**
- `ConflictsWith` prevents both fields from being specified simultaneously at plan time
- Priority model is more flexible: users can configure both and the provider handles the precedence
- Allows gradual migration: users can add `disaster_recover_group_ids` while keeping `placement_group_id` as fallback

### Decision 8: MaxItems: 3 for `disaster_recover_group_ids`

**Choice:** Use `MaxItems: 3` on the `disaster_recover_group_ids` TypeSet field.

**Rationale:**
- The API restricts `DisasterRecoverGroupIds` to a maximum of 3 group IDs
- Terraform SDK `MaxItems` provides plan-time validation before any API call

### Decision 9: Conditional `placement_group_id` Readback

**Choice:** Only read `placement_group_id` from API when `disaster_recover_group_ids` is NOT in state.

**Rationale:**
- When `disaster_recover_group_ids` is set, the API may return a single `DisasterRecoverGroupId` that differs from the user's configured `placement_group_id`
- Unconditionally reading back would overwrite state and cause false plan diffs on subsequent `terraform plan`
- Guarding the readback with `d.GetOk("disaster_recover_group_ids")` preserves state integrity
