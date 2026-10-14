---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_time_window"
sidebar_current: "docs-tencentcloud-resource-postgresql_time_window"
description: |-
  Provides a resource to create a postgres postgresql_time_window
---

# tencentcloud_postgresql_time_window

Provides a resource to create a postgres postgresql_time_window

## Example Usage

```hcl
resource "tencentcloud_postgresql_time_window" "postgresql_time_window" {
  db_instance_id      = "postgres-45b0vlmr"
  maintain_duration   = 2
  maintain_start_time = "04:00"
  maintain_week_days = [
    "friday",
    "monday",
    "saturday",
    "sunday",
    "thursday",
    "tuesday",
    "wednesday",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Instance id.
* `maintain_duration` - (Optional, Int) Maintenance duration, Unit: hours.
* `maintain_start_time` - (Optional, String) Maintenance start time. Time zone is UTC+8.
* `maintain_week_days` - (Optional, Set: [`String`]) Maintenance cycle.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyMaintainTimeWindow](https://www.tencentcloud.com/document/product/409/78801) - This API is used to modify the instance maintenance window.
* Read: [DescribeMaintainTimeWindow](https://www.tencentcloud.com/document/product/409/78803) - This API is used to query instance maintenance window.
* Update: [ModifyMaintainTimeWindow](https://www.tencentcloud.com/document/product/409/78801) - This API is used to modify the instance maintenance window.

## Import

postgres postgresql_time_window can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_time_window.postgresql_time_window instance_id
```

