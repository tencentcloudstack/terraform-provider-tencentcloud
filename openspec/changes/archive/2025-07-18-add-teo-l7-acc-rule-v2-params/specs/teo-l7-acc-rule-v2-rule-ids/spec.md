## ADDED Requirements

### Requirement: rule_ids computed attribute
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL expose a `rule_ids` computed attribute of type list of strings, populated from the `CreateL7AccRules` API response field `RuleIds`.

#### Scenario: rule_ids populated after resource creation
- **WHEN** the `tencentcloud_teo_l7_acc_rule_v2` resource is created successfully via the `CreateL7AccRules` API
- **THEN** the `rule_ids` attribute SHALL contain all rule IDs returned in `response.Response.RuleIds`

#### Scenario: rule_ids populated during resource read
- **WHEN** the `tencentcloud_teo_l7_acc_rule_v2` resource is read via the `DescribeL7AccRules` API
- **THEN** the `rule_ids` attribute SHALL be populated from the `RuleId` field of each `RuleEngineItem` in the response

#### Scenario: rule_ids is computed only
- **WHEN** a user writes a Terraform configuration for `tencentcloud_teo_l7_acc_rule_v2`
- **THEN** the `rule_ids` attribute SHALL NOT be settable by the user (Computed: true, not Optional or Required)

### Requirement: rule_ids schema definition
The `rule_ids` field SHALL be defined in the resource schema as a TypeList with TypeString elements, marked as Computed.

#### Scenario: schema definition correctness
- **WHEN** the resource schema is defined
- **THEN** `rule_ids` SHALL have `Type: schema.TypeList`, `Computed: true`, and `Elem: &schema.Schema{Type: schema.TypeString}`
