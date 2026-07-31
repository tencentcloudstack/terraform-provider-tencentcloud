Use this data source to query DB Custom cluster resource summary info from TencentCloud DBDC product.

Example Usage

Query dbdc db custom cluster resources by cluster_id

```hcl
data "tencentcloud_dbdc_db_custom_cluster_resources" "example" {
  cluster_id = "dbcc-b2arjlee"
}
```
