Use this data source to query detailed information of sqlserver upload_backup_info

Example Usage

```hcl
data "tencentcloud_sqlserver_upload_backup_info" "example" {
  instance_id         = "mssql-qelbzgwf"
  backup_migration_id = "mssql-backup-migration-8a0f3eht"
}
```