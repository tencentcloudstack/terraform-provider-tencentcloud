Use this data source to query detailed information of mariadb log_files

Example Usage

```hcl
data "tencentcloud_mariadb_log_files" "log_files" {
  instance_id = "tdsql-9vqvls95"
  type        = 1
}
```