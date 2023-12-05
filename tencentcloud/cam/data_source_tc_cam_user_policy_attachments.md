Use this data source to query detailed information of CAM user policy attachments

Example Usage

```hcl
# query by user_id
data "tencentcloud_cam_user_policy_attachments" "foo" {
  user_id = tencentcloud_cam_user.foo.id
}

# query by user_id and policy_id
data "tencentcloud_cam_user_policy_attachments" "bar" {
  user_id   = tencentcloud_cam_user.foo.id
  policy_id = tencentcloud_cam_policy.foo.id
}
```