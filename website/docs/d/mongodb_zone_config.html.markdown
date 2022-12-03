---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_zone_config"
sidebar_current: "docs-tencentcloud-datasource-mongodb_zone_config"
description: |-
  Use this data source to query the available mongodb specifications for different zone.
---

# tencentcloud_mongodb_zone_config

Use this data source to query the available mongodb specifications for different zone.

## Example Usage

```hcl
data "tencentcloud_mongodb_zone_config" "mongodb" {
  available_zone = "ap-guangzhou-2"
}
```

## Argument Reference

The following arguments are supported:

* `available_zone` - (Optional, String) The available zone of the Mongodb.
* `result_output_file` - (Optional, String) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of zone config. Each element contains the following attributes:
  * `available_zone` - The available zone of the Mongodb.
  * `cluster_type` - Type of Mongodb cluster.
  * `cpu` - Number of cpu's core.
  * `default_storage` - Default disk size.
  * `engine_version` - Version of the Mongodb version.
  * `machine_type` - Type of Mongodb instance.
  * `max_replicate_set_num` - Maximum replicate set num for sharding cluster.
  * `max_storage` - Maximum size of the disk.
  * `memory` - Memory size.
  * `min_replicate_set_num` - Minimum replicate set num for sharding cluster.
  * `min_storage` - Minimum sie of the disk.


