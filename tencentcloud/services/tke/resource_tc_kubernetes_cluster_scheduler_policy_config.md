Provides a resource to manage TKE cluster scheduler policy configuration (singleton per cluster).

Example Usage

```hcl
resource "tencentcloud_kubernetes_cluster_scheduler_policy_config" "example" {
  cluster_id       = "cls-man1vvi2"
  high_performance = true
  client_connection {
    burst = 100
    qps   = 50
  }

  scheduler_policy_config {
    scheduler_name = "default-scheduler"
    plugin_configs {
      name = "NodeResourcesFit"
      args = jsonencode({
        apiVersion = "kubescheduler.config.k8s.io/v1"
        kind       = "NodeResourcesFitArgs"
        scoringStrategy = {
          resources = [
            {
              name   = "cpu"
              weight = 1
            },
            {
              name   = "memory"
              weight = 1
            },
          ]
          type = "LeastAllocated"
        }
      })
    }

    plugin_set {
      disabled {
        name   = "Coscheduling"
        weight = 0
      }
      disabled {
        name   = "PlacementPolicy"
        weight = 0
      }
      disabled {
        name   = "SafeBalance"
        weight = 0
      }
      disabled {
        name   = "qGPU"
        weight = 0
      }
      enabled {
        name   = "DefaultPreemption"
        weight = 0
      }
      enabled {
        name   = "ImageLocality"
        weight = 1
      }
      enabled {
        name   = "InterPodAffinity"
        weight = 2
      }
      enabled {
        name   = "NodeAffinity"
        weight = 2
      }
      enabled {
        name   = "NodeResourcesBalancedAllocation"
        weight = 1
      }
      enabled {
        name   = "NodeResourcesFit"
        weight = 1
      }
      enabled {
        name   = "PodTopologySpread"
        weight = 2
      }
      enabled {
        name   = "TaintToleration"
        weight = 3
      }
    }
  }
}
```

Import

TKE cluster scheduler policy config can be imported using the cluster ID, e.g.

```
terraform import tencentcloud_kubernetes_cluster_scheduler_policy_config.example cls-man1vvi2
```
