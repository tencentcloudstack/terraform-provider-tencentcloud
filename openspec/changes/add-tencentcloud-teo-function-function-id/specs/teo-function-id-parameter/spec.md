## ADDED Requirements

### Requirement: Optional FunctionId Parameter
The tencentcloud_teo_function resource SHALL support an optional `function_id` parameter that users can specify during resource creation.

#### Scenario: User creates function with custom FunctionId
- **WHEN** user provides a `function_id` parameter in the resource configuration
- **THEN** the provider SHALL include the specified `function_id` in the CreateFunction API request
- **AND** the API SHALL create the function with the specified ID
- **AND** the resource state SHALL contain the provided `function_id`

#### Scenario: User creates function without FunctionId
- **WHEN** user does not provide a `function_id` parameter in the resource configuration
- **THEN** the provider SHALL NOT include the `function_id` field in the CreateFunction API request
- **AND** the API SHALL generate a unique FunctionId automatically
- **AND** the resource state SHALL contain the API-generated `function_id` from the response

### Requirement: Backward Compatibility
The change SHALL maintain full backward compatibility with existing tencentcloud_teo_function resources and configurations.

#### Scenario: Existing resource without function_id in state
- **WHEN** a resource created before this change (without `function_id` in config) is refreshed or updated
- **THEN** the resource SHALL continue to function normally
- **AND** the provider SHALL read the `function_id` from the API response
- **AND** no configuration changes SHALL be required

#### Scenario: Existing resource configuration without function_id
- **WHEN** a user applies an existing Terraform configuration without `function_id` parameter
- **THEN** the resource SHALL be created or updated successfully
- **AND** the behavior SHALL be identical to the previous implementation

### Requirement: Schema Definition
The resource schema SHALL define `function_id` as both Optional and Computed.

#### Scenario: Schema field configuration
- **WHEN** the resource schema is defined
- **THEN** the `function_id` field SHALL have `Optional: true` attribute
- **AND** the `function_id` field SHALL have `Computed: true` attribute
- **AND** the field SHALL be of TypeString
- **AND** the field SHALL have a description indicating it is optional

### Requirement: API Request Handling
The provider SHALL conditionally include the FunctionId parameter in the CreateFunction API request based on user input.

#### Scenario: Include FunctionId when provided
- **WHEN** user specifies a value for `function_id` in the resource configuration
- **THEN** the CreateFunction request SHALL include the `FunctionId` field with the user-specified value
- **AND** the field SHALL be set using the helper.String() function

#### Scenario: Omit FunctionId when not provided
- **WHEN** user does not specify a value for `function_id` in the resource configuration
- **THEN** the CreateFunction request SHALL NOT include the `FunctionId` field
- **AND** the API SHALL handle the request as if FunctionId was not specified

### Requirement: Resource ID Consistency
The resource ID SHALL maintain the composite structure of `zone_id#function_id`.

#### Scenario: Resource ID construction
- **WHEN** a function resource is created or imported
- **THEN** the resource ID SHALL be constructed as `zone_id#function_id`
- **AND** the ID SHALL use the FILED_SP constant as separator
- **AND** the function_id SHALL be either user-specified or API-generated

#### Scenario: Import with known FunctionId
- **WHEN** a user imports an existing function using `terraform import tencentcloud_teo_function.example zone_id#function_id`
- **THEN** the import SHALL succeed
- **AND** the resource state SHALL contain the imported `function_id`
- **AND** the `function_id` field SHALL be set to the value from the import ID

### Requirement: Read Operation Consistency
The read operation SHALL continue to function correctly for both scenarios (user-specified and API-generated FunctionId).

#### Scenario: Read function with API-generated ID
- **WHEN** the provider reads a function that was created without specifying `function_id`
- **THEN** the read operation SHALL retrieve the function details from the API
- **AND** the `function_id` field SHALL be set from the API response
- **AND** the resource state SHALL be updated correctly

#### Scenario: Read function with user-specified ID
- **WHEN** the provider reads a function that was created with a user-specified `function_id`
- **THEN** the read operation SHALL retrieve the function details from the API
- **AND** the `function_id` field SHALL be set from the API response
- **AND** the resource state SHALL be updated correctly

### Requirement: Error Handling
The provider SHALL handle API errors appropriately when user-specified FunctionId is invalid or conflicts.

#### Scenario: API rejects invalid FunctionId
- **WHEN** user specifies a `function_id` that does not meet API validation requirements
- **THEN** the API SHALL return an error
- **AND** the provider SHALL return the error to the user without modifying state
- **AND** the error message SHALL be preserved for debugging

#### Scenario: API rejects duplicate FunctionId
- **WHEN** user specifies a `function_id` that already exists for the same zone
- **THEN** the API SHALL return a conflict error
- **AND** the provider SHALL return the error to the user without modifying state
- **AND** the error message SHALL indicate the duplicate ID conflict

### Requirement: Update Operation
The update operation SHALL NOT allow modification of the `function_id` field.

#### Scenario: Attempt to update function_id
- **WHEN** user tries to change the `function_id` value in the resource configuration
- **THEN** the update operation SHALL reject the change
- **AND** the provider SHALL return an error indicating that `function_id` cannot be changed
- **OR** the provider SHALL treat `function_id` as immutable (like other identifier fields)
