---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_readonly_instance"
sidebar_current: "docs-tencentcloud-resource-cynosdb_readonly_instance"
description: |-
  Provide a resource to create a CynosDB readonly instance.
---

# tencentcloud_cynosdb_readonly_instance

Provide a resource to create a CynosDB readonly instance.

## Example Usage

```hcl
resource "tencentcloud_cynosdb_readonly_instance" "foo" {
  cluster_id           = cynosdbmysql-dzj5l8gz
  instance_name        = "tf-cynosdb-readonly-instance"
  force_delete         = true
  instance_cpu_core    = 2
  instance_memory_size = 4

  instance_maintain_duration   = 7200
  instance_maintain_start_time = 21600
  instance_maintain_weekdays = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID which the readonly instance belongs to.
* `instance_cpu_core` - (Required, Int) The number of CPU cores of read-write type instance in the CynosDB cluster. Note: modification of this field will take effect immediately, if want to upgrade on maintenance window, please upgrade from console.
* `instance_memory_size` - (Required, Int) Memory capacity of read-write type instance, unit in GB. Note: modification of this field will take effect immediately, if want to upgrade on maintenance window, please upgrade from console.
* `instance_name` - (Required, String, ForceNew) Name of instance.
* `force_delete` - (Optional, Bool) Indicate whether to delete readonly instance directly or not. Default is false. If set true, instance will be deleted instead of staying recycle bin. Note: works for both `PREPAID` and `POSTPAID_BY_HOUR` cluster.
* `instance_maintain_duration` - (Optional, Int) Duration time for maintenance, unit in second. `3600` by default.
* `instance_maintain_start_time` - (Optional, Int) Offset time from 00:00, unit in second. For example, 03:00am should be `10800`. `10800` by default.
* `instance_maintain_weekdays` - (Optional, Set: [`String`]) Weekdays for maintenance. `["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]` by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_status` - Status of the instance.
* `instance_storage_size` - Storage size of the instance, unit in GB.


## Import

CynosDB readonly instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_cynosdb_readonly_instance.foo cynosdbmysql-ins-dhwynib6
```

