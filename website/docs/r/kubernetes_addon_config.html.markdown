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
resource "tencentcloud_kubernetes_addon_config" "example" {
  cluster_id = "cls-5yezvaxo"
  addon_name = "cluster-autoscaler"
  raw_values = jsonencode({
    "autoDiscovery" : {
      "labels" : [
        {
          "node.tke.cloud.tencent.com/autoscaling-enabled" : "true"
        }
      ]
    },
    "extraArgs" : {
      "expander" : "random",
      "ignore-daemonsets-utilization" : false,
      "ignore-taint_1" : "tke.cloud.tencent.com/direct-eni-unavailable",
      "ignore-taint_2" : "tke.cloud.tencent.com/eni-ip-unavailable",
      "ignore-taint_3" : "tke.cloud.tencent.com/uninitialized",
      "ignore-taint_4" : "tke.cloud.tencent.com/no-aia-ip",
      "scale-down-unready-time" : "20m0s",
      "scale-down-utilization-threshold" : 0.005,
      "skip-nodes-with-local-storage" : true,
      "scale-down-delay-after-add" : "10mm",
      "scale-down-enabled" : true,
      "scale-down-unneeded-time" : "10mm",
      "skip-nodes-with-system-pods" : true
      "max-empty-bulk-delete" : 11,
      "max-nodes-total" : 5,
      "max-total-unready-percentage" : 33,
      "ok-total-unready-count" : 3,
    },
    "image" : {
      "repository" : "ccr.ccs.tencentyun.com/tkeimages/cluster-autoscaler"
    },
    "resources" : {
      "limits" : {
        "cpu" : "2",
        "memory" : "4Gi"
      },
      "requests" : {
        "cpu" : "200m",
        "memory" : "256Mi"
      }
    }
  })
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


