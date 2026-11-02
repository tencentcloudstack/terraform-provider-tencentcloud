Provides a resource to create a TEO load balancer instance.

Example Usage

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

Example Usage - HTTP_V2 with probe_cluster

```hcl
resource "tencentcloud_teo_load_balancer" "example" {
  zone_id = "zone-3fkff38fyw8s"
  name    = "tf-example-v2"
  type    = "HTTP_V2"

  origin_groups {
    priority        = "priority_1"
    origin_group_id = "og-3w0nrwiw4c6v"
  }

  health_checker {
    type               = "HTTP"
    path               = "www.example.com/test"
    port               = 80
    interval           = 30
    timeout            = 5
    health_threshold   = 3
    critical_threshold = 2
    method             = "HEAD"
    expected_codes     = ["200"]
    follow_redirect    = "false"
    probe_cluster      = "global"
  }

  steering_policy = "Pritory"
  failover_policy = "OtherOriginGroup"
}
```

Import

TEO load balancer can be imported using the zoneId#instanceId, e.g.

```
terraform import tencentcloud_teo_load_balancer.example zone-3fkff38fyw8s#lb-3pfzdob8hh3d
```
