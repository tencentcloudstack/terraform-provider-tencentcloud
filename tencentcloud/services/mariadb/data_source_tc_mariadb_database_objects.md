Use this data source to query detailed information of mariadb database_objects

Example Usage

```hcl
data "tencentcloud_mariadb_database_objects" "database_objects" {
	instance_id = "tdsql-n2fw7pn3"
	db_name = "mysql"
}
```