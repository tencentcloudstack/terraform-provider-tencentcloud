---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_backup_time"
sidebar_current: "docs-tencentcloud-resource-mariadb_backup_time"
description: |-
  Provides a resource to create a mariadb backup_time
---

# tencentcloud_mariadb_backup_time

Provides a resource to create a mariadb backup_time

## Example Usage

```hcl
resource "tencentcloud_mariadb_backup_time" "backup_time" {
  instance_id       = "tdsql-9vqvls95"
  start_backup_time = "01:00"
  end_backup_time   = "04:00"
}
```

## Argument Reference

The following arguments are supported:

* `end_backup_time` - (Required, String) End time of daily backup window in the format of `mm:ss`, such as 23:59.
* `instance_id` - (Required, String) instance id.
* `start_backup_time` - (Required, String) Start time of daily backup window in the format of `mm:ss`, such as 22:00.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mariadb backup_time can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_backup_time.backup_time backup_time_id
```

