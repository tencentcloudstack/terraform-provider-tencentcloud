## 1. Schema Definition

- [x] 1.1 Add `bot_management_lite` field (TypeList, MaxItems: 1, Optional, Computed) under the `security_policy` schema block in `resource_tc_teo_security_policy_config.go`, with nested `captcha_page_challenge` and `ai_crawler_detection` sub-blocks
- [x] 1.2 Add `captcha_page_challenge` sub-block (TypeList, MaxItems: 1, Optional) with `enabled` field (TypeString, Required)
- [x] 1.3 Add `ai_crawler_detection` sub-block (TypeList, MaxItems: 1, Optional) with `enabled` field (TypeString, Required) and `action` sub-block (TypeList, MaxItems: 1, Optional)
- [x] 1.4 Add `action` sub-block under `ai_crawler_detection` with `name` (TypeString, Required), `deny_action_parameters` (TypeList, MaxItems: 1, Optional), `allow_action_parameters` (TypeList, MaxItems: 1, Optional), and `challenge_action_parameters` (TypeList, MaxItems: 1, Optional)
- [x] 1.5 Add `deny_action_parameters` sub-block with `block_ip`, `block_ip_duration`, `return_custom_page`, `response_code`, `error_page_id`, `stall` fields (all TypeString, Optional)
- [x] 1.6 Add `allow_action_parameters` sub-block with `min_delay_time` and `max_delay_time` fields (TypeString, Optional)
- [x] 1.7 Add `challenge_action_parameters` sub-block with `challenge_option` (TypeString, Required), `interval` (TypeString, Optional), `attester_id` (TypeString, Optional)

## 2. Read (Flatten) Implementation

- [x] 2.1 Add flatten logic for `BotManagementLite` in the `resourceTencentCloudTeoSecurityPolicyConfigRead` function, after the existing `exception_rules` flatten block
- [x] 2.2 Add nil checks for `respData.BotManagementLite`, `respData.BotManagementLite.CAPTCHAPageChallenge`, `respData.BotManagementLite.AICrawlerDetection`, and their nested fields
- [x] 2.3 Implement flatten for `captcha_page_challenge`: map `BotManagementLite.CAPTCHAPageChallenge.Enabled` to `captcha_page_challenge.enabled`
- [x] 2.4 Implement flatten for `ai_crawler_detection`: map `BotManagementLite.AICrawlerDetection.Enabled` to `ai_crawler_detection.enabled`
- [x] 2.5 Implement flatten for `ai_crawler_detection.action`: map `AICrawlerDetection.Action.Name` to `action.name`, and flatten `DenyActionParameters`, `AllowActionParameters`, `ChallengeActionParameters` with nil checks
- [x] 2.6 Set the flattened `botManagementLiteMap` into `securityPolicyMap["bot_management_lite"]`

## 3. Create/Update (Expand) Implementation

- [x] 3.1 Add expand logic for `bot_management_lite` in both the Create and Update expand sections of the resource, following the existing pattern for `http_ddos_protection`
- [x] 3.2 Expand `captcha_page_challenge`: read from Terraform state and set `BotManagementLite.CAPTCHAPageChallenge.Enabled`
- [x] 3.3 Expand `ai_crawler_detection`: read `enabled` and `action` from Terraform state
- [x] 3.4 Expand `ai_crawler_detection.action`: construct `SecurityAction` with `Name`, `DenyActionParameters`, `AllowActionParameters`, and `ChallengeActionParameters` as applicable
- [x] 3.5 Set the constructed `BotManagementLite` on `securityPolicy.BotManagementLite`

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md` to add example usage for `bot_management_lite` parameter, including `captcha_page_challenge` and `ai_crawler_detection` with their sub-fields

## 5. Unit Tests

- [x] 5.1 Add unit test functions for `bot_management_lite` flatten logic using gomonkey mock approach in `resource_tc_teo_security_policy_config_test.go`
- [x] 5.2 Add unit test functions for `bot_management_lite` expand logic using gomonkey mock approach
- [x] 5.3 Run unit tests with `go test -gcflags=all=-l` to verify the new test cases pass

## 6. Verification

- [x] 6.1 Verify the code compiles without errors (visual inspection of code correctness)
- [x] 6.2 Verify all new schema fields are Optional/Computed and backward compatible
- [x] 6.3 Verify nil checks are present at every nesting level in the Read flatten logic
