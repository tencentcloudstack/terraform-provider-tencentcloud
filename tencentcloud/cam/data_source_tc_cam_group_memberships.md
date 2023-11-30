Use this data source to query detailed information of CAM group memberships

Example Usage

```hcl
data "tencentcloud_cam_group_memberships" "foo" {
  group_id = tencentcloud_cam_group.foo.id
}
```