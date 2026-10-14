## ADDED Requirements

### Requirement: CreateL7AccRules parameter mapping
The `tencentcloud_teo_l7_acc_rule_v2` resource's Create function SHALL map the `zone_id` schema parameter to `request.ZoneId` and the rule fields (`status`, `rule_name`, `description`, `branches`) to a `RuleEngineItem` within `request.Rules` when calling the `CreateL7AccRules` API.

#### Scenario: Create maps zone_id to request.ZoneId
- **WHEN** a user creates a `tencentcloud_teo_l7_acc_rule_v2` resource with `zone_id = "zone-abc123"`
- **THEN** the `CreateL7AccRules` API SHALL be called with `request.ZoneId = "zone-abc123"`

#### Scenario: Create maps rule fields to request.Rules
- **WHEN** a user creates a resource with `status = "enable"`, `rule_name = "test"`, `description = ["desc"]`, `branches = [...]`
- **THEN** the `CreateL7AccRules` API SHALL be called with `request.Rules` containing a `RuleEngineItem` with the corresponding `Status`, `RuleName`, `Description`, and `Branches` fields

### Requirement: ModifyL7AccRule parameter mapping
The `tencentcloud_teo_l7_acc_rule_v2` resource's Update function SHALL map the `zone_id` schema parameter to `request.ZoneId` and the rule fields to `request.Rule` (a `RuleEngineItem`) when calling the `ModifyL7AccRule` API.

#### Scenario: Update maps zone_id to request.ZoneId
- **WHEN** a user updates a `tencentcloud_teo_l7_acc_rule_v2` resource
- **THEN** the `ModifyL7AccRule` API SHALL be called with `request.ZoneId` set from the composite ID's zone_id part

#### Scenario: Update maps rule_id to request.Rule.RuleId
- **WHEN** a user updates a resource with composite ID `zone-abc#rule-xyz`
- **THEN** the `ModifyL7AccRule` API SHALL be called with `request.Rule.RuleId = "rule-xyz"`

#### Scenario: Update maps mutable fields to request.Rule
- **WHEN** a user changes `status`, `rule_name`, `description`, or `branches`
- **THEN** the `ModifyL7AccRule` API SHALL be called with `request.Rule` containing the updated `Status`, `RuleName`, `Description`, and/or `Branches` fields

### Requirement: DeleteL7AccRules parameter mapping
The `tencentcloud_teo_l7_acc_rule_v2` resource's Delete function SHALL map the `zone_id` to `request.ZoneId` and `rule_id` to `request.RuleIds` when calling the `DeleteL7AccRules` API.

#### Scenario: Delete maps zone_id to request.ZoneId
- **WHEN** a user deletes a `tencentcloud_teo_l7_acc_rule_v2` resource with composite ID `zone-abc#rule-xyz`
- **THEN** the `DeleteL7AccRules` API SHALL be called with `request.ZoneId = "zone-abc"`

#### Scenario: Delete maps rule_id to request.RuleIds
- **WHEN** a user deletes a resource with composite ID `zone-abc#rule-xyz`
- **THEN** the `DeleteL7AccRules` API SHALL be called with `request.RuleIds = ["rule-xyz"]`

### Requirement: DescribeL7AccRules parameter mapping
The `tencentcloud_teo_l7_acc_rule_v2` resource's Read function SHALL use `Filters` with `Name: "rule-id"` and `Values: [ruleId]` derived from the `rule_id` parameter when calling the `DescribeL7AccRules` API.

#### Scenario: Read uses Filters with rule-id
- **WHEN** a user reads a `tencentcloud_teo_l7_acc_rule_v2` resource with composite ID `zone-abc#rule-xyz`
- **THEN** the `DescribeL7AccRules` API SHALL be called with `request.ZoneId = "zone-abc"` and `request.Filters` containing `{Name: "rule-id", Values: ["rule-xyz"]}`
