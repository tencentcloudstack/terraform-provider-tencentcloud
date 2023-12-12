Provides a resource to create a chdfs life_cycle_rule

Example Usage

```hcl
resource "tencentcloud_chdfs_life_cycle_rule" "life_cycle_rule" {
  file_system_id = "f14mpfy5lh4e"

  life_cycle_rule {
    life_cycle_rule_name = "terraform-test"
    path                 = "/test"
    status               = 1

    transitions {
      days = 30
      type = 1
    }
  }
}
```

Import

chdfs life_cycle_rule can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_life_cycle_rule.life_cycle_rule file_system_id#life_cycle_rule_id
```