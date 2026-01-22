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
resource "tencentcloud_mongodb_instance" "example" {
  instance_name  = "tf-example"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_40_WT"
  machine_type   = "HIO10G"
  available_zone = "ap-guangzhou-6"
  vpc_id         = "vpc-i5yyodl9"
  subnet_id      = "subnet-hhi88a58"
  project_id     = 0
  password       = "Password@123"
}
```

## Argument Reference

The following arguments are supported:

* `available_zone` - (Required, String, ForceNew) The available zone of the Mongodb.
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
* `volume` - (Required, Int) Disk size. The minimum value is 25, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.
* `add_node_list` - (Optional, List) Add node attribute list.
* `auto_renew_flag` - (Optional, Int) Auto renew flag. Valid values are `0`(NOTIFY_AND_MANUAL_RENEW), `1`(NOTIFY_AND_AUTO_RENEW) and `2`(DISABLE_NOTIFY_AND_MANUAL_RENEW). Default value is `0`. Note: only works for PREPAID instance. Only supports`0` and `1` for creation.
* `availability_zone_list` - (Optional, List: [`String`]) If cloud database instances are deployed in multiple availability zones, specify a list of multiple availability zones.
	- To deploy an instance with multiple availability zones, the parameter Zone specifies the primary availability zone information of the instance; Availability ZoneList specifies all availability zone information, including the primary availability zone. The input format is as follows: [ap-Guangzhou-2,ap-Guangzhou-3,ap-Guangzhou-4].
	- You can obtain availability zone information planned in different regions of the cloud database through the interface DescribeSpecInfo, so as to specify effective availability zones.
	- Multiple availability zone deployment nodes can only be deployed in 3 different availability zones. Deploying most nodes of a cluster in the same availability zone is not supported. For example, a 3-node cluster does not support 2 nodes deployed in the same zone.
* `charge_type` - (Optional, String, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. Default value is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`. Caution that update operation on this field will delete old instances and create new one with new charge type.
* `hidden_zone` - (Optional, String) The availability zone to which the Hidden node belongs. This parameter is required in cross-AZ instance deployment.
* `in_maintenance` - (Optional, Int) Switch time for instance configuration changes.
	- 0: When the adjustment is completed, perform the configuration task immediately. Default is 0.
	- 1: Perform reconfiguration tasks within the maintenance time window.
Note: Adjusting the number of nodes and slices does not support changes within the maintenance window.
* `maintenance_end` - (Optional, String) Maintenance window end time.
	- The value range is any full point or half point from `00:00-23:00`, and the maintenance time duration is at least 30 minutes and at most 3 hours.
	- The end time must be based on the start time backwards.
* `maintenance_start` - (Optional, String) Maintenance window start time. The value range is any full point or half point from `00:00-23:00`, such as 00:00 or 00:30.
* `node_num` - (Optional, Int) The number of nodes in each replica set. Default value: 3.
* `password` - (Optional, String) Password of this Mongodb account.
* `prepaid_period` - (Optional, Int) The tenancy (time unit is month) of the prepaid instance. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. NOTE: it only works when charge_type is set to `PREPAID`.
* `project_id` - (Optional, Int) ID of the project which the instance belongs.
* `remove_node_list` - (Optional, List) Add node attribute list.
* `security_groups` - (Optional, Set: [`String`]) ID of the security group.
* `subnet_id` - (Optional, String, ForceNew) ID of the subnet within this VPC. The value is required if `vpc_id` is set.
* `tags` - (Optional, Map) The tags of the Mongodb. Key name `project` is system reserved and can't be used.
* `vpc_id` - (Optional, String, ForceNew) ID of the VPC.

The `add_node_list` object supports the following:

* `role` - (Required, String) The node role that needs to be added.
- SECONDARY: Mongod node;
- READONLY: read-only node;
- MONGOS: Mongos node.
* `zone` - (Required, String) The availability zone corresponding to the node.
- single availability zone, where all nodes are in the same availability zone;
- multiple availability zones: the current standard specification is the distribution of three availability zones, and the master and slave nodes are not in the same availability zone. You should pay attention to configuring the availability zone corresponding to the new node, and the rule that the number of nodes in any two availability zones is greater than the third availability zone must be met after the addition.

The `remove_node_list` object supports the following:

* `node_name` - (Required, String) The node ID to delete. The shard cluster must specify the name of the node to be deleted by a group of shards, and the rest of the shards should be grouped and aligned.
* `role` - (Required, String) The node role that needs to be deleted.
- SECONDARY: Mongod node;
- READONLY: read-only node;
- MONGOS: Mongos node.
* `zone` - (Required, String) The availability zone corresponding to the node.
- single availability zone, where all nodes are in the same availability zone;
- multiple availability zones: the current standard specification is the distribution of three availability zones, and the master and slave nodes are not in the same availability zone. You should pay attention to configuring the availability zone corresponding to the new node, and the rule that the number of nodes in any two availability zones is greater than the third availability zone must be met after the addition.

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
terraform import tencentcloud_mongodb_instance.example cmgo-41s6jwy4
```

