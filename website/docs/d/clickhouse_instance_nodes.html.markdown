---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_instance_nodes"
sidebar_current: "docs-tencentcloud-datasource-clickhouse_instance_nodes"
description: |-
  Use this data source to query detailed information of clickhouse instance_nodes
---

# tencentcloud_clickhouse_instance_nodes

Use this data source to query detailed information of clickhouse instance_nodes

## Example Usage

```hcl
data "tencentcloud_clickhouse_instance_nodes" "instance_nodes" {
  instance_id    = "cdwch-mvfjh373"
  node_role      = "data"
  display_policy = "all"
  force_all      = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) InstanceId.
* `display_policy` - (Optional, String) Display strategy, display all when All.
* `force_all` - (Optional, Bool) When true, returns all nodes, that is, the Limit is infinitely large.
* `node_role` - (Optional, String) Cluster role type, default is `data` data node.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_nodes_list` - Total number of instance nodes.
  * `cluster` - Name of the clickhouse cluster to which it belongs.
  * `core` - CPU cores.
  * `disk_size` - Disk size.
  * `disk_type` - Disk type.
  * `ip` - IP Address.
  * `is_ch_proxy` - When true, it indicates that the chproxy process has been deployed on the node.
  * `memory` - Memory size.
  * `node_groups` - Group information to which the node belongs.
    * `group_name` - Group Name.
    * `replica_name` - Copy variable name.
    * `shard_name` - Fragmented variable name.
  * `rip` - VPC IP.
  * `spec` - Model, such as S1.


