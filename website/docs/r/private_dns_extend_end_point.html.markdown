---
subcategory: "PrivateDNS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_private_dns_extend_end_point"
sidebar_current: "docs-tencentcloud-resource-private_dns_extend_end_point"
description: |-
  Provides a resource to create a Private Dns extend end point
---

# tencentcloud_private_dns_extend_end_point

Provides a resource to create a Private Dns extend end point

## Example Usage

### If access_type is CLB

```hcl
resource "tencentcloud_private_dns_extend_end_point" "example" {
  end_point_name   = "tf-example"
  end_point_region = "ap-jakarta"
  forward_ip {
    access_type = "CLB"
    host        = "10.0.1.12"
    port        = 9000
    vpc_id      = "vpc-1v2i79fc"
  }
}
```

### If access_type is CCN

```hcl
resource "tencentcloud_private_dns_extend_end_point" "example" {
  end_point_name   = "tf-example"
  end_point_region = "ap-jakarta"
  forward_ip {
    access_type       = "CCN"
    host              = "1.1.1.1"
    port              = 8080
    vpc_id            = "vpc-2qjckjg2"
    access_gateway_id = "ccn-eo13f8ub"
  }
}
```

## Argument Reference

The following arguments are supported:

* `end_point_name` - (Required, String, ForceNew) Outbound endpoint name.
* `end_point_region` - (Required, String, ForceNew) The region of the outbound endpoint must be consistent with the region of the forwarding target VIP.
* `forward_ip` - (Optional, List, ForceNew) Forwarding target.

The `forward_ip` object supports the following:

* `access_type` - (Required, String, ForceNew) Forwarding target IP network access type. CLB: The forwarding IP is the internal CLB VIP. CCN: Forwarding IP through CCN routing.
* `host` - (Required, String, ForceNew) Forwarding target IP address.
* `port` - (Required, Int, ForceNew) Specifies the forwarding IP port number.
* `vpc_id` - (Required, String, ForceNew) Unique VPC ID.
* `access_gateway_id` - (Optional, String, ForceNew) CCN id. Required when the access type is CCN.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Private Dns extend end point can be imported using the id, e.g.

```
terraform import tencentcloud_private_dns_extend_end_point.example eid-960fb0ee9677
```

