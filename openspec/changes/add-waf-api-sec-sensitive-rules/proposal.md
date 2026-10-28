## Why

The WAF API security "sensitive data" feature is configured through a single multiplexed API (`ModifyApiSecSensitiveRule`) that accepts seven different rule structures in one request body. Today none of these rule types are manageable via Terraform. Exposing the whole multiplexed API as one resource would be confusing and error-prone, so we split it into seven focused, independent resources — one per rule structure — giving users a clean, declarative way to manage each rule type.

## What Changes

- Add **7 new CRUD resources** under the `waf` service, each mapping to exactly one input field of `ModifyApiSecSensitiveRule`:
  - `tencentcloud_waf_api_sec_sensitive_custom_rule` → `CustomRule` (`ApiSecCustomSensitiveRule`)
  - `tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule` → `CustomApiExtractRule` (`ApiSecExtractRule`)
  - `tencentcloud_waf_api_sec_sensitive_privilege_rule` → `ApiSecPrivilegeRule` (`ApiSecPrivilegeRule`)
  - `tencentcloud_waf_api_sec_sensitive_scene_rule` → `ApiSecSceneRule` (`ApiSecSceneRule`)
  - `tencentcloud_waf_api_sec_sensitive_custom_event_rule` → `ApiSecCustomEventRuleRule` (`ApiSecCustomEventRule`)
  - `tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule` → `CustomApiExcludeRule` (`ApiSecExcludeRule`)
  - `tencentcloud_waf_api_sec_sensitive_white_rule` → `ApiSecSensitiveWhiteRuleRule` (`ApiSecSensitiveWhiteRule`)
- Each resource's schema contains **only** the fields of its corresponding struct (strictly validated, no extra fields).
- Create / Read / Update / Delete are all backed by `ModifyApiSecSensitiveRule` (write) and `DescribeApiSecSensitiveRuleList` (read).
- Resource ID is the composite `Domain` + `tccommon.FILED_SP` (`#`) + `RuleName`.
- The `RuleName` carried inside each sub-struct is unified with the top-level `RuleName` parameter; the sub-struct's own `RuleName` is not exposed as a separate field.
- The `Status` field is exposed to users only as `0` (off) / `1` (on); the value `3` (delete) is used internally as the default `Status` in the Delete path.
- Add accompanying `.md` example docs, acceptance tests, service-layer query helpers, generated website docs, and provider registration.

## Capabilities

### New Capabilities
- `waf-api-sec-sensitive-custom-rule`: Manage the custom sensitive-data detection rule (`CustomRule`).
- `waf-api-sec-sensitive-custom-api-extract-rule`: Manage the API extraction rule (`CustomApiExtractRule`).
- `waf-api-sec-sensitive-privilege-rule`: Manage the API privilege (auth) rule (`ApiSecPrivilegeRule`).
- `waf-api-sec-sensitive-scene-rule`: Manage the custom API scene rule (`ApiSecSceneRule`).
- `waf-api-sec-sensitive-custom-event-rule`: Manage the custom event rule (`ApiSecCustomEventRuleRule`).
- `waf-api-sec-sensitive-custom-api-exclude-rule`: Manage the invalid-API exclude rule (`CustomApiExcludeRule`).
- `waf-api-sec-sensitive-white-rule`: Manage the sensitive-data allowlist (white) rule (`ApiSecSensitiveWhiteRuleRule`).

### Modified Capabilities
<!-- None: these are all new resources; no existing spec requirements change. -->

## Impact

- **New code**: 7 × (`resource_tc_waf_api_sec_sensitive_*.go`, `.md`, `_test.go`) under `tencentcloud/services/waf/`.
- **Service layer**: new query helpers in `service_tencentcloud_waf.go` wrapping `DescribeApiSecSensitiveRuleList`.
- **Provider registration**: 7 new entries in `tencentcloud/provider.go`.
- **Docs**: 7 new pages under `website/docs/r/` plus entries in `website/tencentcloud.erb`.
- **SDK**: uses existing `waf/v20180125` SDK (`ModifyApiSecSensitiveRule`, `DescribeApiSecSensitiveRuleList`); no SDK changes required.
- **Backward compatibility**: additive only — no existing resource/schema/state is modified.
