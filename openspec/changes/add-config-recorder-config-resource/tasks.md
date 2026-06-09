## 1. Service Layer

- [x] 1.1 Append `DescribeConfigRecorder()` to `service_tencentcloud_config.go` — wraps `DescribeConfigRecorder` (no params)

## 2. Resource Implementation

- [x] 2.1 Create `resource_tc_config_recorder_config.go` with schema
- [x] 2.2 Create handler: `d.SetId(helper.BuildToken())` then call Update handler
- [x] 2.3 Read handler: call `DescribeConfigRecorder`, map Status→bool, map Items→resource_types list
- [x] 2.4 Update handler: if `status` changed → call Open/Close; if `resource_types` changed → call `UpdateConfigRecorder`
- [x] 2.5 Delete handler: no-op

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_config_recorder_config` in `provider.go` ResourcesMap

## 4. Documentation

- [x] 4.1 Create `resource_tc_config_recorder_config.md` with usage examples

## 5. Tests

- [x] 5.1 Create `resource_tc_config_recorder_config_test.go` with Create/Update/Import test

