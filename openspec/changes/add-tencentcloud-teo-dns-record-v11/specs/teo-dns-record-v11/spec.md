## ADDED Requirements

### Requirement: Create DNS Record
The system SHALL allow users to create a DNS record in TEO service by providing zone ID, domain, record type, record value, and other optional parameters.

#### Scenario: Successful DNS record creation
- **WHEN** user provides valid zone_id, domain, record_type, and record_value
- **THEN** system creates a DNS record in TEO service
- **AND** system returns the record ID
- **AND** system waits for the record to be created successfully

#### Scenario: DNS record creation with invalid parameters
- **WHEN** user provides invalid zone_id or missing required parameters
- **THEN** system returns an error message
- **AND** no DNS record is created

#### Scenario: DNS record creation with duplicate record
- **WHEN** user provides parameters that match an existing DNS record
- **THEN** system returns an error message indicating the record already exists

### Requirement: Read DNS Record
The system SHALL allow users to read DNS record details by providing the record ID and zone ID.

#### Scenario: Successful DNS record read
- **WHEN** user provides valid zone_id and record_id
- **THEN** system retrieves and returns the DNS record details
- **AND** returned details include domain, record_type, record_value, TTL, and other parameters

#### Scenario: Read non-existent DNS record
- **WHEN** user provides zone_id and record_id that do not exist
- **THEN** system returns a not found error

#### Scenario: Read DNS record with invalid zone ID
- **WHEN** user provides invalid zone_id
- **THEN** system returns an error message

### Requirement: Update DNS Record
The system SHALL allow users to update an existing DNS record's parameters such as record value, TTL, and other modifiable fields.

#### Scenario: Successful DNS record update
- **WHEN** user provides valid zone_id, record_id, and updated parameters
- **THEN** system updates the DNS record in TEO service
- **AND** system waits for the update to be applied successfully
- **AND** updated record reflects the new parameter values

#### Scenario: DNS record update with read-only parameters
- **WHEN** user attempts to update read-only parameters (e.g., domain, record_type)
- **THEN** system ignores these parameters or returns an error
- **AND** only modifiable parameters are updated

#### Scenario: Update non-existent DNS record
- **WHEN** user provides zone_id and record_id that do not exist
- **THEN** system returns a not found error
- **AND** no update is performed

### Requirement: Delete DNS Record
The system SHALL allow users to delete an existing DNS record by providing zone ID and record ID.

#### Scenario: Successful DNS record deletion
- **WHEN** user provides valid zone_id and record_id
- **THEN** system deletes the DNS record from TEO service
- **AND** system waits for the deletion to be applied successfully
- **AND** subsequent read operations return not found error

#### Scenario: Delete non-existent DNS record
- **WHEN** user provides zone_id and record_id that do not exist
- **THEN** system returns a not found error
- **AND** no deletion is performed

#### Scenario: Delete already deleted DNS record
- **WHEN** user attempts to delete a DNS record that has already been deleted
- **THEN** system returns success (idempotent operation)
- **AND** no error is raised

### Requirement: Handle Async Operations
The system SHALL handle asynchronous DNS record operations by polling the status until completion or timeout.

#### Scenario: Create operation waits for completion
- **WHEN** user creates a DNS record
- **THEN** system initiates the creation
- **AND** system polls the record status until it becomes active
- **AND** system returns success when the record is active
- **AND** system returns an error if polling times out

#### Scenario: Update operation waits for completion
- **WHEN** user updates a DNS record
- **THEN** system initiates the update
- **AND** system polls the record status until the update is reflected
- **AND** system returns success when the update is complete
- **AND** system returns an error if polling times out

#### Scenario: Delete operation waits for completion
- **WHEN** user deletes a DNS record
- **THEN** system initiates the deletion
- **AND** system polls the record status until it is removed
- **AND** system returns success when the record is deleted
- **AND** system returns an error if polling times out

### Requirement: Configurable Timeouts
The system SHALL allow users to configure custom timeout values for create, update, and delete operations.

#### Scenario: Default timeouts
- **WHEN** user does not specify custom timeouts
- **THEN** system uses default timeout values (10 minutes)

#### Scenario: Custom create timeout
- **WHEN** user specifies a custom create timeout
- **THEN** system uses the specified timeout for create operations
- **AND** system returns timeout error if operation does not complete within the specified time

#### Scenario: Custom update timeout
- **WHEN** user specifies a custom update timeout
- **THEN** system uses the specified timeout for update operations
- **AND** system returns timeout error if operation does not complete within the specified time

#### Scenario: Custom delete timeout
- **WHEN** user specifies a custom delete timeout
- **THEN** system uses the specified timeout for delete operations
- **AND** system returns timeout error if operation does not complete within the specified time

### Requirement: Error Handling and Logging
The system SHALL provide comprehensive error handling and logging for all DNS record operations.

#### Scenario: API error handling
- **WHEN** TEO API returns an error
- **THEN** system translates the error to a user-friendly message
- **AND** system logs the error details for debugging

#### Scenario: Network error handling
- **WHEN** network connectivity issues occur
- **THEN** system retries the operation with exponential backoff
- **AND** system returns an error if retries are exhausted

#### Scenario: Operation logging
- **WHEN** any DNS record operation is performed
- **THEN** system logs the operation details (start time, end time, duration)
- **AND** system logs the operation result (success/failure)

### Requirement: Resource Schema
The system SHALL define a complete resource schema that maps to TEO DNS record parameters.

#### Scenario: Required fields
- **WHEN** user creates a DNS record
- **THEN** system requires zone_id, domain, record_type, and record_value

#### Scenario: Optional fields
- **WHEN** user creates or updates a DNS record
- **THEN** system accepts optional parameters such as TTL, priority, weight, etc.

#### Scenario: Computed fields
- **WHEN** user reads a DNS record
- **THEN** system returns computed fields such as record_id, status, created_at, updated_at

#### Scenario: Schema validation
- **WHEN** user provides invalid data types or values
- **THEN** system validates the input
- **AND** system returns an error message for invalid inputs
