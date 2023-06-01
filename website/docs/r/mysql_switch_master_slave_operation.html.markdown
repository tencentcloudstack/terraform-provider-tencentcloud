---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_switch_master_slave_operation"
sidebar_current: "docs-tencentcloud-resource-mysql_switch_master_slave_operation"
description: |-
  Provides a resource to create a mysql switch_master_slave_operation
---

# tencentcloud_mysql_switch_master_slave_operation

Provides a resource to create a mysql switch_master_slave_operation

## Example Usage

```hcl
resource "tencentcloud_mysql_switch_master_slave_operation" "switch_master_slave_operation" {
  instance_id  = "cdb-d9gbh7lt"
  dst_slave    = "first"
  force_switch = true
  wait_switch  = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) instance id.
* `dst_slave` - (Optional, String, ForceNew) target instance. Possible values: `first` - first standby; `second` - second standby. The default value is `first`, and only multi-AZ instances support setting it to `second`.
* `force_switch` - (Optional, Bool, ForceNew) Whether to force switch. Default is False. Note that if you set the mandatory switch to True, there is a risk of data loss on the instance, so use it with caution.
* `wait_switch` - (Optional, Bool, ForceNew) Whether to switch within the time window. The default is False, i.e. do not switch within the time window. Note that if the ForceSwitch parameter is set to True, this parameter will not take effect.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mysql switch_master_slave_operation can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_switch_master_slave_operation.switch_master_slave_operation switch_master_slave_operation_id
```

