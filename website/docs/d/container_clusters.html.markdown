---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_container_clusters"
sidebar_current: "docs-tencentcloud-datasource-container-clusters-x"
description: |-
  Get container clusters in the current region.
---

# tencentcloud_container_clusters

Use this data source to get container clusters in the current region. 
By default every clusters in current region will be returned. 

~> **NOTE:** It has been deprecated and replaced by tencentcloud_kubernetes_clusters.
## Example Usage

```hcl
data "tencentcloud_container_clusters" "foo" {}
```

## Argument Reference

 * `cluster_id` - (Optional) An id identify the cluster, like `cls-xxxxxx`.
 * `limit` - (Optional) An int variable describe how many cluster in return at most  .


## Attributes Reference

A list of clusters will be exported and its every element contains the following attributes:

 * `cluster_id` - An id identify the cluster, like `cls-xxxxxx`.
 * `security_certification_authority` - Describe the certificate string needed for using kubectl to access to kubernetes.
 * `security_cluster_external_endpoint` - Describe the address needed for using kubectl to access to kubernetes.
 * `security_username` - Describe the username needed for using kubectl to access to kubernetes.
 * `security_password` - Describe the password needed for using kubectl to access to kubernetes.
 * `description` - The description of the cluster.
 * `kubernetes_version` - Describe the running kubernetes version on the cluster.
 * `nodes_num` - Describe how many cluster instances in the cluster.
 * `nodes_status` - Describe the current status of the instances in the cluster.
 * `total_cpu` - Describe the total cpu of each instance in the cluster.
 * `total_mem` - Describe the total memory of each instance in the cluster.
