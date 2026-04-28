---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_load_balancer"
sidebar_current: "docs-tencentcloud-resource-teo_load_balancer"
description: |-
  Provides a resource to create a TEO load balancer instance.
---

# tencentcloud_teo_load_balancer

Provides a resource to create a TEO load balancer instance.

## Example Usage

```hcl
resource "tencentcloud_teo_load_balancer" "example" {
  zone_id = "zone-3fkff38fyw8s"
  name    = "tf-example"
  type    = "HTTP"

  origin_groups {
    priority        = "priority_1"
    origin_group_id = "og-3pfz5626nmbb"
  }

  origin_groups {
    priority        = "priority_2"
    origin_group_id = "og-3pfz1ztltzo0"
  }

  health_checker {
    type               = "ICMP Ping"
    interval           = 30
    timeout            = 5
    health_threshold   = 3
    critical_threshold = 2
  }

  steering_policy = "Pritory"
  failover_policy = "OtherOriginGroup"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Load balancer instance name, 1-200 characters, allowed characters: `a-z`, `A-Z`, `0-9`, `_`, `-`.
* `origin_groups` - (Required, List) Source origin group list with failover priority.
* `type` - (Required, String, ForceNew) Instance type. Valid values: `HTTP` (HTTP-specific, supports HTTP-specific and general origin groups, only referenced by site acceleration services); `GENERAL` (general, only supports general origin groups, can be referenced by site acceleration and Layer-4 proxy).
* `zone_id` - (Required, String, ForceNew) Site ID.
* `failover_policy` - (Optional, String) Retry policy on request failure. Valid values: `OtherOriginGroup` (retry next priority origin group); `OtherRecordInOriginGroup` (retry another origin in the same group). Default: `OtherRecordInOriginGroup`.
* `health_checker` - (Optional, List) Health check policy. If not set, health check is disabled by default.
* `steering_policy` - (Optional, String) Traffic steering policy between origin groups. Valid value: `Pritory` (failover by priority). Default: `Pritory`.

The `headers` object of `health_checker` supports the following:

* `key` - (Required, String) Custom header key.
* `value` - (Required, String) Custom header value.

The `health_checker` object supports the following:

* `type` - (Required, String) Health check type. Valid values: `HTTP`, `HTTPS`, `TCP`, `UDP`, `ICMP Ping`, `NoCheck`. `NoCheck` means health check is disabled.
* `critical_threshold` - (Optional, Int) Number of consecutive unhealthy results to mark the origin as unhealthy. Default: 2.
* `expected_codes` - (Optional, List) Expected response status codes, valid when `type` is `HTTP` or `HTTPS`. e.g., `["200", "301"]`.
* `follow_redirect` - (Optional, String) Whether to follow 301/302 redirects, valid when `type` is `HTTP` or `HTTPS`. Valid values: `true`, `false`.
* `headers` - (Optional, List) Custom HTTP request headers (up to 10), valid when `type` is `HTTP` or `HTTPS`.
* `health_threshold` - (Optional, Int) Number of consecutive healthy results to mark the origin as healthy. Default: 3, minimum: 1.
* `interval` - (Optional, Int) Check interval in seconds. Valid values: `30`, `60`, `180`, `300`, `600`.
* `method` - (Optional, String) Request method, valid when `type` is `HTTP` or `HTTPS`. Valid values: `GET`, `HEAD`.
* `path` - (Optional, String) Probe path, valid when `type` is `HTTP` or `HTTPS`. Full host/path without protocol, e.g., `www.example.com/test`.
* `port` - (Optional, Int) Check port. Required when `type` is `HTTP`, `HTTPS`, `TCP`, or `UDP`.
* `recv_context` - (Optional, String) Expected response content from origin, valid when `type` is `UDP`. Only ASCII visible characters allowed, max length 500.
* `send_context` - (Optional, String) Content to send during health check, valid when `type` is `UDP`. Only ASCII visible characters allowed, max length 500.
* `timeout` - (Optional, Int) Timeout in seconds for each health check. Must be less than `interval`. Default: 5.

The `origin_groups` object supports the following:

* `origin_group_id` - (Required, String) Origin group ID.
* `priority` - (Required, String) Priority, format: `priority_` + number, highest priority is `priority_1`. Valid values: `priority_1`, `priority_2`..., `priority_10`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_id` - Load balancer instance ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `10m`) Used when creating the resource.
* `update` - (Defaults to `10m`) Used when updating the resource.

## Import

TEO load balancer can be imported using the zoneId#instanceId, e.g.

```
terraform import tencentcloud_teo_load_balancer.example zone-3fkff38fyw8s#lb-3pfzdob8hh3d
```

