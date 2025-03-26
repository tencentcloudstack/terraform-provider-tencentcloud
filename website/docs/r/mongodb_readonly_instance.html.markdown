---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_readonly_instance"
sidebar_current: "docs-tencentcloud-resource-mongodb_readonly_instance"
description: |-
  Provide a resource to create a Readonly mongodb instance.
---

# tencentcloud_mongodb_readonly_instance

Provide a resource to create a Readonly mongodb instance.

## Example Usage

### Replset readonly instance

```hcl
resource "tencentcloud_mongodb_readonly_instance" "mongodb" {
  instance_name          = "tf-mongodb-readonly-test"
  memory                 = 4
  volume                 = 100
  engine_version         = "MONGO_44_WT"
  machine_type           = "HIO10G"
  available_zone         = "ap-guangzhou-3"
  project_id             = 0
  father_instance_id     = "cmgo-xxxxxx"
  father_instance_region = "ap-guangzhou"
  vpc_id                 = "vpc-xxxxxx"
  subnet_id              = "subnet-xxxxxx"
  security_groups        = ["sg-xxxxxx"]
  cluster_type           = "REPLSET"
}
```

### Shard readonly instance

```hcl
resource "tencentcloud_mongodb_readonly_instance" "sharding_mongodb" {
  instance_name          = "tf-mongodb-readonly-shard"
  memory                 = 4
  volume                 = 100
  engine_version         = "MONGO_44_WT"
  machine_type           = "HIO10G"
  available_zone         = "ap-guangzhou-3"
  project_id             = 0
  father_instance_id     = "cmgo-xxxxxx"
  father_instance_region = "ap-guangzhou"
  vpc_id                 = "vpc-xxxxxx"
  subnet_id              = "subnet-xxxxxx"
  security_groups        = ["sg-xxxxxx"]
  cluster_type           = "SHARD"
  mongos_cpu             = 1
  mongos_memory          = 2
  mongos_node_num        = 3
}
```

## Argument Reference

The following arguments are supported:

* `available_zone` - (Required, String, ForceNew) The available zone of the Mongodb.
* `cluster_type` - (Required, String, ForceNew) Instance schema type.	- REPLSET: Replset cluster;	- SHARD: Shard cluster.
* `engine_version` - (Required, String, ForceNew) Version of the Mongodb, and available values include `MONGO_36_WT` (MongoDB 3.6 WiredTiger Edition), `MONGO_40_WT` (MongoDB 4.0 WiredTiger Edition) and `MONGO_42_WT`  (MongoDB 4.2 WiredTiger Edition). NOTE: `MONGO_3_WT` (MongoDB 3.2 WiredTiger Edition) and `MONGO_3_ROCKS` (MongoDB 3.2 RocksDB Edition) will deprecated.
* `father_instance_id` - (Required, String, ForceNew) Indicates the main instance ID of readonly instances.
* `father_instance_region` - (Required, String, ForceNew) Indicates the region of main instance.
* `instance_name` - (Required, String) Name of the Mongodb instance.
* `machine_type` - (Required, String, ForceNew) Type of Mongodb instance, and available values include `HIO`(or `GIO` which will be deprecated, represents high IO) and `HIO10G`(or `TGIO` which will be deprecated, represents 10-gigabit high IO).
* `memory` - (Required, Int) Memory size. The minimum value is 2, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.
* `volume` - (Required, Int) Disk size. The minimum value is 25, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.
* `auto_renew_flag` - (Optional, Int) Auto renew flag. Valid values are `0`(NOTIFY_AND_MANUAL_RENEW), `1`(NOTIFY_AND_AUTO_RENEW) and `2`(DISABLE_NOTIFY_AND_MANUAL_RENEW). Default value is `0`. Note: only works for PREPAID instance. Only supports`0` and `1` for creation.
* `charge_type` - (Optional, String, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. Default value is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`. Caution that update operation on this field will delete old instances and create new one with new charge type.
* `mongos_cpu` - (Optional, Int) Number of mongos cpu.
* `mongos_memory` - (Optional, Int) Mongos memory size in GB.
* `mongos_node_num` - (Optional, Int) Number of mongos.
* `node_num` - (Optional, Int) The number of nodes in each replica set. Default value: 3.
* `nodes_per_shard` - (Optional, Int, ForceNew) Number of nodes per shard, at least 3(one master and two slaves).
* `prepaid_period` - (Optional, Int) The tenancy (time unit is month) of the prepaid instance. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. NOTE: it only works when charge_type is set to `PREPAID`.
* `project_id` - (Optional, Int) ID of the project which the instance belongs.
* `security_groups` - (Optional, Set: [`String`]) ID of the security group.
* `shard_quantity` - (Optional, Int, ForceNew) Number of sharding.
* `subnet_id` - (Optional, String, ForceNew) ID of the subnet within this VPC. The value is required if `vpc_id` is set.
* `tags` - (Optional, Map) The tags of the Mongodb. Key name `project` is system reserved and can't be used.
* `vpc_id` - (Optional, String, ForceNew) ID of the VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the Mongodb instance.
* `status` - Status of the Mongodb instance, and available values include pending initialization(expressed with 0),  processing(expressed with 1), running(expressed with 2) and expired(expressed with -2).
* `vip` - IP of the Mongodb instance.
* `vport` - IP port of the Mongodb instance.


## Import

Mongodb instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_mongodb_instance.mongodb cmgo-xxxxxx
```

