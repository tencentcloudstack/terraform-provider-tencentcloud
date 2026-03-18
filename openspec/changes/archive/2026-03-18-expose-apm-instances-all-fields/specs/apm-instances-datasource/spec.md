## ADDED Requirements

### Requirement: APM Instances Data Source Fields
The `tencentcloud_apm_instances` data source SHALL expose all fields returned by the `DescribeApmInstances` API in the `instance_list` output attribute, including storage usage, billing info, log configuration, security detection switches, metric duration, dashboard association, and URL convergence thresholds.

#### Scenario: Query instance returns all API fields
- **WHEN** user queries APM instances via `tencentcloud_apm_instances`
- **THEN** the `instance_list` output contains all fields from `ApmInstanceDetail`, including previously missing fields such as `amount_of_used_storage`, `log_region`, `metric_duration`, `is_sql_injection_analysis`, `token`, etc.

#### Scenario: Backward compatibility
- **WHEN** existing Terraform configurations reference `tencentcloud_apm_instances`
- **THEN** all previously available fields continue to work unchanged
- **AND** new fields are available as additional Computed attributes
