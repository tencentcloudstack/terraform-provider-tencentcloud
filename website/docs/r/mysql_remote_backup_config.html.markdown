---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_remote_backup_config"
sidebar_current: "docs-tencentcloud-resource-mysql_remote_backup_config"
description: |-
  Provides a resource to create a mysql remote_backup_config
---

# tencentcloud_mysql_remote_backup_config

Provides a resource to create a mysql remote_backup_config

## Example Usage

```hcl
resource "tencentcloud_mysql_remote_backup_config" "remote_backup_config" {
  instance_id        = "cdb-c1nl9rpv"
  remote_backup_save = "on"
  remote_binlog_save = "on"
  remote_region      = ["ap-shanghai"]
  expire_days        = 7
}
```

## Argument Reference

The following arguments are supported:

* `expire_days` - (Required, Int) Remote backup retention time, in days.
* `instance_id` - (Required, String) Instance ID, in the format: cdb-c1nl9rpv. Same instance ID as displayed in the ApsaraDB for Console page.
* `remote_backup_save` - (Required, String) Remote data backup switch, off - disable remote backup, on - enable remote backup.
* `remote_binlog_save` - (Required, String) Off-site log backup switch, off - off off-site backup, on-on off-site backup, only when the parameter RemoteBackupSave is on, the RemoteBinlogSave parameter can be set to on.
* `remote_region` - (Required, Set: [`String`]) User settings off-site backup region list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mysql remote_backup_config can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_remote_backup_config.remote_backup_config remote_backup_config_id
```

