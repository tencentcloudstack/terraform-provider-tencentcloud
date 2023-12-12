Use this data source to query detailed information of CAM group policy attachments

Example Usage

```hcl
# query by group_id
data "tencentcloud_cam_group_policy_attachments" "foo" {
  group_id = tencentcloud_cam_group.foo.id
}

# query by group_id and policy_id
data "tencentcloud_cam_group_policy_attachments" "bar" {
  group_id  = tencentcloud_cam_group.foo.id
  policy_id = tencentcloud_cam_policy.foo.id
}
```