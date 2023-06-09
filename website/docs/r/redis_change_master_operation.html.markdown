---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_change_master_operation"
sidebar_current: "docs-tencentcloud-resource-redis_change_master_operation"
description: |-
  Provides a resource to create a redis change_master_operation
---

# tencentcloud_redis_change_master_operation

Provides a resource to create a redis change_master_operation

## Example Usage

```hcl
resource "tencentcloud_redis_change_master_operation" "change_master_operation" {
  instance_id = "crs-c1nl9rpv"
  group_id    = "crs-rpl-c1nl9rpv"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, ForceNew) The ID of group.
* `instance_id` - (Required, String, ForceNew) The ID of instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



