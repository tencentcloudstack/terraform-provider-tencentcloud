Provides a resource to create a rum whitelist

Example Usage

```hcl
resource "tencentcloud_rum_whitelist" "whitelist" {
  instance_id = "rum-pasZKEI3RLgakj"
  remark = "white list remark"
  whitelist_uin = "20221122"
  # aid = ""
}

```
Import

rum whitelist can be imported using the id, e.g.
```
$ terraform import tencentcloud_rum_whitelist.whitelist whitelist_id
```