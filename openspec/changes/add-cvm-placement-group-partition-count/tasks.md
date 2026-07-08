## 1. Service Layer Changes

- [x] 1.1 Update `CreatePlacementGroup` function signature in `tencentcloud/services/cvm/service_tencentcloud_cvm.go` to accept `strategy string` and `partitionCount int` parameters
- [x] 1.2 Pass `Strategy` to `CreateDisasterRecoverGroup` API request when `strategy != ""`
- [x] 1.3 Pass `PartitionCount` to `CreateDisasterRecoverGroup` API request when `partitionCount > 0`
- [x] 1.4 Refactor `CreatePlacementGroup` to return full `*cvm.CreateDisasterRecoverGroupResponse` instead of individual fields

## 2. Resource Schema Changes

- [x] 2.1 Add `strategy` schema field (TypeString, Optional, Computed, validate SPREAD/PARTITION) to `ResourceTencentCloudPlacementGroup()` in `tencentcloud/services/cvm/resource_tc_placement_group.go`
- [x] 2.2 Add `partition_count` schema field (TypeInt, Optional, Computed, validate 2-30, description mentions strategy dependency) to `ResourceTencentCloudPlacementGroup()`
- [x] 2.3 Add `CVM_PLACEMENT_GROUP_STRATEGY` constants and var array in `tencentcloud/services/cvm/extension_cvm.go`
- [x] 2.4 Add `CVM_PLACEMENT_GROUP_IMMUTABLE_ARGS` variable in `extension_cvm.go` containing `["type", "strategy", "affinity", "partition_count"]`

## 3. Create Function Changes

- [x] 3.1 Read `strategy` and `partition_count` from schema data
- [x] 3.2 Validate `partition_count` only set when `strategy == "PARTITION"`, return error otherwise
- [x] 3.3 Pass values to `CreatePlacementGroup` service function
- [x] 3.4 Extract `DisasterRecoverGroupId` from full response and call `d.SetId()`

## 4. Read Function Changes

- [x] 4.1 Read `placement.Strategy` from Describe response and set to `strategy` in state
- [x] 4.2 Read `placement.PartitionCount` from Describe response and set to `partition_count` in state

## 5. Update Function Changes

- [x] 5.1 Add `immutableArgs` array with `["type", "strategy", "affinity", "partition_count"]`
- [x] 5.2 Iterate `immutableArgs` and return error if any field has changed

## 6. SDK Update

- [x] 6.1 Update CVM SDK to version that includes `Strategy` and `PartitionCount` in `DisasterRecoverGroup` struct

## 7. Documentation

- [x] 7.1 Update `tencentcloud/services/cvm/resource_tc_placement_group.md` with `partition_count` and `strategy` usage example

## 8. Validation

- [x] 8.1 Verify the code compiles successfully
- [x] 8.2 Verify no lint errors