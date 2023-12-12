Provides a resource to create a lighthouse disk_backup

Example Usage

```hcl
resource "tencentcloud_lighthouse_disk_backup" "disk_backup" {
  disk_id = "lhdisk-xxxxx"
  disk_backup_name = "disk-backup"
}
```

Import

lighthouse disk_backup can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_disk_backup.disk_backup disk_backup_id
```