Use this data source to get information about a parameter group of a database instance.

Example Usage

```hcl
data "tencentcloud_mysql_parameter_list" "mysql" {
  mysql_id           = "terraform-test-local-database"
  engine_version     = "5.5"
  result_output_file = "mytestpath"
}
```