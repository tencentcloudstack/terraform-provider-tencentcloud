## ADDED Requirements

### Requirement: User can configure SNAT enable on target group creation

The `tencentcloud_clb_target_group` resource SHALL support an optional `snat_enable` parameter of type bool. When set to `true`, SNAT (source IP replacement) SHALL be enabled on the target group during creation. When set to `false` or not specified, SNAT SHALL remain disabled (default behavior). The parameter SHALL be passed as `SnatEnable` to the `CreateTargetGroup` API.

#### Scenario: Create target group with SNAT enabled
- **WHEN** user specifies `snat_enable = true` in the resource configuration
- **THEN** the system SHALL pass `SnatEnable = true` to the `CreateTargetGroup` API and the target group SHALL be created with SNAT enabled

#### Scenario: Create target group with SNAT disabled explicitly
- **WHEN** user specifies `snat_enable = false` in the resource configuration
- **THEN** the system SHALL pass `SnatEnable = false` to the `CreateTargetGroup` API and the target group SHALL be created with SNAT disabled

#### Scenario: Create target group without specifying snat_enable
- **WHEN** user does not specify `snat_enable` in the resource configuration
- **THEN** the system SHALL NOT pass `SnatEnable` to the `CreateTargetGroup` API and the target group SHALL use the API default (SNAT disabled)

### Requirement: User can update SNAT enable on existing target group

The `tencentcloud_clb_target_group` resource SHALL support updating the `snat_enable` parameter after creation via the `ModifyTargetGroupAttribute` API.

#### Scenario: Update target group to enable SNAT
- **WHEN** user changes `snat_enable` from `false` to `true` on an existing target group
- **THEN** the system SHALL call `ModifyTargetGroupAttribute` API with `SnatEnable = true`

#### Scenario: Update target group to disable SNAT
- **WHEN** user changes `snat_enable` from `true` to `false` on an existing target group
- **THEN** the system SHALL call `ModifyTargetGroupAttribute` API with `SnatEnable = false`

### Requirement: SNAT enable is a write-only attribute

Since the `DescribeTargetGroups` API response does not include the `SnatEnable` field, the resource Read method SHALL NOT attempt to read or set this value from the API. The value SHALL be preserved in Terraform state from the user's configuration.

#### Scenario: Read does not overwrite snat_enable state
- **WHEN** the resource Read method is called (e.g., during refresh)
- **THEN** the system SHALL NOT modify the `snat_enable` value in Terraform state (it remains as configured by the user)

#### Scenario: Import does not populate snat_enable
- **WHEN** user imports an existing target group resource
- **THEN** the `snat_enable` field SHALL NOT be populated in the imported state (it will be empty/unset)
