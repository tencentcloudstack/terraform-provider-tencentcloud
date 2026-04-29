## Why

The `tencentcloud_teo_security_policy_config` resource currently lacks support for the `BotManagementLite` field in the `security_policy` block. The cloud API's `SecurityPolicy` structure supports this field through both `DescribeSecurityPolicy` and `ModifySecurityPolicy` interfaces, but the Terraform resource does not expose it. Users who need to configure basic Bot management features (CAPTCHA page challenge and AI crawler detection) cannot do so through Terraform.

## What Changes

- Add a new `bot_management_lite` parameter (TypeList, MaxItems: 1, Optional) under the `security_policy` block of the `tencentcloud_teo_security_policy_config` resource
- The `bot_management_lite` parameter will contain two sub-fields:
  - `captcha_page_challenge` (TypeList, MaxItems: 1, Optional): CAPTCHA page challenge configuration with an `enabled` field (TypeString, Required)
  - `ai_crawler_detection` (TypeList, MaxItems: 1, Optional): AI crawler detection configuration with an `enabled` field (TypeString, Required) and an `action` field (TypeList, MaxItems: 1, Optional) using the existing SecurityAction schema pattern
- Implement Read logic in `resourceTencentCloudTeoSecurityPolicyConfigRead` to flatten `BotManagementLite` from the `DescribeSecurityPolicy` response
- Implement Create/Update logic in `resourceTencentCloudTeoSecurityPolicyConfigCreate`/`Update` to expand `bot_management_lite` into the `ModifySecurityPolicy` request's `SecurityPolicy.BotManagementLite` field
- Add unit tests for the new parameter

## Capabilities

### New Capabilities
- `teo-security-policy-bot-management-lite`: Adds support for the BotManagementLite field in the tencentcloud_teo_security_policy_config resource, enabling Terraform users to configure basic Bot management features including CAPTCHA page challenge and AI crawler detection through the SecurityPolicy API.

### Modified Capabilities

## Impact

- **Affected files**: `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go`, `tencentcloud/services/teo/resource_tc_teo_security_policy_config_test.go`, `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md`
- **Cloud API**: Uses existing `DescribeSecurityPolicy` and `ModifySecurityPolicy` APIs (no new API dependencies)
- **Backward compatibility**: Adding a new Optional parameter is backward compatible; existing Terraform configurations will continue to work without changes
