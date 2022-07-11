---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_alb_server_attachment"
sidebar_current: "docs-tencentcloud-resource-alb_server_attachment"
description: |-
  Provides an tencentcloud application load balancer servers attachment as a resource, to attach and detach instances from load balancer.
---

# tencentcloud_alb_server_attachment

Provides an tencentcloud application load balancer servers attachment as a resource, to attach and detach instances from load balancer.

~> **NOTE:** It has been deprecated and replaced by `tencentcloud_clb_attachment`.

~> **NOTE:** Currently only support existing `loadbalancer_id` `listener_id` `location_id` and Application layer 7 load balancer

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `backends` - (Required, Set) list of backend server.
* `listener_id` - (Required, String, ForceNew) listener ID.
* `loadbalancer_id` - (Required, String, ForceNew) loadbalancer ID.
* `location_id` - (Optional, String, ForceNew) location ID, only support for layer 7 loadbalancer.

The `backends` object supports the following:

* `instance_id` - (Required, String) A list backend instance ID (CVM instance ID).
* `port` - (Required, Int) The port used by the backend server. Valid value range: [1-65535].
* `weight` - (Optional, Int) Weight of the backend server. Valid value range: [0-100]. Default to 10.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `protocol_type` - The protocol type, http or tcp.


