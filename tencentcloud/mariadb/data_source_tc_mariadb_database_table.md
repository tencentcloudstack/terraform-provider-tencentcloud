Use this data source to query detailed information of mariadb database_table

Example Usage

```hcl
data "tencentcloud_mariadb_database_table" "database_table" {
  instance_id = "tdsql-e9tklsgz"
  db_name = "mysql"
  table = "server_cost"
}
```