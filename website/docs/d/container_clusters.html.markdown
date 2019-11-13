---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_container_clusters"
sidebar_current: "docs-tencentcloud-datasource-container_clusters"
description: |-
  Get container clusters in the current region.
---

# tencentcloud_container_clusters

Get container clusters in the current region.

Use this data source to get container clusters in the current region. By default every clusters in current region will be returned.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_kubernetes_clusters.

## Example Usage

```hcl
data "tencentcloud_container_clusters" "foo" {
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional) An id identify the cluster, like `cls-xxxxxx`.
* `limit` - (Optional) An int variable describe how many cluster in return at most.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `clusters` - An information list of kubernetes clusters.
  * `cluster_id` - An id identify the cluster, like `cls-xxxxxx`.
  * `cluster_name` - Name the cluster.
  * `description` - The description of the cluster.
  * `kubernetes_version` - Describe the running kubernetes version on the cluster.
  * `nodes_num` - Describe how many cluster instances in the cluster.
  * `nodes_status` - Describe the current status of the instances in the cluster.
  * `security_certification_authority` - Describe the certificate string needed for using kubectl to access to kubernetes.
  * `security_cluster_external_endpoint` - Describe the address needed for using kubectl to access to kubernetes.
  * `security_password` - Describe the password needed for using kubectl to access to kubernetes.
  * `security_username` - Describe the username needed for using kubectl to access to kubernetes.
  * `total_cpu` - Describe the total cpu of each instance in the cluster.
  * `total_mem` - Describe the total memory of each instance in the cluster.
* `total_count` - Number of clusters.


