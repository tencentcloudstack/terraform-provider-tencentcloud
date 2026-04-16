---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_scheduler_policy_config"
sidebar_current: "docs-tencentcloud-resource-kubernetes_cluster_scheduler_policy_config"
description: |-
  Provides a resource to manage TKE cluster scheduler policy configuration (singleton per cluster).
---

# tencentcloud_kubernetes_cluster_scheduler_policy_config

Provides a resource to manage TKE cluster scheduler policy configuration (singleton per cluster).

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `client_connection` - (Optional, List) Client connection configuration.
* `extenders` - (Optional, List) Extender scheduler configuration list.
* `high_performance` - (Optional, Bool) High performance mode switch.
* `scheduler_policy_config` - (Optional, List) Scheduler policy configuration list.

The `client_connection` object supports the following:

* `burst` - (Optional, Int) Burst request limit.
* `qps` - (Optional, Float64) Maximum queries per second.

The `disabled` object of `plugin_set` supports the following:

* `name` - (Required, String) Plugin name.
* `weight` - (Optional, Int) Plugin weight.

The `enabled` object of `plugin_set` supports the following:

* `name` - (Required, String) Plugin name.
* `weight` - (Optional, Int) Plugin weight.

The `extender_client_config` object of `extenders` supports the following:

* `service` - (Optional, List) Service reference configuration.

The `extenders` object supports the following:

* `extender_client_config` - (Optional, List) Extender client configuration.
* `filter_verb` - (Optional, String) Filter stage interface.
* `node_cache_capable` - (Optional, Bool) Whether node cache capability is enabled.
* `preempt_verb` - (Optional, String) Preempt stage interface.
* `prioritize_verb` - (Optional, String) Prioritize stage interface.
* `weight` - (Optional, Int) Weight for prioritize stage.

The `plugin_configs` object of `scheduler_policy_config` supports the following:

* `args` - (Optional, String) Plugin args in raw JSON format. Terraform will automatically base64-encode it before calling the API and decode it on read.
* `name` - (Optional, String) Plugin name.

The `plugin_set` object of `scheduler_policy_config` supports the following:

* `disabled` - (Optional, Set) List of plugins to disable.
* `enabled` - (Optional, Set) List of plugins to enable.

The `scheduler_policy_config` object supports the following:

* `plugin_configs` - (Optional, List) Scheduler plugin configuration list.
* `plugin_set` - (Optional, List) Plugin set configuration.
* `scheduler_name` - (Optional, String) Scheduler name.

The `service` object of `extender_client_config` supports the following:

* `name` - (Optional, String) Service name.
* `namespace` - (Optional, String) Service namespace.
* `path` - (Optional, String) Service path.
* `port` - (Optional, Int) Service port.
* `scheme` - (Optional, String) Service protocol scheme (e.g. http, https).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `policy` - Raw scheduler policy JSON string.


## Import

TKE cluster scheduler policy config can be imported using the cluster ID, e.g.

```
terraform import tencentcloud_kubernetes_cluster_scheduler_policy_config.example cls-man1vvi2
```

