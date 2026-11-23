Use this data source to query detailed information of dcdb database tables

Example Usage

```hcl
data "tencentcloud_dcdb_database_tables" "example" {
  instance_id = "dcdbt-ow7t8lmc"
  db_name     = "tf_test_db"
  table       = "tf_test_table"
}
```
