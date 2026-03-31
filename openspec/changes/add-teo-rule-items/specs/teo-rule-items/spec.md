## ADDED Requirements

### Requirement: RuleItems parameter support in schema
The tencentcloud_teo_rule_engine resource SHALL support a `rule_items` parameter in the schema. This parameter SHALL be a list of rule item objects, where each item contains conditions, actions, and priority configuration.

#### Scenario: Schema definition includes rule_items
- **WHEN** defining the resource schema
- **THEN** the schema SHALL include a `rule_items` attribute of type List
- **THEN** the `rule_items` attribute SHALL be Optional to maintain backward compatibility

#### Scenario: Rule item structure validation
- **WHEN** a rule item is defined in the configuration
- **THEN** the rule item SHALL contain a `conditions` list
- **THEN** the rule item SHALL contain an `actions` list
- **THEN** the rule item SHALL contain a `priority` integer field

### Requirement: Read RuleItems from DescribeRules API
The provider SHALL read RuleItems data from the DescribeRules API when refreshing the tencentcloud_teo_rule_engine resource state.

#### Scenario: Successfully read rule items from API
- **WHEN** the DescribeRules API is called
- **THEN** the provider SHALL extract RuleItems from the API response
- **THEN** each rule item SHALL be mapped to the corresponding Terraform state structure
- **THEN** the state SHALL be updated with the retrieved rule items

#### Scenario: Handle empty rule items
- **WHEN** the DescribeRules API returns an empty RuleItems list
- **THEN** the provider SHALL set the state rule_items to an empty list
- **THEN** no error SHALL be raised

### Requirement: Create resource with RuleItems
The provider SHALL support creating tencentcloud_teo_rule_engine resources with RuleItems specified in the configuration.

#### Scenario: Create resource with single rule item
- **WHEN** creating a tencentcloud_teo_rule_engine resource with one rule item
- **THEN** the provider SHALL send the rule item configuration to the cloud API
- **THEN** the API SHALL accept and store the rule item
- **THEN** the provider SHALL read back the created rule items to verify the state

#### Scenario: Create resource with multiple rule items
- **WHEN** creating a tencentcloud_teo_rule_engine resource with multiple rule items
- **THEN** the provider SHALL send all rule items to the cloud API
- **THEN** the order of rule items SHALL be preserved according to priority
- **THEN** all rule items SHALL be successfully created

### Requirement: Update RuleItems
The provider SHALL support updating RuleItems when the tencentcloud_teo_rule_engine resource configuration changes.

#### Scenario: Add new rule items
- **WHEN** new rule items are added to the resource configuration
- **THEN** the provider SHALL call the update API with the complete set of rule items
- **THEN** the new rule items SHALL be added to the existing rules
- **THEN** the provider SHALL verify the update by reading the rule items back

#### Scenario: Remove existing rule items
- **WHEN** rule items are removed from the resource configuration
- **THEN** the provider SHALL call the update API with the reduced set of rule items
- **THEN** the removed rule items SHALL be deleted from the cloud resource
- **THEN** the provider SHALL verify the removal by reading the rule items back

#### Scenario: Modify existing rule items
- **WHEN** existing rule items are modified in the resource configuration
- **THEN** the provider SHALL call the update API with the modified rule items
- **THEN** the modified rule items SHALL be updated in the cloud resource
- **THEN** the provider SHALL verify the modification by reading the rule items back

### Requirement: Delete resource with RuleItems
The provider SHALL support deleting tencentcloud_teo_rule_engine resources that contain RuleItems.

#### Scenario: Delete resource with rule items
- **WHEN** a tencentcloud_teo_rule_engine resource with rule items is deleted
- **THEN** the provider SHALL call the delete API
- **THEN** all associated rule items SHALL be deleted along with the resource
- **THEN** the resource and rule items SHALL be removed from state

### Requirement: Data source support for RuleItems
The provider SHALL support reading RuleItems from the tencentcloud_teo_rule_engine data source.

#### Scenario: Query rule items via data source
- **WHEN** querying the tencentcloud_teo_rule_engine data source
- **THEN** the data source SHALL return the rule_items list
- **THEN** each rule item SHALL contain all its properties (conditions, actions, priority)

### Requirement: Backward compatibility
The provider SHALL maintain backward compatibility with existing tencentcloud_teo_rule_engine configurations that do not use RuleItems.

#### Scenario: Existing configuration without rule_items
- **WHEN** an existing tencentcloud_teo_rule_engine configuration without rule_items is applied
- **THEN** the resource SHALL continue to function normally
- **THEN** no changes SHALL be made to the resource's rule configuration
- **THEN** no state migration SHALL be required

#### Scenario: State compatibility
- **WHEN** importing an existing resource state created before rule_items support
- **THEN** the import SHALL succeed without errors
- **THEN** the rule_items attribute SHALL be set to an empty list in state
- **THEN** the resource SHALL continue to operate normally
