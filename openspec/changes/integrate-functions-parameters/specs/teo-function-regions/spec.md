## ADDED Requirements

### Requirement: TEO Function supports region selection configuration
The system SHALL allow users to configure region selection for TEO functions through the `region_selection` field in the `tencentcloud_teo_function` resource. Region selection SHALL include:
- `regions`: List of region codes where the function should be deployed (optional, if not specified, function is deployed globally)

#### Scenario: Create function with region selection
- **WHEN** user creates a `tencentcloud_teo_function` resource with `region_selection` containing specific region codes
- **THEN** system SHALL create the function with the specified region deployment
- **AND** the function SHALL be deployed only in the specified regions
- **AND** the function SHALL be accessible in the Read operation with region configuration

#### Scenario: Update function region selection
- **WHEN** user updates the `region_selection` field in an existing `tencentcloud_teo_function` resource
- **THEN** system SHALL update the function's region deployment accordingly
- **AND** the function SHALL be deployed in the newly specified regions
- **AND** any region not in the new list SHALL have the function removed

#### Scenario: Read function with region selection
- **WHEN** user reads an existing `tencentcloud_teo_function` resource that has region selection configured
- **THEN** system SHALL return all configured regions
- **AND** the regions SHALL be returned as a list of region codes

#### Scenario: Validate region codes
- **WHEN** user provides invalid region codes in `region_selection`
- **THEN** system SHALL return a validation error
- **AND** the error message SHALL indicate which region codes are invalid
- **AND** valid region codes SHALL follow the format defined by TEO (e.g., "CN", "CN.GD", "US")

### Requirement: TEO Function supports global deployment when region selection is not specified
The system SHALL support global deployment for functions when `region_selection` is not specified or when the regions list is empty.

#### Scenario: Create function without region selection
- **WHEN** user creates a `tencentcloud_teo_function` resource without specifying `region_selection`
- **THEN** system SHALL create the function with global deployment
- **AND** the function SHALL be deployed in all available regions
- **AND** the function SHALL be accessible globally

#### Scenario: Create function with empty region list
- **WHEN** user creates a `tencentcloud_teo_function` resource with `region_selection` containing an empty regions list
- **THEN** system SHALL treat this as global deployment
- **AND** the function SHALL be deployed in all available regions

### Requirement: TEO Function supports single and multiple region deployments
The system SHALL support both single-region and multi-region deployments based on the `region_selection` configuration.

#### Scenario: Create function with single region
- **WHEN** user creates a `tencentcloud_teo_function` resource with `region_selection` containing only one region code
- **THEN** system SHALL deploy the function only in that specified region
- **AND** the function SHALL only handle requests from that region

#### Scenario: Create function with multiple regions
- **WHEN** user creates a `tencentcloud_teo_function` resource with `region_selection` containing multiple region codes
- **THEN** system SHALL deploy the function in all specified regions
- **AND** the function SHALL handle requests from any of the specified regions
- **AND** requests SHALL be routed to the nearest available region

### Requirement: TEO Function maintains backward compatibility without region selection
The system SHALL continue to support functions without explicit region selection to maintain backward compatibility.

#### Scenario: Update existing function without adding region selection
- **WHEN** user updates an existing `tencentcloud_teo_function` resource that was created before this feature
- **THEN** system SHALL maintain the function's existing behavior
- **AND** the function SHALL continue to be deployed globally
- **AND** no changes to the function's region configuration SHALL occur

#### Scenario: Read existing function without region selection
- **WHEN** user reads an existing `tencentcloud_teo_function` resource that was created before this feature
- **THEN** system SHALL return the function without `region_selection` field
- **AND** the function SHALL behave identically to functions created before this feature
