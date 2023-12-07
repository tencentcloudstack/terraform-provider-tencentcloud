Provides a resource to create a CAM group.

Example Usage

```hcl
resource "tencentcloud_cam_group" "foo" {
  name   = "tf_cam_group"
  remark = "tf_group_remark"
}
```

Import

CAM group can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_group.foo 90496
```