---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_addon"
sidebar_current: "docs-tencentcloud-resource-kubernetes_addon"
description: |-
  Provide a resource to configure kubernetes cluster app addons.
---

# tencentcloud_kubernetes_addon

Provide a resource to configure kubernetes cluster app addons.

## Example Usage

### Install cos addon

```hcl
resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                  = "vpc-xxxxxxxx"
  cluster_cidr            = "10.31.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "tf_example_cluster"
  cluster_desc            = "example for tke cluster"
  cluster_max_service_num = 32
  cluster_internet        = false # (can be ignored) open it after the nodes added
  cluster_version         = "1.22.5"
  cluster_deploy_type     = "MANAGED_CLUSTER"
  # without any worker config
}

resource "tencentcloud_kubernetes_addon" "kubernetes_addon" {
  cluster_id    = tencentcloud_kubernetes_cluster.example.id
  addon_name    = "cos"
  addon_version = "2018-05-25"
  raw_values    = "e30="
}
```

## Argument Reference

The following arguments are supported:

* `addon_name` - (Required, String, ForceNew) Name of addon.
* `cluster_id` - (Required, String, ForceNew) ID of cluster.
* `addon_version` - (Optional, String) Version of addon.
* `raw_values` - (Optional, String) Params of addon, base64 encoded json format.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `phase` - Status of addon.
* `reason` - Reason of addon failed.


## Import

Addon can be imported by using cluster_id#addon_name
```
$ terraform import tencentcloud_kubernetes_addon.addon_cos cls-xxx#addon_name
```

