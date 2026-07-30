## Context

The `tencentcloud_placement_group` resource currently supports creating placement groups with `name`, `type`, and `affinity` parameters. The CVM `CreateDisasterRecoverGroup` API also accepts `Strategy` (SPREAD/PARTITION) and `PartitionCount` (for PARTITION strategy), but the Terraform resource does not expose these parameters.

**Current state:**
- Resource file: `tencentcloud/services/cvm/resource_tc_placement_group.go`
- Service layer: `tencentcloud/services/cvm/service_tencentcloud_cvm.go` (`CreatePlacementGroup` function)
- Constants: `tencentcloud/services/cvm/extension_cvm.go`
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312`

**API behavior analysis (after SDK update):**

| API | Strategy in Request | PartitionCount in Request | Strategy in Response | PartitionCount in Response |
|-----|--------------------|--------------------------|---------------------|---------------------------|
| `CreateDisasterRecoverGroup` | Yes (`Strategy *string`, SPREAD/PARTITION) | Yes (`PartitionCount *int64`, range 2-30) | No | Yes (`PartitionCount *int64`) |
| `DescribeDisasterRecoverGroups` | N/A | N/A | Yes (`Strategy *string`) | Yes (`PartitionCount *int64`) |
| `ModifyDisasterRecoverGroupAttribute` | No | No | N/A | N/A |
| `DeleteDisasterRecoverGroups` | No | No | N/A | N/A |

**Key constraint:** `Strategy` and `PartitionCount` are only available in the Create request for writing. The Describe API's `DisasterRecoverGroup` struct now includes both fields (after SDK update), so they can be refreshed on Read.

## Goals / Non-Goals

**Goals:**
- Add `strategy` (Optional, Computed, immutable after creation) parameter to `tencentcloud_placement_group` with valid values `SPREAD` and `PARTITION`
- Add `partition_count` (Optional, Computed, immutable after creation) parameter to `tencentcloud_placement_group` with valid range 2-30
- Pass `Strategy` to `CreateDisasterRecoverGroup` API when specified by user
- Pass `PartitionCount` to `CreateDisasterRecoverGroup` API when `strategy` is `PARTITION` and user specifies `partition_count`
- Validate that `partition_count` can only be set when `strategy` is `PARTITION`
- Implement immutable args check using a shared array (`CVM_PLACEMENT_GROUP_IMMUTABLE_ARGS`) in the Update function (consistent with ckafka resource pattern), including `type`, `strategy`, `affinity`, and `partition_count`
- Read `Strategy` and `PartitionCount` from `DescribeDisasterRecoverGroups` API response to support state refresh and import
- `CreatePlacementGroup` service function returns full response for extensibility
- Maintain full backward compatibility — existing configurations continue to work unchanged

**Non-Goals:**
- Making `Strategy` or `PartitionCount` updatable (API does not support it; must be immutable)
- Adding `strategy` or `partition_count` to the datasource `tencentcloud_placement_groups` (out of scope)

## Decisions

### Decision 1: `strategy` and `partition_count` are immutable (not ForceNew)

**Rationale:** The `ModifyDisasterRecoverGroupAttribute` API does not accept `Strategy` or `PartitionCount`. Instead of using `ForceNew: true` (which silently destroys and recreates the resource), we use an `immutableArgs` array pattern in the Update function that returns a clear error when these values change. This gives users a better error message than silently destroying resources.

### Decision 2: Read `Strategy` and `PartitionCount` from Describe API

**Rationale:** After SDK update, the `DisasterRecoverGroup` struct now includes `Strategy` and `PartitionCount` fields. Reading these values enables proper state refresh and supports imported resources.

### Decision 3: Use `int` type with validation range 2-30 for `partition_count`

**Rationale:** Consistent with existing `affinity` field (also `TypeInt` with `ValidateIntegerInRange`). The API documentation states the valid range is 2-30.

### Decision 4: Validate `partition_count` requires `strategy=PARTITION`

**Rationale:** The API documentation states `PartitionCount` is only valid when the placement group type is partition (`Strategy=PARTITION`). Adding provider-side validation gives users clearer error messages than relying solely on API error responses.

### Decision 5: Return full response from `CreatePlacementGroup`

**Rationale:** Instead of returning individual fields from the response (`placementId string, respPartitionCount int`), return the complete `*cvm.CreateDisasterRecoverGroupResponse`. This provides better extensibility — callers can access any response field without modifying the service function signature.

### Decision 6: Immutable args array pattern

**Rationale:** Define `CVM_PLACEMENT_GROUP_IMMUTABLE_ARGS` in `extension_cvm.go` as a package-level variable containing `["type", "strategy", "affinity", "partition_count"]`. The Update function iterates this array and returns an error for any changed field. This pattern is consistent with the ckafka resource and is easily extensible.

## Risks / Trade-offs

- **[Risk] Imported resources have `strategy` and `partition_count` populated**: After SDK update, `DescribeDisasterRecoverGroups` returns these fields, so imported resources will have correct values.
  - **Mitigation:** No special handling needed; values are read from API response.

- **[Risk] Changing strategy or partition_count destroys the resource**: Using immutable args pattern (not ForceNew)
  - **Mitigation:** Users receive a clear error message explaining which field cannot be changed, rather than silent destruction.