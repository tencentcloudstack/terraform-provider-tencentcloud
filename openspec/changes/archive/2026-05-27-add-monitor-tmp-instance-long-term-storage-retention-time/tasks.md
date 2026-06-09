## 1. Schema and CRUD Implementation

- [x] 1.1 Add `long_term_storage_retention_time` parameter (Optional, Computed, int) to the resource schema in `tencentcloud/services/tmp/resource_tc_monitor_tmp_instance.go`
- [x] 1.2 Wire `long_term_storage_retention_time` to `request.InstanceAttributes` with key `LongTermStorageRetentionTime` in the Create function
- [x] 1.3 Wire `long_term_storage_retention_time` to `request.InstanceAttributes` with key `LongTermStorageRetentionTime` in the Update function (when changed)
- [x] 1.4 Read `LongTermStorageRetentionTime` from `PrometheusInstancesItem.InstanceAttributes` in the Read function, parsing the string value to int

## 2. Unit Tests

- [x] 2.1 Update unit test cases in `tencentcloud/services/tmp/resource_tc_monitor_tmp_instance_test.go` to verify Create with `long_term_storage_retention_time`, Read populating the field from `InstanceAttributes`, Read with nil `InstanceAttributes`, and Update modifying the value

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/tmp/resource_tc_monitor_tmp_instance.md` to include `long_term_storage_retention_time` in the parameter documentation
