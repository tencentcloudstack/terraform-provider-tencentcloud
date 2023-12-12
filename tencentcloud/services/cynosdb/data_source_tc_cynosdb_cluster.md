Use this data source to query detailed information of cynosdb cluster

Example Usage

```hcl
data "tencentcloud_cynosdb_cluster" "cluster" {
  cluster_id = "cynosdbmysql-bws8h88b"
  database   = "users"
  table      = "tb_user_name"
  table_type = "all"
}
```