# Design: Config Rules Data Sources

## Architecture

Follows the standard pattern (reference: `tencentcloud_igtm_instance_list`):

```
provider.go (registration)
    ├─ data_source_tc_config_system_rules.go
    ├─ data_source_tc_config_rule_evaluation_results.go
    └─ data_source_tc_config_rules.go
           └─ service_tencentcloud_config.go (appended methods)
                  └─ config SDK v20220802
```

## File Layout

| File | Action |
|---|---|
| `tencentcloud/services/config/data_source_tc_config_system_rules.go` | New |
| `tencentcloud/services/config/data_source_tc_config_system_rules.md` | New |
| `tencentcloud/services/config/data_source_tc_config_system_rules_test.go` | New |
| `tencentcloud/services/config/data_source_tc_config_rule_evaluation_results.go` | New |
| `tencentcloud/services/config/data_source_tc_config_rule_evaluation_results.md` | New |
| `tencentcloud/services/config/data_source_tc_config_rule_evaluation_results_test.go` | New |
| `tencentcloud/services/config/data_source_tc_config_rules.go` | New |
| `tencentcloud/services/config/data_source_tc_config_rules.md` | New |
| `tencentcloud/services/config/data_source_tc_config_rules_test.go` | New |
| `tencentcloud/services/config/service_tencentcloud_config.go` | Modified — append 3 service methods |
| `tencentcloud/provider.go` | Modified — register 3 data sources |

## Pagination Strategy

All three APIs use `Limit` + `Offset` with `Total` in response.
- Loop until `len(items) >= total`
- Each page call wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)`
- Limit values: `ListSystemRules`=100, `ListConfigRuleEvaluationResults`=100, `ListConfigRules`=200

## Schema Design

### `tencentcloud_config_system_rules`

**Filters (Optional):** `keyword` (String), `risk_level` (Int)

**Result:** `rule_list` (List of `SystemConfigRule`)
- Fields: `identifier`, `rule_name`, `description`, `risk_level`, `service_function`, `create_time`, `update_time`, `trigger_type` (list), `resource_type` (list), `label` (list), `reference_count`, `identifier_type`

### `tencentcloud_config_rule_evaluation_results`

**Filters:** `config_rule_id` (String, Required), `resource_type` (List), `compliance_type` (List)

**Result:** `result_list` (List of `EvaluationResult`)
- Fields: `resource_id`, `resource_type`, `resource_region`, `resource_name`, `config_rule_id`, `config_rule_name`, `compliance_pack_id`, `risk_level`, `compliance_type`, `invoking_event_message_type`, `config_rule_invoked_time`, `result_recorded_time`, `annotation` (nested: `configuration`, `desired_value`, `operator`, `property`)

### `tencentcloud_config_rules`

**Filters (Optional):** `rule_name` (String), `risk_level` (List of Int), `state` (String), `compliance_result` (List of String), `order_type` (String)

**Result:** `rule_list` (List of `ConfigRule`)
- Fields: `identifier`, `rule_name`, `risk_level`, `service_function`, `create_time`, `description`, `status`, `compliance_result`, `config_rule_id`, `identifier_type`, `compliance_pack_id`, `compliance_pack_name`, `resource_type` (list), `labels` (list)

## Service Methods

```go
DescribeConfigSystemRulesByFilter(ctx, paramMap) ([]*SystemConfigRule, error)
DescribeConfigRuleEvaluationResultsByFilter(ctx, paramMap) ([]*EvaluationResult, error)
DescribeConfigRulesByFilter(ctx, paramMap) ([]*ConfigRule, error)
```
