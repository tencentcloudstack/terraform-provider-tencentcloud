---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_endpoint"
sidebar_current: "docs-tencentcloud-resource-kubernetes_cluster_endpoint"
description: |-
  Provide a resource to create a KubernetesClusterEndpoint.
This resource allows you to create an empty cluster first without any workers. Only all attached node depends create complete, cluster endpoint will finally be enabled.
---

# tencentcloud_kubernetes_cluster_endpoint

Provide a resource to create a KubernetesClusterEndpoint.
This resource allows you to create an empty cluster first without any workers. Only all attached node depends create complete, cluster endpoint will finally be enabled.

~> **NOTE:** Recommend using `depends_on` to make sure endpoint create after node pools or workers does.

## Example Usage

```hcl
resource "tencentcloud_kubernetes_node_pool" "pool1" {}

resource "tencentcloud_kubernetes_cluster_endpoint" "foo" {
  cluster_id       = "cls-xxxxxxxx"
  cluster_internet = true
  cluster_intranet = true
  # managed_cluster_internet_security_policies = [
  "192.168.0.0/24"
]
cluster_intranet_subnet_id = "subnet-xxxxxxxx"
depends_on = [
  tencentcloud_kubernetes_node_pool.pool1
]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Specify cluster ID.
* `cluster_internet_security_group` - (Optional, String) Specify security group, NOTE: This argument must not be empty if cluster internet enabled.
* `cluster_internet` - (Optional, Bool) Open internet access or not.
* `cluster_intranet_subnet_id` - (Optional, String) Subnet id who can access this independent cluster, this field must and can only set  when `cluster_intranet` is true. `cluster_intranet_subnet_id` can not modify once be set.
* `cluster_intranet` - (Optional, Bool) Open intranet access or not.
* `managed_cluster_internet_security_policies` - (Optional, List: [`String`], **Deprecated**) this argument was deprecated, use `cluster_internet_security_group` instead. Security policies for managed cluster internet, like:'192.168.1.0/24' or '113.116.51.27', '0.0.0.0/0' means all. This field can only set when field `cluster_deploy_type` is 'MANAGED_CLUSTER' and `cluster_internet` is true. `managed_cluster_internet_security_policies` can not delete or empty once be set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `certification_authority` - The certificate used for access.
* `cluster_deploy_type` - Cluster deploy type of `MANAGED_CLUSTER` or `INDEPENDENT_CLUSTER`.
* `cluster_external_endpoint` - External network address to access.
* `domain` - Domain name for access.
* `password` - Password of account.
* `pgw_endpoint` - The Intranet address used for access.
* `user_name` - User name of account.


## Import

KubernetesClusterEndpoint instance can be imported by passing cluster id, e.g.
```
$ terraform import tencentcloud_kubernetes_cluster_endpoint.test cluster-id
```

