---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_clusters"
sidebar_current: "docs-tencentcloud-datasource-kubernetes_clusters"
description: |-
  Use this data source to query detailed information of kubernetes clusters.
---

# tencentcloud_kubernetes_clusters

Use this data source to query detailed information of kubernetes clusters.

## Example Usage

```hcl
data "tencentcloud_kubernetes_clusters" "name" {
  cluster_name = "terraform"
}

data "tencentcloud_kubernetes_clusters" "id" {
  cluster_id = "cls-godovr32"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional) ID of the cluster. Conflict with cluster_name, can not be set at the same time.
* `cluster_name` - (Optional) Name of the cluster. Conflict with cluster_id, can not be set at the same time.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - An information list of kubernetes clusters. Each element contains the following attributes:
  * `certification_authority` - The certificate used for access.
  * `cluster_cidr` - A network address block of the cluster. Different from vpc cidr and cidr of other clusters within this vpc.
  * `cluster_deploy_type` - Deployment type of the cluster.
  * `cluster_desc` - Description of the cluster.
  * `cluster_external_endpoint` - External network address to access.
  * `cluster_ipvs` - Indicates whether ipvs is enabled.
  * `cluster_max_pod_num` - The maximum number of Pods per node in the cluster.
  * `cluster_max_service_num` - The maximum number of services in the cluster.
  * `cluster_name` - Name of the cluster.
  * `cluster_node_num` - Number of nodes in the  cluster.
  * `cluster_os` - Operating system of the cluster.
  * `cluster_version` - Version of the cluster.
  * `container_runtime` - (**Deprecated**) It has been deprecated from version 1.18.1. Container runtime of the cluster.
  * `domain` - Domain name for access.
  * `ignore_cluster_cidr_conflict` - Indicates whether to ignore the cluster cidr conflict error.
  * `password` - Password of account.
  * `pgw_endpoint` - The Intranet address used for access.
  * `project_id` - Project Id of the cluster.
  * `security_policy` - Access policy.
  * `user_name` - User name of account.
  * `vpc_id` - Vpc Id of the cluster.
  * `worker_instances_list` - An information list of cvm within the WORKER clusters. Each element contains the following attributes.
    * `failed_reason` - Information of the cvm when it is failed.
    * `instance_id` - ID of the cvm.
    * `instance_role` - Role of the cvm.
    * `instance_state` - State of the cvm.


