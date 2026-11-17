## ADDED Requirements

### Requirement: CLS Config supports input_type parameter
The `tencentcloud_cls_config` resource SHALL support an optional `input_type` parameter of type `TypeString` that specifies the log input source type for collection configurations.

#### Scenario: Create a config with input_type set to windows_event
- **WHEN** user creates a `tencentcloud_cls_config` resource with `input_type = "windows_event"`
- **THEN** the `CreateConfig` API is called with `InputType` set to `"windows_event"`
- **AND** the resource is created successfully

#### Scenario: Create a config without input_type
- **WHEN** user creates a `tencentcloud_cls_config` resource without specifying `input_type`
- **THEN** the `CreateConfig` API is called without the `InputType` field set
- **AND** the resource is created successfully (backward compatible)

#### Scenario: Read a config that has input_type
- **WHEN** user reads a `tencentcloud_cls_config` resource that has `InputType` returned from the API
- **THEN** the `input_type` attribute is set in the Terraform state

#### Scenario: Read a config without input_type
- **WHEN** user reads a `tencentcloud_cls_config` resource that has `InputType` returned as nil from the API
- **THEN** the `input_type` attribute is not set in the Terraform state (no unnecessary diff)

#### Scenario: Update input_type on an existing config
- **WHEN** user changes `input_type` from `"file"` to `"syslog"` on an existing `tencentcloud_cls_config` resource
- **THEN** the `ModifyConfig` API is called with `InputType` set to `"syslog"`
- **AND** the update succeeds without recreating the resource