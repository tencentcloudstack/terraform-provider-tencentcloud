# Add tencentcloud_teo_security_js_injection_rule Resource

## What

Add a new Terraform resource `tencentcloud_teo_security_js_injection_rule` for managing Tencent Cloud EdgeOne (TEO) JavaScript injection rules. This resource manages a single JS injection rule entry under a site zone, supporting full CRUD lifecycle.

## Why

TEO JavaScript injection rules allow users to configure per-request JS injection strategies (e.g., SDK injection for TC-RCE / TC-CAPTCHA authentication) based on match conditions. Currently no Terraform resource exists to manage these rules, requiring manual portal operations. This resource enables infrastructure-as-code management of TEO JS injection rules.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create | `CreateSecurityJSInjectionRule` | `JSInjectionRules` array length limited to **1**; returns `JSInjectionRuleIds` array |
| Read | `DescribeSecurityJSInjectionRule` | Paginate by zone; match by `RuleId` field in response |
| Update | `ModifySecurityJSInjectionRule` | Pass `RuleId` + updated fields in `JSInjectionRules[0]` |
| Delete | `DeleteSecurityJSInjectionRule` | Pass `ZoneId` + `JSInjectionRuleIds=[<id>]` |

## Resource ID

Composite: `<zone_id>#<js_injection_rule_id>` (using `tccommon.FILED_SP`), e.g. `zone-123123322#injection-2184008405`.
`js_injection_rule_id` is taken from `JSInjectionRuleIds[0]` returned by `CreateSecurityJSInjectionRule`.
