## ADDED Requirements

### Requirement: TEO Rule Engine returns RuleItems in read operation
The tencentcloud_teo_rule_engine resource SHALL return RuleItems parameter when reading the resource configuration from DescribeRules API.

#### Scenario: Read resource with RuleItems
- **WHEN** user reads tencentcloud_teo_rule_engine resource
- **THEN** the resource SHALL include RuleItems parameter with all rule items from the API response

### Requirement: RuleItems schema structure matches API response
The RuleItems parameter SHALL have a schema structure that matches the DescribeRules API response structure.

#### Scenario: RuleItems contains rule details
- **WHEN** RuleItems is populated from API response
- **THEN** each rule item SHALL include all required fields such as rule type, conditions, and actions

### Requirement: RuleItems is optional and backward compatible
The RuleItems parameter SHALL be optional and not break existing configurations.

#### Scenario: Existing resource without RuleItems
- **WHEN** existing resources are read
- **THEN** RuleItems SHALL be empty or null without breaking the resource read operation

#### Scenario: New resource with RuleItems
- **WHEN** resources are created after the change
- **THEN** resources SHALL include RuleItems parameter if available in API response
