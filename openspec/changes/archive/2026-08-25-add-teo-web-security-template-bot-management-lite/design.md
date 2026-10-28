## Context

The `tencentcloud_teo_web_security_template` resource manages web security policy templates for TEO (Tencent EdgeOne). The existing resource already supports `bot_management` (Bot 管理) under the `security_policy` block, which provides full bot management with custom rules, basic bot settings, and client attestation. The TEO cloud API has introduced `BotManagementLite` (基础 Bot 管理) as a lighter-weight alternative, currently available in the SDK's `SecurityPolicy` struct.

The existing `bot_management` schema in this resource provides a well-established pattern for handling nested `SecurityAction` structures with their various parameters (DenyActionParameters, AllowActionParameters, etc.). The new `bot_management_lite` block follows the same patterns but with a simpler structure.

## Goals / Non-Goals

**Goals:**
- Add `bot_management_lite` block to the `security_policy` schema of `tencentcloud_teo_web_security_template`
- Support `captcha_page_challenge` (CAPTCHA page challenge) with `enabled` field
- Support `ai_crawler_detection` (AI crawler detection) with `enabled` and `action` (SecurityAction) fields
- Wire up Create, Read, and Update operations for the new parameters
- Maintain full backward compatibility

**Non-Goals:**
- Adding `bot_management_lite` to other resources (e.g., `tencentcloud_teo_security_policy_config`)
- Modifying existing `bot_management` schema or behavior
- Adding new cloud API interfaces

## Decisions

### Decision 1: Schema placement under `security_policy` alongside `bot_management`

**Rationale:** The SDK `SecurityPolicy` struct has `BotManagement` and `BotManagementLite` as sibling fields. Placing `bot_management_lite` as a sibling of `bot_management` in the Terraform schema mirrors the API structure exactly, making the mapping straightforward and intuitive.

**Alternatives considered:**
- Placing it inside `bot_management` — rejected because the API treats them as separate fields
- Creating a separate top-level resource — rejected because it's a sub-configuration of the security policy

### Decision 2: Reuse existing SecurityAction sub-block patterns

**Rationale:** The existing `bot_management` blocks already define the complete SecurityAction sub-block structure (name, deny_action_parameters, allow_action_parameters, challenge_action_parameters, redirect_action_parameters, block_ip_action_parameters, return_custom_page_action_parameters). The `ai_crawler_detection.action` block uses the same SDK `SecurityAction` type and should follow the identical schema pattern.

### Decision 3: All fields Optional+Computed

**Rationale:** All fields in `bot_management_lite` are new and Optional. The `Computed` flag is set on the parent blocks (TypeList containers) to allow the API to return values that weren't explicitly configured, consistent with the existing `bot_management` pattern.

### Decision 4: Expand in Create/Update, flatten in Read

**Rationale:** Follow the established pattern in this resource:
- **Create/Update (expand):** Convert Terraform schema maps into SDK request structs (`BotManagementLite`, `CAPTCHAPageChallenge`, `AICrawlerDetection`, `SecurityAction`)
- **Read (flatten):** Convert SDK response structs back into Terraform state maps with nil checks at each level

## Risks / Trade-offs

- **Risk:** Schema size growth — the file is already ~1MB. Adding `bot_management_lite` with its nested SecurityAction sub-blocks will add significant lines.
  → **Mitigation:** Follow the exact same pattern as existing `bot_management` to keep code consistent and maintainable

- **Risk:** The `SecurityAction` sub-blocks are duplicated across multiple locations in the schema (custom_rules, bot_management, and now bot_management_lite).
  → **Mitigation:** This is an existing pattern in the codebase. A future refactor could extract shared schema components, but that's out of scope for this change.

## Implementation Plan

1. **Schema changes** in `resource_tc_teo_web_security_template.go`:
   - Add `bot_management_lite` block under `security_policy` (after `bot_management`)
   - Contains `captcha_page_challenge` (TypeList, MaxItems:1, Optional) with `enabled` (TypeString, Required)
   - Contains `ai_crawler_detection` (TypeList, MaxItems:1, Optional) with `enabled` (TypeString, Required) and `action` (TypeList, MaxItems:1, Optional)
   - The `action` block contains the standard SecurityAction fields: `name`, `deny_action_parameters`, `allow_action_parameters`, `challenge_action_parameters`, `redirect_action_parameters`, `block_ip_action_parameters`, `return_custom_page_action_parameters`

2. **Create/Update function changes:**
   - After processing `bot_management`, process `bot_management_lite` from the schema map
   - Expand `captcha_page_challenge` into `CAPTCHAPageChallenge` struct
   - Expand `ai_crawler_detection` into `AICrawlerDetection` struct, including the `Action` → `SecurityAction` expansion
   - Set `request.SecurityPolicy.BotManagementLite = &botManagementLite`

3. **Read function changes:**
   - After processing `bot_management` response data, process `BotManagementLite`
   - Check `respData.BotManagementLite != nil`
   - Flatten `CAPTCHAPageChallenge.Enabled`
   - Flatten `AICrawlerDetection.Enabled` and the full `Action` → `SecurityAction` chain
   - Set `securityPolicyMap["bot_management_lite"]`

4. **Test changes:**
   - Add test cases in `resource_tc_teo_web_security_template_test.go`

5. **Documentation:**
   - Update `resource_tc_teo_web_security_template.md` with new `bot_management_lite` example usage