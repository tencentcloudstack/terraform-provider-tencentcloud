Provides a resource to query TEO l7 acc rule

~> **NOTE:** This feature only supports the sites in the plans of Standard Edition and Enterprise Edition.

Example Usage

```hcl
data "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-36bjhygh1bxe"
}

output "total_count" {
  value = data.tencentcloud_teo_l7_acc_rule.example.total_count
}
```

Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) Zone id.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total_count` - (Integer) Total count of rules.
* `rules` - (List) Rules content.
  * `rule_id` - (String) Rule ID. Unique identifier of the rule.
  * `rule_name` - (String) Rule name. The name length limit is 255 characters.
  * `description` - (List) Rule annotation. multiple annotations can be added.
  * `rule_priority` - (Integer) Rule priority. only used as an output parameter.
  * `branches` - (List) Sub-Rule branch. this list currently supports filling in only one rule; multiple entries are invalid.
