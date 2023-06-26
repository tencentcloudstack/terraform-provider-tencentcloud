---
subcategory: "TencentCloud ServiceMesh(TCM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcm_mesh"
sidebar_current: "docs-tencentcloud-resource-tcm_mesh"
description: |-
  Provides a resource to create a tcm mesh
---

# tencentcloud_tcm_mesh

Provides a resource to create a tcm mesh

## Example Usage

```hcl
resource "tencentcloud_tcm_mesh" "mesh" {
  display_name = "test_mesh"
  mesh_version = "1.12.5"
  type         = "HOSTED"
  config {
    istio {
      outbound_traffic_policy = "ALLOW_ANY"
      disable_policy_checks   = true
      enable_pilot_http       = true
      disable_http_retry      = true
      smart_dns {
        istio_meta_dns_capture       = true
        istio_meta_dns_auto_allocate = true
      }
      tracing {
        enable = false
      }
    }
    tracing {
      enable   = true
      sampling = 1
      apm {
        enable = true
        region = "ap-guangzhou"
      }
    }
    prometheus {
      custom_prom {
        url       = "https://10.0.0.1:1000"
        auth_type = "none"
        vpc_id    = "vpc-j9yhbzpn"
      }
    }
    inject {
      exclude_ip_ranges                   = ["172.16.0.0/16"]
      hold_application_until_proxy_starts = true
      hold_proxy_until_application_ends   = true
    }

    sidecar_resources {
      limits {
        name     = "cpu"
        quantity = "2"
      }
      limits {
        name     = "memory"
        quantity = "1Gi"
      }
      requests {
        name     = "cpu"
        quantity = "100m"
      }
      requests {
        name     = "memory"
        quantity = "128Mi"
      }
    }
  }
  tag_list {
    key         = "key"
    value       = "value"
    passthrough = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `config` - (Required, List) Mesh configuration.
* `display_name` - (Required, String) Mesh name.
* `mesh_version` - (Required, String) Mesh version.
* `type` - (Required, String) Mesh type.
* `mesh_id` - (Optional, String) Mesh ID.
* `tag_list` - (Optional, List) A list of associated tags.

The `apm` object supports the following:

* `enable` - (Optional, Bool) Whether enable APM.
* `instance_id` - (Optional, String) Instance id of the APM.
* `region` - (Optional, String) Region.

The `apm` object supports the following:

* `enable` - (Required, Bool) Whether enable APM.
* `instance_id` - (Optional, String) Instance id of the APM.
* `region` - (Optional, String) Region.

The `config` object supports the following:

* `inject` - (Optional, List) Sidecar inject configuration.
* `istio` - (Optional, List) Istio configuration.
* `prometheus` - (Optional, List) Prometheus configuration.
* `sidecar_resources` - (Optional, List) Default sidecar requests and limits.
* `tracing` - (Optional, List) Tracing config.

The `custom_prom` object supports the following:

* `auth_type` - (Required, String) Authentication type of the prometheus.
* `url` - (Required, String) Url of the prometheus.
* `is_public_addr` - (Optional, Bool) Whether it is public address, default false.
* `password` - (Optional, String) Password of the prometheus, used in basic authentication type.
* `username` - (Optional, String) Username of the prometheus, used in basic authentication type.
* `vpc_id` - (Optional, String) Vpc id.

The `inject` object supports the following:

* `exclude_ip_ranges` - (Optional, Set) IP ranges that should not be proxied.
* `hold_application_until_proxy_starts` - (Optional, Bool) Let istio-proxy(sidecar) start first, before app container.
* `hold_proxy_until_application_ends` - (Optional, Bool) Let istio-proxy(sidecar) stop last, after app container.

The `istio` object supports the following:

* `outbound_traffic_policy` - (Required, String) Outbound traffic policy, REGISTRY_ONLY or ALLOW_ANY, see https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#MeshConfig-OutboundTrafficPolicy-Mode.
* `disable_http_retry` - (Optional, Bool) Disable http retry.
* `disable_policy_checks` - (Optional, Bool) Disable policy checks.
* `enable_pilot_http` - (Optional, Bool) Enable HTTP/1.0 support.
* `smart_dns` - (Optional, List) SmartDNS configuration.
* `tracing` - (Optional, List) Tracing config(Deprecated, please use MeshConfig.Tracing for configuration).

The `limits` object supports the following:

* `name` - (Optional, String) Resource type name, `cpu/memory`.
* `quantity` - (Optional, String) Resource quantity, example: cpu-`100m`, memory-`1Gi`.

The `prometheus` object supports the following:

* `custom_prom` - (Optional, List) Custom prometheus.
* `instance_id` - (Optional, String) Instance id.
* `region` - (Optional, String) Region.
* `subnet_id` - (Optional, String) Subnet id.
* `vpc_id` - (Optional, String) Vpc id.

The `requests` object supports the following:

* `name` - (Optional, String) Resource type name, `cpu/memory`.
* `quantity` - (Optional, String) Resource quantity, example: cpu-`100m`, memory-`1Gi`.

The `sidecar_resources` object supports the following:

* `limits` - (Optional, Set) Sidecar limits.
* `requests` - (Optional, Set) Sidecar requests.

The `smart_dns` object supports the following:

* `istio_meta_dns_auto_allocate` - (Optional, Bool) Enable auto allocate address.
* `istio_meta_dns_capture` - (Optional, Bool) Enable dns proxy.

The `tag_list` object supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, String) Tag value.
* `passthrough` - (Optional, Bool) Passthrough to other related product.

The `tracing` object supports the following:

* `apm` - (Optional, List) APM config.
* `enable` - (Optional, Bool) Whether enable tracing.
* `sampling` - (Optional, Float64) Tracing sampling, 0.0-1.0.
* `zipkin` - (Optional, List) Third party zipkin config.

The `zipkin` object supports the following:

* `address` - (Required, String) Zipkin address.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tcm mesh can be imported using the id, e.g.
```
$ terraform import tencentcloud_tcm_mesh.mesh mesh_id
```

