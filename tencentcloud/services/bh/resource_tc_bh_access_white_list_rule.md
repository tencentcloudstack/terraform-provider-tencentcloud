Provides a resource to create a BH access white list rule

Example Usage

```hcl
resource "tencentcloud_bh_access_white_list_rule" "example" {
  source = "1.1.1.1"
  remark = "remark."
}
```

Import

BH access white list rule can be imported using the id, e.g.

```
terraform import tencentcloud_bh_access_white_list_rule.example 1235
```
