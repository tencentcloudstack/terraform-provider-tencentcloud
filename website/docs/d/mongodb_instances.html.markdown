---
subcategory: "MongoDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instances"
sidebar_current: "docs-tencentcloud-datasource-mongodb_instances"
description: |-
  Use this data source to query detailed information of Mongodb instances.
---

# tencentcloud_mongodb_instances

Use this data source to query detailed information of Mongodb instances.

## Example Usage

```hcl
data "tencentcloud_mongodb_instances" "mongodb" {
  instance_id  = "cmgo-l6lwdsel"
  cluster_type = "REPLSET"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_type` - (Optional, String) Type of Mongodb cluster, and available values include replica set cluster(expressed with `REPLSET`), sharding cluster(expressed with `SHARD`).
* `instance_id` - (Optional, String) ID of the Mongodb instance to be queried.
* `instance_name_prefix` - (Optional, String) Name prefix of the Mongodb instance.
* `result_output_file` - (Optional, String) Used to store results.
* `tags` - (Optional, Map) Tags of the Mongodb instance to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - A list of instances. Each element contains the following attributes:
  * `auto_renew_flag` - Auto renew flag.
  * `available_zone` - The available zone of the Mongodb.
  * `charge_type` - The charge type of instance.
  * `cluster_type` - Type of Mongodb cluster.
  * `cpu` - Number of cpu's core.
  * `create_time` - Creation time of the Mongodb instance.
  * `engine_version` - Version of the Mongodb engine.
  * `instance_id` - ID of the Mongodb instance.
  * `instance_name` - Name of the Mongodb instance.
  * `machine_type` - Type of Mongodb instance.
  * `memory` - Memory size.
  * `project_id` - ID of the project which the instance belongs.
  * `shard_quantity` - Number of sharding.
  * `status` - Status of the Mongodb, and available values include pending initialization(expressed with 0),  processing(expressed with 1), running(expressed with 2) and expired(expressed with -2).
  * `subnet_id` - ID of the subnet.
  * `tags` - Tags of the Mongodb instance.
  * `vip` - IP of the Mongodb instance.
  * `volume` - Disk size.
  * `vpc_id` - ID of the VPC.
  * `vport` - IP port of the Mongodb instance.


