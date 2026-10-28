## 1. Schema Definition

- [x] 1.1 Add `bot_management_lite` block (TypeList, MaxItems: 1, Optional, Computed) under `security_policy` in the resource schema, after the existing `bot_management` block
- [x] 1.2 Add `captcha_page_challenge` sub-block (TypeList, MaxItems: 1, Optional, Computed) with `enabled` field (TypeString, Required, "on"/"off")
- [x] 1.3 Add `ai_crawler_detection` sub-block (TypeList, MaxItems: 1, Optional, Computed) with `enabled` field (TypeString, Required, "on"/"off") and `action` field (TypeList, MaxItems: 1, Optional)
- [x] 1.4 Add `action` sub-fields: `name` (Required), `deny_action_parameters`, `allow_action_parameters`, `challenge_action_parameters`, `redirect_action_parameters`, `block_ip_action_parameters`, `return_custom_page_action_parameters` (all Optional, TypeList, MaxItems: 1), following the same pattern as existing `bot_management` action blocks

## 2. Create Operation (Expand)

- [x] 2.1 In the Create function, after processing `bot_management`, process `bot_management_lite` map from the Terraform config
- [x] 2.2 Expand `captcha_page_challenge` into `CAPTCHAPageChallenge` struct with `Enabled` field
- [x] 2.3 Expand `ai_crawler_detection` into `AICrawlerDetection` struct with `Enabled` and `Action` (SecurityAction) fields
- [x] 2.4 Expand the `action` sub-block into `SecurityAction` struct with all nested parameters (DenyActionParameters, AllowActionParameters, ChallengeActionParameters, RedirectActionParameters, BlockIPActionParameters, ReturnCustomPageActionParameters)
- [x] 2.5 Set `request.SecurityPolicy.BotManagementLite = &botManagementLite`

## 3. Update Operation (Expand)

- [x] 3.1 In the Update function, apply the same `bot_management_lite` expansion logic as in Create
- [x] 3.2 Set `request.SecurityPolicy.BotManagementLite = &botManagementLite` in the ModifyWebSecurityTemplate request

## 4. Read Operation (Flatten)

- [x] 4.1 In the Read function, after processing `bot_management` response data, process `BotManagementLite` from the response
- [x] 4.2 Check `respData.BotManagementLite != nil` before accessing sub-fields
- [x] 4.3 Flatten `CAPTCHAPageChallenge.Enabled` into the `captcha_page_challenge` map
- [x] 4.4 Flatten `AICrawlerDetection.Enabled` and the full `Action` → `SecurityAction` chain into the `ai_crawler_detection` map
- [x] 4.5 Set `securityPolicyMap["bot_management_lite"]` with the flattened map

## 5. Unit Tests

- [x] 5.1 Add test cases for `bot_management_lite` in `resource_tc_teo_web_security_template_test.go` covering basic create/read/update scenarios with CAPTCHA page challenge and AI crawler detection

## 6. Documentation

- [x] 6.1 Update `tencentcloud/services/teo/resource_tc_teo_web_security_template.md` with `bot_management_lite` example usage in the Example Usage section