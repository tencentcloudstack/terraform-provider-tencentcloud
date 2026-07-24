## MODIFIED Requirements

### Requirement: No rule_ids computed attribute
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL NOT have a `rule_ids` computed attribute. The `rule_id` attribute provides the single rule ID managed by this resource.

#### Scenario: rule_ids is not in schema
- **WHEN** a user inspects the `tencentcloud_teo_l7_acc_rule_v2` resource schema
- **THEN** there SHALL NOT be a `rule_ids` attribute
