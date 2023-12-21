Provides a resource to create a mongodb instance_backup_download_task

Example Usage

```hcl
resource "tencentcloud_mongodb_instance_backup_download_task" "instance_backup_download_task" {
  instance_id = "cmgo-b43i3wkj"
  backup_name = "cmgo-b43i3wkj_2023-05-09 14:54"
  backup_sets {
    replica_set_id = "cmgo-b43i3wkj_0"
  }
}
```

Import

mongodb instance_backup_download_task can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance_backup_download_task.instance_backup_download_task instanceId#backupName
```
