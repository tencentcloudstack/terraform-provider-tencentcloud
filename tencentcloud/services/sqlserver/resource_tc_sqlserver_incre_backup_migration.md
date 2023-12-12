Provides a resource to create a sqlserver incre_backup_migration

Example Usage

```hcl
resource "tencentcloud_sqlserver_incre_backup_migration" "example" {
  instance_id         = "mssql-4gmc5805"
  backup_migration_id = "mssql-backup-migration-9tj0sxnz"
  backup_files        = []
  is_recovery         = "YES"
}
```

Import

sqlserver incre_backup_migration can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_incre_backup_migration.incre_backup_migration incre_backup_migration_id
```