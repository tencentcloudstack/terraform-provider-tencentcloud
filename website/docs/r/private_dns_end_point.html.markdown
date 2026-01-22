---
subcategory: "PrivateDNS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_private_dns_end_point"
sidebar_current: "docs-tencentcloud-resource-private_dns_end_point"
description: |-
  Provides a resource to create a Private Dns end point
---

# tencentcloud_private_dns_end_point

Provides a resource to create a Private Dns end point

## Example Usage

```hcl
resource "tencentcloud_private_dns_end_point" "example" {
  end_point_name       = "tf-example"
  end_point_service_id = "vpcsvc-61wcwmar"
  end_point_region     = "ap-guangzhou"
  ip_num               = 1
}
```

## Argument Reference

The following arguments are supported:

* `end_point_name` - (Required, String, ForceNew) Endpoint name.
* `end_point_region` - (Required, String, ForceNew) Endpoint region, which should be consistent with the region of the endpoint service.
* `end_point_service_id` - (Required, String, ForceNew) Endpoint service ID (namely, VPC endpoint service ID).
* `ip_num` - (Optional, Int, ForceNew) Number of endpoint IP addresses.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `end_point_vip_set` - Vip list of endpoint.


## Import

Private Dns end point can be imported using the id, e.g.

```
terraform import tencentcloud_private_dns_end_point.example eid-77a246c867
```

