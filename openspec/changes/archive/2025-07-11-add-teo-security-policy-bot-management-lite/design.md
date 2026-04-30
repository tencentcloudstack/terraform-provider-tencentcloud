## Context

The `tencentcloud_teo_security_policy_config` resource manages TEO (TencentCloud EdgeOne) security policies. The resource currently supports 5 out of 8 fields in the cloud API's `SecurityPolicy` struct: `custom_rules`, `managed_rules`, `http_ddos_protection`, `rate_limiting_rules`, and `exception_rules`. The missing fields are `BotManagement`, `BotManagementLite`, and `DefaultDenySecurityActionParameters`.

This change adds the `bot_management_lite` field, which corresponds to the `BotManagementLite` struct in the cloud API. `BotManagementLite` is the lightweight version of Bot management, providing basic Bot protection capabilities including CAPTCHA page challenge and AI crawler detection.

The existing resource file is large (~3660 lines) and follows consistent patterns for schema definition, read (flatten), and create/update (expand) operations.

## Goals / Non-Goals

**Goals:**
- Add `bot_management_lite` as an Optional parameter under `security_policy` in the resource schema
- Support `captcha_page_challenge` (with `enabled` field) and `ai_crawler_detection` (with `enabled` and `action` fields) sub-configurations
- Implement Read (flatten) logic to populate `bot_management_lite` from `DescribeSecurityPolicy` response
- Implement Create/Update (expand) logic to set `SecurityPolicy.BotManagementLite` in `ModifySecurityPolicy` request
- Add unit tests using gomonkey mock approach
- Update the resource `.md` documentation

**Non-Goals:**
- Adding `BotManagement` (full Bot management) or `DefaultDenySecurityActionParameters` - these are separate changes
- Modifying any existing schema fields or behavior
- Changing the resource's ID format or import logic

## Decisions

### 1. Schema structure for `bot_management_lite`

**Decision**: Use TypeList with MaxItems:1 for both `bot_management_lite` and its nested blocks, consistent with existing fields like `managed_rules`, `http_ddos_protection`, and `exception_rules`.

**Rationale**: The existing resource consistently uses TypeList with MaxItems:1 for all single-object nested blocks under `security_policy`. This maintains consistency with the codebase patterns.

### 2. Schema structure for `ai_crawler_detection.action`

**Decision**: Include `deny_action_parameters`, `allow_action_parameters`, and `challenge_action_parameters` as Optional sub-blocks of the `action` field, but NOT `redirect_action_parameters`, `block_ip_action_parameters`, or `return_custom_page_action_parameters`.

**Rationale**: The cloud API documentation for `AICrawlerDetection.Action` states: "SecurityAction 的 Name 取值仅支持：Deny, Monitor, Allow, Challenge". Since the API only supports Deny, Monitor, Allow, and Challenge actions, only the corresponding action parameter blocks should be included. Monitor has no additional parameters.

### 3. AllowActionParameters schema

**Decision**: Include `min_delay_time` and `max_delay_time` fields in `allow_action_parameters`.

**Rationale**: The cloud API's `AllowActionParameters` struct has two fields: `MinDelayTime` and `MaxDelayTime`. These must be mapped to the Terraform schema.

### 4. DenyActionParameters schema for AICrawlerDetection

**Decision**: Reuse the same `DenyActionParameters` schema pattern as used in `http_ddos_protection` sub-blocks, with `block_ip`, `block_ip_duration`, `return_custom_page`, `response_code`, `error_page_id`, and `stall` fields.

**Rationale**: The `AICrawlerDetection.Action.DenyActionParameters` uses the same `DenyActionParameters` struct as other security actions in the API. Consistency with existing patterns ensures maintainability.

### 5. Unit testing approach

**Decision**: Use gomonkey mock approach for unit tests, as this is a new resource parameter addition.

**Rationale**: Per project rules, new terraform resources should use mock-based unit tests rather than terraform test suites.

## Risks / Trade-offs

- [Risk: API field naming mismatch] → The cloud API uses `CAPTCHAPageChallenge` with all-caps "CAPTCHA", which converts to `captcha_page_challenge` in Terraform schema. Need to ensure correct mapping in flatten/expand logic.
- [Risk: Large file modification] → The resource file is already 3660 lines. Adding new schema, flatten, and expand code will increase it further. Mitigation: follow existing patterns closely to minimize review overhead.
- [Risk: Nil pointer checks] → Must follow existing pattern of checking for nil before accessing nested fields in flatten logic, as cloud API responses may return nil for optional fields.
