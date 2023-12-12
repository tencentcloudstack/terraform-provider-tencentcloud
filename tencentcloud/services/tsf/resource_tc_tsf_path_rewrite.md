Provides a resource to create a tsf path_rewrite

Example Usage

```hcl
resource "tencentcloud_tsf_path_rewrite" "path_rewrite" {
  gateway_group_id = "group-a2j9zxpv"
  regex = "/test"
  replacement = "/tt"
  blocked = "N"
  order = 2
}
```

Import

tsf path_rewrite can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_path_rewrite.path_rewrite rewrite-nygq33v2
```