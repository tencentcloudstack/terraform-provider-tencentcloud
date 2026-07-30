## Why

Users need the ability to assume IAM roles via Terraform to obtain temporary security credentials (STS AssumeRole). This enables cross-account access, service role assumption, and temporary credential management directly within Terraform workflows, without requiring external scripts or manual API calls.

## What Changes

- Add a new Terraform resource `tencentcloud_sts_assume_role_operation` of type RESOURCE_KIND_OPERATION
- The resource calls the STS `AssumeRole` API to obtain temporary credentials (Token, TmpSecretId, TmpSecretKey)
- All input parameters are ForceNew (one-time operation, no state tracking)
- Read and Delete methods are empty (operation resource pattern)
- Output parameters (credentials, expired_time, expiration) are Computed and stored in state for downstream consumption
- Register the resource in `provider.go` and `provider.md`

## Capabilities

### New Capabilities
- `sts-assume-role-operation`: A one-time operation resource that calls the STS AssumeRole API to obtain temporary security credentials for role assumption

### Modified Capabilities

## Impact

- New files:
  - `tencentcloud/services/sts/resource_tc_sts_assume_role_operation.go` - Resource implementation
  - `tencentcloud/services/sts/resource_tc_sts_assume_role_operation_test.go` - Unit tests with gomonkey mocks
  - `tencentcloud/services/sts/resource_tc_sts_assume_role_operation.md` - Example documentation
- Modified files:
  - `tencentcloud/provider.go` - Register the new resource
  - `tencentcloud/provider.md` - Add resource to provider documentation
- Dependencies: Uses existing `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813` package (already vendored)
- No breaking changes to existing resources or APIs
