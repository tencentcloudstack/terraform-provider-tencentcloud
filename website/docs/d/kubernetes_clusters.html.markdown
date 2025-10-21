---
subcategory: "Tencent Kubernetes Engine(TKE)"
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

* `cluster_id` - (Optional, String) ID of the cluster. Conflict with cluster_name, can not be set at the same time.
* `cluster_name` - (Optional, String) Name of the cluster. Conflict with cluster_id, can not be set at the same time.
* `kube_config_file_prefix` - (Optional, String) The path prefix of kube config. You can store KubeConfig in a specified directory by specifying this field, such as ~/.kube/k8s, then public network access will use ~/.kube/k8s-clusterID-kubeconfig naming, and intranet access will use ~/.kube /k8s-clusterID-kubeconfig-intranet naming. If this field is not set, the KubeConfig will not be exported.
* `result_output_file` - (Optional, String) Used to save results.
* `tags` - (Optional, Map) Tags of the cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - An information list of kubernetes clusters. Each element contains the following attributes:
  * `cdc_id` - CDC ID.
  * `certification_authority` - The certificate used for access.
  * `claim_expired_seconds` - The expired seconds to recycle ENI.
  * `cluster_as_enabled` - Indicates whether to enable cluster node auto scaler.
  * `cluster_cidr` - A network address block of the cluster. Different from vpc cidr and cidr of other clusters within this VPC.
  * `cluster_deploy_type` - Deployment type of the cluster.
  * `cluster_desc` - Description of the cluster.
  * `cluster_external_endpoint` - External network address to access.
  * `cluster_extra_args` - Customized parameters for master component.
    * `kube_apiserver` - The customized parameters for kube-apiserver.
    * `kube_controller_manager` - The customized parameters for kube-controller-manager.
    * `kube_scheduler` - The customized parameters for kube-scheduler.
  * `cluster_id` - ID of cluster.
  * `cluster_ipvs` - Indicates whether ipvs is enabled.
  * `cluster_max_pod_num` - The maximum number of Pods per node in the cluster.
  * `cluster_max_service_num` - The maximum number of services in the cluster.
  * `cluster_name` - Name of the cluster.
  * `cluster_node_num` - Number of nodes in the cluster.
  * `cluster_os` - Operating system of the cluster.
  * `cluster_version` - Version of the cluster.
  * `container_runtime` - (**Deprecated**) It has been deprecated from version 1.18.1. Container runtime of the cluster.
  * `deletion_protection` - Indicates whether cluster deletion protection is enabled.
  * `domain` - Domain name for access.
  * `eni_subnet_ids` - Subnet IDs for cluster with VPC-CNI network mode.
  * `ignore_cluster_cidr_conflict` - Indicates whether to ignore the cluster cidr conflict error.
  * `is_non_static_ip_mode` - Indicates whether non-static ip mode is enabled.
  * `kube_config_intranet` - Kubernetes config of private network.
  * `kube_config` - Kubernetes config.
  * `kube_proxy_mode` - Cluster kube-proxy mode.
  * `network_type` - Cluster network type.
  * `node_name_type` - Node name type of cluster.
  * `password` - Password of account.
  * `pgw_endpoint` - The Intranet address used for access.
  * `project_id` - Project ID of the cluster.
  * `security_policy` - Access policy.
  * `service_cidr` - The network address block of the cluster.
  * `tags` - Tags of the cluster.
  * `user_name` - User name of account.
  * `vpc_cni_type` - Distinguish between shared network card multi-IP mode and independent network card mode.
  * `vpc_id` - Vpc ID of the cluster.
  * `worker_instances_list` - An information list of cvm within the WORKER clusters. Each element contains the following attributes.
    * `failed_reason` - Information of the cvm when it is failed.
    * `instance_id` - ID of the cvm.
    * `instance_role` - Role of the cvm.
    * `instance_state` - State of the cvm.
    * `lan_ip` - LAN IP of the cvm.


