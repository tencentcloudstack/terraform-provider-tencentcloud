Provides a resource to create a lighthouse apply_disk_backup

Example Usage

```hcl
resource "tencentcloud_lighthouse_apply_disk_backup" "apply_disk_backup" {
  disk_id = "lhdisk-xxxxxx"
  disk_backup_id = "lhbak-xxxxxx"
}
```