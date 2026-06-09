## 1. Service Layer

- [x] 1.1 Append `DescribeConfigRemediationById()` to `service_tencentcloud_config.go` — calls `ListRemediations` with `RuleIds` filter and finds by `RemediationId`

## 2. Resource Implementation

- [x] 2.1 Create `resource_tc_config_remediation.go` with `ResourceTencentCloudConfigRemediation()` schema
- [x] 2.2 Implement `resourceTencentCloudConfigRemediationCreate()` calling `CreateRemediation`, set ID to `RemediationId`
- [x] 2.3 Implement `resourceTencentCloudConfigRemediationRead()` calling service layer
- [x] 2.4 Implement `resourceTencentCloudConfigRemediationUpdate()` calling `UpdateRemediation` on changed fields
- [x] 2.5 Implement `resourceTencentCloudConfigRemediationDelete()` calling `DeleteRemediations`

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_config_remediation` in `provider.go` ResourcesMap

## 4. Documentation

- [x] 4.1 Create `resource_tc_config_remediation.md` with usage example and import section

## 5. Tests

- [x] 5.1 Create `resource_tc_config_remediation_test.go` with basic CRUD acceptance test

