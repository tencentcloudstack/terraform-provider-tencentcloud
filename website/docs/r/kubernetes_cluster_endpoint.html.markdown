---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_endpoint"
sidebar_current: "docs-tencentcloud-resource-kubernetes_cluster_endpoint"
description: |-
  Provide a resource to create a kubernetes cluster endpoint. This resource allows you to create an empty cluster first without any workers. Only all attached node depends create complete, cluster endpoint will finally be enabled.
---

# tencentcloud_kubernetes_cluster_endpoint

Provide a resource to create a kubernetes cluster endpoint. This resource allows you to create an empty cluster first without any workers. Only all attached node depends create complete, cluster endpoint will finally be enabled.

~> **NOTE:** Recommend using `depends_on` to make sure endpoint create after node pools or workers does.

~> **NOTE:** Please do not use this resource and resource `tencentcloud_kubernetes_cluster` to operate cluster public network/intranet access at the same time.

## Example Usage

### Open intranet access for kubernetes cluster

```hcl
resource "tencentcloud_kubernetes_cluster_endpoint" "example" {
  cluster_id                 = "cls-fdy7hm1q"
  cluster_intranet           = true
  cluster_intranet_subnet_id = "subnet-7nl0sswi"
  cluster_intranet_domain    = "intranet_demo.com"
}
```

### Open internet access for kubernetes cluster

```hcl
resource "tencentcloud_kubernetes_cluster_endpoint" "example" {
  cluster_id                      = "cls-fdy7hm1q"
  cluster_internet                = true
  cluster_internet_security_group = "sg-e6a8xxib"
  cluster_internet_domain         = "internet_demo.com"
  extensive_parameters = jsonencode({
    "AddressIPVersion" : "IPV4",
    "InternetAccessible" : {
      "InternetChargeType" : "TRAFFIC_POSTPAID_BY_HOUR",
      "InternetMaxBandwidthOut" : 10
    }
  })
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Specify cluster ID.
* `cluster_internet_domain` - (Optional, String) Domain name for cluster Kube-apiserver internet access.  Be careful if you modify value of this parameter, the cluster_external_endpoint value may be changed automatically too.
* `cluster_internet_security_group` - (Optional, String) Security group ID for internet cluster endpoint. NOTE: This argument must not be empty if cluster internet enabled.
* `cluster_internet` - (Optional, Bool) Open internet access or not.
* `cluster_intranet_domain` - (Optional, String) Domain name for cluster Kube-apiserver intranet access. Be careful if you modify value of this parameter, the pgw_endpoint value may be changed automatically too.
* `cluster_intranet_security_group` - (Optional, String, ForceNew) Security group ID for intranet cluster endpoint.
* `cluster_intranet_subnet_id` - (Optional, String) Subnet id who can access this independent cluster, this field must and can only set  when `cluster_intranet` is true. `cluster_intranet_subnet_id` can not modify once be set.
* `cluster_intranet` - (Optional, Bool) Open intranet access or not.
* `existed_load_balancer_id` - (Optional, String, ForceNew) Enable internal or external access using an existing CLB.
* `extensive_parameters` - (Optional, String, ForceNew) The LB parameter. Only used for public network access.
* `managed_cluster_internet_security_policies` - (Optional, List: [`String`], **Deprecated**) this argument was deprecated, use `cluster_internet_security_group` instead. Security policies for managed cluster internet, like:'192.168.1.0/24' or '113.116.51.27', '0.0.0.0/0' means all. This field can only set when field `cluster_deploy_type` is 'MANAGED_CLUSTER' and `cluster_internet` is true. `managed_cluster_internet_security_policies` can not delete or empty once be set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `certification_authority` - The certificate used for access.
* `cluster_deploy_type` - Cluster deploy type of `MANAGED_CLUSTER` or `INDEPENDENT_CLUSTER`.
* `cluster_external_endpoint` - External network address to access.
* `domain` - Domain name for access.
* `kube_config_intranet` - Kubernetes config of private network.
* `kube_config` - The Intranet address used for access.
* `password` - Password of account.
* `pgw_endpoint` - The Intranet address used for access.
* `user_name` - User name of account.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeClusterInstances](https://cloud.tencent.com/document/api/457/31863) - This API is used to query information of node instances in a cluster.
* Create: [DescribeClusterVirtualNode](https://cloud.tencent.com/document/api/457/85351) - This API is used to view the Super Node list.
* Read: [DescribeClusterKubeconfig](https://cloud.tencent.com/document/api/457/48149) - This API is used to obtain the cluster kubeconfig file. Different sub-accounts have their own kubeconfig files. The kubeconfig file contains the kube-apiserver client certificate of the corresponding sub-account. By default, the client certificate is created when this API is called for the first time, and the certificate is valid for 20 years with no permissions granted. For the cluster owner or primary account, the cluster-admin permission is granted by default.
* Read: [DescribeClusterSecurity](https://cloud.tencent.com/document/api/457/36703) - This API is used to query the key information of a cluster.
* Read: [DescribeClusters](https://cloud.tencent.com/document/api/457/31862) - This API is used to query clusters list.
* Update: [ModifyClusterEndpointSP](https://cloud.tencent.com/document/api/457/39408) - Modify the security policy of the external port of the managed cluster (the old way, only the external port of the managed cluster is supported)
* Delete: [DescribeClusterSecurity](https://cloud.tencent.com/document/api/457/36703) - This API is used to query the key information of a cluster.

