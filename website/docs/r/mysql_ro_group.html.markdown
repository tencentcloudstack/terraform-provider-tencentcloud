---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_ro_group"
sidebar_current: "docs-tencentcloud-resource-mysql_ro_group"
description: |-
  Provides a resource to create a mysql ro_group
---

# tencentcloud_mysql_ro_group

Provides a resource to create a mysql ro_group

## Example Usage

```hcl
resource "tencentcloud_mysql_ro_group" "ro_group" {
  instance_id = "cdb-e8i766hx"
  ro_group_id = "cdbrg-f49t0gnj"
  ro_group_info {
    ro_group_name          = "keep-ro"
    ro_max_delay_time      = 1
    ro_offline_delay       = 1
    min_ro_in_group        = 1
    weight_mode            = "custom"
    replication_delay_time = 1
  }
  ro_weight_values {
    instance_id = "cdbro-f49t0gnj"
    weight      = 10
  }
  is_balance_ro_load = 1
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID, in the format: cdbro-3i70uj0k.
* `ro_group_id` - (Required, String) The ID of the RO group.
* `is_balance_ro_load` - (Optional, Int) Whether to rebalance the load of RO instances in the RO group. Supported values include: 1 - rebalance load; 0 - do not rebalance load. The default value is 0. Note that when it is set to rebalance the load, the RO instance in the RO group will have a momentary disconnection of the database connection, please ensure that the application can reconnect to the database.
* `ro_group_info` - (Optional, List) Details of the RO group.
* `ro_weight_values` - (Optional, List) The weight of the instance within the RO group. If the weight mode of the RO group is changed to user-defined mode (custom), this parameter must be set, and the weight value of each RO instance needs to be set.

The `ro_group_info` object supports the following:

* `min_ro_in_group` - (Optional, Int) The minimum number of reserved instances. It can be set to any value less than or equal to the number of RO instances under this RO group. Note that if the setting value is greater than the number of RO instances, it will not be removed; if it is set to 0, all instances whose latency exceeds the limit will be removed.
* `replication_delay_time` - (Optional, Int) Delayed replication time.
* `ro_group_name` - (Optional, String) RO group name.
* `ro_max_delay_time` - (Optional, Int) RO instance maximum latency threshold. The unit is seconds, the minimum value is 1. Note that the RO group must have enabled instance delay culling policy for this value to be valid.
* `ro_offline_delay` - (Optional, Int) Whether to enable delayed culling of instances. Supported values are: 1 - on; 0 - not on. Note that if you enable instance delay culling, you must set the delay threshold (RoMaxDelayTime) parameter.
* `weight_mode` - (Optional, String) weight mode. Supported values include: `system` - automatically assigned by the system; `custom` - user-defined settings. Note that if the `custom` mode is set, the RO instance weight configuration (RoWeightValues) parameter must be set.

The `ro_weight_values` object supports the following:

* `instance_id` - (Required, String) RO instance ID.
* `weight` - (Required, Int) Weights. The value range is [0, 100].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



