Use this data source to query detailed information of sqlserver datasource_backup_command

Example Usage

```hcl
data "tencentcloud_sqlserver_backup_commands" "example" {
  backup_file_type = "FULL"
  data_base_name   = "keep-publish-instance"
  is_recovery      = "NO"
}
```