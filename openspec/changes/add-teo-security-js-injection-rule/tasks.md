## 1. Service Layer

- [x] 1.1 Append `DescribeTeoSecurityJSInjectionRuleById(ctx, zoneId, jsInjectionRuleId)` to `tencentcloud/services/teo/service_tencentcloud_teo.go` — paginates `DescribeSecurityJSInjectionRule` (Limit=100) until the entry with `RuleId == jsInjectionRuleId` is found; returns `*teo.JSInjectionRule` or nil

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_security_js_injection_rule.go` with full schema following `tencentcloud_igtm_strategy` style:
  - Top-level fields: `zone_id` (Required, ForceNew)
  - `js_injection_rules` (Required, TypeList, MaxItems:1) with sub-fields 100% mapping SDK `JSInjectionRule` struct:
    - `name` (Required, String)
    - `priority` (Required, Int)
    - `condition` (Required, String)
    - `inject_j_s` (Required, String)
    - `rule_id` (Computed, String)

- [x] 2.2 Implement Create:
  - Build `CreateSecurityJSInjectionRuleRequest` with `ZoneId` and `JSInjectionRules` (single-element, no `RuleId`)
  - Call `CreateSecurityJSInjectionRuleWithContext`; extract `JSInjectionRuleIds[0]`
  - Set resource ID to `strings.Join([]string{zoneId, jsInjectionRuleId}, tccommon.FILED_SP)`
  - Call Read

- [x] 2.3 Implement Read:
  - Split ID → `zoneId`, `jsInjectionRuleId`
  - Call `DescribeTeoSecurityJSInjectionRuleById`; if nil → `d.SetId("")`
  - Populate `zone_id` and `js_injection_rules` block from response

- [x] 2.4 Implement Update:
  - Build `ModifySecurityJSInjectionRuleRequest` with `ZoneId` and `JSInjectionRules=[{RuleId, Name, Priority, Condition, InjectJS}]`
  - Call `ModifySecurityJSInjectionRuleWithContext`
  - Call Read

- [x] 2.5 Implement Delete:
  - Build `DeleteSecurityJSInjectionRuleRequest` with `ZoneId` and `JSInjectionRuleIds=[jsInjectionRuleId]`
  - Call `DeleteSecurityJSInjectionRuleWithContext` with Retry

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_security_js_injection_rule` in `tencentcloud/provider.go` ResourcesMap, pointing to `teo.ResourceTencentCloudTeoSecurityJSInjectionRule()`

## 4. Documentation & Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_security_js_injection_rule.md` — document all arguments, attributes, and import syntax with example HCL
- [x] 4.2 Create `tencentcloud/services/teo/resource_tc_teo_security_js_injection_rule_test.go` — basic acceptance test covering create/update/import/delete following `resource_tc_igtm_strategy_test.go` style
