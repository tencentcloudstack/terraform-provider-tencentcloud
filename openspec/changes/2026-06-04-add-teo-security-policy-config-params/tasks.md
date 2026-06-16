## 1. Resource Schema and CRUD Implementation

- [x] 1.1 Add/verify schema definitions for `zone_id`, `entity`, `host`, `template_id`, `security_config`, and `security_policy` parameters in `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go`
- [x] 1.2 Implement Create function using `ModifySecurityPolicy` API with all parameters and composite ID generation
- [x] 1.3 Implement Read function using `DescribeSecurityPolicy` API with `zone_id`, `entity`, `host`, `template_id` inputs and `security_policy` output mapping
- [x] 1.4 Implement Update function using `ModifySecurityPolicy` API with all parameters
- [x] 1.5 Implement Delete function (no-op or reset to defaults)

## 2. Service Layer

- [x] 2.1 Add/verify `DescribeTeoSecurityPolicyConfigById` function in `tencentcloud/services/teo/service_tencentcloud_teo.go` with retry logic using `tccommon.ReadRetryTimeout`

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_security_policy_config` resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md` with Example Usage for all entity types and Import section with composite ID format

## 5. Unit Tests

- [x] 5.1 Add unit tests in `tencentcloud/services/teo/resource_tc_teo_security_policy_config_test.go` using gomonkey to mock cloud API calls, covering Create/Read/Update/Delete operations
