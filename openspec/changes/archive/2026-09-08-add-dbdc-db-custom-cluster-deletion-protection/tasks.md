## 1. Resource Schema Changes

- [x] 1.1 Add `deletion_protection` schema field (TypeBool, Optional, description: whether to enable cluster deletion protection) to `ResourceTencentCloudDbdcDbCustomCluster()` in `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_cluster.go`

## 2. Create Function Changes

- [x] 2.1 Read `deletion_protection` from schema data using `d.GetOk("deletion_protection")` in `resourceTencentCloudDbdcDbCustomClusterCreate`
- [x] 2.2 Set `request.DeletionProtection = helper.Bool(v.(bool))` when the parameter is specified

## 3. Read Function Changes

- [x] 3.1 Add nil-check for `respData.DeletionProtection` and set `deletion_protection` in state via `_ = d.Set("deletion_protection", respData.DeletionProtection)` in `resourceTencentCloudDbdcDbCustomClusterRead`

## 4. Update Function Changes

- [x] 4.1 Add `if d.HasChange("deletion_protection")` block in `resourceTencentCloudDbdcDbCustomClusterUpdate`
- [x] 4.2 Create `ModifyDBCustomClusterAttributesRequest` with `ClusterId` and `DeletionProtection` fields
- [x] 4.3 Call `ModifyDBCustomClusterAttributesWithContext` wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with proper error handling and logging
- [x] 4.4 Check response is not nil, log errors on failure

## 5. Unit Test Changes

- [x] 5.1 Add unit test cases for the `deletion_protection` parameter in `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_cluster_test.go` using gomonkey mock approach for the cloud API

## 6. Documentation

- [x] 6.1 Update `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_cluster.md` with `deletion_protection` usage example

## 7. Validation

- [x] 7.1 Verify the code compiles successfully
- [x] 7.2 Verify no lint errors
