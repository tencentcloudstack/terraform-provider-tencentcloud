Use this data source to query the Kubernetes scheduling configuration (labels and taints) of nodes inside a DB Custom cluster from TencentCloud DBDC product.

Example Usage

Query dbdc db custom cluster node config list by cluster_id

```hcl
data "tencentcloud_dbdc_db_custom_cluster_node_config_list" "example" {
  cluster_id = "dbcc-nmtmsew8"
}
```

Query dbdc db custom cluster node config list by node_ids

```hcl
data "tencentcloud_dbdc_db_custom_cluster_node_config_list" "example" {
  cluster_id = "dbcc-nmtmsew8"
  node_ids   = ["node-abc123", "node-def456"]
}
```
