Provides a resource to create a cls logset

Example Usage

```hcl
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "demo"
  tags = {
    "createdBy" = "terraform"
  }
}

```
Import

cls logset can be imported using the id, e.g.
```
$ terraform import tencentcloud_cls_logset.logset logset_id
```