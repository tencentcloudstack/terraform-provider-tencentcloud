---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_target_group"
sidebar_current: "docs-tencentcloud-resource-clb_target_group"
description: |-
  Provides a resource to create a CLB target group.
---

# tencentcloud_clb_target_group

Provides a resource to create a CLB target group.

## Example Usage

### Create V1 target group with health check and tags

```hcl
resource "tencentcloud_clb_target_group" "test" {
  target_group_name = "test-v1"
  port              = 80
  type              = "v1"

  health_check {
    health_switch = true
    protocol      = "TCP"
    timeout       = 5
    gap_time      = 10
    good_limit    = 3
    bad_limit     = 3
  }

  tags = {
    "createdBy" = "terraform"
  }
}
```

### Create V2 target group with advanced features

```hcl
resource "tencentcloud_clb_target_group" "test_v2" {
  target_group_name   = "test-v2"
  vpc_id              = "vpc-xxxxxx"
  port                = 80
  type                = "v2"
  protocol            = "HTTP"
  schedule_algorithm  = "WRR"
  session_expire_time = 1800
  keepalive_enable    = true
  weight              = 50

  health_check {
    health_switch     = true
    protocol          = "HTTP"
    port              = 8080
    timeout           = 5
    gap_time          = 10
    good_limit        = 3
    bad_limit         = 3
    http_check_path   = "/health"
    http_check_method = "GET"
    http_code         = 2 # 2xx
  }

  tags = {
    "createdBy" = "terraform"
    "env"       = "production"
  }
}
```

### Create V2 HTTP target group with IP hash scheduling

```hcl
resource "tencentcloud_clb_target_group" "ip_hash" {
  target_group_name  = "ip-hash-tg"
  vpc_id             = "vpc-xxxxxx"
  type               = "v2"
  protocol           = "HTTP"
  schedule_algorithm = "IP_HASH"
  ip_version         = "IPv4"

  health_check {
    health_switch = true
    protocol      = "HTTP"
  }
}
```

### Create V2 full listener target group

```hcl
resource "tencentcloud_clb_target_group" "full_listener" {
  target_group_name  = "full-listener-tg"
  vpc_id             = "vpc-xxxxxx"
  type               = "v2"
  protocol           = "TCP"
  full_listen_switch = true

  health_check {
    health_switch   = true
    protocol        = "HTTP"
    http_version    = "HTTP/1.1"
    http_check_path = "/healthz"
  }
}
```

### Create IPv6 target group

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
  }
}
```

## Argument Reference

The following arguments are supported:

* `full_listen_switch` - (Optional, Bool) Whether this is a full listener target group. Only valid for v2 target groups. true: full listener target group, false: normal target group.
* `health_check` - (Optional, List) Health check configuration.
* `ip_version` - (Optional, String) IP version type. Common values: IPv4, IPv6, IPv6FullChain.
* `keepalive_enable` - (Optional, Bool) Enable keep-alive connections. Only valid for HTTP/HTTPS target groups. true: enable, false: disable. Default: false.
* `port` - (Optional, Int) The default port of target group, add server after can use it.
* `protocol` - (Optional, String, ForceNew) Backend forwarding protocol of the target group. this field is required for the new version (v2) target group. currently supports TCP, UDP, HTTP, HTTPS, GRPC.
* `schedule_algorithm` - (Optional, String) Scheduling algorithm. Only valid for v2 target groups with HTTP/HTTPS/GRPC protocols. Valid values: WRR (weighted round robin), LEAST_CONN (least connections), IP_HASH (IP hash). Default: WRR.
* `session_expire_time` - (Optional, Int) Session persistence time in seconds. Only valid for v2 target groups with HTTP/HTTPS/GRPC protocols. Range: 30-3600 or 0 (disabled). Default: 0 (disabled).
* `tags` - (Optional, Map) Resource tags for the target group.
* `target_group_instances` - (Optional, List, **Deprecated**) It has been deprecated from version 1.77.3. please use `tencentcloud_clb_target_group_instance_attachment` instead. The backend server of target group bind.
* `target_group_name` - (Optional, String) Target group name.
* `type` - (Optional, String, ForceNew) Target group type, currently supported v1 (legacy version target group) and v2 (new version target group), defaults to v1 (legacy version target group).
* `vpc_id` - (Optional, String, ForceNew) VPC ID, default is based on the network.
* `weight` - (Optional, Int) Default backend server weight. Range: [0, 100]. Only valid for v2 target groups. When set, backend servers added to the target group will use this default weight if not specified.

The `health_check` object supports the following:

* `health_switch` - (Required, Bool) Whether to enable health check. true: enable, false: disable.
* `bad_limit` - (Optional, Int) Unhealthy threshold. Number of consecutive failed health checks required before marking the backend as unhealthy. Range: [2, 10]. Default: 3.
* `extended_code` - (Optional, String) Extended status code for health check.
* `gap_time` - (Optional, Int) Health check interval in seconds. Range: [2, 300]. Default: 5.
* `good_limit` - (Optional, Int) Healthy threshold. Number of consecutive successful health checks required before marking the backend as healthy. Range: [2, 10]. Default: 3.
* `http_check_domain` - (Optional, String) Health check domain. For HTTP/HTTPS protocol.
* `http_check_method` - (Optional, String) Health check HTTP method. For HTTP/HTTPS protocol. Valid values: HEAD, GET. Default: HEAD.
* `http_check_path` - (Optional, String) Health check path. For HTTP/HTTPS protocol. Must start with /. If not specified, / is used by default.
* `http_code` - (Optional, Int) HTTP status codes indicating health. For HTTP/HTTPS protocol. Example: 1 (1xx), 2 (2xx), 4 (3xx), 8 (4xx), 16 (5xx). Multiple values can be combined, e.g., 7 (1xx, 2xx, 3xx).
* `http_version` - (Optional, String) HTTP version for health check. Required when health check protocol is HTTP. Valid values: HTTP/1.0, HTTP/1.1. Only valid for TCP target groups.
* `port` - (Optional, Int) Health check port. If not specified, the backend server port is used by default.
* `protocol` - (Optional, String) Health check protocol. Valid values: TCP, HTTP, HTTPS, PING, CUSTOM, GRPC. Valid for v2 target groups.
* `timeout` - (Optional, Int) Health check response timeout in seconds. Range: [2, 60]. Default: 2.

The `target_group_instances` object supports the following:

* `bind_ip` - (Required, String) The internal ip of target group instance.
* `port` - (Required, Int) The port of target group instance.
* `new_port` - (Optional, Int) The new port of target group instance.
* `weight` - (Optional, Int) The weight of target group instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CLB target group can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group.test lbtg-3k3io0i0
```

