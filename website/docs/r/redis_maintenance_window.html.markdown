---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_maintenance_window"
sidebar_current: "docs-tencentcloud-resource-redis_maintenance_window"
description: |-
  Provides a resource to create a redis maintenance_window
---

# tencentcloud_redis_maintenance_window

Provides a resource to create a redis maintenance_window

## Example Usage

```hcl
resource "tencentcloud_redis_maintenance_window" "maintenance_window" {
  instance_id = "crs-c1nl9rpv"
  start_time  = "17:00"
  end_time    = "19:00"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) The end time of the maintenance window, e.g. 19:00.
* `instance_id` - (Required, String) The ID of instance.
* `start_time` - (Required, String) Maintenance window start time, e.g. 17:00.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

redis maintenance_window can be imported using the id, e.g.

```
terraform import tencentcloud_redis_maintenance_window.maintenance_window maintenance_window_id
```

