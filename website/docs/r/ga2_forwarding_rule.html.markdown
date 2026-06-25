---
subcategory: "Global Accelerator(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_forwarding_rule"
sidebar_current: "docs-tencentcloud-resource-ga2_forwarding_rule"
description: |-
  Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) layer-7 forwarding rule.
---

# tencentcloud_ga2_forwarding_rule

Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) layer-7 forwarding rule.

## Example Usage

```hcl
resource "tencentcloud_ga2_forwarding_rule" "example" {
  global_accelerator_id = "ga-fhhs8w84"
  listener_id           = "lsr-dyy8jhzp"
  forwarding_policy_id  = "dm-rjssxr8k"

  rule_conditions {
    rule_condition_type  = "Path"
    rule_condition_value = ["/path"]
  }

  rule_actions {
    rule_action_type  = "ForwardGroup"
    rule_action_value = "epg-nt4iwozo"
  }
}
```

### with origin headers

```hcl
resource "tencentcloud_ga2_forwarding_rule" "example1" {
  global_accelerator_id = "ga-fhhs8w84"
  listener_id           = "lsr-dyy8jhzp"
  forwarding_policy_id  = "dm-rjssxr8k"
  origin_host           = "2.2.2.2"

  rule_conditions {
    rule_condition_type  = "Path"
    rule_condition_value = ["/path"]
  }

  rule_actions {
    rule_action_type  = "ForwardGroup"
    rule_action_value = "epg-nt4iwozo"
  }

  origin_headers {
    key   = "key"
    value = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `forwarding_policy_id` - (Required, String, ForceNew) Forwarding policy ID this forwarding rule belongs to.
* `global_accelerator_id` - (Required, String, ForceNew) Global accelerator instance ID this forwarding rule belongs to.
* `listener_id` - (Required, String, ForceNew) Listener ID this forwarding rule belongs to.
* `rule_actions` - (Required, Set) Layer-7 forwarding rule action list. Treated as an unordered set; HCL element order has no semantic meaning.
* `rule_conditions` - (Required, Set) Layer-7 forwarding rule condition list. Treated as an unordered set; HCL element order has no semantic meaning.
* `enable_origin_sni` - (Optional, Bool) Whether to enable origin SNI.
* `origin_headers` - (Optional, Set) Origin request header list. Treated as an unordered set; HCL element order has no semantic meaning.
* `origin_host` - (Optional, String) Origin host value.
* `origin_sni` - (Optional, String) Origin SNI value.

The `origin_headers` object supports the following:

* `key` - (Required, String) Origin request header key.
* `value` - (Required, String) Origin request header value.

The `rule_actions` object supports the following:

* `rule_action_type` - (Required, String) Layer-7 forwarding rule action type.
* `rule_action_value` - (Required, String) Layer-7 forwarding rule action value.

The `rule_conditions` object supports the following:

* `rule_condition_type` - (Required, String) Layer-7 forwarding rule condition type.
* `rule_condition_value` - (Required, Set) Layer-7 forwarding rule condition values. Treated as an unordered set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `forwarding_rule_id` - Layer-7 forwarding rule ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `5m`) Used when creating the resource.
* `update` - (Defaults to `5m`) Used when updating the resource.
* `delete` - (Defaults to `5m`) Used when deleting the resource.

## Import

GA2 forwarding rule can be imported using the composite id `<global_accelerator_id>#<listener_id>#<forwarding_policy_id>#<forwarding_rule_id>`, e.g.

```
terraform import tencentcloud_ga2_forwarding_rule.example ga-fhhs8w84#lsr-dyy8jhzp#dm-rjssxr8k#rule-757r3bk2
```

