Use this data source to query detailed information of cynosdb cluster_detail_databases

Example Usage

```hcl
data "tencentcloud_cynosdb_cluster_detail_databases" "cluster_detail_databases" {
  cluster_id = "cynosdbmysql-bws8h88b"
  db_name    = "users"
}
```