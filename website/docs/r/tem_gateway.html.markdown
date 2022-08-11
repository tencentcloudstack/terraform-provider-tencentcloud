---
subcategory: "TEM"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tem_gateway"
sidebar_current: "docs-tencentcloud-resource-tem_gateway"
description: |-
  Provides a resource to create a tem gateway
---

# tencentcloud_tem_gateway

Provides a resource to create a tem gateway

## Example Usage

```hcl
resource "tencentcloud_tem_gateway" "gateway" {
  ingress {
    ingress_name        = "demo"
    environment_id      = "en-853mggjm"
    cluster_namespace   = "default"
    address_i_p_version = "IPV4"
    rewrite_type        = "NONE"
    mixed               = false
    rules {
      host     = "test.com"
      protocol = "http"
      http {
        paths {
          path = "/"
          backend {
            service_name = "demo"
            service_port = 80
          }
        }
      }
    }
    rules {
      host     = "hello.com"
      protocol = "http"
      http {
        paths {
          path = "/"
          backend {
            service_name = "hello"
            service_port = 36000
          }
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `ingress` - (Optional, List) gateway properties.

The `backend` object supports the following:

* `service_name` - (Required, String) backend name.
* `service_port` - (Required, Int) backend port.

The `http` object supports the following:

* `paths` - (Required, List) path payload.

The `ingress` object supports the following:

* `address_i_p_version` - (Required, String) ip version, IPV4 only.
* `cluster_namespace` - (Required, String) inner namespace, default only.
* `environment_id` - (Required, String) environment ID.
* `ingress_name` - (Required, String) gateway name.
* `mixed` - (Required, Bool) vpc ID.
* `rules` - (Required, List) proxy rules.
* `rewrite_type` - (Optional, String) redirect http to https.
* `tls` - (Optional, List) ingress TLS configurations.

The `paths` object supports the following:

* `backend` - (Required, List) backend payload.
* `path` - (Required, String) path.

The `rules` object supports the following:

* `http` - (Required, List) rule payload.
* `host` - (Optional, String) host name.
* `protocol` - (Optional, String) protocol.

The `tls` object supports the following:

* `certificate_id` - (Required, String) certificate ID.
* `hosts` - (Required, Set) host names.
* `secret_name` - (Optional, String) secret name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tem gateway can be imported using the id, e.g.
```
$ terraform import tencentcloud_tem_gateway.gateway gateway_id
```

