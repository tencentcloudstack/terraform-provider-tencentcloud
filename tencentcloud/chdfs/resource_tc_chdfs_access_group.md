Provides a resource to create a chdfs access_group

Example Usage

```hcl
resource "tencentcloud_chdfs_access_group" "access_group" {
  access_group_name = "testAccessGroup"
  vpc_type          = 1
  vpc_id            = "vpc-4owdpnwr"
  description       = "test access group"
}
```

Import

chdfs access_group can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_access_group.access_group access_group_id
```