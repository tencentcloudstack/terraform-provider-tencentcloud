---
subcategory: "Web Application Firewall(WAF)"
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

* `action_type` - (Required, String) Action type, 1 represents blocking, 2 represents captcha, 3 represents observation, and 4 represents redirection.
* `domain` - (Required, String) Domain name that needs to add policy.
* `expire_time` - (Required, String) Expiration time, measured in seconds, such as 1677254399, which means the expiration time is 2023-02-24 23:59:59 0 means never expires.
* `name` - (Required, String) Rule Name.
* `sort_id` - (Required, String) Priority, value range 0-100.
* `strategies` - (Required, List) Strategies detail.
* `redirect` - (Optional, String) If the action is a redirect, it represents the redirect address; Other situations can be left blank.
* `status` - (Optional, String) The status of the switch, 1 is on, 0 is off, default 1.

The `strategies` object supports the following:

* `arg` - (Required, String) Matching parameters.
* `compare_func` - (Required, String) Logical symbol.
* `content` - (Required, String) Matching Content.
* `field` - (Required, String) Matching Fields.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - rule ID.


## Import

waf clb_domain can be imported using the id, e.g.

```
terraform import tencentcloud_waf_clb_domain.example waf_2kxtlbky00b2v1fn#test.com#waf-0FSehoRU
```

