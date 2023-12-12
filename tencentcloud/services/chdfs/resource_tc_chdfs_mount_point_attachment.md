Provides a resource to create a chdfs mount_point_attachment

Example Usage

```hcl
resource "tencentcloud_chdfs_mount_point_attachment" "mount_point_attachment" {
  access_group_ids = [
    "ag-bvmzrbsm",
    "ag-lairqrgr",
  ]
  mount_point_id   = "f14mpfy5lh4e-KuiL"
}
```

Import

chdfs mount_point_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_mount_point_attachment.mount_point_attachment mount_point_id
```