Use this data source to query detailed information of DB Custom cluster node resources from TencentCloud DBDC product.

Example Usage

Query dbdc db custom cluster node resources by cluster_id

```hcl
data "tencentcloud_dbdc_db_custom_cluster_node_resources" "example" {
  cluster_id = "dbcc-nmtmsew8"
}
```

Query dbdc db custom cluster node resources by cluster_id and node_ids

```hcl
data "tencentcloud_dbdc_db_custom_cluster_node_resources" "example" {
  cluster_id = "dbcc-nmtmsew8"

  node_ids = ["node-abc123", "node-def456"]
}
```
