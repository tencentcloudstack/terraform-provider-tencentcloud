---
layout: "tencentcloud"
page_title: "tencentcloud: tencentcloud_alb_server_attachment"
sidebar_current: "docs-tencentcloud-resource-lb-server-attachment"
description: |-
  Provides an tencentcloud application load balancer servers attachment as a resource, to attach and detach instances from load balancer.
---

# tencentcloud_alb_server_attachment

Provides Load Balancer server attachment resource.

~> **NOTE:** Currently only support existing `loadbalancer_id` `listener_id` `location_id` and Application layer 7 load balancer

## Example Usage

```hcl
resource "tencentcloud_alb_server_attachment" "service1" {
  loadbalancer_id = "lb-qk1dqox5"
  listener_id = "lbl-ghoke4tl"
  location_id = "loc-i858qv1l"
  backends = [
    {
      instance_id = "ins-4j30i5pe"
      port = 80
      weight = 50
    },
    {
      instance_id = "ins-4j30i5pe"
      port = 8080
      weight = 50
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `loadbalancer_id` - (Required, Forces new resource) loadbalancer ID.
* `listener_id` - (Required, Forces new resource) listener ID.
* `location_id` - (Optional) location ID only support for layer 7 loadbalancer
* `backends` - (Required) list of backend server. Valid value range [1-100].

### Block backends

The backends mapping supports the following:

* `instance_id` - (Required) A list backend instance ID (CVM instance ID).
* `port` - (Required) The port used by the backend server. Valid value range: [1-65535].
* `weight` - (Optional) Weight of the backend server. Valid value range: [0-100]. Default to 10.

## Attributes Reference

The following attributes are exported:

* `loadbalancer_id` - loadbalancer ID.
* `listener_id` - listener ID.
* `location_id` - location ID (only support for layer 7 loadbalancer)
* `protocol_type` - http or tcp
