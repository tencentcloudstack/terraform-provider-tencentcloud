---
subcategory: "Container Cluster"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_container_cluster_instances"
sidebar_current: "docs-tencentcloud-datasource-container_cluster_instances"
description: |-
  Get all instances of the specific cluster.
---

# tencentcloud_container_cluster_instances

Get all instances of the specific cluster.

Use this data source to get all instances in a specific cluster.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_kubernetes_clusters.

## Example Usage

```hcl
data "tencentcloud_container_cluster_instances" "foo_instance" {
  cluster_id = "cls-abcdefg"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) An id identify the cluster, like cls-xxxxxx.
* `limit` - (Optional) An int variable describe how many instances in return at most.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `nodes` - An information list of kubernetes instances.
  * `abnormal_reason` - Describe the reason when node is in abnormal state(if it was).
  * `cpu` - Describe the cpu of the node.
  * `instance_id` - An id identify the node, provided by cvm.
  * `is_normal` - Describe whether the node is normal.
  * `lan_ip` - Describe the lan ip of the node.
  * `mem` - Describe the memory of the node.
  * `wan_ip` - Describe the wan ip of the node.
* `total_count` - Number of instances.


