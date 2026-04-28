---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_load_balancer"
sidebar_current: "docs-tencentcloud-resource-teo_load_balancer"
description: |-
  Provides a resource to create a TEO load balancer instance
---

# tencentcloud_teo_load_balancer

Provides a resource to create a TEO load balancer instance

## Example Usage

### Create a load balancer with HTTP health checker

```hcl
resource "tencentcloud_teo_load_balancer" "example" {
  zone_id = "zone-197z8rf93cfw"
  name    = "test-lb"
  type    = "HTTP"

  origin_groups {
    priority        = "priority_1"
    origin_group_id = "og-aaa"
  }

  origin_groups {
    priority        = "priority_2"
    origin_group_id = "og-bbb"
  }

  health_checker {
    type            = "HTTP"
    port            = 80
    interval        = 30
    timeout         = 5
    path            = "/health"
    method          = "GET"
    follow_redirect = "on"

    expected_codes = ["200", "301"]

    headers {
      key   = "X-Custom-Header"
      value = "health-check"
    }
  }

  steering_policy = "Pritory"
  failover_policy = "OtherOriginGroup"
}
```

### Create a load balancer with UDP health checker

```hcl
resource "tencentcloud_teo_load_balancer" "example" {
  zone_id = "zone-197z8rf93cfw"
  name    = "test-lb-udp"
  type    = "GENERAL"

  origin_groups {
    priority        = "priority_1"
    origin_group_id = "og-aaa"
  }

  health_checker {
    type         = "UDP"
    port         = 53
    interval     = 30
    timeout      = 5
    send_context = "health_check"
    recv_context = "ok"
  }

  steering_policy = "Pritory"
  failover_policy = "OtherRecordInOriginGroup"
}
```

### Create a load balancer without health checker

```hcl
resource "tencentcloud_teo_load_balancer" "example" {
  zone_id = "zone-197z8rf93cfw"
  name    = "test-lb-nocheck"
  type    = "GENERAL"

  origin_groups {
    priority        = "priority_1"
    origin_group_id = "og-aaa"
  }

  steering_policy = "Pritory"
  failover_policy = "OtherRecordInOriginGroup"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Instance name, can be 1-200 characters, allowed characters are a-z, A-Z, 0-9, _, -.
* `origin_groups` - (Required, List) Origin group list and corresponding failover scheduling priority.
* `type` - (Required, String, ForceNew) Instance type, valid values: `HTTP` (HTTP dedicated type), `GENERAL` (general type).
* `zone_id` - (Required, String, ForceNew) Site ID.
* `failover_policy` - (Optional, String) Request retry policy when accessing an origin fails, valid values: `OtherOriginGroup` (retry next priority origin group), `OtherRecordInOriginGroup` (retry other origins in the same group). Default is OtherRecordInOriginGroup.
* `health_checker` - (Optional, List) Health check policy.
* `steering_policy` - (Optional, String) Traffic scheduling policy between origin groups, valid values: `Pritory` (failover by priority order). Default is Pritory.

The `headers` object of `health_checker` supports the following:

* `key` - (Required, String) Custom header Key.
* `value` - (Required, String) Custom header Value.

The `health_checker` object supports the following:

* `type` - (Required, String) Health check policy type, valid values: `HTTP`, `HTTPS`, `TCP`, `UDP`, `ICMP Ping`, `NoCheck`.
* `critical_threshold` - (Optional, Int) Unhealthy threshold, the number of consecutive health checks that are 'unhealthy' before judging the origin as 'unhealthy', default 2.
* `expected_codes` - (Optional, List) Expected status codes for health determination, only valid when Type=HTTP or HTTPS.
* `follow_redirect` - (Optional, String) Whether to enable 301/302 redirect following, only valid when Type=HTTP or HTTPS.
* `headers` - (Optional, List) Custom HTTP request headers for detection, only valid when Type=HTTP or HTTPS, up to 10.
* `health_threshold` - (Optional, Int) Health threshold, the number of consecutive health checks that are 'healthy' before judging the origin as 'healthy', default 3, minimum 1.
* `interval` - (Optional, Int) Check frequency, how often to initiate a health check task, in seconds. Valid values: 30, 60, 180, 300, 600.
* `method` - (Optional, String) Request method, only valid when Type=HTTP or HTTPS. Valid values: `GET`, `HEAD`.
* `path` - (Optional, String) Detection path, only valid when Type=HTTP or HTTPS. Need to fill in the complete host/path, excluding the protocol part.
* `port` - (Optional, Int) Check port. Required when Type=HTTP, HTTPS, TCP or UDP.
* `recv_context` - (Optional, String) Expected response from origin for health check, only valid when Type=UDP. Only ASCII visible characters, max 500 characters.
* `send_context` - (Optional, String) Content sent by health check, only valid when Type=UDP. Only ASCII visible characters, max 500 characters.
* `timeout` - (Optional, Int) Timeout for each health check, in seconds, default is 5s, must be less than Interval.

The `origin_groups` object supports the following:

* `origin_group_id` - (Required, String) Origin group ID.
* `priority` - (Required, String) Priority, format is 'priority_' + 'number', the highest priority is 'priority_1'. Valid values: priority_1 to priority_10.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_id` - Load balancer instance ID.
* `l4_used_list` - List of L4 proxy instances bound to this load balancer instance.
* `l7_used_list` - List of L7 domain names bound to this load balancer instance.
* `origin_group_health_status` - Origin group health status.
  * `origin_group_id` - Origin group ID.
  * `origin_group_name` - Origin group name.
  * `origin_health_status` - Health status of origins in the origin group.
    * `healthy` - Origin health status, valid values: `Healthy`, `Unhealthy`, `Undetected`.
    * `origin` - Origin.
  * `origin_type` - Origin group type, valid values: `HTTP`, `GENERAL`.
  * `priority` - Priority.
* `references` - List of instances that reference this load balancer.
  * `instance_id` - Reference instance ID.
  * `instance_name` - Reference instance name.
  * `instance_type` - Reference service type.
* `status` - Load balancer status, valid values: `Pending` (deploying), `Deleting` (deleting), `Running` (effective).


## Import

teo load_balancer can be imported using the zone_id#instance_id, e.g.

```
terraform import tencentcloud_teo_load_balancer.example zone-297z8rf93cfw#lb-12345678
```

