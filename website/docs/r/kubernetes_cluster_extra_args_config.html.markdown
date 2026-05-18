---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_extra_args_config"
sidebar_current: "docs-tencentcloud-resource-kubernetes_cluster_extra_args_config"
description: |-
  Provides a resource to manage TKE cluster extra args configuration.
---

# tencentcloud_kubernetes_cluster_extra_args_config

Provides a resource to manage TKE cluster extra args configuration.

~> **NOTE:** This resource must exclusive in one cluster, do not declare additional args resources of this extra args elsewhere.

## Example Usage

```hcl
resource "tencentcloud_kubernetes_cluster_extra_args_config" "example" {
  cluster_id = "cls-man1vvi2"
  kube_apiserver = [
    "goaway-chance=0",
    "kubelet-preferred-address-types=Hostname"
  ]

  kube_controller_manager = [
    "concurrent-serviceaccount-token-syncs=5"
  ]

  kube_scheduler = [
    "kube-api-qps=50"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID. Only managed clusters are supported.
* `etcd` - (Optional, List: [`String`]) Custom args for etcd. Only standalone clusters are supported, format: ["k1=v1", "k2=v2"].
* `kube_apiserver` - (Optional, List: [`String`]) Custom args for kube-apiserver, format: ["k1=v1", "k2=v2"], e.g. ["max-requests-inflight=500","feature-gates=PodShareProcessNamespace=true"].
* `kube_controller_manager` - (Optional, List: [`String`]) Custom args for kube-controller-manager, format: ["k1=v1", "k2=v2"].
* `kube_scheduler` - (Optional, List: [`String`]) Custom args for kube-scheduler, format: ["k1=v1", "k2=v2"].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `10m`) Used when creating the resource.
* `update` - (Defaults to `10m`) Used when updating the resource.

## Import

TKE cluster extra args config can be imported using the clusterId, e.g.

```
terraform import tencentcloud_kubernetes_cluster_extra_args_config.example cls-man1vvi2
```

