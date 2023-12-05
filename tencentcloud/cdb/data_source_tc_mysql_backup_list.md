Use this data source to query the list of backup databases.

Example Usage

```hcl
data "tencentcloud_mysql_backup_list" "default" {
  mysql_id           = "terraform-test-local-database"
  max_number         = 10
  result_output_file = "mytestpath"
}
```