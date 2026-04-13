## 1. SDK & Infrastructure

- [x] 1.1 Add Config SDK v20220802 to vendor directory
- [x] 1.2 Register `UseConfigV20220802Client()` in `tencentcloud/connectivity/client.go`

## 2. Service Layer

- [x] 2.1 Create `tencentcloud/services/cfg/service_tencentcloud_config.go` with `ConfigService` struct
- [x] 2.2 Implement `DescribeConfigCompliancePackById()` wrapping `DescribeCompliancePack` API
- [x] 2.3 Implement ratelimit and error handling in service methods

## 3. Resource Implementation

- [x] 3.1 Create `tencentcloud/services/cfg/resource_tc_config_compliance_pack.go`
- [x] 3.2 Define `ResourceTencentCloudConfigCompliancePack()` with full schema
- [x] 3.3 Implement `resourceTencentCloudConfigCompliancePackCreate()` calling `AddCompliancePack`
- [x] 3.4 Implement `resourceTencentCloudConfigCompliancePackRead()` calling service layer
- [x] 3.5 Implement `resourceTencentCloudConfigCompliancePackUpdate()` calling `UpdateCompliancePack` and/or `UpdateCompliancePackStatus`
- [x] 3.6 Implement `resourceTencentCloudConfigCompliancePackDelete()` — disable first, then delete

## 4. Provider Registration

- [x] 4.1 Register `tencentcloud_config_compliance_pack` in `tencentcloud/provider.go` ResourcesMap

## 5. Test Implementation

- [x] 5.1 Create `tencentcloud/services/cfg/resource_tc_config_compliance_pack_test.go`
- [x] 5.2 Implement `TestAccTencentCloudConfigCompliancePackResource_basic` test

## 6. Documentation

- [x] 6.1 Create `tencentcloud/services/cfg/resource_tc_config_compliance_pack.md`
- [x] 6.2 Add usage example with required fields
- [x] 6.3 Add import section
