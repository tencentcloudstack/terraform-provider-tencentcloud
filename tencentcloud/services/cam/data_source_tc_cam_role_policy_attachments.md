Use this data source to query detailed information of CAM role policy attachments

Example Usage

```hcl
# query by role_id
data "tencentcloud_cam_role_policy_attachments" "foo" {
  role_id = tencentcloud_cam_role.foo.id
}

# query by role_id and policy_id
data "tencentcloud_cam_role_policy_attachments" "bar" {
  role_id   = tencentcloud_cam_role.foo.id
  policy_id = tencentcloud_cam_policy.foo.id
}
```