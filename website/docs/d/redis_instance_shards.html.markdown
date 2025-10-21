---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_instance_shards"
sidebar_current: "docs-tencentcloud-datasource-redis_instance_shards"
description: |-
  Use this data source to query detailed information of redis instance_shards
---

# tencentcloud_redis_instance_shards

Use this data source to query detailed information of redis instance_shards

## Example Usage

```hcl
data "tencentcloud_redis_instance_shards" "instance_shards" {
  instance_id  = "crs-c1nl9rpv"
  filter_slave = false
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of instance.
* `filter_slave` - (Optional, Bool) Whether to filter out slave information.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_shards` - Instance shard list information.
  * `connected` - Service status: 0-down;1-on.
  * `keys` - Number of keys.
  * `role` - role.
  * `runid` - The node ID of the instance runtime.
  * `shard_id` - Shard node ID.
  * `shard_name` - Shard node name.
  * `slots` - Slot information.
  * `storage_slope` - Capacity tilt.
  * `storage` - Used capacity.


