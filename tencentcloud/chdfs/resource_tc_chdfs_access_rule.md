Provides a resource to create a chdfs access_rule

Example Usage

```hcl
resource "tencentcloud_chdfs_access_rule" "access_rule" {
  access_group_id = "ag-bvmzrbsm"

  access_rule {
    access_mode    = 2
    address        = "10.0.1.1"
    priority       = 12
  }
}
```

Import

chdfs access_rule can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_access_rule.access_rule access_group_id#access_rule_id
```