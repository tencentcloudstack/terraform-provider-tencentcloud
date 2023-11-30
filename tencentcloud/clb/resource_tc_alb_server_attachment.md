Provides an tencentcloud application load balancer servers attachment as a resource, to attach and detach instances from load balancer.

~> **NOTE:** It has been deprecated and replaced by `tencentcloud_clb_attachment`.

~> **NOTE:** Currently only support existing `loadbalancer_id` `listener_id` `location_id` and Application layer 7 load balancer

Example Usage

```hcl
resource "tencentcloud_alb_server_attachment" "service1" {
  loadbalancer_id = "lb-qk1dqox5"
  listener_id     = "lbl-ghoke4tl"
  location_id     = "loc-i858qv1l"

  backends = [
    {
      instance_id = "ins-4j30i5pe"
      port        = 80
      weight      = 50
    },
    {
      instance_id = "ins-4j30i5pe"
      port        = 8080
      weight      = 50
    },
  ]
}
```