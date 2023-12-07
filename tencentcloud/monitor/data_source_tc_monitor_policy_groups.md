Use this data source to query monitor policy groups (There is a lot of data and it is recommended to output to a file)

Example Usage

```hcl
data "tencentcloud_monitor_policy_groups" "groups" {
  policy_view_names = ["REDIS-CLUSTER", "cvm_device"]
}

data "tencentcloud_monitor_policy_groups" "name" {
  name = "test"
}
```