---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_backup_encryption_status"
sidebar_current: "docs-tencentcloud-resource-mysql_backup_encryption_status"
description: |-
  Provides a resource to create a mysql backup_encryption_status
---

# tencentcloud_mysql_backup_encryption_status

Provides a resource to create a mysql backup_encryption_status

## Example Usage

```hcl
resource "tencentcloud_mysql_backup_encryption_status" "backup_encryption_status" {
  instance_id       = "cdb-c1nl9rpv"
  encryption_status = "on"
}
```

## Argument Reference

The following arguments are supported:

* `encryption_status` - (Required, String) Whether physical backup encryption is enabled for the instance. Possible values are `on`, `off`.
* `instance_id` - (Required, String) Instance ID, in the format: cdb-XXXX. Same instance ID as displayed in the ApsaraDB for Console page.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mysql backup_encryption_status can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_backup_encryption_status.backup_encryption_status backup_encryption_status_id
```

