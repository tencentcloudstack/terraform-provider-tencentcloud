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

resource "tencentcloud_ga2_endpoint_group" "example" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
  listener_id           = tencentcloud_ga2_listener.example.listener_id
  endpoint_group_type   = "VIRTUAL"

  endpoint_group_configuration {
    name                  = "tf-example"
    endpoint_group_region = "ap-guangzhou"
    description           = "tf example endpoint group"
    enable_health_check   = true
    forward_protocol      = "HTTP"
    check_type            = "HTTP"
    check_domain          = "check.com"
    check_method          = "GET"
    check_path            = "/path"
    connect_timeout       = 2
    health_check_interval = 30
    healthy_threshold     = 3
    unhealthy_threshold   = 3
    status_mask = [
      "http_2xx",
      "http_3xx",
      "http_4xx"
    ]

    endpoint_configurations {
      endpoint_type    = "CustomPublicIp"
      endpoint_service = "1.1.1.1"
      weight           = 10
    }

    endpoint_configurations {
      endpoint_type    = "CustomDomain"
      endpoint_service = "example.com"
      weight           = 20
    }

    port_overrides {
      listener_port = 90
      endpoint_port = 9090
    }
  }
}

resource "tencentcloud_ga2_forwarding_policy" "example" {
  global_accelerator_id = tencentcloud_ga2_accelerate_area.example.global_accelerator_id
  listener_id           = tencentcloud_ga2_listener.example.listener_id
  host                  = "example.com"
}

resource "tencentcloud_ga2_forwarding_rule" "example" {
  global_accelerator_id = tencentcloud_ga2_accelerate_area.example.global_accelerator_id
  listener_id           = tencentcloud_ga2_listener.example.listener_id
  forwarding_policy_id  = tencentcloud_ga2_forwarding_policy.example.forwarding_policy_id

  rule_conditions {
    rule_condition_type  = "Path"
    rule_condition_value = ["/path"]
  }

  rule_actions {
    rule_action_type  = "ForwardGroup"
    rule_action_value = tencentcloud_ga2_endpoint_group.example.endpoint_group_id
  }
}
```

### with origin headers

```hcl
resource "tencentcloud_ga2_forwarding_rule" "example" {
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
* `rule_conditions` - (Required, Set) Layer-7 forwarding rule condition list. Maximum of 1 element. Treated as an unordered set; HCL element order has no semantic meaning.
* `enable_origin_sni` - (Optional, Bool) Whether to enable origin SNI. Default: `false`. Required when `rule_actions.rule_action_type` is `ForwardGroup`.
* `origin_headers` - (Optional, Set) Origin request header list. Maximum of 5 elements. Required when `rule_actions.rule_action_type` is `ForwardGroup`. Treated as an unordered set; HCL element order has no semantic meaning.
* `origin_host` - (Optional, String) Origin host value. Maximum length is 80 characters. Required when `rule_actions.rule_action_type` is `ForwardGroup`.
* `origin_sni` - (Optional, String) Origin SNI value. Maximum length is 80 characters. Required when `enable_origin_sni` is `true`, and also required when `rule_actions.rule_action_type` is `ForwardGroup`.

The `origin_headers` object supports the following:

* `key` - (Required, String) Origin request header key. Must contain only printable ASCII characters and must not contain `()<>@,;:\"/[ ]?={}`. Length must be between 1 and 40 characters.
* `value` - (Required, String) Origin request header value. Maximum length is 128 characters. If the value contains `$`, only `$remote_addr` or `$remote_port` are supported.

The `rule_actions` object supports the following:

* `rule_action_type` - (Required, String) Layer-7 forwarding rule action type. Valid values: `ForwardGroup` (forward to an endpoint group), `Drop` (drop the request).
* `rule_action_value` - (Required, String) Layer-7 forwarding rule action value. Not required when `rule_action_type` is `Drop`. Required when `rule_action_type` is `ForwardGroup`, in which case it must be a custom endpoint group ID (the default endpoint group is not supported).

The `rule_conditions` object supports the following:

* `rule_condition_type` - (Required, String) Layer-7 forwarding rule condition type. Valid values: `Path`.
* `rule_condition_value` - (Required, Set) Layer-7 forwarding rule condition values. Each value must match the regular expression `^[a-zA-Z0-9_.-/]{1,80}$`. Maximum of 1 element. Treated as an unordered set.

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

