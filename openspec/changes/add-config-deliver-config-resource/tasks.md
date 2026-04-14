## 1. Service Layer

- [x] 1.1 Append `DescribeConfigDeliver()` to `service_tencentcloud_config.go` — calls `DescribeConfigDeliver` (no params), returns `*DescribeConfigDeliverResponseParams`

## 2. Resource Implementation

- [x] 2.1 Create `resource_tc_config_deliver_config.go` with `ResourceTencentCloudConfigDeliverConfig()` schema
- [x] 2.2 Implement Create handler: call `UpdateConfigDeliver`; set ID to `helper.BuildToken()`
- [x] 2.3 Implement Read handler: call service, populate all fields
- [x] 2.4 Implement Update handler: call `UpdateConfigDeliver` on change
- [x] 2.5 Implement Delete handler: no-op (return nil)

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_config_deliver_config` in `provider.go` ResourcesMap

## 4. Documentation

- [x] 4.1 Create `resource_tc_config_deliver_config.md` with usage example

## 5. Tests

- [x] 5.1 Create `resource_tc_config_deliver_config_test.go` with basic CRUD acceptance test

