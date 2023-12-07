Use this data source to query detailed information of mysql databases

Example Usage

```hcl
data "tencentcloud_mysql_databases" "databases" {
  instance_id = "cdb-c1nl9rpv"
  database_regexp = ""
}
```