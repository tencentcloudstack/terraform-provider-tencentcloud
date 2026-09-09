# monitor-alarm-policy-is-bind-all Specification

## Purpose
TBD - created by archiving change add-monitor-alarm-policy-is-bind-all. Update Purpose after archive.
## Requirements
### Requirement: IsBindAll Parameter for Alarm Policy
The `tencentcloud_monitor_alarm_policy` resource SHALL support the `is_bind_all` parameter to specify whether the alarm policy binds to all objects. The parameter SHALL be an optional integer with valid values `0` (no, default) and `1` (yes). The parameter SHALL be immutable after creation (ForceNew) because no Modify API supports updating `IsBindAll`.

#### Scenario: Create alarm policy with is_bind_all set to 1
- **WHEN** user creates a `tencentcloud_monitor_alarm_policy` resource with `is_bind_all` set to `1`
- **THEN** the provider SHALL pass `IsBindAll` with value `1` to the `CreateAlarmPolicy` API request
- **AND** the alarm policy SHALL be created as a "bind all objects" policy

#### Scenario: Create alarm policy without is_bind_all
- **WHEN** user creates a `tencentcloud_monitor_alarm_policy` resource without specifying `is_bind_all`
- **THEN** the provider SHALL NOT pass `IsBindAll` to the `CreateAlarmPolicy` API request (or pass the default value)
- **AND** the alarm policy SHALL be created with the API default behavior (not bind all)

#### Scenario: Read alarm policy is_bind_all
- **WHEN** Terraform performs a read operation on the alarm policy resource
- **THEN** the provider SHALL read `IsBindAll` from the `DescribeAlarmPolicy` API response (`response.Policy.IsBindAll`)
- **AND** if `IsBindAll` is nil, the provider SHALL NOT set the `is_bind_all` field
- **AND** if `IsBindAll` is not nil, the provider SHALL set `is_bind_all` in Terraform state to the returned value

#### Scenario: Update is_bind_all triggers resource recreation
- **WHEN** user modifies `is_bind_all` in an existing alarm policy resource
- **THEN** the provider SHALL treat the change as a ForceNew operation
- **AND** the existing alarm policy SHALL be destroyed and a new one SHALL be created with the updated `is_bind_all` value

#### Scenario: Validate is_bind_all accepts only 0 or 1
- **WHEN** user sets `is_bind_all` to a value other than `0` or `1`
- **THEN** the provider SHALL reject the configuration with a validation error

### Requirement: Backward Compatibility for IsBindAll
The new `is_bind_all` parameter SHALL be fully backward compatible with existing alarm policy configurations.

#### Scenario: Existing configuration without is_bind_all
- **WHEN** user has an existing alarm policy without `is_bind_all`
- **THEN** the provider SHALL continue to work without errors
- **AND** existing Terraform state SHALL remain valid

