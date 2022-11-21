---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_shards"
sidebar_current: "docs-tencentcloud-datasource-dcdb_shards"
description: |-
  Use this data source to query detailed information of dcdb shards
---

# tencentcloud_dcdb_shards

Use this data source to query detailed information of dcdb shards

## Example Usage

```hcl
data "tencentcloud_dcdb_shards" "shards" {
  instance_id        = "your_instance_id"
  shard_instance_ids = ["shard1_id"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `result_output_file` - (Optional, String) Used to save results.
* `shard_instance_ids` - (Optional, Set: [`String`]) shard instance ids.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - shard list.
  * `cpu` - cpu cores.
  * `create_time` - create time.
  * `instance_id` - instance id.
  * `memory_usage` - memory usage.
  * `memory` - memory, the unit is GB.
  * `node_count` - node count.
  * `paymode` - pay mode.
  * `period_end_time` - expired time.
  * `project_id` - project id.
  * `proxy_version` - proxy version.
  * `range` - the range of shard key.
  * `region` - region.
  * `shard_instance_id` - shard instance id.
  * `shard_master_zone` - shard master zone.
  * `shard_serial_id` - shard serial id.
  * `shard_slave_zones` - shard slave zones.
  * `status_desc` - status description.
  * `status` - status.
  * `storage_usage` - storage usage.
  * `storage` - memory, the unit is GB.
  * `subnet_id` - subnet id.
  * `vpc_id` - vpc id.
  * `zone` - zone.


