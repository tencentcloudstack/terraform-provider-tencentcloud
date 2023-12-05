Use this data source to query detailed information of as last_activity

Example Usage

```hcl
data "tencentcloud_as_last_activity" "last_activity" {
  auto_scaling_group_ids = ["asc-lo0b94oy"]
}
```