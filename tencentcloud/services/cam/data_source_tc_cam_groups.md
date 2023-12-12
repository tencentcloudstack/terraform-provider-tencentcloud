Use this data source to query detailed information of CAM groups

Example Usage

```hcl
# query by group_id
data "tencentcloud_cam_groups" "foo" {
  group_id = tencentcloud_cam_group.foo.id
}

# query by name
data "tencentcloud_cam_groups" "bar" {
  name = "cam-group-test"
}
```