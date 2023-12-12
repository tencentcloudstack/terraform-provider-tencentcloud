Provides a resource to create a CFS access group.

Example Usage

```hcl
resource "tencentcloud_cfs_access_group" "example" {
  name        = "tx_example"
  description = "desc."
}
```

Import

CFS access group can be imported using the id, e.g.

```
$ terraform import tencentcloud_cfs_access_group.example pgroup-7nx89k7l
```