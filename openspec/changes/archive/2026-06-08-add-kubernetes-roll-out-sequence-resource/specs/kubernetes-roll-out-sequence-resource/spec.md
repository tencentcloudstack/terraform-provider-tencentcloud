## ADDED Requirements

### Requirement: Create roll-out sequence
The system SHALL create a TKE cluster roll-out sequence by calling the `CreateRollOutSequence` API with the provided name, sequence_flows, and enabled parameters. Upon successful creation, the system SHALL store the returned ID as the Terraform resource ID (converted from int64 to string). The system SHALL verify that the API response and the returned ID are not nil, returning a non-retryable error if they are.

#### Scenario: Successful creation
- **WHEN** user applies a `tencentcloud_kubernetes_roll_out_sequence` resource with valid name, sequence_flows, and enabled fields
- **THEN** the system calls `CreateRollOutSequence` API with the mapped parameters and sets the resource ID from the response

#### Scenario: API returns nil response or nil ID
- **WHEN** the `CreateRollOutSequence` API returns a nil response or nil ID field
- **THEN** the system returns a non-retryable error indicating the response is invalid

#### Scenario: API call fails with retryable error
- **WHEN** the `CreateRollOutSequence` API call fails
- **THEN** the system retries the call using `tccommon.ReadRetryTimeout` and wraps the error with `tccommon.RetryError()`

### Requirement: Read roll-out sequence
The system SHALL read a TKE cluster roll-out sequence by calling the `DescribeRollOutSequences` API with pagination (Limit=20), iterating through all pages until the sequence with the matching ID is found. If the sequence is not found, the system SHALL remove the resource from state (indicating it was deleted externally).

#### Scenario: Successful read
- **WHEN** the system reads a `tencentcloud_kubernetes_roll_out_sequence` resource
- **THEN** the system paginates through `DescribeRollOutSequences` results, finds the sequence matching the resource ID, and sets all schema fields (name, sequence_flows, enabled) from the response

#### Scenario: Sequence not found
- **WHEN** the system reads a resource but the sequence ID is not found in any page of results
- **THEN** the system removes the resource from Terraform state (calls `d.SetId("")`)

#### Scenario: Setting fields from response
- **WHEN** the system finds the matching sequence in the API response
- **THEN** the system SHALL check each field (Name, SequenceFlows, Enabled) for nil before calling `d.Set()`, only setting non-nil fields

### Requirement: Update roll-out sequence
The system SHALL update a TKE cluster roll-out sequence by calling the `ModifyRollOutSequence` API with the resource ID and all updatable fields (name, sequence_flows, enabled).

#### Scenario: Successful update
- **WHEN** user modifies any field (name, sequence_flows, enabled) of an existing `tencentcloud_kubernetes_roll_out_sequence` resource
- **THEN** the system calls `ModifyRollOutSequence` API with the ID and all current field values, then calls the Read function to refresh state

#### Scenario: API call fails with retryable error
- **WHEN** the `ModifyRollOutSequence` API call fails
- **THEN** the system retries the call using `tccommon.ReadRetryTimeout` and wraps the error with `tccommon.RetryError()`

### Requirement: Delete roll-out sequence
The system SHALL delete a TKE cluster roll-out sequence by calling the `DeleteRollOutSequence` API with the resource ID.

#### Scenario: Successful deletion
- **WHEN** user destroys a `tencentcloud_kubernetes_roll_out_sequence` resource
- **THEN** the system calls `DeleteRollOutSequence` API with the ID parsed from the resource state

#### Scenario: API call fails with retryable error
- **WHEN** the `DeleteRollOutSequence` API call fails
- **THEN** the system retries the call using `tccommon.ReadRetryTimeout` and wraps the error with `tccommon.RetryError()`

### Requirement: Resource schema definition
The system SHALL define the following Terraform schema for `tencentcloud_kubernetes_roll_out_sequence`:
- `name` (Required, String): The roll-out sequence name
- `sequence_flows` (Required, List): List of sequence flow steps, each containing:
  - `tags` (Required, List): List of tags, each containing:
    - `key` (Required, String): Tag key
    - `value` (Required, List of String): Tag values
  - `soak_time` (Required, Int): Wait time in seconds between steps
- `enabled` (Required, Bool): Whether the sequence is enabled

#### Scenario: Valid resource configuration
- **WHEN** user defines a `tencentcloud_kubernetes_roll_out_sequence` resource with all required fields
- **THEN** the system accepts the configuration and proceeds with the CRUD operation

#### Scenario: Missing required field
- **WHEN** user defines a resource without a required field (name, sequence_flows, or enabled)
- **THEN** Terraform reports a validation error before any API call

### Requirement: Resource registration
The system SHALL register `tencentcloud_kubernetes_roll_out_sequence` in `tencentcloud/provider.go` resource map and add the corresponding entry in `tencentcloud/provider.md`.

#### Scenario: Resource is available after registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_kubernetes_roll_out_sequence` is available as a valid resource type

### Requirement: Resource import support
The system SHALL support importing existing roll-out sequences using the sequence ID (as a string representation of the int64 ID).

#### Scenario: Import by ID
- **WHEN** user runs `terraform import tencentcloud_kubernetes_roll_out_sequence.example 123`
- **THEN** the system reads the sequence with ID 123 and populates the state

### Requirement: Unit tests with gomonkey mocks
The system SHALL provide unit tests that mock the TKE cloud API client using gomonkey. Tests SHALL cover Create, Read, Update, and Delete operations and SHALL pass with `go test -gcflags=all=-l`.

#### Scenario: Test create operation
- **WHEN** the create unit test runs
- **THEN** it mocks `CreateRollOutSequence` to return a valid ID and verifies the resource state is set correctly

#### Scenario: Test read operation
- **WHEN** the read unit test runs
- **THEN** it mocks `DescribeRollOutSequences` to return a list containing the target sequence and verifies all fields are set

#### Scenario: Test update operation
- **WHEN** the update unit test runs
- **THEN** it mocks `ModifyRollOutSequence` and verifies the API is called with correct parameters

#### Scenario: Test delete operation
- **WHEN** the delete unit test runs
- **THEN** it mocks `DeleteRollOutSequence` and verifies the API is called with the correct ID
