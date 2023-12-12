Use this data source to query detailed information of CAM roles

Example Usage

```hcl
# query by role_id
data "tencentcloud_cam_roles" "foo" {
  role_id = tencentcloud_cam_role.foo.id
}

# query by name
data "tencentcloud_cam_roles" "bar" {
  name = "cam-role-test"
}
```