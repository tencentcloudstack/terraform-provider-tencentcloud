## 1. Service Layer Changes

- [x] 1.1 Update `CreatePlacementGroup` function signature in `tencentcloud/services/cvm/service_tencentcloud_cvm.go` to accept `partitionCount int` parameter
- [x] 1.2 Pass `PartitionCount` to `CreateDisasterRecoverGroup` API request when `partitionCount > 0`

## 2. Resource Schema Changes

- [x] 2.1 Add `partition_count` schema field (TypeInt, Optional, ForceNew, Computed, validate 2-30) to `ResourceTencentCloudPlacementGroup()` in `tencentcloud/services/cvm/resource_tc_placement_group.go`
- [x] 2.2 Update `resourceTencentCloudPlacementGroupCreate` to read `partition_count` from schema and pass to `CreatePlacementGroup` service call

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/cvm/resource_tc_placement_group.md` with `partition_count` usage example

## 4. Validation

- [x] 4.1 Verify the code compiles successfully (`go vet` equivalent check)
- [x] 4.2 Run `gofmt` on all modified `.go` files (deferred to tfpacer-finalize skill)