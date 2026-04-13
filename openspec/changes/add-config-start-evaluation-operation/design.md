# Design: tencentcloud_config_start_config_rule_evaluation_operation

## Architecture

Follows the one-shot operation resource pattern (reference: `tencentcloud_advisor_authorization_operation`), combined with igtm_strategy code style for variable declarations:

```
provider.go (registration)
    └─ resource_tc_config_start_config_rule_evaluation_operation.go
           └─ config SDK v20220802 (StartConfigRuleEvaluationWithContext)
```

No service layer method needed — direct API call in the resource file.

## File Layout

| File | Action |
|---|---|
| `tencentcloud/services/config/resource_tc_config_start_config_rule_evaluation_operation.go` | New — operation resource |
| `tencentcloud/services/config/resource_tc_config_start_config_rule_evaluation_operation.md` | New — inline doc |
| `tencentcloud/services/config/resource_tc_config_start_config_rule_evaluation_operation_test.go` | New — acceptance test |
| `tencentcloud/provider.go` | Modified — register resource |

## Schema

### Optional

| Field | Type | Description |
|---|---|---|
| `rule_id` | String | Config rule ID to trigger evaluation for |
| `compliance_pack_id` | String | Compliance pack ID to trigger evaluation for |

> At least one of `rule_id` or `compliance_pack_id` should be provided (API allows both optional, but practically one is needed).

## Resource Lifecycle

| Handler | Behavior |
|---|---|
| Create | Call `StartConfigRuleEvaluation`; `d.SetId(helper.BuildToken())` |
| Read | No-op — return nil |
| Delete | No-op — return nil |

No `Update`, no `Importer` (one-shot resources are not importable).
