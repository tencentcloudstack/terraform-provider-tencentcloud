---
subcategory: "TDSQL Boundless(tdmysql)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmysql_dbs_backup"
sidebar_current: "docs-tencentcloud-resource-tdmysql_dbs_backup"
description: |-
  Provides a resource to create a TDSQL Boundless(tdmysql) manual backup set.
---

# tencentcloud_tdmysql_dbs_backup

Provides a resource to create a TDSQL Boundless(tdmysql) manual backup set.

## Example Usage

```hcl
resource "tencentcloud_tdmysql_dbs_backup" "example" {
  instance_id   = "tdsql3-ibpyogg7tg"
  backup_method = "physical"
  backup_type   = "full"
  backup_name   = "tf-example"

  timeouts {
    create = "60m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `backup_method` - (Required, String) Backup method. Valid values: `physical` (physical backup), `snapshot` (snapshot backup). This value should stay consistent with the `backupMethod` returned by the DescribeDBSBackupPolicy API.
* `backup_type` - (Required, String) Backup type. Only `full` is supported for now.
* `instance_id` - (Required, String) Instance ID.
* `backup_name` - (Optional, String) Backup comment. Modifiable through the ModifyDBSBackupSetComment API. If omitted, the cloud side generates one automatically.
* `timeouts` - (Optional, Object) The timeouts block allows you to specify the timeout for the create operation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource, in the form of the backup set ID.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateDBSBackup](https://www.tencentcloud.com/document/api/1292/78258)
  - [DescribeDBSBackupSets](https://www.tencentcloud.com/document/api/1292/80116)
* Read: [DescribeDBSBackupSets](https://www.tencentcloud.com/document/api/1292/80116)
* Update:
  - [DescribeDBSBackupSets](https://www.tencentcloud.com/document/api/1292/80116)
  - [ModifyDBSBackupSetComment](https://www.tencentcloud.com/document/api/1292/78275)
* Delete: [DeleteDBSBackupSets](https://www.tencentcloud.com/document/api/1292/78277)

## Import

TDSQL Boundless(tdmysql) manual backup set can be imported using the id, e.g.

```
terraform import tencentcloud_tdmysql_dbs_backup.example 94403665169021
```

