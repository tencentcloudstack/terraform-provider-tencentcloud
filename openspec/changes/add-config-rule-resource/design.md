# Design: tencentcloud_config_rule Resource

## Architecture

Follows `tencentcloud_igtm_strategy` style:

```
provider.go
    └─ resource_tc_config_rule.go (CRUD handlers)
           └─ service_tencentcloud_config.go (DescribeConfigRuleById)
                  └─ config SDK v20220802
```

## Schema

### Required (Create)

| Field | Type | ForceNew | Description |
|---|---|---|---|
| `identifier` | String | Yes | Rule template identifier |
| `identifier_type` | String | Yes | Template type: SYSTEM, CUSTOMIZE |
| `rule_name` | String | No | Rule name |
| `resource_type` | List of String | Yes | Supported resource types |
| `trigger_type` | List (object) | No | Trigger type list (max 2); each: `message_type`, `maximum_execution_frequency` |
| `risk_level` | Int | No | Risk level: 1/2/3 |

### Optional

| Field | Type | Description |
|---|---|---|
| `input_parameter` | List (object) | Input params: `parameter_key`, `type`, `value` |
| `description` | String | Rule description |
| `regions_scope` | List of String | Region scope filter |
| `tags_scope` | List (object) | Tag scope filter: `tag_key`, `tag_value` |
| `exclude_resource_ids_scope` | List of String | Excluded resource IDs |
| `status` | String | Rule status: ACTIVE / UN_ACTIVE (maps to Open/Close APIs) |

### Computed

| Field | Type | Description |
|---|---|---|
| `config_rule_id` | String | Rule ID (same as resource ID) |
| `create_time` | String | Creation time |
| `compliance_result` | String | Compliance result |
| `config_rule_invoked_time` | String | Last evaluation time |

## Update Logic

```
if content fields changed (rule_name, trigger_type, risk_level, ...):
    call UpdateConfigRule
if status changed:
    if status == "ACTIVE"    → call OpenConfigRule
    if status == "UN_ACTIVE" → call CloseConfigRule
```
