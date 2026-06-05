## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/teo/resource_tc_teo_apply_free_certificate_operation.go` with schema definition (zone_id, domain, verification_method as Required/ForceNew inputs; dns_verification, file_verification as Computed outputs) and CRUD functions (Create calls ApplyFreeCertificate API with retry, Read/Delete are no-ops)
- [x] 1.2 Create unit test file `tencentcloud/services/teo/resource_tc_teo_apply_free_certificate_operation_test.go` using gomonkey to mock the ApplyFreeCertificate API call and verify business logic

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_apply_free_certificate` resource in `tencentcloud/provider.go` ResourcesMap
- [x] 2.2 Add `tencentcloud_teo_apply_free_certificate` entry in `tencentcloud/provider.md` under the teo service section

## 3. Documentation

- [x] 3.1 Create example usage documentation file `tencentcloud/services/teo/resource_tc_teo_apply_free_certificate_operation.md` with Example Usage section showing both DNS and HTTP verification methods
