Provides data source for TEO l7 acc rule

~> **NOTE:** This feature only supports sites in the plans of the Standard Edition and Enterprise Edition.

Example Usage

```hcl
data "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-36bjhygh1bxe"
}

output "rule_count" {
  description = "Total number of rules in the zone"
  value       = data.tencentcloud_teo_l7_acc_rule.example.total_count
}

output "rule_ids" {
  description = "List of rule IDs"
  value       = data.tencentcloud_teo_l7_acc_rule.example.rules[*].rule_id
}
```

Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `rules` - List of Layer 7 acceleration rules.
  * `rule_id` - Rule ID. Unique identifier of the rule.
  * `rule_name` - Rule name. The name length limit is 255 characters.
  * `description` - Rule annotation. Multiple annotations can be added.
  * `rule_priority` - Rule priority. Only used as an output parameter.
  * `branches` - Sub-Rule branch. This list currently supports filling in only one rule; multiple entries are invalid.
* `total_count` - Total number of Layer 7 acceleration rules matching the query criteria.
