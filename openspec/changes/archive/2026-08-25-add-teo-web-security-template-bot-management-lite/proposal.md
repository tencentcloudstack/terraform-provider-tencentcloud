## Why

The `tencentcloud_teo_web_security_template` resource currently lacks support for the Bot Management Lite (基础 Bot 管理) configuration within the security policy. The TEO cloud API has already added `BotManagementLite` to the `SecurityPolicy` struct in `CreateWebSecurityTemplate`, `ModifyWebSecurityTemplate`, and `DescribeWebSecurityTemplate` APIs. This change enables Terraform users to configure CAPTCHA page challenge and AI crawler detection settings in their web security templates.

## What Changes

- Add a new optional `bot_management_lite` block (TypeList, MaxItems: 1) under `security_policy` in the `tencentcloud_teo_web_security_template` resource schema
- The `bot_management_lite` block contains two sub-blocks:
  - `captcha_page_challenge` (TypeList, MaxItems: 1, Optional): CAPTCHA page challenge configuration with `enabled` field
  - `ai_crawler_detection` (TypeList, MaxItems: 1, Optional): AI crawler detection configuration with `enabled` and `action` fields
- The `action` sub-block under `ai_crawler_detection` follows the same `SecurityAction` pattern used by existing `bot_management` blocks, supporting:
  - `name` (Required): Deny, Monitor, Allow, Challenge
  - `deny_action_parameters`, `allow_action_parameters`, `challenge_action_parameters`, `redirect_action_parameters`, `block_ip_action_parameters`, `return_custom_page_action_parameters`
- Support Create, Read, Update operations for the new parameters through the respective cloud API calls

## Capabilities

### New Capabilities
- `teo-web-security-template-bot-management-lite`: Bot Management Lite (CAPTCHA page challenge and AI crawler detection) configuration in the `tencentcloud_teo_web_security_template` resource

### Modified Capabilities
- `teo-security-policy-bot-management-lite`: The same capability is being extended to the `tencentcloud_teo_web_security_template` resource (in addition to the existing `tencentcloud_teo_security_policy_config` resource). The schema definition and API interaction patterns remain the same.

## Impact

- Affected code: `tencentcloud/services/teo/resource_tc_teo_web_security_template.go`
- Affected tests: `tencentcloud/services/teo/resource_tc_teo_web_security_template_test.go`
- Affected docs: `tencentcloud/services/teo/resource_tc_teo_web_security_template.md`
- No breaking changes: all new fields are Optional, existing configurations continue to work
- SDK dependency: Uses existing `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` types (`BotManagementLite`, `CAPTCHAPageChallenge`, `AICrawlerDetection`, `SecurityAction`)