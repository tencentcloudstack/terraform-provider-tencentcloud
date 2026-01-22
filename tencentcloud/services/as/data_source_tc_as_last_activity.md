Use this data source to query detailed information of AS last activity

Example Usage

```hcl
data "tencentcloud_as_last_activity" "example" {
  auto_scaling_group_ids     = ["asg-3st9wq9m"]
  exclude_cancelled_activity = true
}
```