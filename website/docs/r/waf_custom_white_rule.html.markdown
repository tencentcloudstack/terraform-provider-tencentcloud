---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_custom_white_rule"
sidebar_current: "docs-tencentcloud-resource-waf_custom_white_rule"
description: |-
  Provides a resource to create a waf custom_white_rule
---

# tencentcloud_waf_custom_white_rule

Provides a resource to create a waf custom_white_rule

## Example Usage

```hcl
resource "tencentcloud_waf_custom_white_rule" "example" {
  name        = "tf-example"
  sort_id     = "30"
  expire_time = "0"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "1.1.1.1"
    arg          = ""
  }

  status = "1"
  domain = "test.com"
  bypass = "geoip,cc,owasp"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain name.
* `instance_id` - (Required, String) Instance unique ID.
* `region` - (Required, String) Regions of LB bound by domain.
* `alb_type` - (Optional, String) Load balancer type: clb, apisix or tsegw, default clb.
* `api_safe_status` - (Optional, Int) Whether to enable api safe, 1 enable, 0 disable.
* `bot_status` - (Optional, Int) Whether to enable bot, 1 enable, 0 disable.
* `cls_status` - (Optional, Int) Whether to enable access logs, 1 enable, 0 disable.
* `engine` - (Optional, Int) Protection Status: 10: Rule Observation&&AI Off Mode, 11: Rule Observation&&AI Observation Mode, 12: Rule Observation&&AI Interception Mode, 20: Rule Interception&&AI Off Mode, 21: Rule Interception&&AI Observation Mode, 22: Rule Interception&&AI Interception Mode, Default 20.
* `flow_mode` - (Optional, Int) WAF traffic mode, 1 cleaning mode, 0 mirroring mode.
* `ip_headers` - (Optional, List: [`String`]) When is_cdn=3, this parameter needs to be filled in to indicate a custom header.
* `is_cdn` - (Optional, Int) Whether a proxy has been enabled before WAF, 0 no deployment, 1 deployment and use first IP in X-Forwarded-For as client IP, 2 deployment and use remote_addr as client IP, 3 deployment and use values of custom headers as client IP.
* `load_balancer_set` - (Optional, List) List of bound LB.
* `status` - (Optional, Int) Binding status between waf and LB, 0:not bind, 1:binding.

The `load_balancer_set` object supports the following:

* `listener_id` - (Required, String) Unique ID of listener in LB.
* `listener_name` - (Required, String) Listener name.
* `load_balancer_id` - (Required, String) LoadBalancer unique ID.
* `load_balancer_name` - (Required, String) LoadBalancer name.
* `protocol` - (Required, String) Protocol of listener, http or https.
* `region` - (Required, String) LoadBalancer region.
* `vip` - (Required, String) LoadBalancer IP.
* `vport` - (Required, Int) LoadBalancer port.
* `zone` - (Required, String) LoadBalancer zone.
* `load_balancer_type` - (Optional, String) Network type for load balancer.
* `numerical_vpc_id` - (Optional, Int) VPCID for load balancer, public network is -1, and internal network is filled in according to actual conditions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `domain_id` - Domain id.


## Import

waf custom_white_rule can be imported using the id, e.g.

```
terraform import tencentcloud_waf_custom_white_rule.example test.com#1100310837
```

