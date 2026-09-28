Use this data source to query detailed information of cat node groups

Example Usage

```hcl
data "tencentcloud_cat_node_groups" "groups" {
  node_type        = [1]
  ip_type          = 1
  node_group_type  = 1
}
```
