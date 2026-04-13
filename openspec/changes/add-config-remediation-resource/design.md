# Design: tencentcloud_config_remediation Resource

## Architecture

Follows the standard terraform-provider-tencentcloud resource pattern (reference: `tencentcloud_igtm_strategy`):

```
provider.go (registration)
    └─ resource_tc_config_remediation.go   (CRUD handlers)
           └─ service_tencentcloud_config.go (DescribeConfigRemediationById)
                  └─ config SDK v20220802
```

## File Layout

| File | Action |
|---|---|
| `tencentcloud/services/config/resource_tc_config_remediation.go` | New — resource CRUD |
| `tencentcloud/services/config/resource_tc_config_remediation.md` | New — inline doc |
| `tencentcloud/services/config/resource_tc_config_remediation_test.go` | New — acceptance test |
| `tencentcloud/services/config/service_tencentcloud_config.go` | Modified — append `DescribeConfigRemediationById` |
| `tencentcloud/provider.go` | Modified — register resource |

## Schema

### Required (Create-time, ForceNew where applicable)

| Field | Type | ForceNew | Description |
|---|---|---|---|
| `rule_id` | String | Yes | Config rule ID the remediation is bound to |
| `remediation_type` | String | No | Remediation type. Valid value: `SCF` |
| `remediation_template_id` | String | No | Remediation template ID (SCF function ARN) |
| `invoke_type` | String | No | Execution mode: `MANUAL_EXECUTION`, `AUTO_EXECUTION`, `NON_EXECUTION`, `NOT_CONFIG` |

### Optional

| Field | Type | Description |
|---|---|---|
| `source_type` | String | Template source. Valid value: `CUSTOM` |

### Computed

| Field | Type | Description |
|---|---|---|
| `remediation_id` | String | Remediation setting ID (same as resource ID) |
| `owner_uin` | String | Owner account UIN |
| `remediation_source_type` | String | Source type returned from API |

## Read Strategy

`ListRemediations` does not support querying by `RemediationId` directly. The service method will:
1. Call `ListRemediations` with `RuleIds: [ruleId]` (from state)
2. Iterate the `Remediations` array to find the entry matching `RemediationId`
3. Return nil if not found (resource deleted externally)

## Mutable Fields

All fields except `rule_id` can be updated via `UpdateRemediation`.

## Delete

Call `DeleteRemediations` with `RemediationIds: [remediationId]`.
