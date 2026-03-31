## ADDED Requirements

### Requirement: Expose RuleItems from DescribeRules API
The tencentcloud_teo_rule_engine resource SHALL expose the complete RuleItems array returned by the DescribeRules API as a computed field named "rule_items" in the Terraform schema.

#### Scenario: Read operation returns all rule items
- **WHEN** a user performs a read operation on a tencentcloud_teo_rule_engine resource
- **THEN** the system SHALL populate the "rule_items" computed field with all rule items returned by the DescribeRules API for the specified zone

### Requirement: Maintain backward compatibility for existing schema fields
The tencentcloud_teo_rule_engine resource SHALL maintain all existing schema fields and their behavior when adding the "rule_items" computed field.

#### Scenario: Existing resource read works unchanged
- **WHEN** a user reads an existing tencentcloud_teo_rule_engine resource
- **THEN** all existing fields (rule_id, rule_name, status, rule_priority, tags, rules) SHALL continue to function as before
- **AND** the new "rule_items" field SHALL be available as additional computed output

### Requirement: RuleItems field structure matches API response
The "rule_items" computed field SHALL have a schema structure that represents the RuleItems array structure returned by the DescribeRules API, matching the existing "rules" field structure for consistency.

#### Scenario: RuleItems field contains nested rule data
- **WHEN** the DescribeRules API returns multiple rule items
- **THEN** the "rule_items" field SHALL contain a list of rule items
- **AND** each rule item SHALL contain the same nested structure as the existing "rules" field
- **AND** the structure SHALL include or, and, actions, and sub_rules as appropriate

### Requirement: Service method returns complete RuleItems array
The TeoService SHALL provide a method to fetch and return the complete RuleItems array from the DescribeRules API response.

#### Scenario: Service method fetches all rule items
- **WHEN** the service method is called with a zone_id
- **THEN** the method SHALL call the DescribeRules API
- **AND** the method SHALL return the complete RuleItems array from the API response
- **AND** the method SHALL handle API rate limiting and retries according to existing patterns

### Requirement: New field is computed and read-only
The "rule_items" field SHALL be defined as a computed field in the Terraform schema, making it read-only and not affecting create or update operations.

#### Scenario: Computed field cannot be set by user
- **WHEN** a user attempts to set a value for the "rule_items" field in their Terraform configuration
- **THEN** the system SHALL ignore any user-provided value for this field
- **AND** the field SHALL only be populated during read operations from the API

### Requirement: Documentation includes new field
The tencentcloud_teo_rule_engine resource documentation SHALL describe the new "rule_items" computed field and its purpose.

#### Scenario: Documentation describes new field
- **WHEN** users view the resource documentation
- **THEN** the documentation SHALL include a description of the "rule_items" field
- **AND** the documentation SHALL explain that it contains all rule items from the DescribeRules API
- **AND** the documentation SHALL note that it is a computed, read-only field
