## ADDED Requirements

### Requirement: Resource shall support Rules parameter in schema
The `tencentcloud_teo_l7_acc_rule` resource SHALL include a Rules parameter in its schema definition that represents the access control rule configuration. The Rules parameter MUST be a nested structure that matches the API response structure of `DescribeL7AccRules`.

#### Scenario: Define Rules parameter in resource schema
- **WHEN** defining the resource schema for `tencentcloud_teo_l7_acc_rule`
- **THEN** the schema MUST include a Rules field with appropriate type (List/Set of nested objects)
- **AND** the Rules field MUST support all properties from the API response

### Requirement: Resource shall read Rules from DescribeL7AccRules API
The resource SHALL call `DescribeL7AccRules` API during read operations and extract the Rules field from the API response.

#### Scenario: Read operation retrieves Rules parameter
- **WHEN** Terraform performs a read operation on `tencentcloud_teo_l7_acc_rule` resource
- **THEN** the provider MUST call `DescribeL7AccRules` API
- **AND** the provider MUST parse the Rules field from the API response
- **AND** the provider MUST set the parsed Rules value in the Terraform state

### Requirement: Rules parameter must maintain backward compatibility
The addition of the Rules parameter MUST NOT break existing Terraform configurations or state files.

#### Scenario: Existing resource continues to work without Rules
- **WHEN** an existing `tencentcloud_teo_l7_acc_rule` resource is updated with provider that includes Rules parameter
- **THEN** the resource MUST continue to function without requiring configuration changes
- **AND** the state file MUST remain compatible

### Requirement: Rules parameter must support all required fields
The Rules parameter MUST support all fields returned by the `DescribeL7AccRules` API response.

#### Scenario: All API fields mapped to schema
- **WHEN** mapping API response to Terraform schema
- **THEN** every field in the API's Rules response MUST have a corresponding field in the Terraform schema
- **AND** field types MUST match (e.g., string, number, boolean, list, nested objects)
