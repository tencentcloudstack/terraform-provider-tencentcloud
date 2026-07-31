---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_db_custom_cluster_node_resources"
sidebar_current: "docs-tencentcloud-datasource-dbdc_db_custom_cluster_node_resources"
description: |-
  Use this data source to query detailed information of DB Custom cluster node resources from TencentCloud DBDC product.
---

# tencentcloud_dbdc_db_custom_cluster_node_resources

Use this data source to query detailed information of DB Custom cluster node resources from TencentCloud DBDC product.

## Example Usage

### Query dbdc db custom cluster node resources by cluster_id

```hcl
data "tencentcloud_dbdc_db_custom_cluster_node_resources" "example" {
  cluster_id = "dbcc-nmtmsew8"
}
```

### Query dbdc db custom cluster node resources by cluster_id and node_ids

```hcl
data "tencentcloud_dbdc_db_custom_cluster_node_resources" "example" {
  cluster_id = "dbcc-nmtmsew8"

  node_ids = ["node-abc123", "node-def456"]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) DB Custom cluster ID.
* `node_ids` - (Optional, List: [`String`]) Node ID list. Up to 50 node IDs per request (enforced by the cloud API).
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `node_set` - DB Custom cluster node resource list.
  * `allocatable` - Node allocatable capacity = Capacity - system reserved.
    * `cpu` - CPU cores.
    * `memory` - Memory, in GiB.
    * `pods` - Number of pods.
  * `available` - Node schedulable remainder = max(0, Allocatable - Requests).
    * `cpu` - CPU cores.
    * `memory` - Memory, in GiB.
    * `pods` - Number of pods.
  * `capacity` - Node physical resource total capacity.
    * `cpu` - CPU cores.
    * `memory` - Memory, in GiB.
    * `pods` - Number of pods.
  * `limits` - Sum of limits of all non-terminal pods on the node (including system pods).
    * `cpu` - CPU cores.
    * `memory` - Memory, in GiB.
    * `pods` - Number of pods.
  * `node_id` - Node ID.
  * `requests` - Sum of requests of all non-terminal pods on the node (including system pods).
    * `cpu` - CPU cores.
    * `memory` - Memory, in GiB.
    * `pods` - Number of pods.


