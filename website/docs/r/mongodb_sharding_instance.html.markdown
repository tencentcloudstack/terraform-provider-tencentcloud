---
subcategory: "MongoDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_sharding_instance"
sidebar_current: "docs-tencentcloud-resource-mongodb_sharding_instance"
description: |-
  Provide a resource to create a Mongodb sharding instance.
---

# tencentcloud_mongodb_sharding_instance

Provide a resource to create a Mongodb sharding instance.

## Example Usage

```hcl
resource "tencentcloud_mongodb_sharding_instance" "mongodb" {
  instance_name   = "mongodb"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 4
  volume          = 100
  engine_version  = "MONGO_3_WT"
  machine_type    = "GIO"
  available_zone  = "ap-guangzhou-3"
  vpc_id          = "vpc-mz3efvbw"
  subnet_id       = "subnet-lk0svi3p"
  project_id      = 0
  password        = "password1234"
}
```

## Argument Reference

The following arguments are supported:

* `available_zone` - (Required, ForceNew) The available zone of the Mongodb.
* `engine_version` - (Required, ForceNew) Version of the Mongodb, and available values include `MONGO_3_WT`, `MONGO_3_ROCKS`, `MONGO_36_WT` and `MONGO_40_WT`.
* `instance_name` - (Required) Name of the Mongodb instance.
* `machine_type` - (Required, ForceNew) Type of Mongodb instance, and available values include `HIO`(or `GIO` which will be deprecated) and `HIO10G`(or `TGIO` which will be deprecated).
* `memory` - (Required) Memory size. The minimum value is 2, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.
* `nodes_per_shard` - (Required, ForceNew) Number of nodes per shard, at least 3(one master and two slaves).
* `password` - (Required) Password of this Mongodb account.
* `shard_quantity` - (Required, ForceNew) Number of sharding.
* `volume` - (Required) Disk size. The minimum value is 25, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.
* `auto_renew_flag` - (Optional) Auto renew flag. Valid values are `0`(NOTIFY_AND_MANUAL_RENEW), `1`(NOTIFY_AND_AUTO_RENEW) and `2`(DISABLE_NOTIFY_AND_MANUAL_RENEW). Default value is `0`. Note: only works for PREPAID instance. Only supports`0` and `1` for creation.
* `charge_type` - (Optional, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. Default value is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`. Caution that update operation on this field will delete old instances and create new one with new charge type.
* `prepaid_period` - (Optional) The tenancy (time unit is month) of the prepaid instance. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. NOTE: it only works when charge_type is set to `PREPAID`.
* `project_id` - (Optional) ID of the project which the instance belongs.
* `security_groups` - (Optional, ForceNew) ID of the security group. NOTE: for instance which `engine_version` is `MONGO_40_WT`, `security_groups` is not supported.
* `subnet_id` - (Optional, ForceNew) ID of the subnet within this VPC. The value is required if `vpc_id` is set.
* `tags` - (Optional) The tags of the Mongodb.
* `vpc_id` - (Optional, ForceNew) ID of the VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the Mongodb instance.
* `status` - Status of the Mongodb instance, and available values include pending initialization(expressed with 0),  processing(expressed with 1), running(expressed with 2) and expired(expressed with -2).
* `vip` - IP of the Mongodb instance.
* `vport` - IP port of the Mongodb instance.


## Import

Mongodb sharding instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_mongodb_sharding_instance.mongodb cmgo-41s6jwy4
```

