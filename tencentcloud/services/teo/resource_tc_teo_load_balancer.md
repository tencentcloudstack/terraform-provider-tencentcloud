Provides a resource to create a TEO load balancer instance

Example Usage

Create a load balancer with HTTP health checker

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
    type        = "HTTP"
    port        = 80
    interval    = 30
    timeout     = 5
    path        = "/health"
    method      = "GET"
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

Create a load balancer with UDP health checker

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
    type          = "UDP"
    port          = 53
    interval      = 30
    timeout       = 5
    send_context  = "health_check"
    recv_context  = "ok"
  }

  steering_policy = "Pritory"
  failover_policy = "OtherRecordInOriginGroup"
}
```

Create a load balancer without health checker

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

Import

teo load_balancer can be imported using the zone_id#instance_id, e.g.

```
terraform import tencentcloud_teo_load_balancer.example zone-297z8rf93cfw#lb-12345678
```
