Provides a resource to rollback cbs disk backup.

Example Usage

```hcl
resource "tencentcloud_cbs_disk_backup_rollback_operation" "operation" {
  disk_backup_id  = "dbp-xxx"
  disk_id = "disk-xxx"
}
```