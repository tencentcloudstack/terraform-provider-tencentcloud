---
subcategory: "Global Accelerator(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_listener"
sidebar_current: "docs-tencentcloud-resource-ga2_listener"
description: |-
  Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) listener.
---

# tencentcloud_ga2_listener

Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) listener.

## Example Usage

```hcl
resource "tencentcloud_ga2_global_accelerator" "example" {
  name                 = "tf-example"
  instance_charge_type = "POSTPAID"
  description          = "tf example global accelerator"

  tags = {
    createdBy = "Terraform"
  }
}

resource "tencentcloud_ga2_accelerate_area" "example" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
  accelerate_region     = "ap-guangzhou"
  bandwidth             = 10
  isp_type              = "BGP"
  ip_version            = "IPv4"
}

resource "tencentcloud_ga2_listener" "example1" {
  global_accelerator_id = tencentcloud_ga2_accelerate_area.example.global_accelerator_id
  name                  = "tf-example-tcp"
  protocol              = "TCP"

  port_ranges {
    from_port = 80
    to_port   = 80
  }

  description      = "tf example listener"
  get_real_ip_type = "ProxyProtocol"
  client_affinity  = "Open"
  listener_type    = "Standard"
  idle_timeout     = 900
}

resource "tencentcloud_ga2_listener" "example2" {
  global_accelerator_id = tencentcloud_ga2_accelerate_area.example.global_accelerator_id
  name                  = "tf-example-http"
  protocol              = "HTTP"

  port_ranges {
    from_port = 90
    to_port   = 90
  }

  description             = "tf example listener"
  idle_timeout            = 15
  request_timeout         = 60
  listener_type           = "Standard"
  x_forwarded_for_real_ip = true

  depends_on = [tencentcloud_ga2_listener.example1]
}
```

## Argument Reference

The following arguments are supported:

* `global_accelerator_id` - (Required, String, ForceNew) Global accelerator instance ID this listener belongs to.
* `port_ranges` - (Required, List, ForceNew) Listening port range. Cannot be modified after creation; modifying it forces a new resource.
* `certification_type` - (Optional, String) SSL authentication mode. Valid values: `UNIDIRECTIONAL`, `MUTUAL`.
* `cipher_policy_id` - (Optional, String) TLS cipher policy ID.
* `client_affinity_time` - (Optional, Int) Session-stickiness duration in seconds. NOTE: this field is silently ignored on Create (the SDK CreateListener API has no equivalent slot) and forwarded only on Update via ModifyListener.
* `client_affinity` - (Optional, String) Whether to enable session stickiness.
* `client_ca_certificates` - (Optional, Set: [`String`]) Client CA certificate ID list. Treated as an unordered set; HCL element order has no semantic meaning.
* `description` - (Optional, String) Listener description. Maximum length is 100 bytes.
* `get_real_ip_type` - (Optional, String) Layer-4 real-IP method. Valid values: `TOA`, `ProxyProtocol`.
* `idle_timeout` - (Optional, Int) Connection idle timeout in seconds.
* `listener_type` - (Optional, String, ForceNew) Listener routing type. Defaults to smart routing. Cannot be modified after creation.
* `name` - (Optional, String) Listener name. Maximum length is 60 bytes.
* `protocol` - (Optional, String, ForceNew) Listener protocol. Valid values: `TCP`, `UDP`, `HTTP`, `HTTPS`. Default: `TCP`. Cannot be modified after creation.
* `request_timeout` - (Optional, Int) Request timeout in seconds.
* `server_certificates` - (Optional, Set: [`String`]) Server certificate ID list. Treated as an unordered set; HCL element order has no semantic meaning.
* `x_forwarded_for_real_ip` - (Optional, Bool) Whether to enable layer-7 real-IP forwarding.

The `port_ranges` object supports the following:

* `from_port` - (Required, Int) Inclusive start port.
* `to_port` - (Required, Int) Inclusive end port.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Listener creation time.
* `endpoint_group_counts` - Number of endpoint groups attached to this listener.
* `http_version` - HTTP version negotiated for this listener.
* `listener_id` - Listener instance ID.
* `status` - Listener operational status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `5m`) Used when creating the resource.
* `update` - (Defaults to `5m`) Used when updating the resource.
* `delete` - (Defaults to `5m`) Used when deleting the resource.

## Import

GA2 listener can be imported using the composite id `<global_accelerator_id>#<listener_id>`, e.g.

```
terraform import tencentcloud_ga2_listener.example ga-4mredmiu#lsr-llr0dng1
```

