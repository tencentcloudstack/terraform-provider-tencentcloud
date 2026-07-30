## 1. Schema and Resource Code

- [x] 1.1 Add `init_params` schema field (TypeList, Optional, ForceNew) with nested `param` (string, Required) and `value` (string, Required) to `ResourceTencentCloudMariadbHourDbInstance()` in `tencentcloud/services/mariadb/resource_tc_mariadb_hour_db_instance.go`
- [x] 1.2 Modify `resourceTencentCloudMariadbHourDbInstanceCreate` to read `init_params` from config: if provided, set `request.InitParams` on the `CreateHourDBInstance` request and skip the hardcoded `InitDBInstances` call; if not provided, keep existing hardcoded `InitDBInstances` behavior

## 2. Tests

- [x] 2.1 Add unit test cases in `tencentcloud/services/mariadb/resource_tc_mariadb_hour_db_instance_test.go` using gomonkey to mock the cloud API, verifying that `init_params` is correctly passed to the `CreateHourDBInstance` request

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/mariadb/resource_tc_mariadb_hour_db_instance.md` to add `init_params` usage example in the Example Usage section
