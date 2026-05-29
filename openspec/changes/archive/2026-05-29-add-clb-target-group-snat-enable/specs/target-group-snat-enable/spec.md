## ADDED Requirements

### Requirement: Support snat_enable parameter in tencentcloud_clb_target_group resource

The `tencentcloud_clb_target_group` resource SHALL support an optional `snat_enable` parameter of type bool that controls whether SNAT (source IP replacement) is enabled for the target group. When `snat_enable` is true, the system replaces the client source IP; when false or not set, SNAT is disabled (default behavior).

#### Scenario: Create target group with snat_enable set to true
- **WHEN** user creates a `tencentcloud_clb_target_group` resource with `snat_enable = true`
- **THEN** the system SHALL pass `SnatEnable = true` to the `CreateTargetGroup` API request and the target group is created with SNAT enabled

#### Scenario: Create target group with snat_enable set to false
- **WHEN** user creates a `tencentcloud_clb_target_group` resource with `snat_enable = false`
- **THEN** the system SHALL pass `SnatEnable = false` to the `CreateTargetGroup` API request and the target group is created with SNAT disabled

#### Scenario: Create target group without snat_enable specified
- **WHEN** user creates a `tencentcloud_clb_target_group` resource without specifying `snat_enable`
- **THEN** the system SHALL NOT pass `SnatEnable` to the `CreateTargetGroup` API request, and the API uses its default (disabled)

#### Scenario: Read target group with snat_enable
- **WHEN** the system reads a target group via `DescribeTargetGroups` API
- **THEN** the system SHALL read the `SnatEnable` field from the `TargetGroupInfo` response and set it in the Terraform state if it is not nil

#### Scenario: Update snat_enable from false to true
- **WHEN** user changes `snat_enable` from `false` to `true` on an existing `tencentcloud_clb_target_group` resource
- **THEN** the system SHALL call `ModifyTargetGroupAttribute` API with `SnatEnable = true` to update the target group in-place without recreation

#### Scenario: Update snat_enable from true to false
- **WHEN** user changes `snat_enable` from `true` to `false` on an existing `tencentcloud_clb_target_group` resource
- **THEN** the system SHALL call `ModifyTargetGroupAttribute` API with `SnatEnable = false` to update the target group in-place without recreation
