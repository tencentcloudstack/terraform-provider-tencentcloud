---
subcategory: "Global Accelerator(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_forwarding_policy"
sidebar_current: "docs-tencentcloud-resource-ga2_forwarding_policy"
description: |-
  Provides a resource to create a GA2 forwarding policy
---

# tencentcloud_ga2_forwarding_policy

Provides a resource to create a GA2 forwarding policy

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

resource "tencentcloud_ga2_listener" "example" {
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
}

resource "tencentcloud_ga2_forwarding_policy" "example" {
  global_accelerator_id = tencentcloud_ga2_accelerate_area.example.global_accelerator_id
  listener_id           = tencentcloud_ga2_listener.example.listener_id
  host                  = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `global_accelerator_id` - (Required, String, ForceNew) Global accelerator instance ID this forwarding policy belongs to.
* `host` - (Required, String) Domain name. Must match the regular expression `^(?:[a-z0-9-]{0,61}[a-z0-9]?.)+[a-z]{2,}$`. Length must be between 1 and 80 characters.
* `listener_id` - (Required, String, ForceNew) Listener ID this forwarding policy belongs to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `default_host_flag` - Whether this is the default domain name for the listener.
* `forwarding_policy_id` - Forwarding policy ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `5m`) Used when creating the resource.
* `update` - (Defaults to `5m`) Used when updating the resource.
* `delete` - (Defaults to `5m`) Used when deleting the resource.

## Import

GA2 forwarding policy can be imported using the composite id `<global_accelerator_id>#<listener_id>#<forwarding_policy_id>`, e.g.

```
terraform import tencentcloud_ga2_forwarding_policy.example ga-jnyfyyss#lsr-hzht200v#dm-kvassops
```

