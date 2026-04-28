## 1. Resource Schema & CRUD Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_security_client_attester.go` with schema definition including `zone_id` (Required, TypeString), `client_attesters` (Required, TypeList with nested Resource), `client_attester_ids` (Computed, TypeList), and nested option blocks (`tc_rce_option`, `tc_captcha_option`, `tc_eo_captcha_option`)
- [x] 1.2 Implement `resourceTencentCloudTeoSecurityClientAttesterCreate` - call `CreateSecurityClientAttester` with ZoneId and ClientAttesters, store composite ID as `zone_id#client_attester_ids` using FILED_SP separator
- [x] 1.3 Implement `resourceTencentCloudTeoSecurityClientAttesterRead` - call `DescribeSecurityClientAttester` with ZoneId, handle pagination with Limit=100, match by client_attester_ids from composite ID, flatten response into Terraform state
- [x] 1.4 Implement `resourceTencentCloudTeoSecurityClientAttesterUpdate` - check immutable args (`zone_id`), call `ModifySecurityClientAttester` with ZoneId and updated ClientAttesters
- [x] 1.5 Implement `resourceTencentCloudTeoSecurityClientAttesterDelete` - call `DeleteSecurityClientAttester` with ZoneId and ClientAttesterIds extracted from d.Get() fields

## 2. Service Layer

- [x] 2.1 Add `DescribeTeoSecurityClientAttesterById` method to `TeoService` in `tencentcloud/services/teo/service_tencentcloud_teo.go` with pagination support (Limit=100) and retry logic using `tccommon.ReadRetryTimeout`

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_security_client_attester` resource in `tencentcloud/provider.go` ResourcesMap
- [x] 3.2 Add `tencentcloud_teo_security_client_attester` entry in `tencentcloud/provider.md`

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_security_client_attester_test.go` with gomonkey-based unit tests for Create, Read, Update, Delete functions

## 5. Documentation

- [x] 5.1 Create `tencentcloud/services/teo/resource_tc_teo_security_client_attester.md` with resource description, example usage (covering TC-RCE, TC-CAPTCHA, TC-EO-CAPTCHA scenarios), and import section
