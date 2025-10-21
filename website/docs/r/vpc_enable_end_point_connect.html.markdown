---
subcategory: "Private Link(PLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_enable_end_point_connect"
sidebar_current: "docs-tencentcloud-resource-vpc_enable_end_point_connect"
description: |-
  Provides a resource to create a vpc enable_end_point_connect
---

# tencentcloud_vpc_enable_end_point_connect

Provides a resource to create a vpc enable_end_point_connect

## Example Usage

```hcl
resource "tencentcloud_vpc_enable_end_point_connect" "enable_end_point_connect" {
  end_point_service_id = "vpcsvc-98jddhcz"
  end_point_id         = ["vpce-6q0ftmke"]
  accept_flag          = true
}
```

## Argument Reference

The following arguments are supported:

* `accept_flag` - (Required, Bool, ForceNew) Whether to accept endpoint connection requests. `true`: Accept automatically. `false`: Do not automatically accept.
* `end_point_id` - (Required, Set: [`String`], ForceNew) Endpoint ID.
* `end_point_service_id` - (Required, String, ForceNew) Endpoint service ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



