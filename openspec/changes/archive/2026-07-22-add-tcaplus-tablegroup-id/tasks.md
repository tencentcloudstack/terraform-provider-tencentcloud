## 1. Service Layer Changes

- [x] 1.1 Update `CreateGroup` function signature in `tencentcloud/services/tcaplusdb/service_tencentcloud_tcaplus.go` (line 277) to accept an additional `tableGroupId string` parameter
- [x] 1.2 Pass `tableGroupId` to the `CreateTableGroup` API request as `request.TableGroupId` when the value is non-empty (leave unset when empty so the API auto-increments)
- [x] 1.3 Keep returning the created `groupId` from `response.Response.TableGroupId` (existing behavior, no change to return value)

## 2. Resource Schema Changes

- [x] 2.1 Add `table_group_id` schema field (TypeString, Optional, immutable after creation, description explaining it is user-specified or auto-incremented) to `ResourceTencentCloudTcaplusTableGroup()` in `tencentcloud/services/tcaplusdb/resource_tc_tcaplus_tablegroup.go`

## 3. Create Function Changes

- [x] 3.1 Read `table_group_id` from schema data in `resourceTencentCloudTcaplusTableGroupCreate`
- [x] 3.2 Pass `table_group_id` to the updated `CreateGroup` service function
- [x] 3.3 Preserve existing resource ID format (`clusterId:tableGroupId`) via `d.SetId()`

## 4. Read Function Changes

- [x] 4.1 Set `table_group_id` into state from `DescribeTableGroups` API response (`info.TableGroupId`) in `resourceTencentCloudTcaplusTableGroupRead`, with a nil check before calling `d.Set()`

## 5. Update Function Changes

- [x] 5.1 Add `immutableArgs` array containing `table_group_id` in `resourceTencentCloudTcaplusTableGroupUpdate`
- [x] 5.2 Iterate `immutableArgs` and return a formatted error if `table_group_id` has changed (reject changes with a clear message instead of ForceNew)
- [x] 5.3 Keep existing `tablegroup_name` change handling unchanged

## 6. Documentation

- [x] 6.1 Update `tencentcloud/services/tcaplusdb/resource_tc_tcaplus_tablegroup.md` with `table_group_id` usage example
- [x] 6.2 Run `make doc` during finalize phase to regenerate `website/docs/r/tcaplus_tablegroup.html.markdown`

## 7. Tests

- [x] 7.1 Update `tencentcloud/services/tcaplusdb/resource_tc_tcaplus_tablegroup_test.go` to add test cases covering `table_group_id` (create with user-specified id, and immutable-update error scenario), using the existing terraform test suite pattern

## 8. Validation

- [x] 8.1 Verify the code compiles successfully
- [x] 8.2 Run gofmt during the finalize phase
