---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_switck_master"
sidebar_current: "docs-tencentcloud-resource-redis_switck_master"
description: |-
  Provides a resource to create a redis switck_master
---

# tencentcloud_redis_switck_master

Provides a resource to create a redis switck_master

## Example Usage

```hcl
resource "tencentcloud_redis_switck_master" "switck_master" {
  instance_id = "crs-kfdkirid"
  group_id    = 29369
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of instance.
* `group_id` - (Optional, Int) Replication group ID, required for multi-AZ instances.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



