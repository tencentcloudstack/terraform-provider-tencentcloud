Use this data source to query placement groups.

Example Usage

```hcl
data "tencentcloud_placement_groups" "foo" {
  placement_group_id = "ps-21q9ibvr"
  name               = "test"
}
```