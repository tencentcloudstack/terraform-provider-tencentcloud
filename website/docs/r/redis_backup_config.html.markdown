---
subcategory: "TencentDB for Redis"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_backup_config"
sidebar_current: "docs-tencentcloud-resource-redis_backup_config"
description: |-
  Use this resource to create a backup config of redis.
---

# tencentcloud_redis_backup_config

Use this resource to create a backup config of redis.

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

* `backup_time` - (Required, String) Specifys what time the backup action should take place. And the time interval should be one hour.
* `redis_id` - (Required, String, ForceNew) ID of a redis instance to which the policy will be applied.
* `backup_period` - (Optional, Set: [`String`], **Deprecated**) It has been deprecated from version 1.58.2. It makes no difference to online config at all Specifys which day the backup action should take place. Valid values: `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday` and `Sunday`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Redis  backup config can be imported, e.g.

```
$ terraform import tencentcloud_redis_backup_config.redisconfig redis-id
```

