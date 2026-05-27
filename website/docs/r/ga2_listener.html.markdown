---
subcategory: "Global Accelerator 2(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_listener"
sidebar_current: "docs-tencentcloud-resource-ga2_listener"
description: |-
  Provides a resource to create a GA2 (Global Accelerator 2) listener.
---

# tencentcloud_ga2_listener

Provides a resource to create a GA2 (Global Accelerator 2) listener.

## Example Usage

```hcl
resource "tencentcloud_ga2_listener" "example" {
  global_accelerator_id = "ga2-xxxxxxxx"
  name                  = "tf-example"
  protocol              = "TCP"
  listener_type         = "INTELLIGENT"

  port_ranges {
    from_port = 80
    to_port   = 80
  }

  description      = "tf example listener"
  idle_timeout     = 900
  client_affinity  = "SOURCE_IP"
  get_real_ip_type = "TOA"
}
```

## Argument Reference

The following arguments are supported:

* `global_accelerator_id` - (Required, String, ForceNew) Global accelerator instance ID.
* `name` - (Required, String) Listener name. Maximum length is 60 bytes.
* `port_ranges` - (Required, List, ForceNew) Port ranges.
* `certification_type` - (Optional, String) Certificate type. `UNIDIRECTIONAL`: one-way. `MUTUAL`: two-way.
* `cipher_policy_id` - (Optional, String) Cipher policy ID.
* `client_affinity_time` - (Optional, Int) Session persistence time.
* `client_affinity` - (Optional, String) Whether to enable session persistence.
* `client_ca_certificates` - (Optional, List: [`String`]) Client CA certificates.
* `description` - (Optional, String) Description. Maximum length is 100 bytes.
* `get_real_ip_type` - (Optional, String) Layer 4 method to get real IP. Valid values: `TOA`, `ProxyProtocol`.
* `idle_timeout` - (Optional, Int) Connection idle timeout.
* `listener_type` - (Optional, String, ForceNew) Listener type, default is intelligent routing.
* `protocol` - (Optional, String, ForceNew) Protocol, default is TCP.
* `request_timeout` - (Optional, Int) Request timeout.
* `server_certificates` - (Optional, List: [`String`]) Server certificates.
* `x_forwarded_for_real_ip` - (Optional, Bool) Whether to enable layer 7 method to get real IP.

The `port_ranges` object supports the following:

* `from_port` - (Required, Int) Start port.
* `to_port` - (Required, Int) End port.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `listener_id` - Listener ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.
* `delete` - (Defaults to `20m`) Used when deleting the resource.

## Import

GA2 listener can be imported using the composite ID `<global_accelerator_id>#<listener_id>`, e.g.

```
terraform import tencentcloud_ga2_listener.example ga2-xxxxxxxx#lbl-xxxxxxxx
```

