Provides a resource to create a FWM rule group

Example Usage

```hcl
resource "tencentcloud_fwm_rule_group" "example" {
  group_name = "tf-example"
  product    = "enterprise_sg"
}
```

Import

FWM rule group can be imported using the id, e.g.

```
terraform import tencentcloud_fwm_rule_group.example fwmrg-wsrrhuvrib
```
