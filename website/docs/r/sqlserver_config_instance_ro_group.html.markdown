---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_config_instance_ro_group"
sidebar_current: "docs-tencentcloud-resource-sqlserver_config_instance_ro_group"
description: |-
  Provides a resource to create a sqlserver config_instance_ro_group
---

# tencentcloud_sqlserver_config_instance_ro_group

Provides a resource to create a sqlserver config_instance_ro_group

## Example Usage

```hcl
resource "tencentcloud_sqlserver_config_instance_ro_group" "config_instance_ro_group" {
  instance_id              = "mssql-ds1xhnt9"
  read_only_group_id       = "mssqlrg-cbya44fb"
  read_only_group_name     = "keep-ro-group-customize"
  is_offline_delay         = 1
  read_only_max_delay_time = 10
  min_read_only_in_group   = 1
  weight_pairs {
    read_only_instance_id = "mssqlro-o6dv2ugx"
    read_only_weight      = 50
  }
  auto_weight = 0
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `read_only_group_id` - (Required, String) Read-only group ID.
* `auto_weight` - (Optional, Int) 0-user-defined weight (adjusted according to WeightPairs), 1-system automatically assigns weight (WeightPairs is invalid), the default is 0.
* `balance_weight` - (Optional, Int) 0-do not rebalance the load, 1-rebalance the load, the default is 0.
* `is_offline_delay` - (Optional, Int) Whether to enable timeout culling function. 0- Disable the culling function. 1- Enable the culling function.
* `min_read_only_in_group` - (Optional, Int) After the timeout removal function is enabled, the number of read-only copies retained by the read-only group at least, if this parameter is not filled, it will not be modified.
* `read_only_group_name` - (Optional, String) Read-only group name. If this parameter is not specified, it is not modified.
* `read_only_max_delay_time` - (Optional, Int) After the timeout elimination function is enabled, the timeout threshold used, if this parameter is not filled, it will not be modified.
* `weight_pairs` - (Optional, List) Read-only group instance weight modification set, if this parameter is not filled, it will not be modified.

The `weight_pairs` object supports the following:

* `read_only_instance_id` - (Required, String) Read-only instance ID, in the format: mssqlro-3l3fgqn7.
* `read_only_weight` - (Required, Int) Read-only instance weight, the range is 0-100.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver config_instance_ro_group can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_instance_ro_group.config_instance_ro_group config_instance_ro_group_id
```

