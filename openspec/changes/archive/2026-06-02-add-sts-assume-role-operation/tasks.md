## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/sts/resource_tc_sts_assume_role_operation.go` with schema definition (all input fields ForceNew, output fields Computed/Sensitive) and CRUD functions (Create calls AssumeRole API with retry, Read/Delete are no-ops)
- [x] 1.2 Register `tencentcloud_sts_assume_role_operation` in `tencentcloud/provider.go` under the STS service section
- [x] 1.3 Add `tencentcloud_sts_assume_role_operation` entry to `tencentcloud/provider.md`

## 2. Documentation

- [x] 2.1 Create example documentation file `tencentcloud/services/sts/resource_tc_sts_assume_role_operation.md` with Example Usage and resource description

## 3. Testing

- [x] 3.1 Create unit test file `tencentcloud/services/sts/resource_tc_sts_assume_role_operation_test.go` using gomonkey to mock the STS client AssumeRoleWithContext method, testing the Create logic with required and optional parameters
- [x] 3.2 Run unit tests with `go test -gcflags=all=-l` to verify tests pass
