## 1. Schema Changes

- [x] 1.1 Add `disaster_recover_group_ids` schema field (TypeList of TypeString, Optional, ForceNew, MaxItems 1) to `ResourceTencentCloudDbdcDbCustomNode()` in `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node.go`, with a description noting it maps to `CreateDBCustomNodes` `DisasterRecoverGroupIds` and that only one ID is supported
- [x] 1.2 Add `disaster_recover_group_id` schema field (TypeString, Computed) to `ResourceTencentCloudDbdcDbCustomNode()`, with a description noting it is refreshed from `DescribeDBCustomNodes` `DBCustomNode.DisasterRecoverGroupId`

## 2. Create Function Changes

- [x] 2.1 In `resourceTencentCloudDbdcDbCustomNodeCreate`, read `disaster_recover_group_ids` from schema data and populate `request.DisasterRecoverGroupIds` (append each `*string` element), placing the logic alongside the other create-time request fields and before the `resource.Retry` call

## 3. Read Function Changes

- [x] 3.1 In `resourceTencentCloudDbdcDbCustomNodeRead`, add a nil-guarded `d.Set("disaster_recover_group_id", respData.DisasterRecoverGroupId)` block alongside the existing computed field reads (after `respData` is fetched via `DescribeDBCustomNodeById`)

## 4. Tests

- [x] 4.1 Add unit test coverage in `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node_test.go` for the new `disaster_recover_group_ids` input (Create request mapping) and `disaster_recover_group_id` computed output (Read), using gomonkey mocks for the cloud API per the project convention for newly added resources

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node.md` example to include `disaster_recover_group_ids` usage (website docs are generated via `make doc` in the finalize phase; update the source `.md` only)

## 6. Validation

- [x] 6.1 Verify the code compiles successfully
- [x] 6.2 Verify no lint errors
