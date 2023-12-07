Provides a resource to create a tcm mesh

Example Usage

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
Import

tcm mesh can be imported using the id, e.g.
```
$ terraform import tencentcloud_tcm_mesh.mesh mesh_id
```