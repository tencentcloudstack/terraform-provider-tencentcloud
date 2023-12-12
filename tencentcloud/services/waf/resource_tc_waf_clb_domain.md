Provides a resource to create a waf clb_domain

Example Usage

Create a basic waf clb domain

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

  region          = "gz"
  alb_type        = "clb"
}
```

Create a complete waf clb domain

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
  ip_headers      = [
    "headers_1",
    "headers_2",
    "headers_3",
  ]
}
```

Create a complete waf tsegw domain

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

Create a complete waf apisix domain

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

Import

waf clb_domain can be imported using the id, e.g.

```
terraform import tencentcloud_waf_clb_domain.example waf_2kxtlbky00b2v1fn#test.com#waf-0FSehoRU
```