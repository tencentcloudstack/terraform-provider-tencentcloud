## ADDED Requirements

### Requirement: Create sandbox tool
The system SHALL support creating an AGS sandbox tool via `CreateSandboxTool` API with all required and optional parameters.

#### Scenario: Create sandbox tool with minimum required fields
- **WHEN** user provides `tool_name`, `tool_type`, and `network_configuration` with `network_mode`
- **THEN** system calls `CreateSandboxTool` API and returns the created `tool_id`

#### Scenario: Create sandbox tool with all optional fields
- **WHEN** user provides all parameters including `description`, `default_timeout`, `tags`, `client_token`, `role_arn`, `storage_mounts`, `custom_configuration`, `computer_configuration`, `log_configuration`, `persistent`
- **THEN** system calls `CreateSandboxTool` API with all fields and returns the created `tool_id`

#### Scenario: Create sandbox tool fails with nil response
- **WHEN** `CreateSandboxTool` API returns nil response or nil ToolId
- **THEN** system returns a non-retryable error

### Requirement: Read sandbox tool
The system SHALL support reading an AGS sandbox tool via `DescribeSandboxToolList` API using the `tool_id`.

#### Scenario: Read existing sandbox tool
- **WHEN** system reads a sandbox tool with a valid `tool_id`
- **THEN** system returns all tool attributes including `tool_name`, `tool_type`, `status`, `description`, `network_configuration`, `tags`, `create_time`, `update_time`, `role_arn`, `storage_mounts`, `custom_configuration`, `computer_configuration`, `log_configuration`, `persistent`, `status_reason`

#### Scenario: Read non-existent sandbox tool
- **WHEN** `DescribeSandboxToolList` returns empty `SandboxToolSet`
- **THEN** system logs a warning, sets id to empty, and returns nil error

### Requirement: Update sandbox tool
The system SHALL support updating an AGS sandbox tool via `UpdateSandboxTool` API for mutable fields: `description`, `network_configuration`, `tags`, `custom_configuration`, `computer_configuration`.

#### Scenario: Update sandbox tool description
- **WHEN** user changes `description`
- **THEN** system calls `UpdateSandboxTool` API with new description

#### Scenario: Update sandbox tool network configuration
- **WHEN** user changes `network_configuration` (network_mode or vpc_config)
- **THEN** system calls `UpdateSandboxTool` API with new network configuration

#### Scenario: Update sandbox tool tags
- **WHEN** user changes `tags`
- **THEN** system calls `UpdateSandboxTool` API with new tags

### Requirement: Delete sandbox tool
The system SHALL support deleting an AGS sandbox tool via `DeleteSandboxTool` API.

#### Scenario: Delete existing sandbox tool
- **WHEN** user destroys the sandbox tool resource
- **THEN** system calls `DeleteSandboxTool` API with the tool's `tool_id`

### Requirement: Import sandbox tool
The system SHALL support importing an existing sandbox tool by `tool_id`.

#### Scenario: Import sandbox tool by tool_id
- **WHEN** user imports resource with `terraform import tencentcloud_ags_sandbox_tool.xxx tool-xxxxxxxx`
- **THEN** system reads the tool attributes and populates the state

### Requirement: ForceNew for immutable fields
The system SHALL mark immutable fields with `ForceNew: true` so that changes require resource recreation.

#### Scenario: Change tool_name forces recreation
- **WHEN** user modifies `tool_name` on an existing resource
- **THEN** Terraform destroys the existing resource and creates a new one

#### Scenario: Change tool_type forces recreation
- **WHEN** user modifies `tool_type` on an existing resource
- **THEN** Terraform destroys the existing resource and creates a new one

#### Scenario: Change storage_mounts forces recreation
- **WHEN** user modifies `storage_mounts` on an existing resource
- **THEN** Terraform destroys the existing resource and creates a new one

### Requirement: Schema completeness for nested objects
The system SHALL flatten nested SDK structs into Terraform `TypeList` schemas following existing provider conventions.

#### Scenario: network_configuration with VPC config
- **WHEN** user provides `network_configuration` with `network_mode: "VPC"` and `vpc_config`
- **THEN** system properly maps subnet_ids and security_group_ids to the SDK request

#### Scenario: custom_configuration with all sub-fields
- **WHEN** user provides `custom_configuration` with `image`, `env`, `ports`, `resources`, `probe`, `dns_config`, `security_context`
- **THEN** system builds the complete `CustomConfiguration` SDK struct

### Requirement: Tags handling
The system SHALL support tags as `TypeMap` in the Terraform schema, converting to/from `[]*Tag` for SDK calls.

#### Scenario: Create with tags
- **WHEN** user provides `tags = { "env" = "prod" }`
- **THEN** system sends `Tags: [{Key: "env", Value: "prod"}]` to the API

#### Scenario: Read returns tags
- **WHEN** API returns `Tags: [{Key: "env", Value: "prod"}]`
- **THEN** system sets `tags = { "env" = "prod" }` in state