---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_sharding_instance"
sidebar_current: "docs-tencentcloud-resource-mongodb_sharding_instance"
description: |-
  Provide a resource to create a Mongodb sharding instance.
---

# tencentcloud_mongodb_sharding_instance

Provide a resource to create a Mongodb sharding instance.

~> **NOTE:** The `add_node_list` and `remove_node_list` arguments are used to submit node change actions. When updating the resource, only newly added items in these lists will be sent to the API. If an existing item is removed from the Terraform configuration, Terraform only updates the local state and does not submit a repeated add or remove request. To add or remove another read-only node, append a new block instead of modifying an existing one. After the change is completed, obsolete action records can be removed from the configuration, and this cleanup does not trigger a new node operation when the remaining list is a subset of the previous list. In general, it is recommended to keep these action records in the configuration and avoid cleanup unless necessary.

~> **NOTE:** The `cpu` parameter takes effect only when the configuration is changed. Changing the `cpu` triggers the `ModifyDBInstanceSpec` API to adjust the CPU specification of the running MongoDB instance in-place. The supported CPU specifications can be obtained through the `DescribeSpecInfo` API.

## Example Usage

```hcl
resource "tencentcloud_mongodb_sharding_instance" "example" {
  instance_name   = "tf-example"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 4
  volume          = 100
  engine_version  = "MONGO_40_WT"
  machine_type    = "HIO10G"
  available_zone  = "ap-guangzhou-6"
  vpc_id          = "vpc-i5yyodl9"
  subnet_id       = "subnet-hhi88a58"
  project_id      = 0
  password        = "Password@123"
  mongos_cpu      = 1
  mongos_memory   = 2
  mongos_node_num = 3
}
```

### Add a read-only node

```hcl
resource "tencentcloud_mongodb_sharding_instance" "example" {
  instance_name   = "tf-example"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 4
  volume          = 100
  engine_version  = "MONGO_40_WT"
  machine_type    = "HIO10G"
  available_zone  = "ap-guangzhou-6"
  vpc_id          = "vpc-i5yyodl9"
  subnet_id       = "subnet-hhi88a58"
  project_id      = 0
  password        = "Password@123"
  mongos_cpu      = 1
  mongos_memory   = 2
  mongos_node_num = 3

  add_node_list {
    role = "READONLY"
    zone = "ap-guangzhou-6"
  }
}
```

### Remove a read-only node

```hcl
resource "tencentcloud_mongodb_sharding_instance" "example" {
  instance_name   = "tf-example"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 4
  volume          = 100
  engine_version  = "MONGO_40_WT"
  machine_type    = "HIO10G"
  available_zone  = "ap-guangzhou-6"
  vpc_id          = "vpc-i5yyodl9"
  subnet_id       = "subnet-hhi88a58"
  project_id      = 0
  password        = "Password@123"
  mongos_cpu      = 1
  mongos_memory   = 2
  mongos_node_num = 3

  remove_node_list {
    role      = "READONLY"
    node_name = "cmgo-xxxx_0-node-readonly0"
    zone      = "ap-guangzhou-6"
  }
}
```

## Argument Reference

The following arguments are supported:

* `available_zone` - (Required, String) The available zone of the Mongodb.
* `engine_version` - (Required, String) Refers to version information. The DescribeSpecInfo API can be called to obtain detailed information about the supported versions.
- MONGO_40_WT: version of the MongoDB 4.0 WiredTiger storage engine.
- MONGO_42_WT: version of the MongoDB 4.2 WiredTiger storage engine.
- MONGO_44_WT: version of the MongoDB 4.4 WiredTiger storage engine.
- MONGO_50_WT: version of the MongoDB 5.0 WiredTiger storage engine.
- MONGO_60_WT: version of the MongoDB 6.0 WiredTiger storage engine.
- MONGO_70_WT: version of the MongoDB 7.0 WiredTiger storage engine.
- MONGO_80_WT: version of the MongoDB 8.0 WiredTiger storage engine.
* `instance_name` - (Required, String) Name of the Mongodb instance.
* `machine_type` - (Required, String, ForceNew) Type of Mongodb instance, and available values include `HIO`(or `GIO` which will be deprecated, represents high IO) and `HIO10G`(or `TGIO` which will be deprecated, represents 10-gigabit high IO).
* `memory` - (Required, Int) Memory size. The minimum value is 2, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.
* `nodes_per_shard` - (Required, Int, ForceNew) Number of nodes per shard, at least 3(one master and two slaves). Allow value[3, 5, 7].
* `shard_quantity` - (Required, Int, ForceNew) Number of sharding.
* `volume` - (Required, Int) Disk size. The minimum value is 25, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.
* `add_node_list` - (Optional, List) Add node list. Node type and availability zone information.
* `auto_renew_flag` - (Optional, Int) Auto renew flag. Valid values are `0`(NOTIFY_AND_MANUAL_RENEW), `1`(NOTIFY_AND_AUTO_RENEW) and `2`(DISABLE_NOTIFY_AND_MANUAL_RENEW). Default value is `0`. Note: only works for PREPAID instance. Only supports`0` and `1` for creation.
* `availability_zone_list` - (Optional, List: [`String`]) A list of nodes deployed in multiple availability zones. For more information, please use the API DescribeSpecInfo.
			- Multi-availability zone deployment nodes can only be deployed in 3 different availability zones. It is not supported to deploy most nodes of the cluster in the same availability zone. For example, a 3-node cluster does not support the deployment of 2 nodes in the same zone.
			- Version 4.2 and above are not supported.
			- Read-only disaster recovery instances are not supported.
			- Basic network cannot be selected.
* `charge_type` - (Optional, String, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. Default value is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`. Caution that update operation on this field will delete old instances and create new one with new charge type.
* `cpu` - (Optional, Int) The CPU core count of the MongoDB instance after the configuration change. Unit: C. When this parameter is empty, the current CPU size of the instance is used by default. The supported CPU specifications can be obtained through the DescribeSpecInfo API.
* `hidden_zone` - (Optional, String) The availability zone to which the Hidden node belongs. This parameter is required in cross-AZ instance deployment.
* `in_maintenance` - (Optional, Int) Switch time for instance configuration changes.
	- 0: When the adjustment is completed, perform the configuration task immediately. Default is 0.
	- 1: Perform reconfiguration tasks within the maintenance time window.
Note: Adjusting the number of nodes and slices does not support changes within the maintenance window.
* `mongos_cpu` - (Optional, Int) Number of mongos cpu.
* `mongos_memory` - (Optional, Int) Mongos memory size in GB.
* `mongos_node_num` - (Optional, Int) Number of mongos.
* `password` - (Optional, String) Password of this Mongodb account.
* `prepaid_period` - (Optional, Int) The tenancy (time unit is month) of the prepaid instance. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. NOTE: it only works when charge_type is set to `PREPAID`.
* `project_id` - (Optional, Int) ID of the project which the instance belongs.
* `remove_node_list` - (Optional, List) Remove node list. Node type, node name, and availability zone information. Note: Based on the consistency principle of each shard node in a sharding instance, when removing nodes, you only need to specify the node corresponding to shard 0, e.g., `cmgo-xxxx_0-node-readonly0` will remove the first readonly node of each shard.
* `security_groups` - (Optional, Set: [`String`]) ID of the security group.
* `subnet_id` - (Optional, String, ForceNew) ID of the subnet within this VPC. The value is required if `vpc_id` is set.
* `tags` - (Optional, Map) The tags of the Mongodb. Key name `project` is system reserved and can't be used.
* `vpc_id` - (Optional, String, ForceNew) ID of the VPC.

The `add_node_list` object supports the following:

* `role` - (Required, String) Node role to add. Valid values: `SECONDARY` (Mongod node), `READONLY` (read-only node), `MONGOS` (Mongos node).
* `zone` - (Required, String) The availability zone for the new node.

The `remove_node_list` object supports the following:

* `node_name` - (Required, String) Node ID to remove. For sharding cluster, specify the node name corresponding to one shard group. For example: `cmgo-xxxx_0-node-readonly0`.
* `role` - (Required, String) Node role to remove. Valid values: `SECONDARY` (Mongod secondary node), `READONLY` (read-only node), `MONGOS` (Mongos node).
* `zone` - (Required, String) The availability zone of the node to remove.

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
terraform import tencentcloud_mongodb_sharding_instance.example cmgo-41s6jwy4
```

