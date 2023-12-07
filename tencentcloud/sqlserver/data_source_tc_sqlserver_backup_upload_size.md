Use this data source to query detailed information of sqlserver datasource_backup_upload_size

Example Usage

```hcl
data "tencentcloud_sqlserver_backup_upload_size" "example" {
  instance_id         = "mssql-4gmc5805"
  backup_migration_id = "mssql-backup-migration-9tj0sxnz"
}
```