Use this data source to query detailed information of cynosdb proxy_node

Example Usage

```hcl
data "tencentcloud_cynosdb_proxy_node" "proxy_node" {
  order_by      = "CREATETIME"
  order_by_type = "DESC"
  filters {
    names       = "ClusterId"
    values      = "cynosdbmysql-cgd2gpwr"
    exact_match = false
    name        = "ClusterId"
  }
}
```