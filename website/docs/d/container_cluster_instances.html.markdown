---
layout: "tencentcloud"
page_title: "TencentCloud: container_cluster-instances"
sidebar_current: "docs-tencentcloud-container-cluster-instances"
description: |-
  Get all instances of the specific cluster.
---

# tencentcloud_container_cluster_instances

Use this data source to get all instances in a specific cluster. 

## Example Usage

```hcl
data "tencentcloud_container_cluster_instances" "foo_instance" {
    cluster_id = "cls-abcdefg"
}
```

## Argument Reference

 * `cluster_id` - (Required) An id identify the cluster, like cls-xxxxxx.
 * `limit` - (Optional) An int variable describe how many instances in return at most.

## Attributes Reference
* `total_count` - Describe how many nodes in the cluster.

A list of nodes will be exported and its every element contains the following attributes:

 * `abnormal_reason` - Describe the reason when node is in abnormal state(if it was).
 * `cpu` - Describe the cpu of the node.
 * `mem` - Describe the memory of the node.
 * `instance_id` - An id identify the node, provided by cvm.
 * `is_normal` - Describe whether the node is normal.
 * `wan_ip` - Describe the wan ip of the node.
 * `lan_ip` - Descirbe the lan ip of the node.
