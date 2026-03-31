## ADDED Requirements

### Requirement: TaskId field shall be readable in Read operation
The system SHALL allow users to read the `task_id` field from tencentcloud_teo_l7_acc_rule resource. This field corresponds to the TaskId returned by the ImportZoneConfig API and represents the task ID for the configuration import operation.

#### Scenario: Read resource returns TaskId when available
- **WHEN** user reads an existing tencentcloud_teo_l7_acc_rule resource
- **AND** the ImportZoneConfig API returns a TaskId in the response
- **THEN** the resource's `task_id` field SHALL be populated with the TaskId value
- **AND** the field SHALL be accessible in the Terraform state

#### Scenario: Read resource handles null TaskId
- **WHEN** user reads an existing tencentcloud_teo_l7_acc_rule resource
- **AND** the ImportZoneConfig API returns null or no TaskId
- **THEN** the resource's `task_id` field SHALL be empty
- **AND** no error SHALL be raised

### Requirement: TaskId field shall be Optional and Computed
The `task_id` field in tencentcloud_teo_l7_acc_rule resource SHALL be marked as both Optional and Computed to ensure backward compatibility and proper behavior.

#### Scenario: Existing resources upgrade without error
- **WHEN** a user upgrades the provider to a version supporting `task_id`
- **AND** their existing state does not contain `task_id` field
- **THEN** the resource SHALL continue to work without errors
- **AND** the `task_id` field SHALL be automatically added to the schema

#### Scenario: Users cannot set task_id in configuration
- **WHEN** a user attempts to set `task_id` in their Terraform configuration
- **THEN** the field SHALL be ignored by the provider
- **AND** the value SHALL not be sent to the API

### Requirement: TaskId field shall not affect Create/Update/Delete operations
The `task_id` field SHALL only be populated during Read operations and SHALL not be sent to the API during Create, Update, or Delete operations.

#### Scenario: Create operation ignores task_id
- **WHEN** a user creates a new tencentcloud_teo_l7_acc_rule resource
- **THEN** the `task_id` field SHALL NOT be sent to the ImportZoneConfig API
- **AND** the create operation SHALL complete successfully

#### Scenario: Update operation ignores task_id
- **WHEN** a user updates an existing tencentcloud_teo_l7_acc_rule resource
- **AND** the `task_id` field is present in the configuration
- **THEN** the `task_id` field SHALL NOT be sent to the ImportZoneConfig API
- **AND** the update operation SHALL complete successfully

#### Scenario: Delete operation ignores task_id
- **WHEN** a user deletes a tencentcloud_teo_l7_acc_rule resource
- **THEN** the `task_id` field SHALL NOT be used in the delete operation
- **AND** the delete operation SHALL complete successfully
