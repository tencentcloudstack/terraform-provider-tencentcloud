---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eks_clusters"
sidebar_current: "docs-tencentcloud-datasource-eks_clusters"
description: |-
  Use this data source to query elastic kubernetes cluster resource.
---

# tencentcloud_eks_clusters

Use this data source to query elastic kubernetes cluster resource.

## Example Usage

```hcl
data "tencentcloud_eks_clusters" "foo" {
  cluster_id = "cls-xxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional) ID of the cluster. Conflict with cluster_name, can not be set at the same time.
* `cluster_name` - (Optional) Name of the cluster. Conflict with cluster_id, can not be set at the same time.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - EKS cluster list.
  * `cluster_desc` - Description of the cluster.
  * `cluster_id` - ID of the cluster.
  * `cluster_name` - Name of the cluster.
  * `created_time` - Create time of the clusters.
  * `dns_servers` - List of cluster custom DNS Server info.
    * `domain` - DNS Server domain. Empty indicates all domain.
    * `servers` - List of DNS Server IP address.
  * `enable_vpc_core_dns` - Indicates whether to enable dns in user cluster, default value is `true`.
  * `k8s_version` - EKS cluster kubernetes version.
  * `need_delete_cbs` - Indicates whether to delete CBS after EKS cluster remove.
  * `service_subnet_id` - Subnet id of service.
  * `status` - EKS status.
  * `subnet_ids` - Subnet id list.
  * `tags` - Tags of EKS cluster.
  * `vpc_id` - Vpc id.


