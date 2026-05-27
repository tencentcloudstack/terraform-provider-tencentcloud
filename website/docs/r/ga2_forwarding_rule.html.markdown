---
subcategory: "Global Accelerator 2(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_forwarding_rule"
sidebar_current: "docs-tencentcloud-resource-ga2_forwarding_rule"
description: |-
  Provides a resource to create a GA2 (Global Accelerator 2) forwarding rule.
---

# tencentcloud_ga2_forwarding_rule

Provides a resource to create a GA2 (Global Accelerator 2) forwarding rule.

## Example Usage

```hcl
resource "tencentcloud_ga2_forwarding_rule" "example" {
  global_accelerator_id = "ga2-xxxxxxxx"
  listener_id           = "lis-xxxxxxxx"
  forwarding_policy_id  = "fp-xxxxxxxx"

  rule_conditions {
    rule_condition_type  = "HOST"
    rule_condition_value = ["example.com"]
  }

  rule_actions {
    rule_action_type  = "ForwardGroup"
    rule_action_value = "eg-xxxxxxxx"
  }

  origin_headers {
    key   = "X-Custom-Header"
    value = "custom-value"
  }

  enable_origin_sni = true
  origin_sni        = "example.com"
  origin_host       = "origin.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `forwarding_policy_id` - (Required, String, ForceNew) Forwarding policy ID.
* `global_accelerator_id` - (Required, String, ForceNew) Global accelerator instance ID.
* `listener_id` - (Required, String, ForceNew) Listener ID.
* `rule_actions` - (Required, List) Layer-7 forwarding rule actions.
* `rule_conditions` - (Required, List) Layer-7 forwarding rule conditions.
* `enable_origin_sni` - (Optional, Bool) Whether to enable origin SNI.
* `origin_headers` - (Optional, List) Origin headers.
* `origin_host` - (Optional, String) Origin host value.
* `origin_sni` - (Optional, String) Origin SNI value.

The `origin_headers` object supports the following:

* `key` - (Optional, String) Header key.
* `value` - (Optional, String) Header value.

The `rule_actions` object supports the following:

* `rule_action_type` - (Optional, String) Rule action type.
* `rule_action_value` - (Optional, String) Rule action value.

The `rule_conditions` object supports the following:

* `rule_condition_type` - (Optional, String) Rule condition type.
* `rule_condition_value` - (Optional, List) Rule condition values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `forwarding_rule_id` - Forwarding rule ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.
* `delete` - (Defaults to `20m`) Used when deleting the resource.

## Import

GA2 forwarding rule can be imported using the composite ID `<global_accelerator_id>#<listener_id>#<forwarding_policy_id>#<forwarding_rule_id>`, e.g.

```
terraform import tencentcloud_ga2_forwarding_rule.example ga2-xxxxxxxx#lis-xxxxxxxx#fp-xxxxxxxx#fr-xxxxxxxx
```

