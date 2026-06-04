## 1. Resource Implementation

- [x] 1.1 Verify and ensure the `security_policy` schema field in `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go` properly maps all sub-fields of the `SecurityPolicy` struct from the vendor SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`)
- [x] 1.2 Verify the Update function correctly constructs `request.SecurityPolicy` from the `security_policy` schema field and passes it to `ModifySecurityPolicy` API
- [x] 1.3 Verify the Read function correctly calls `DescribeSecurityPolicy` with `ZoneId`, `Entity`, `Host`, `TemplateId` input parameters and flattens `response.Response.SecurityPolicy` into the `security_policy` attribute

## 2. Unit Tests

- [x] 2.1 Add unit tests in `tencentcloud/services/teo/resource_tc_teo_security_policy_config_test.go` using gomonkey to mock the cloud API calls, verifying the `security_policy` parameter is correctly handled in create/update/read operations

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md` with example usage demonstrating the `security_policy` parameter configuration
