---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_backup_config"
sidebar_current: "docs-tencentcloud-resource-redis_backup_config"
description: |-
  Use this data source to query which instance types of Redis are available in a specific region.
---

# tencentcloud_redis_backup_config

Use this data source to query which instance types of Redis are available in a specific region.

## Example Usage

```hcl
resource "tencentcloud_redis_backup_config" "redislab" {
  redis_id      = "crs-7yl0q0dd"
  backup_time   = "04:00-05:00"
  backup_period = ["Monday"]
}
```

## Argument Reference

The following arguments are supported:

* `backup_period` - (Required) Specifys which day the backup action should take place. Supported values include: Monday，Tuesday, Wednesday, Thursday, Friday, Saturday and Sunday.
* `backup_time` - (Required) Specifys what time the backup action should take place.
* `redis_id` - (Required, ForceNew) ID of a Redis instance to which the policy will be applied.


## Import

Redis  backup config can be imported, e.g.

```hcl
$ terraform import tencentcloud_redis_backup_config.redisconfig redis-id
```

