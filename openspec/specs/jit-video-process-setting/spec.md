## ADDED Requirements

### Requirement: JIT video process parameter schema definition

The `tencentcloud_teo_zone_setting` resource SHALL expose a `jit_video_process` parameter as an Optional, Computed, TypeList field with MaxItems=1. The inner schema SHALL contain a `switch` field of type string that is Required, accepting values `on` (enable) or `off` (disable).

#### Scenario: Schema definition is correct
- **WHEN** the resource schema is loaded
- **THEN** the `jit_video_process` field SHALL be present with type TypeList, Optional=true, Computed=true, MaxItems=1, and contain a nested `switch` field of type TypeString with Required=true

### Requirement: Read JIT video process from API response

The resource Read method SHALL read the `JITVideoProcess.Switch` value from the `DescribeZoneSetting` API response and set it into the Terraform state as `jit_video_process`.

#### Scenario: API returns JITVideoProcess with Switch=on
- **WHEN** the `DescribeZoneSetting` API response contains `ZoneSetting.JITVideoProcess` with `Switch` = "on"
- **THEN** the resource state SHALL contain `jit_video_process.0.switch` = "on"

#### Scenario: API returns JITVideoProcess as nil
- **WHEN** the `DescribeZoneSetting` API response contains `ZoneSetting.JITVideoProcess` as nil
- **THEN** the resource state SHALL NOT set the `jit_video_process` field (no error)

### Requirement: Update JIT video process via API request

The resource Update method SHALL include `JITVideoProcess` in the `ModifyZoneSetting` API request when the `jit_video_process` field has changed.

#### Scenario: User enables JIT video processing
- **WHEN** the user sets `jit_video_process { switch = "on" }` in their Terraform configuration and applies
- **THEN** the `ModifyZoneSetting` API request SHALL include `JITVideoProcess` with `Switch` = "on"

#### Scenario: User disables JIT video processing
- **WHEN** the user sets `jit_video_process { switch = "off" }` in their Terraform configuration and applies
- **THEN** the `ModifyZoneSetting` API request SHALL include `JITVideoProcess` with `Switch` = "off"

#### Scenario: User does not specify jit_video_process
- **WHEN** the user does not include `jit_video_process` in their Terraform configuration
- **THEN** the `ModifyZoneSetting` API request SHALL NOT include `JITVideoProcess` (preserving existing server-side configuration)

### Requirement: Backward compatibility

Adding the `jit_video_process` parameter SHALL NOT break existing Terraform configurations that do not use this field.

#### Scenario: Existing configuration without jit_video_process
- **WHEN** an existing Terraform configuration for `tencentcloud_teo_zone_setting` does not include `jit_video_process`
- **THEN** the resource SHALL continue to function correctly without errors, and the field SHALL be populated from the API response via the Computed attribute
