Provides a resource to create a cfw sg_rule

Example Usage

```hcl
resource "tencentcloud_sg_rule" "sg_rule" {
  data = {
  }
}
```

Import

cfw sg_rule can be imported using the id, e.g.

```
terraform import tencentcloud_sg_rule.sg_rule sg_rule_id
```
