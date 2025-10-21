---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_addon_config"
sidebar_current: "docs-tencentcloud-resource-kubernetes_addon_config"
description: |-
  Provide a resource to configure addon that kubernetes comes with.
---

# tencentcloud_kubernetes_addon_config

Provide a resource to configure addon that kubernetes comes with.

## Example Usage

### Update cluster-autoscaler addon

```hcl
resource "tencentcloud_kubernetes_addon_config" "kubernetes_addon_config" {
  cluster_id = "cls-xxxxxx"
  addon_name = "cluster-autoscaler"
  raw_values = "{\"extraArgs\":{\"scale-down-enabled\":true,\"max-empty-bulk-delete\":11,\"scale-down-delay-after-add\":\"10mm\",\"scale-down-unneeded-time\":\"10mm\",\"scale-down-utilization-threshold\":0.005,\"ignore-daemonsets-utilization\":false,\"skip-nodes-with-local-storage\":true,\"skip-nodes-with-system-pods\":true}}"
}
`
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


