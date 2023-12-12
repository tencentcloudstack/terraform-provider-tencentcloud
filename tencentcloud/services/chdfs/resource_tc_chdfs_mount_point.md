Provides a resource to create a chdfs mount_point

Example Usage

```hcl
resource "tencentcloud_chdfs_mount_point" "mount_point" {
  file_system_id     = "f14mpfy5lh4e"
  mount_point_name   = "terraform-test"
  mount_point_status = 1
}
```

Import

chdfs mount_point can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_mount_point.mount_point mount_point_id
```