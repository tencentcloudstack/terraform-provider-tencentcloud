---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_db_custom_cluster_resources"
sidebar_current: "docs-tencentcloud-datasource-dbdc_db_custom_cluster_resources"
description: |-
  Use this data source to query DB Custom cluster resource summary info from TencentCloud DBDC product.
---

# tencentcloud_dbdc_db_custom_cluster_resources

Use this data source to query DB Custom cluster resource summary info from TencentCloud DBDC product.

## Example Usage

### Query dbdc db custom cluster resources by cluster_id

```hcl
data "tencentcloud_dbdc_db_custom_cluster_resources" "example" {
  cluster_id = "dbcc-nmtmsew8"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `allocatable` - The sum of the allocatable capacity of all nodes in the cluster (= Capacity - system reservation).
  * `cpu` - Number of CPU cores.
  * `memory` - Memory size in GiB.
  * `pods` - Number of pods.
* `available` - Cluster schedulable remaining capacity (the sum of max(0, Allocatable - Requests) for all nodes).
  * `cpu` - Number of CPU cores.
  * `memory` - Memory size in GiB.
  * `pods` - Number of pods.
* `capacity` - The sum of the physical total resource capacity of all nodes in the cluster.
  * `cpu` - Number of CPU cores.
  * `memory` - Memory size in GiB.
  * `pods` - Number of pods.
* `limits` - The sum of the limits of all non-terminal pods in the cluster (including system pods, the Pods field has no semantics and is fixed to 0).
  * `cpu` - Number of CPU cores.
  * `memory` - Memory size in GiB.
  * `pods` - Number of pods.
* `node_count` - Total number of worker nodes participating in the aggregation (excluding control plane nodes).
* `requests` - The sum of the requests of all non-terminal pods in the cluster (including system pods).
  * `cpu` - Number of CPU cores.
  * `memory` - Memory size in GiB.
  * `pods` - Number of pods.


