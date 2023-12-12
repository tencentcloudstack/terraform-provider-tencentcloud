Use this data source to query detailed information of CAM policies

Example Usage

```hcl
# query by policy_id
data "tencentcloud_cam_policies" "foo" {
  policy_id = tencentcloud_cam_policy.foo.id
}

# query by policy_id and name
data "tencentcloud_cam_policies" "bar" {
  policy_id = tencentcloud_cam_policy.foo.id
  name      = "tf-auto-test"
}
```