---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_backup_download_task"
sidebar_current: "docs-tencentcloud-resource-mongodb_instance_backup_download_task"
description: |-
  Provides a resource to create a mongodb instance_backup_download_task
---

# tencentcloud_mongodb_instance_backup_download_task

Provides a resource to create a mongodb instance_backup_download_task

## Example Usage

```hcl
resource "tencentcloud_mongodb_instance_backup_download_task" "instance_backup_download_task" {
  instance_id = "cmgo-b43i3wkj"
  backup_name = "cmgo-b43i3wkj_2023-05-09 14:54"
  backup_sets {
    replica_set_id = "cmgo-b43i3wkj_0"
  }
}
```

## Argument Reference

The following arguments are supported:

* `backup_name` - (Required, String, ForceNew) The name of the backup file to be downloaded can be obtained through the DescribeDBBackups interface.
* `backup_sets` - (Required, List, ForceNew) Specifies the node names of replica sets to download or a list of shard names for sharded clusters.For example, the replica set cmgo-p8vnipr5, example (fixed value): BackupSets.0=cmgo-p8vnipr5_0, the full amount of data can be downloaded.For example, the sharded cluster cmgo-p8vnipr5, for example: BackupSets.0=cmgo-p8vnipr5_0&amp;amp;BackupSets.1=cmgo-p8vnipr5_1, that is, to download the data of shard 0 and 1. If the sharded cluster needs to be downloaded in full, please pass in the example. Full slice name.
* `instance_id` - (Required, String, ForceNew) Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.

The `backup_sets` object supports the following:

* `replica_set_id` - (Required, String) Replication Id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateBackupDownloadTask](https://www.tencentcloud.com/document/api/240/40343)
* Read: [DescribeBackupDownloadTask](https://www.tencentcloud.com/document/api/240/40342)

## Import

mongodb instance_backup_download_task can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance_backup_download_task.instance_backup_download_task instanceId#backupName
```

