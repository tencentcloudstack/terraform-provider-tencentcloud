Use this data source to query policy group binding objects.

Example Usage

```hcl
data "tencentcloud_monitor_policy_groups" "name" {
  name = "test"
}

data "tencentcloud_monitor_binding_objects" "objects" {
  group_id = data.tencentcloud_monitor_policy_groups.name.list[0].group_id
}
```