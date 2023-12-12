Provides a resource to create a tsf enable_unit_rule

Example Usage

```hcl
resource "tencentcloud_tsf_enable_unit_rule" "enable_unit_rule" {
  rule_id = "unit-rl-is9m4nxz"
  switch = "enabled"
}
```

Import

tsf enable_unit_rule can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_enable_unit_rule.enable_unit_rule enable_unit_rule_id
```