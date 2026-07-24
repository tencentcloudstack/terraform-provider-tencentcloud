## 1. Resource Read Function Fix

- [x] 1.1 In `resourceTencentCloudGa2ForwardingPolicyRead`, after the `DescribeGa2ForwardingPolicyById` call, add a `ResourceNotFound` error check using `sdkErrors.TencentCloudSDKError` type assertion. When the error code is `"ResourceNotFound"` and the resource is not new (`!d.IsNewResource()`), log a WARN message, call `d.SetId("")`, and return nil.

## 2. Unit Test

- [ ] 2.1 Add unit test cases in `resource_tc_ga2_forwarding_policy_test.go` to verify that the Read function handles `ResourceNotFound` errors correctly (using gomonkey to mock the service layer).

## 3. Verification

- [ ] 3.1 Run `go vet` on the modified file to ensure no compilation errors.
- [ ] 3.2 Run `go test -gcflags=all=-l` on the unit test file to verify the mock-based tests pass.