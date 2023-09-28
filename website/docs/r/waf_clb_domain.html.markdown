---
subcategory: "Waf"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_clb_domain"
sidebar_current: "docs-tencentcloud-resource-waf_clb_domain"
description: |-
  Provides a resource to create a waf clb_domain
---

# tencentcloud_waf_clb_domain

Provides a resource to create a waf clb_domain

## Example Usage

### Create a basic waf clb domain

```hcl
resource "tencentcloud_waf_clb_domain" "example" {
  instance_id = "waf_2kxtlbky00b2v1fn"
  domain      = "test.com"

  load_balancer_set {
    load_balancer_id   = "lb-5dnrkgry"
    load_balancer_name = "keep-listener-clb"
    listener_id        = "lbl-nonkgvc2"
    listener_name      = "dsadasd"
    vip                = "106.55.220.8"
    vport              = "80"
    region             = "gz"
    protocol           = "HTTP"
    zone               = "ap-guangzhou-6"
    numerical_vpc_id   = "5232945"
    load_balancer_type = "OPEN"
  }

  region   = "gz"
  alb_type = "clb"
}
```

### Create a complete waf clb domain

```hcl
resource "tencentcloud_waf_clb_domain" "example" {
  instance_id = "waf_2kxtlbky00b2v1fn"
  domain      = "test.com"
  is_cdn      = 3
  status      = 1
  engine      = 21

  load_balancer_set {
    load_balancer_id   = "lb-5dnrkgry"
    load_balancer_name = "keep-listener-clb"
    listener_id        = "lbl-nonkgvc2"
    listener_name      = "dsadasd"
    vip                = "106.55.220.8"
    vport              = "80"
    region             = "gz"
    protocol           = "HTTP"
    zone               = "ap-guangzhou-6"
    numerical_vpc_id   = "5232945"
    load_balancer_type = "OPEN"
  }

  region          = "gz"
  flow_mode       = 1
  alb_type        = "clb"
  bot_status      = 1
  api_safe_status = 1
  ip_headers = [
    "headers_1",
    "headers_2",
    "headers_3",
  ]
}
```

### Create a complete waf tsegw domain

```hcl
resource "tencentcloud_waf_clb_domain" "example" {
  instance_id     = "waf_2kxtlbky00b2v1fn"
  domain          = "xxx.com"
  is_cdn          = 0
  status          = 1
  engine          = 12
  region          = "gz"
  flow_mode       = 0
  alb_type        = "tsegw"
  bot_status      = 0
  api_safe_status = 0
}
```

### Create a complete waf apisix domain

```hcl
resource "tencentcloud_waf_clb_domain" "example" {
  instance_id     = "waf_2kxtlbky00b2v1fn"
  domain          = "xxx.com"
  is_cdn          = 0
  status          = 1
  engine          = 12
  region          = "gz"
  flow_mode       = 0
  alb_type        = "apisix"
  bot_status      = 0
  api_safe_status = 0
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

waf clb_domain can be imported using the id, e.g.

```
terraform import tencentcloud_waf_clb_domain.example waf_2kxtlbky00b2v1fn#test.com#waf-0FSehoRU
```

