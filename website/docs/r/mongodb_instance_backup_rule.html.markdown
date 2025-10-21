---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_backup_rule"
sidebar_current: "docs-tencentcloud-resource-mongodb_instance_backup_rule"
description: |-
  Provides a resource to create mongodb instance backup rule
---

# tencentcloud_mongodb_instance_backup_rule

Provides a resource to create mongodb instance backup rule

## Example Usage

```hcl
resource "tencentcloud_mongodb_instance_backup" "backup_rule" {
  instance_id   = "cmgo-xxxxxx"
  backup_method = 0
  backup_time   = 10
}
```

## Argument Reference

The following arguments are supported:

* `backup_method` - (Required, Int) Set automatic backup method. Valid values:
- 0: Logical backup;
- 1: Physical backup;
- 3: Snapshot backup (supported only in cloud disk version).
* `backup_time` - (Required, Int) Set the start time for automatic backup. The value range is: [0,23]. For example, setting this parameter to 2 means that backup starts at 02:00.
* `instance_id` - (Required, String, ForceNew) Instance id.
* `backup_retention_period` - (Optional, Int) Specify the number of days to save backup data. The default is 7 days, and the support settings are 7, 30, 90, 180, 365.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mongodb instance backup rule can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance_backup.backup_rule ${instanceId}
```

