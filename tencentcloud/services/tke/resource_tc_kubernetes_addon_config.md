Provide a resource to configure addon that kubernetes comes with.

Example Usage

Update cluster-autoscaler addon

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
