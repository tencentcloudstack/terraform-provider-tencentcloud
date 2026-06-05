## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/teo/resource_tc_teo_check_free_certificate_verification_operation.go` with schema definition (zone_id, domain as Required+ForceNew; common_name, signature_algorithm, expire_time as Computed) and CRUD functions (Create calls CheckFreeCertificateVerification API with retry, Read/Update/Delete are no-op)
- [x] 1.2 Register the resource `tencentcloud_teo_check_free_certificate_verification` in `tencentcloud/provider.go`

## 2. Documentation

- [x] 2.1 Create resource example file `tencentcloud/services/teo/resource_tc_teo_check_free_certificate_verification_operation.md` with Example Usage and description
- [x] 2.2 Add resource entry to `tencentcloud/provider.md`

## 3. Testing

- [x] 3.1 Create unit test file `tencentcloud/services/teo/resource_tc_teo_check_free_certificate_verification_operation_test.go` using gomonkey to mock the CheckFreeCertificateVerification API call, and verify with `go test -gcflags=all=-l`
