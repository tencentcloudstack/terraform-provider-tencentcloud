# Add Config Rules Data Sources

## What

Add three new Terraform data sources for the Tencent Cloud Config (配置审计) service:

1. `tencentcloud_config_system_rules` — query system-preset rules via `ListSystemRules`
2. `tencentcloud_config_rule_evaluation_results` — query rule evaluation results (rule dimension) via `ListConfigRuleEvaluationResults`
3. `tencentcloud_config_rules` — query user-defined/managed config rules via `ListConfigRules`

## Why

The existing `tencentcloud_config_compliance_pack` resource handles compliance pack CRUD, but there are no data sources to inspect or list individual config rules and their evaluation results. Users need to:

- Enumerate available system-preset rules (templates) to build custom compliance packs
- Query evaluation results by rule ID to understand compliance state of resources
- List and filter active user-created config rules by state, risk level, compliance result, etc.

## APIs Used

| Data Source | API | Pagination Limit |
|---|---|---|
| `tencentcloud_config_system_rules` | `ListSystemRules` | 100 (no explicit max in docs; using practical safe value) |
| `tencentcloud_config_rule_evaluation_results` | `ListConfigRuleEvaluationResults` | 100 |
| `tencentcloud_config_rules` | `ListConfigRules` | 200 (documented max) |

## Resource IDs

All data sources use `helper.BuildToken()` as computed ID (read-only).
