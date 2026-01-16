---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster"
sidebar_current: "docs-tencentcloud-resource-cynosdb_cluster"
description: |-
  Provide a resource to create a CynosDB cluster.
---

# tencentcloud_cynosdb_cluster

Provide a resource to create a CynosDB cluster.

~> **NOTE:** params `instance_count` and `instance_init_infos` only choose one. If neither parameter is set, the CynosDB cluster is created with parameter `instance_count` set to `2` by default(one RW instance + one Ro instance). If you only need to create a master instance, explicitly set the `instance_count` field to `1`, or configure the RW instance information in the `instance_init_infos` field.

## Example Usage

### Create a single availability zone NORMAL CynosDB cluster

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

# create cynosdb cluster
resource "tencentcloud_cynosdb_cluster" "example" {
  available_zone               = var.availability_zone
  vpc_id                       = tencentcloud_vpc.vpc.id
  subnet_id                    = tencentcloud_subnet.subnet.id
  db_mode                      = "NORMAL"
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  port                         = 3306
  cluster_name                 = "tf-example"
  password                     = "cynosDB@123"
  instance_maintain_duration   = 7200
  instance_maintain_start_time = 10800
  instance_cpu_core            = 2
  instance_memory_size         = 4
  force_delete                 = false
  instance_maintain_weekdays = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  param_items {
    name          = "character_set_server"
    current_value = "utf8mb4"
  }

  param_items {
    name          = "lower_case_table_names"
    current_value = "0"
  }

  rw_group_sg = [
    tencentcloud_security_group.example.id,
  ]

  ro_group_sg = [
    tencentcloud_security_group.example.id,
  ]

  instance_init_infos {
    cpu            = 2
    memory         = 4
    instance_type  = "rw"
    instance_count = 1
    device_type    = "common"
  }

  instance_init_infos {
    cpu            = 2
    memory         = 4
    instance_type  = "ro"
    instance_count = 1
    device_type    = "exclusive"
  }

  cynos_version = "2.1.14.001"

  tags = {
    createBy = "terraform"
  }
}
```

### Create a multiple availability zone SERVERLESS CynosDB cluster

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

variable "slave_zone" {
  default = "ap-guangzhou-6"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

# create param template
resource "tencentcloud_cynosdb_param_template" "example" {
  db_mode              = "SERVERLESS"
  engine_version       = "8.0"
  template_name        = "tf-example"
  template_description = "terraform-template"

  param_list {
    current_value = "-1"
    param_name    = "optimizer_trace_offset"
  }
}

# create cynosdb cluster
resource "tencentcloud_cynosdb_cluster" "example" {
  available_zone               = var.availability_zone
  slave_zone                   = var.slave_zone
  vpc_id                       = tencentcloud_vpc.vpc.id
  subnet_id                    = tencentcloud_subnet.subnet.id
  db_mode                      = "SERVERLESS"
  db_type                      = "MYSQL"
  db_version                   = "8.0"
  port                         = 3306
  cluster_name                 = "tf-example"
  password                     = "cynosDB@123"
  instance_maintain_duration   = 7200
  instance_maintain_start_time = 10800
  min_cpu                      = 2
  max_cpu                      = 4
  param_template_id            = tencentcloud_cynosdb_param_template.example.template_id
  force_delete                 = false
  instance_maintain_weekdays = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  rw_group_sg = [
    tencentcloud_security_group.example.id,
  ]

  ro_group_sg = [
    tencentcloud_security_group.example.id,
  ]

  tags = {
    createBy = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `available_zone` - (Required, String, ForceNew) The available zone of the CynosDB Cluster.
* `cluster_name` - (Required, String) Name of CynosDB cluster.
* `db_type` - (Required, String, ForceNew) Type of CynosDB, and available values include `MYSQL`.
* `db_version` - (Required, String, ForceNew) Version of CynosDB, which is related to `db_type`. For `MYSQL`, available value is `5.7`, `8.0`.
* `password` - (Required, String) Password of `root` account.
* `subnet_id` - (Required, String) ID of the subnet within this VPC.
* `vpc_id` - (Required, String) ID of the VPC.
* `auto_pause_delay` - (Optional, Int) Specify auto-pause delay in second while `db_mode` is `SERVERLESS`. Value range: `[600, 691200]`. Default: `600`.
* `auto_pause` - (Optional, String) Specify whether the cluster can auto-pause while `db_mode` is `SERVERLESS`. Values: `yes` (default), `no`.
* `auto_renew_flag` - (Optional, Int) Auto renew flag. Valid values are `0`(MANUAL_RENEW), `1`(AUTO_RENEW). Default value is `0`. Only works for PREPAID cluster.
* `charge_type` - (Optional, String, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. Default value is `POSTPAID_BY_HOUR`.
* `cynos_version` - (Optional, String) Kernel minor version, like `3.1.16.002`.
* `db_mode` - (Optional, String) Specify DB mode, only available when `db_type` is `MYSQL`. Values: `NORMAL` (Default), `SERVERLESS`.
* `force_delete` - (Optional, Bool) Indicate whether to delete cluster instance directly or not. Default is false. If set true, the cluster and its `All RELATED INSTANCES` will be deleted instead of staying recycle bin. Note: works for both `PREPAID` and `POSTPAID_BY_HOUR` cluster.
* `instance_count` - (Optional, Int, ForceNew) The number of instances, the range is (0,16], the default value is 2 (i.e. one RW instance + one Ro instance), the passed n means 1 RW instance + n-1 Ro instances (with the same specifications), if you need a more accurate cluster composition, please use InstanceInitInfos.
* `instance_cpu_core` - (Optional, Int) The number of CPU cores of read-write type instance in the CynosDB cluster. Required while creating normal cluster. Note: modification of this field will take effect immediately, if want to upgrade on maintenance window, please upgrade from console.
* `instance_init_infos` - (Optional, List, ForceNew) Instance initialization configuration information, mainly used to select instances of different specifications when purchasing a cluster.
* `instance_maintain_duration` - (Optional, Int) Duration time for maintenance, unit in second. `3600` by default.
* `instance_maintain_start_time` - (Optional, Int) Offset time from 00:00, unit in second. For example, 03:00am should be `10800`. `10800` by default.
* `instance_maintain_weekdays` - (Optional, Set: [`String`]) Weekdays for maintenance. `["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]` by default.
* `instance_memory_size` - (Optional, Int) Memory capacity of read-write type instance, unit in GB. Required while creating normal cluster. Note: modification of this field will take effect immediately, if want to upgrade on maintenance window, please upgrade from console.
* `max_cpu` - (Optional, Float64) Maximum CPU core count, required while `db_mode` is `SERVERLESS`, request DescribeServerlessInstanceSpecs for more reference.
* `min_cpu` - (Optional, Float64) Minimum CPU core count, required while `db_mode` is `SERVERLESS`, request DescribeServerlessInstanceSpecs for more reference.
* `old_ip_reserve_hours` - (Optional, Int) Recycling time of the old address, must be filled in when modifying the vpcRecycling time of the old address, must be filled in when modifying the vpc.
* `param_items` - (Optional, List) Specify parameter list of database. It is valid when `param_template_id` is set in create cluster. Use `data.tencentcloud_mysql_default_params` to query available parameter details.
* `param_template_id` - (Optional, Int) The ID of the parameter template.
* `port` - (Optional, Int, ForceNew) Port of CynosDB cluster.
* `prarm_template_id` - (Optional, Int, **Deprecated**) It will be deprecated. Use `param_template_id` instead. The ID of the parameter template.
* `prepaid_period` - (Optional, Int, ForceNew) The tenancy (time unit is month) of the prepaid instance. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`. NOTE: it only works when charge_type is set to `PREPAID`.
* `project_id` - (Optional, Int, ForceNew) ID of the project. `0` by default.
* `ro_group_sg` - (Optional, List: [`String`]) IDs of security group for `ro_group`.
* `rw_group_sg` - (Optional, List: [`String`]) IDs of security group for `rw_group`.
* `serverless_status_flag` - (Optional, String) Specify whether to pause or resume serverless cluster. values: `resume`, `pause`.
* `slave_zone` - (Optional, String) Multi zone Addresses of the CynosDB Cluster.
* `storage_limit` - (Optional, Int) Storage limit of CynosDB cluster instance, unit in GB. The maximum storage of a non-serverless instance in GB. NOTE: If db_type is `MYSQL` and charge_type is `PREPAID`, the value cannot exceed the maximum storage corresponding to the CPU and memory specifications, and the transaction mode is `order and pay`. when charge_type is `POSTPAID_BY_HOUR`, this argument is unnecessary.
* `storage_pay_mode` - (Optional, Int) Cluster storage billing mode, pay-as-you-go: `0`-yearly/monthly: `1`-The default is pay-as-you-go. When the DbType is MYSQL, when the cluster computing billing mode is post-paid (including DbMode is SERVERLESS), the storage billing mode can only be billing by volume; rollback and cloning do not support yearly subscriptions monthly storage.
* `tags` - (Optional, Map) The tags of the CynosDB cluster.

The `instance_init_infos` object supports the following:

* `cpu` - (Required, Int, ForceNew) CPU of instance.
* `instance_count` - (Required, Int, ForceNew) Instance count. Range: [1, 15].
* `instance_type` - (Required, String, ForceNew) Instance type. Value: `rw`, `ro`.
* `memory` - (Required, Int, ForceNew) Memory of instance.
* `device_type` - (Optional, String, ForceNew) Instance machine type. Values: `common`, `exclusive`.
* `max_ro_count` - (Optional, Int, ForceNew) Maximum number of Serverless instances. Range [1,15].
* `max_ro_cpu` - (Optional, Float64, ForceNew) Maximum Serverless Instance Specifications.
* `min_ro_count` - (Optional, Int, ForceNew) Minimum number of Serverless instances. Range [1,15].
* `min_ro_cpu` - (Optional, Float64, ForceNew) Minimum Serverless Instance Specifications.

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
* `serverless_status` - Serverless cluster status. NOTE: This is a readonly attribute, to modify, please set `serverless_status_flag`.
* `storage_used` - Used storage of CynosDB cluster, unit in MB.


## Import

CynosDB cluster can be imported using the id, e.g.

```
$ terraform import tencentcloud_cynosdb_cluster.example cynosdbmysql-dzj5l8gz
```

