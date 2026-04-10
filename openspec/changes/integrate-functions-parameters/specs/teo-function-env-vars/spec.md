## ADDED Requirements

### Requirement: TEO Function supports environment variables configuration
The system SHALL allow users to configure environment variables for TEO functions through the `environment_variables` field in the `tencentcloud_teo_function` resource. Each environment variable SHALL include:
- `key`: Variable name (required), limited to 64 bytes, containing only letters, numbers, and special characters @ . - _
- `value`: Variable value (optional), limited to 5000 bytes
- `type`: Variable type (optional, default: "string"), supports "string" or "json"

#### Scenario: Create function with environment variables
- **WHEN** user creates a `tencentcloud_teo_function` resource with `environment_variables` containing multiple variables
- **THEN** system SHALL create the function with specified environment variables
- **AND** each variable's key SHALL be unique within the function
- **AND** the function SHALL be accessible in the Read operation with all environment variables

#### Scenario: Update function environment variables
- **WHEN** user updates the `environment_variables` field in an existing `tencentcloud_teo_function` resource
- **THEN** system SHALL update the function's environment variables accordingly
- **AND** any variable removed from the list SHALL be deleted from the function
- **AND** any new or modified variables SHALL be added or updated

#### Scenario: Read function with environment variables
- **WHEN** user reads an existing `tencentcloud_teo_function` resource that has environment variables
- **THEN** system SHALL return all configured environment variables
- **AND** each variable SHALL include key, value, and type attributes

#### Scenario: Validate environment variable key format
- **WHEN** user provides an environment variable key that exceeds 64 bytes or contains invalid characters
- **THEN** system SHALL return a validation error
- **AND** the error message SHALL indicate which variable key is invalid

### Requirement: TEO Function environment variables support both string and JSON types
The system SHALL support environment variable types of "string" and "json". When type is "json", the value SHALL be a valid JSON object.

#### Scenario: Create function with JSON environment variable
- **WHEN** user creates a function with an environment variable where `type` is "json" and `value` is a valid JSON object
- **THEN** system SHALL accept and store the JSON value
- **AND** the variable SHALL be returned with type "json" in Read operations

#### Scenario: Validate JSON type environment variable
- **WHEN** user provides an environment variable with type "json" but value is not valid JSON
- **THEN** system SHALL return a validation error
- **AND** the error message SHALL indicate that the value must be valid JSON

### Requirement: TEO Function maintains backward compatibility without environment variables
The system SHALL continue to support functions without environment variables to maintain backward compatibility.

#### Scenario: Create function without environment variables
- **WHEN** user creates a `tencentcloud_teo_function` resource without specifying `environment_variables`
- **THEN** system SHALL create the function successfully
- **AND** the function SHALL have no environment variables
- **AND** the function SHALL behave identically to functions created before this feature
