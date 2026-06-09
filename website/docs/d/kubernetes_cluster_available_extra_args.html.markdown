---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_available_extra_args"
sidebar_current: "docs-tencentcloud-datasource-kubernetes_cluster_available_extra_args"
description: |-
  Use this data source to query the available custom extra arguments for TKE cluster components.
---

# tencentcloud_kubernetes_cluster_available_extra_args

Use this data source to query the available custom extra arguments for TKE cluster components.

## Example Usage

### Query available extra args for a managed cluster

```hcl
data "tencentcloud_kubernetes_cluster_available_extra_args" "example" {
  cluster_version = "1.34.1"
  cluster_type    = "MANAGED_CLUSTER"
}
```

### Query available extra args for an independent cluster

```hcl
data "tencentcloud_kubernetes_cluster_available_extra_args" "example" {
  cluster_version = "1.30.0"
  cluster_type    = "INDEPENDENT_CLUSTER"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_type` - (Required, String) Cluster type. Valid values: `MANAGED_CLUSTER`, `INDEPENDENT_CLUSTER`.
* `cluster_version` - (Required, String) Cluster version, e.g. `1.28.3`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `available_extra_args` - Available custom extra arguments for cluster components.
  * `kube_apiserver` - Available custom arguments for kube-apiserver.
    * `constraint` - Valid range or allowed values of the argument.
    * `default` - Default value of the argument.
    * `name` - Argument name.
    * `type` - Argument type.
    * `usage` - Argument description.
  * `kube_controller_manager` - Available custom arguments for kube-controller-manager.
    * `constraint` - Valid range or allowed values of the argument.
    * `default` - Default value of the argument.
    * `name` - Argument name.
    * `type` - Argument type.
    * `usage` - Argument description.
  * `kube_scheduler` - Available custom arguments for kube-scheduler.
    * `constraint` - Valid range or allowed values of the argument.
    * `default` - Default value of the argument.
    * `name` - Argument name.
    * `type` - Argument type.
    * `usage` - Argument description.
  * `kubelet` - Available custom arguments for kubelet.
    * `constraint` - Valid range or allowed values of the argument.
    * `default` - Default value of the argument.
    * `name` - Argument name.
    * `type` - Argument type.
    * `usage` - Argument description.


