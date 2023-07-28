---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_backup_operation"
sidebar_current: "docs-tencentcloud-resource-redis_backup_operation"
description: |-
  Provides a resource to create a redis backup_operation
---

# tencentcloud_redis_backup_operation

Provides a resource to create a redis backup_operation

## Example Usage

### Manually back up the Redis instance, and the backup data is kept for 7 days

```hcl
data "tencentcloud_mysql_instance" "foo" {}

resource "tencentcloud_redis_backup_operation" "backup_operation" {
  instance_id  = data.tencentcloud_mysql_instance.foo.instance_list[0].mysql_id
  remark       = "manually back"
  storage_days = 7
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The ID of instance.
* `remark` - (Optional, String, ForceNew) Notes information for the backup.
* `storage_days` - (Optional, Int, ForceNew) Number of days to store.0 specifies the default retention time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



