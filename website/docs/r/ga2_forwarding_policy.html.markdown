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
resource "tencentcloud_ga2_forwarding_policy" "example" {
  global_accelerator_id = "ga-fhhs8w84"
  listener_id           = "lsr-dyy8jhzp"
  host                  = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `global_accelerator_id` - (Required, String, ForceNew) Global accelerator instance ID this forwarding policy belongs to.
* `host` - (Required, String) The domain/host for the forwarding policy.
* `listener_id` - (Required, String, ForceNew) Listener ID this forwarding policy belongs to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `default_host_flag` - Whether this is the default host policy for the listener.
* `forwarding_policy_id` - Forwarding policy ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `5m`) Used when creating the resource.
* `update` - (Defaults to `5m`) Used when updating the resource.
* `delete` - (Defaults to `5m`) Used when deleting the resource.

## Import

GA2 forwarding policy can be imported using the composite id `<global_accelerator_id>#<listener_id>#<forwarding_policy_id>`, e.g.

```
terraform import tencentcloud_ga2_forwarding_policy.example ga-fhhs8w84#lsr-dyy8jhzp#dm-rjssxr8k
```

