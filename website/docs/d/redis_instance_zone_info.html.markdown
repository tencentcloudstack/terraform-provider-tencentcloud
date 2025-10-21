---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_instance_zone_info"
sidebar_current: "docs-tencentcloud-datasource-redis_instance_zone_info"
description: |-
  Use this data source to query detailed information of redis instance_zone_info
---

# tencentcloud_redis_instance_zone_info

Use this data source to query detailed information of redis instance_zone_info

## Example Usage

```hcl
data "tencentcloud_redis_instance_zone_info" "instance_zone_info" {
  instance_id = "crs-c1nl9rpv"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional, String) The ID of instance.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `replica_groups` - List of instance node groups.
  * `group_id` - Node group ID.
  * `group_name` - Node group Name.
  * `redis_nodes` - Node group node list.
    * `keys` - The number of node keys.
    * `node_id` - Node ID.
    * `role` - Node role.
    * `slot` - Node slot distribution.
    * `status` - Node status.
  * `role` - The node group type, master is the primary node, and replica is the replica node.
  * `zone_id` - he availability zone ID of the node, such as ap-guangzhou-1.


