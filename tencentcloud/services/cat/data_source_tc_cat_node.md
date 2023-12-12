Use this data source to query detailed information of cat node

Example Usage

```hcl
data "tencentcloud_cat_node" "node"{
  node_type = 1
  location = 2
  is_ipv6 = false
}
```