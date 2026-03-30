## ADDED Requirements

### Requirement: User can specify custom function_id during creation
The system SHALL allow users to specify a custom `function_id` parameter when creating a `tencentcloud_teo_function` resource. The `function_id` parameter SHALL be optional, and when provided, it SHALL be passed to the CreateFunction API. When not provided, the system SHALL allow the API to auto-generate the function_id.

#### Scenario: Create resource with custom function_id
- **WHEN** user provides a valid `function_id` in the resource configuration
- **THEN** the provider SHALL pass the `function_id` to the CreateFunction API
- **THEN** the provider SHALL use the returned `function_id` from API response for state management

#### Scenario: Create resource without function_id
- **WHEN** user does not provide a `function_id` in the resource configuration
- **THEN** the provider SHALL NOT pass any `function_id` to the CreateFunction API
- **THEN** the provider SHALL use the auto-generated `function_id` returned from API response for state management

#### Scenario: Backward compatibility
- **WHEN** existing resource configuration does not include `function_id` parameter
- **THEN** the resource SHALL continue to work without any changes
- **THEN** the system SHALL auto-generate `function_id` as before

### Requirement: function_id parameter schema definition
The `function_id` parameter SHALL be defined as both Optional and Computed in the resource schema to support both user-provided and system-generated values.

#### Scenario: Schema supports input and output
- **WHEN** the resource schema is inspected
- **THEN** the `function_id` parameter SHALL have both Optional and Computed flags set
- **THEN** the parameter SHALL allow user input during creation
- **THEN** the parameter SHALL be populated and readable after resource creation

### Requirement: Error handling for duplicate function_id
The system SHALL handle API errors when the user-provided `function_id` conflicts with an existing function and return the error message to the user.

#### Scenario: function_id conflict error
- **WHEN** user provides a `function_id` that already exists
- **THEN** the provider SHALL return the API error message to the user
- **THEN** the error message SHALL clearly indicate the conflict reason

### Requirement: Import support with custom function_id
The system SHALL continue to support importing resources using the composite ID format `zone_id#function_id`, regardless of whether the `function_id` was user-provided or auto-generated.

#### Scenario: Import resource with custom function_id
- **WHEN** user imports a resource with a custom `function_id`
- **THEN** the import SHALL succeed using the composite ID format
- **THEN** the imported resource state SHALL correctly reflect the `function_id`

#### Scenario: Import resource with auto-generated function_id
- **WHEN** user imports a resource with an auto-generated `function_id`
- **THEN** the import SHALL succeed using the composite ID format
- **THEN** the imported resource state SHALL correctly reflect the `function_id`
