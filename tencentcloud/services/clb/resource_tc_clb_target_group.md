Provides a resource to create a CLB target group.

Example Usage

Create V1 target group with health check and tags

```hcl
resource "tencentcloud_clb_target_group" "test" {
  target_group_name = "test-v1"
  port              = 80
  type              = "v1"

  health_check {
    health_switch    = true
    protocol         = "TCP"
    timeout          = 5
    gap_time         = 10
    good_limit       = 3
    bad_limit        = 3
  }

  tags = {
    "createdBy" = "terraform"
  }
}
```

Create V2 TCP target group with TCP health check

```hcl
resource "tencentcloud_clb_target_group" "tcp_tg" {
  target_group_name = "tcp_tg"
  vpc_id            = "vpc-xxxxxx"
  type              = "v2"
  protocol          = "TCP"

  health_check {
    health_switch = true
    protocol      = "TCP"
  }
}
```

Create V2 target group with advanced features

```hcl
resource "tencentcloud_clb_target_group" "test_v2" {
  target_group_name    = "test-v2"
  vpc_id               = "vpc-xxxxxx"
  port                 = 80
  type                 = "v2"
  protocol             = "HTTP"
  schedule_algorithm   = "WRR"
  session_expire_time  = 1800
  keepalive_enable     = true
  weight               = 50

  health_check {
    health_switch      = true
    protocol           = "HTTP"
    port               = 8080
    timeout            = 5
    gap_time           = 11
    good_limit         = 4
    bad_limit          = 4
    http_check_path    = "/health"
    http_check_method  = "GET"
    http_check_domain  = "test.com"
    http_code          = 2  # 2xx
  }

  tags = {
    "createdBy" = "terraform"
    "env"       = "production"
  }
}
```

Create V2 HTTP target group with IP hash scheduling

```hcl
resource "tencentcloud_clb_target_group" "ip_hash" {
  target_group_name  = "ip-hash-tg"
  vpc_id             = "vpc-xxxxxxx"
  type               = "v2"
  protocol           = "HTTP"
  schedule_algorithm = "IP_HASH"
  ip_version         = "IPv4"

  health_check {
    health_switch = true
    protocol      = "HTTP"
    http_check_domain  = "test.com"
    timeout            = 5
    gap_time           = 11
    good_limit         = 4
    bad_limit          = 4
  }
}
```

Create V2 full listener target group

```hcl
resource "tencentcloud_clb_target_group" "full_listener" {
  target_group_name  = "full-listener-tg"
  vpc_id             = "vpc-xxxxxx"
  type               = "v2"
  protocol           = "TCP"
  full_listen_switch = true

  health_check {
    health_switch     = true
    protocol          = "HTTP"
    http_version      = "HTTP/1.1"
    http_check_path   = "/healthz"
    http_check_domain = "test.com"
    timeout           = 5
    gap_time          = 11
    good_limit        = 4
    bad_limit         = 4
  }
}
```

Create IPv6 target group

```hcl
resource "tencentcloud_clb_target_group" "ipv6" {
  target_group_name = "ipv6-tg"
  vpc_id            = "vpc-xxxxxx"
  type              = "v2"
  protocol          = "HTTP"
  ip_version        = "IPv6"

  health_check {
    health_switch = true
    protocol      = "HTTP"
    http_check_domain  = "test.com"
    timeout            = 5
    gap_time           = 11
    good_limit         = 4
    bad_limit          = 4
  }
}
```

Import

CLB target group can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group.test lbtg-3k3io0i0
```