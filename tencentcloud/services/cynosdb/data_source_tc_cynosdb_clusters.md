Use this data source to query detailed information of Cynosdb clusters.

Example Usage

```hcl
data "tencentcloud_cynosdb_clusters" "foo" {
  cluster_id   = "cynosdbmysql-dzj5l8gz"
  project_id   = 0
  db_type      = "MYSQL"
  cluster_name = "test"
}
```