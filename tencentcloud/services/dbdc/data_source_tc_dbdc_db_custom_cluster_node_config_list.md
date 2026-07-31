Use this data source to query the Kubernetes scheduling configuration (labels and taints) of nodes inside a DB Custom cluster from TencentCloud DBDC product.

Example Usage

Query dbdc db custom cluster node config list by cluster id and node ids

```hcl
data "tencentcloud_dbdc_db_custom_cluster_node_config_list" "example" {
  cluster_id = "dbcc-b2arjlee"
  node_ids   = ["dbcn-vvkg2xls"]
}
```
