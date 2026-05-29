## ADDED Requirements

### Requirement: Support snat_enable parameter in CLB target group resource

The `tencentcloud_clb_target_group` resource SHALL support a `snat_enable` parameter (bool, Optional, Computed) that controls whether SNAT (source IP replacement) is enabled for the target group. When `snat_enable` is set to `true`, the system MUST replace the client source IP. When set to `false` or not specified, SNAT MUST be disabled (default behavior).

#### Scenario: Create target group with snat_enable set to true
- **WHEN** user creates a `tencentcloud_clb_target_group` resource with `snat_enable = true`
- **THEN** the system SHALL pass `SnatEnable = true` to the `CreateTargetGroup` API and the created target group SHALL have SNAT enabled

#### Scenario: Create target group with snat_enable set to false
- **WHEN** user creates a `tencentcloud_clb_target_group` resource with `snat_enable = false`
- **THEN** the system SHALL pass `SnatEnable = false` to the `CreateTargetGroup` API and the created target group SHALL have SNAT disabled

#### Scenario: Create target group without specifying snat_enable
- **WHEN** user creates a `tencentcloud_clb_target_group` resource without specifying `snat_enable`
- **THEN** the system SHALL not pass `SnatEnable` to the `CreateTargetGroup` API and the API default (false/disabled) SHALL apply

#### Scenario: Read target group snat_enable state
- **WHEN** the system reads the target group state via `DescribeTargetGroups` API
- **THEN** the system SHALL read the `SnatEnable` field from `TargetGroupInfo` and set it in the Terraform state if it is not nil

#### Scenario: Update snat_enable from false to true
- **WHEN** user changes `snat_enable` from `false` to `true` on an existing target group
- **THEN** the system SHALL call `ModifyTargetGroupAttribute` API with `SnatEnable = true` to update the target group in-place without recreation

#### Scenario: Update snat_enable from true to false
- **WHEN** user changes `snat_enable` from `true` to `false` on an existing target group
- **THEN** the system SHALL call `ModifyTargetGroupAttribute` API with `SnatEnable = false` to update the target group in-place without recreation

#### Scenario: Import target group with snat_enable
- **WHEN** user imports an existing target group that has SNAT enabled
- **THEN** the system SHALL correctly read and populate the `snat_enable` field in the Terraform state
