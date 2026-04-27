# Design: tencentcloud_teo_security_js_injection_rule Resource

## Architecture

Follows `tencentcloud_igtm_strategy` style:

```
provider.go
    └─ tencentcloud/services/teo/resource_tc_teo_security_js_injection_rule.go  (CRUD handlers)
           └─ tencentcloud/services/teo/service_tencentcloud_teo.go (DescribeTeoSecurityJSInjectionRuleById)
                  └─ teo SDK v20220901
```

## Resource ID

Composite: `<zone_id>#<js_injection_rule_id>` (using `tccommon.FILED_SP`), e.g. `zone-123123322#injection-2184008405`.

## Key Constraint

`JSInjectionRules` array in Create/Modify is limited to **1 element**. The `js_injection_rules` field uses `TypeList` with `MaxItems: 1` to align with the SDK `JSInjectionRule` struct.

## Schema

### Top-level Required

| Field | Type | ForceNew | Description |
|---|---|---|---|
| `zone_id` | String | Yes | Site ID |
| `js_injection_rules` | List (MaxItems:1) | No | JS injection rule configuration |

### js_injection_rules sub-fields (100% mapping to SDK JSInjectionRule struct)

| Sub-field | Type | Required | SDK Field | Description |
|---|---|---|---|---|
| `name` | String | Required | `Name` | Rule name |
| `priority` | Int | Required | `Priority` | Rule priority, 0-100, smaller value = higher priority |
| `condition` | String | Required | `Condition` | Match condition expression |
| `inject_j_s` | String | Required | `InjectJS` | JS injection option: `no-injection`, `inject-sdk-only` |
| `rule_id` | String | Computed | `RuleId` | Rule ID returned by API |

## Read Logic

Call `DescribeSecurityJSInjectionRule` with `ZoneId`, paginate (Limit=100) until the entry with `JSInjectionRule.RuleId == jsInjectionRuleId` is found.
If not found → resource deleted → `d.SetId("")`.

## Update Logic

Call `ModifySecurityJSInjectionRule` with `ZoneId` and `JSInjectionRules=[{RuleId: jsInjectionRuleId, Name, Priority, Condition, InjectJS}]`.

## Delete Logic

Call `DeleteSecurityJSInjectionRule` with `ZoneId` and `JSInjectionRuleIds=[jsInjectionRuleId]`.

## Key SDK Types

```go
// teo v20220901
type JSInjectionRule struct {
    RuleId    *string  // Computed
    Name      *string  // Required
    Priority  *int64   // Required
    Condition *string  // Required
    InjectJS  *string  // Required
}
```
