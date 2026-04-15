## ADDED Requirements

### Requirement: Create TEO Function V2
The system SHALL allow users to create a new TEO (EdgeOne) function with specified parameters through the `tencentcloud_teo_function_v2` Terraform resource.

#### Scenario: Successful creation
- **WHEN** user creates a `tencentcloud_teo_function_v2` resource with valid `zone_id`, `name`, and `content`
- **THEN** system should call CreateFunction API with the provided parameters
- **THEN** system should assign a unique `function_id` to the newly created function
- **THEN** system should wait for the function to be fully provisioned (domain is assigned)
- **THEN** Terraform state should contain the function ID and all computed attributes

#### Scenario: Creation with optional remark
- **WHEN** user creates a `tencentcloud_teo_function_v2` resource with `remark` parameter
- **THEN** system should include the remark in the CreateFunction API call
- **THEN** the function's remark should be persisted and retrievable

#### Scenario: Creation with invalid parameters
- **WHEN** user creates a `tencentcloud_teo_function_v2` resource with invalid `name` (contains uppercase letters or special characters other than hyphens)
- **THEN** system should return a validation error
- **THEN** Terraform should report the error to the user

### Requirement: Read TEO Function V2
The system SHALL allow users to read the current state of a TEO function using the `tencentcloud_teo_function_v2` Terraform resource.

#### Scenario: Successful read
- **WHEN** user reads an existing `tencentcloud_teo_function_v2` resource by its ID
- **THEN** system should call DescribeFunctions API with the function ID
- **THEN** system should populate Terraform state with all function attributes including `zone_id`, `function_id`, `name`, `remark`, `content`, `domain`, `create_time`, and `update_time`

#### Scenario: Read of non-existent function
- **WHEN** user reads a `tencentcloud_teo_function_v2` resource that has been deleted
- **THEN** system should recognize the function does not exist
- **THEN** Terraform state should be cleared for this resource

#### Scenario: Read with filters
- **WHEN** system needs to query functions by name or remark
- **THEN** system should use the Filters parameter in DescribeFunctions API
- **THEN** system should support fuzzy matching for name and remark fields

### Requirement: Update TEO Function V2
The system SHALL allow users to update mutable attributes of a TEO function through the `tencentcloud_teo_function_v2` Terraform resource.

#### Scenario: Successful update of remark
- **WHEN** user updates the `remark` field of an existing `tencentcloud_teo_function_v2` resource
- **THEN** system should call ModifyFunction API with the new remark
- **THEN** the function's remark should be updated in the cloud
- **THEN** Terraform state should reflect the new remark value

#### Scenario: Successful update of content
- **WHEN** user updates the `content` field of an existing `tencentcloud_teo_function_v2` resource
- **THEN** system should call ModifyFunction API with the new content
- **THEN** the function's content should be updated in the cloud
- **THEN** Terraform state should reflect the new content value

#### Scenario: Attempt to update immutable name
- **WHEN** user attempts to update the `name` field of an existing `tencentcloud_teo_function_v2` resource
- **THEN** system should return an error indicating that name cannot be changed
- **THEN** Terraform should report the error to the user without modifying the function

#### Scenario: Partial update
- **WHEN** user updates only the `remark` field without changing `content`
- **THEN** system should call ModifyFunction API with only the remark parameter
- **THEN** the function's content should remain unchanged

### Requirement: Delete TEO Function V2
The system SHALL allow users to delete a TEO function through the `tencentcloud_teo_function_v2` Terraform resource.

#### Scenario: Successful deletion
- **WHEN** user deletes an existing `tencentcloud_teo_function_v2` resource
- **THEN** system should call DeleteFunction API with the zone ID and function ID
- **THEN** the function should be removed from the cloud
- **THEN** Terraform state should be cleared for this resource

#### Scenario: Deletion of non-existent function
- **WHEN** user attempts to delete a `tencentcloud_teo_function_v2` resource that does not exist
- **THEN** system should handle the error gracefully
- **THEN** Terraform state should be cleared for this resource

### Requirement: Resource ID Format
The system SHALL use a composite resource ID format for `tencentcloud_teo_function_v2` to uniquely identify each function.

#### Scenario: Resource ID format
- **WHEN** a `tencentcloud_teo_function_v2` resource is created or imported
- **THEN** the resource ID should be in the format `zone_id#function_id`
- **THEN** the system should parse the ID to extract zone_id and function_id for API calls

#### Scenario: Import existing function
- **WHEN** user imports an existing TEO function using `terraform import`
- **THEN** system should accept the composite ID format `zone_id#function_id`
- **THEN** system should populate Terraform state with the function's current attributes

### Requirement: Asynchronous Operation Handling
The system SHALL properly handle asynchronous operations when creating or updating TEO functions.

#### Scenario: Wait for function creation completion
- **WHEN** CreateFunction API is called
- **THEN** system should poll DescribeFunctions API until the function's domain is assigned
- **THEN** system should use a 10-second delay with 3-second minimum timeout and 600-second maximum timeout
- **THEN** system should return success only when the function is fully provisioned

#### Scenario: Handle timeout during creation
- **WHEN** function creation takes longer than the timeout period
- **THEN** system should return a timeout error
- **THEN** Terraform should report the error to the user
- **THEN** the user may retry the operation

### Requirement: Retry Logic
The system SHALL implement retry logic for handling transient errors and eventual consistency.

#### Scenario: Retry on write operation failure
- **WHEN** a write operation (Create, Update, Delete) fails due to transient network issues
- **THEN** system should retry the operation up to the configured timeout
- **THEN** system should use exponential backoff for retries
- **THEN** system should return success if the operation succeeds within the timeout

#### Scenario: Retry on read operation failure
- **WHEN** a read operation fails due to eventual consistency
- **THEN** system should retry the operation up to the configured timeout
- **THEN** system should use exponential backoff for retries
- **THEN** system should return success if the operation succeeds within the timeout

### Requirement: Error Handling and Logging
The system SHALL provide comprehensive error handling and logging for all operations.

#### Scenario: Log successful operations
- **WHEN** any API operation succeeds
- **THEN** system should log the operation details including request and response bodies
- **THEN** log should include a unique log ID for tracking

#### Scenario: Log failed operations
- **WHEN** any API operation fails
- **THEN** system should log the error with context
- **THEN** log should include the request body and error reason
- **THEN** system should return a clear error message to the user

#### Scenario: Inconsistent state check
- **WHEN** any CRUD operation completes
- **THEN** system should check for state inconsistencies
- **THEN** system should log the elapsed time for each operation
