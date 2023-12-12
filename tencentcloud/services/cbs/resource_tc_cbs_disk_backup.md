Provides a resource to create a cbs disk_backup.

~> **NOTE:** Backup quota must greater than 1.

Example Usage

```hcl

	resource "tencentcloud_cbs_disk_backup" "disk_backup" {
	  disk_id = "disk-xxx"
	  disk_backup_name = "xxx"
	}

```

Import

cbs disk_backup can be imported using the id, e.g.

```
terraform import tencentcloud_cbs_disk_backup.disk_backup disk_backup_id
```