## Why

The `tencentcloud_teo_security_policy_config` resource needs to support the full set of parameters for the `ModifySecurityPolicy` and `DescribeSecurityPolicy` cloud APIs, enabling users to manage TEO security policy configurations including zone-level, template-level, and host-level policies through Terraform.

## What Changes

- Add parameters to `tencentcloud_teo_security_policy_config` resource to support the `ModifySecurityPolicy` API inputs: `zone_id`, `entity`, `host`, `template_id`, `security_config`, and `security_policy`.
- Add parameters to support the `DescribeSecurityPolicy` API inputs (`zone_id`, `entity`, `host`, `template_id`) and output (`security_policy`).
- The resource uses `ModifySecurityPolicy` for create/update operations and `DescribeSecurityPolicy` for read operations.

## Capabilities

### New Capabilities
- `teo-security-policy-config-params`: Add full parameter support for the TEO security policy config resource, covering zone_id, entity, host, template_id, security_config, and security_policy fields mapped to the ModifySecurityPolicy and DescribeSecurityPolicy cloud APIs.

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go`: Resource schema and CRUD logic
- `tencentcloud/services/teo/resource_tc_teo_security_policy_config_test.go`: Unit tests
- `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md`: Documentation
- `tencentcloud/services/teo/service_tencentcloud_teo.go`: Service layer for API calls
- `tencentcloud/provider.go`: Resource registration
- Depends on `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` SDK package
