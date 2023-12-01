---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance"
sidebar_current: "docs-tencentcloud-resource-mongodb_instance"
description: |-
  Provide a resource to create a Mongodb instance.
---

# tencentcloud_mongodb_instance

Provide a resource to create a Mongodb instance.

## Example Usage

```hcl
resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "mongodb"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_36_WT"
  machine_type   = "HIO10G"
  available_zone = "ap-guangzhou-2"
  vpc_id         = "vpc-xxxxxx"
  subnet_id      = "subnet-xxxxxx"
  project_id     = 0
  password       = "password1234"
}
```

## Argument Reference

The following arguments are supported:

* `available_zone` - (Required, String, ForceNew) The available zone of the Mongodb.
* `engine_version` - (Required, String, ForceNew) Version of the Mongodb, and available values include `MONGO_36_WT` (MongoDB 3.6 WiredTiger Edition), `MONGO_40_WT` (MongoDB 4.0 WiredTiger Edition) and `MONGO_42_WT`  (MongoDB 4.2 WiredTiger Edition). NOTE: `MONGO_3_WT` (MongoDB 3.2 WiredTiger Edition) and `MONGO_3_ROCKS` (MongoDB 3.2 RocksDB Edition) will deprecated.
* `instance_name` - (Required, String) Name of the Mongodb instance.
* `machine_type` - (Required, String, ForceNew) Type of Mongodb instance, and available values include `HIO`(or `GIO` which will be deprecated, represents high IO) and `HIO10G`(or `TGIO` which will be deprecated, represents 10-gigabit high IO).
* `memory` - (Required, Int) Memory size. The minimum value is 2, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.
* `volume` - (Required, Int) Disk size. The minimum value is 25, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.
* `auto_renew_flag` - (Optional, Int) Auto renew flag. Valid values are `0`(NOTIFY_AND_MANUAL_RENEW), `1`(NOTIFY_AND_AUTO_RENEW) and `2`(DISABLE_NOTIFY_AND_MANUAL_RENEW). Default value is `0`. Note: only works for PREPAID instance. Only supports`0` and `1` for creation.
* `availability_zone_list` - (Optional, List: [`String`]) A list of nodes deployed in multiple availability zones. For more information, please use the API DescribeSpecInfo.
			- Multi-availability zone deployment nodes can only be deployed in 3 different availability zones. It is not supported to deploy most nodes of the cluster in the same availability zone. For example, a 3-node cluster does not support the deployment of 2 nodes in the same zone.
			- Version 4.2 and above are not supported.
			- Read-only disaster recovery instances are not supported.
			- Basic network cannot be selected.
* `charge_type` - (Optional, String, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. Default value is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`. Caution that update operation on this field will delete old instances and create new one with new charge type.
* `hidden_zone` - (Optional, String) The availability zone to which the Hidden node belongs. This parameter must be configured to deploy instances across availability zones.
* `node_num` - (Optional, Int) The number of nodes in each replica set. Default value: 3.
* `password` - (Optional, String) Password of this Mongodb account.
* `prepaid_period` - (Optional, Int) The tenancy (time unit is month) of the prepaid instance. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. NOTE: it only works when charge_type is set to `PREPAID`.
* `project_id` - (Optional, Int) ID of the project which the instance belongs.
* `security_groups` - (Optional, Set: [`String`], ForceNew) ID of the security group.
* `subnet_id` - (Optional, String, ForceNew) ID of the subnet within this VPC. The value is required if `vpc_id` is set.
* `tags` - (Optional, Map) The tags of the Mongodb. Key name `project` is system reserved and can't be used.
* `vpc_id` - (Optional, String, ForceNew) ID of the VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the Mongodb instance.
* `standby_instance_list` - List of standby instances' info.
  * `standby_instance_id` - Indicates the ID of standby instance.
  * `standby_instance_region` - Indicates the region of standby instance.
* `status` - Status of the Mongodb instance, and available values include pending initialization(expressed with 0),  processing(expressed with 1), running(expressed with 2) and expired(expressed with -2).
* `vip` - IP of the Mongodb instance.
* `vport` - IP port of the Mongodb instance.


## Import

Mongodb instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_mongodb_instance.mongodb cmgo-41s6jwy4
```

