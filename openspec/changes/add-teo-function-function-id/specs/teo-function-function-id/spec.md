## ADDED Requirements

### Requirement: Support FunctionId parameter in tencentcloud_teo_function resource
The tencentcloud_teo_function resource SHALL accept an optional FunctionId parameter in its schema. When provided during resource creation, the parameter SHALL be passed to the CreateFunction API call.

#### Scenario: Create resource with FunctionId provided
- **WHEN** user creates a tencentcloud_teo_function resource with FunctionId parameter specified
- **THEN** the provider SHALL pass the FunctionId value to the CreateFunction API
- **AND** the provider SHALL store the FunctionId in the resource state

#### Scenario: Create resource without FunctionId
- **WHEN** user creates a tencentcloud_teo_function resource without FunctionId parameter
- **THEN** the provider SHALL call CreateFunction API without FunctionId
- **AND** the resource SHALL be created successfully as before (backward compatibility)

### Requirement: Maintain backward compatibility
The FunctionId parameter SHALL be optional (Optional attribute in Terraform schema). Existing configurations without FunctionId SHALL continue to work without any changes.

#### Scenario: Existing resource without FunctionId continues to work
- **WHEN** an existing tencentcloud_teo_function resource (created before this change) is applied
- **THEN** the provider SHALL manage the resource without requiring FunctionId
- **AND** no errors SHALL occur

### Requirement: Read FunctionId from resource state
After resource creation, the provider SHALL read back the FunctionId from the API response and store it in the resource state for consistency.

#### Scenario: Read FunctionId after creation
- **WHEN** tencentcloud_teo_function resource is created with FunctionId
- **THEN** the provider SHALL read the FunctionId from the resource's current state
- **AND** the value SHALL match what was provided during creation
