---
subcategory: "CynosDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster"
sidebar_current: "docs-tencentcloud-resource-cynosdb_cluster"
description: |-
  Provide a resource to create a CynosDB cluster.
---

# tencentcloud_cynosdb_cluster

Provide a resource to create a CynosDB cluster.

## Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster" "foo" {
  available_zone               = "ap-guangzhou-4"
  vpc_id                       = "vpc-h70b6b49"
  subnet_id                    = "subnet-q6fhy1mi"
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  storage_limit                = 1000
  cluster_name                 = "tf-cynosdb"
  password                     = "cynos@123"
  instance_maintain_duration   = 7200
  instance_maintain_start_time = 10800
  instance_maintain_weekdays = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  instance_cpu_core    = 1
  instance_memory_size = 2

  param_item {
    name          = "character_set_server"
    current_value = "utf8mb4"
  }

  tags = {
    test = "test"
  }

  force_delete = false

  rw_group_sg = [
    "sg-ibyjkl6r",
  ]
  ro_group_sg = [
    "sg-ibyjkl6r",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `available_zone` - (Required, String, ForceNew) The available zone of the CynosDB Cluster.
* `cluster_name` - (Required, String, ForceNew) Name of CynosDB cluster.
* `db_type` - (Required, String, ForceNew) Type of CynosDB, and available values include `MYSQL`.
* `db_version` - (Required, String, ForceNew) Version of CynosDB, which is related to `db_type`. For `MYSQL`, available value is `5.7`.
* `instance_cpu_core` - (Required, Int) The number of CPU cores of read-write type instance in the CynosDB cluster. Note: modification of this field will take effect immediately, if want to upgrade on maintenance window, please upgrade from console.
* `instance_memory_size` - (Required, Int) Memory capacity of read-write type instance, unit in GB. Note: modification of this field will take effect immediately, if want to upgrade on maintenance window, please upgrade from console.
* `password` - (Required, String, ForceNew) Password of `root` account.
* `subnet_id` - (Required, String, ForceNew) ID of the subnet within this VPC.
* `vpc_id` - (Required, String, ForceNew) ID of the VPC.
* `auto_renew_flag` - (Optional, Int, ForceNew) Auto renew flag. Valid values are `0`(MANUAL_RENEW), `1`(AUTO_RENEW). Default value is `0`. Only works for PREPAID cluster.
* `charge_type` - (Optional, String, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. Default value is `POSTPAID_BY_HOUR`.
* `force_delete` - (Optional, Bool) Indicate whether to delete cluster instance directly or not. Default is false. If set true, the cluster and its `All RELATED INSTANCES` will be deleted instead of staying recycle bin. Note: works for both `PREPAID` and `POSTPAID_BY_HOUR` cluster.
* `instance_maintain_duration` - (Optional, Int) Duration time for maintenance, unit in second. `3600` by default.
* `instance_maintain_start_time` - (Optional, Int) Offset time from 00:00, unit in second. For example, 03:00am should be `10800`. `10800` by default.
* `instance_maintain_weekdays` - (Optional, Set: [`String`]) Weekdays for maintenance. `["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]` by default.
* `param_items` - (Optional, List) Specify parameter list of database. Use `data.tencentcloud_mysql_default_params` to query available parameter details.
* `port` - (Optional, Int, ForceNew) Port of CynosDB cluster.
* `prepaid_period` - (Optional, Int, ForceNew) The tenancy (time unit is month) of the prepaid instance. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`. NOTE: it only works when charge_type is set to `PREPAID`.
* `project_id` - (Optional, Int, ForceNew) ID of the project. `0` by default.
* `ro_group_sg` - (Optional, List: [`String`]) IDs of security group for `ro_group`.
* `rw_group_sg` - (Optional, List: [`String`]) IDs of security group for `rw_group`.
* `storage_limit` - (Optional, Int, ForceNew) Storage limit of CynosDB cluster instance, unit in GB. The maximum storage of a non-serverless instance in GB. NOTE: If db_type is `MYSQL` and charge_type is `PREPAID`, the value cannot exceed the maximum storage corresponding to the CPU and memory specifications, when charge_type is `POSTPAID_BY_HOUR`, this argument is unnecessary.
* `tags` - (Optional, Map) The tags of the CynosDB cluster.

The `param_items` object supports the following:

* `current_value` - (Required, String) Param expected value to set.
* `name` - (Required, String) Name of param, e.g. `character_set_server`.
* `old_value` - (Optional, String) Param old value, indicates the value which already set, this value is required when modifying current_value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `charset` - Charset used by CynosDB cluster.
* `cluster_status` - Status of the Cynosdb cluster.
* `create_time` - Creation time of the CynosDB cluster.
* `instance_id` - ID of instance.
* `instance_name` - Name of instance.
* `instance_status` - Status of the instance.
* `instance_storage_size` - Storage size of the instance, unit in GB.
* `ro_group_addr` - Readonly addresses. Each element contains the following attributes:
  * `ip` - IP address for readonly connection.
  * `port` - Port number for readonly connection.
* `ro_group_id` - ID of read-only instance group.
* `ro_group_instances` - List of instances in the read-only instance group.
  * `instance_id` - ID of instance.
  * `instance_name` - Name of instance.
* `rw_group_addr` - Read-write addresses. Each element contains the following attributes:
  * `ip` - IP address for read-write connection.
  * `port` - Port number for read-write connection.
* `rw_group_id` - ID of read-write instance group.
* `rw_group_instances` - List of instances in the read-write instance group.
  * `instance_id` - ID of instance.
  * `instance_name` - Name of instance.
* `storage_used` - Used storage of CynosDB cluster, unit in MB.


## Import

CynosDB cluster can be imported using the id, e.g.

```
$ terraform import tencentcloud_cynosdb_cluster.foo cynosdbmysql-dzj5l8gz
```

