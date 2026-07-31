Use this data source to query detailed information of DB Custom cluster node resources from TencentCloud DBDC product.

Example Usage

Query dbdc db custom cluster node resources by cluster_id and node_ids

```hcl
data "tencentcloud_dbdc_db_custom_cluster_node_resources" "example" {
  cluster_id = "dbcc-b2arjlee"
  node_ids   = ["dbcn-vvkg2xls"]
}
```
