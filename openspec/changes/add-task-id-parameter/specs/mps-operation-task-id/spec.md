# Spec: MPS Operation Task ID

## ADDED Requirements

### Requirement: MPS operation resources expose task_id computed field

The MPS (Media Processing Service) operation resources SHALL expose a `task_id` computed field in their Terraform schema to provide users with visibility into the task ID returned by the underlying Tencent Cloud API. This field MUST be set when the resource is created and MUST be available for reference in other Terraform resources or outputs.

The `task_id` field MUST be defined as:
- Type: `string`
- Computed: `true`
- Optional: `false` (not user-provided)
- Required: `false` (not user-required)
- Description: "Task ID returned by API, used to track media processing task status."

The `task_id` value MUST match the value used as the Terraform resource ID (`id` attribute).

Affected resources:
- `tencentcloud_mps_process_media_operation`
- `tencentcloud_mps_edit_media_operation`
- `tencentcloud_mps_process_live_stream_operation`
- Any other MPS operation resources that return TaskId from their respective API

#### Scenario: Create process_media_operation returns task_id

- **WHEN** user creates a `tencentcloud_mps_process_media_operation` resource with valid media processing parameters
- **THEN** the resource creation succeeds
- **AND** the `task_id` attribute is populated with the TaskId returned by the ProcessMedia API
- **AND** the `task_id` value matches the resource's `id` attribute
- **AND** the `task_id` value is a non-empty string

#### Scenario: Create edit_media_operation returns task_id

- **WHEN** user creates a `tencentcloud_mps_edit_media_operation` resource with valid edit parameters
- **THEN** the resource creation succeeds
- **AND** the `task_id` attribute is populated with the TaskId returned by the EditMedia API
- **AND** the `task_id` value matches the resource's `id` attribute
- **AND** the `task_id` value is a non-empty string

#### Scenario: Create process_live_stream_operation returns task_id

- **WHEN** user creates a `tencentcloud_mps_process_live_stream_operation` resource with valid stream processing parameters
- **THEN** the resource creation succeeds
- **AND** the `task_id` attribute is populated with the TaskId returned by the ProcessLiveStream API
- **AND** the `task_id` value matches the resource's `id` attribute
- **AND** the `task_id` value is a non-empty string

#### Scenario: User references task_id in other resources

- **WHEN** a user creates an MPS operation resource
- **AND** the user references the `task_id` attribute in another Terraform resource or output
- **THEN** the reference resolves to the correct task ID value
- **AND** the referenced value is consistent across the entire Terraform configuration

#### Scenario: Read operation does not modify task_id

- **WHEN** a user performs a `terraform refresh` or read operation on an MPS operation resource
- **THEN** the `task_id` attribute retains its original value
- **AND** the `task_id` value continues to match the resource's `id` attribute
- **AND** no additional API calls are made to refresh the task_id value

#### Scenario: Existing configurations remain functional

- **WHEN** a user applies a Terraform configuration that uses existing MPS operation resources without referencing `task_id`
- **THEN** the existing resources continue to work without modification
- **AND** the `task_id` field is automatically added to the resource state on next refresh
- **AND** the user does not need to update their configuration to accommodate the new field

#### Scenario: Task ID format validation

- **WHEN** a user creates an MPS operation resource
- **THEN** the `task_id` attribute follows the format returned by the Tencent Cloud API
- **AND** the `task_id` is a valid string identifier used by the MPS service

#### Scenario: Documentation includes task_id

- **WHEN** a user views the documentation for an MPS operation resource
- **THEN** the documentation includes the `task_id` attribute
- **AND** the description clearly states that `task_id` is computed and matches the resource's `id`
- **AND** usage examples show how to reference the `task_id` in other resources or outputs
