Use this data source to query detailed information of DB Custom cluster nodes from TencentCloud DBDC product.

Example Usage

Query dbdc db custom cluster nodes by cluster_id

```hcl
data "tencentcloud_dbdc_db_custom_cluster_nodes" "example" {
  cluster_id = "dbcc-nmtmsew8"
}
```

Query dbdc db custom cluster nodes by filters

```hcl
data "tencentcloud_dbdc_db_custom_cluster_nodes" "example" {
  cluster_id = "dbcc-nmtmsew8"

  filters {
    name   = "node-name"
    values = ["tf-example"]
  }
}
```
