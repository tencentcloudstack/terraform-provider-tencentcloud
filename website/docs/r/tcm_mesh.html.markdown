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
resource "tencentcloud_tcm_mesh" "basic" {
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
    }
  }
  tag_list {
    key         = "key"
    value       = "value"
    passthrough = true
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

The `config` object supports the following:

* `istio` - (Optional, List) Istio configuration.
* `prometheus` - (Optional, List) Prometheus configuration.
* `tracing` - (Optional, List) Tracing config.

The `custom_prom` object supports the following:

* `auth_type` - (Required, String) Authentication type of the prometheus.
* `url` - (Required, String) Url of the prometheus.
* `is_public_addr` - (Optional, Bool) Whether it is public address, default false.
* `password` - (Optional, String) Password of the prometheus, used in basic authentication type.
* `username` - (Optional, String) Username of the prometheus, used in basic authentication type.
* `vpc_id` - (Optional, String) Vpc id.

The `istio` object supports the following:

* `outbound_traffic_policy` - (Required, String) Outbound traffic policy.
* `disable_http_retry` - (Optional, Bool) Disable http retry.
* `disable_policy_checks` - (Optional, Bool) Disable policy checks.
* `enable_pilot_http` - (Optional, Bool) Enable HTTP/1.0 support.
* `smart_dns` - (Optional, List) SmartDNS configuration.

The `prometheus` object supports the following:

* `custom_prom` - (Optional, List) Custom prometheus.
* `instance_id` - (Optional, String) Instance id.
* `region` - (Optional, String) Region.
* `subnet_id` - (Optional, String) Subnet id.
* `vpc_id` - (Optional, String) Vpc id.

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

