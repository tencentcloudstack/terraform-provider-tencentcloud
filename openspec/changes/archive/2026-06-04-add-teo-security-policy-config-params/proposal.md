## Why

The existing `tencentcloud_teo_security_policy_config` resource needs to support the `SecurityPolicy` parameter in the `ModifySecurityPolicy` API and properly read it back via the `DescribeSecurityPolicy` API. This enables users to configure security policies (custom rules, managed rules, HTTP DDoS protection, rate limiting rules, exception rules, bot management, bot management lite, and default deny security action parameters) using the expression-based `SecurityPolicy` structure for TEO (TencentCloud EdgeOne) security policy management.

## What Changes

- Add the `security_policy` parameter to the `tencentcloud_teo_security_policy_config` resource that maps to `request.SecurityPolicy` in the `ModifySecurityPolicy` API for create/update operations.
- Ensure the `DescribeSecurityPolicy` API is called with proper input parameters (`ZoneId`, `Entity`, `Host`, `TemplateId`) and the `response.Response.SecurityPolicy` output is read back into the `security_policy` terraform attribute.

## Capabilities

### New Capabilities
- `teo-security-policy-param`: Add the `SecurityPolicy` parameter support to the `tencentcloud_teo_security_policy_config` resource, enabling configuration of security policies via the `ModifySecurityPolicy` API and reading them back via `DescribeSecurityPolicy` API.

### Modified Capabilities
<!-- No existing capabilities are being modified -->

## Impact

- **Code**: `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go` - Add/verify `security_policy` schema field and CRUD logic
- **Code**: `tencentcloud/services/teo/resource_tc_teo_security_policy_config_test.go` - Add unit tests
- **Code**: `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md` - Update documentation with example usage
- **API**: Uses `ModifySecurityPolicy` and `DescribeSecurityPolicy` from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
- **Dependencies**: No new dependencies required (vendor SDK already includes the TEO v20220901 package)
