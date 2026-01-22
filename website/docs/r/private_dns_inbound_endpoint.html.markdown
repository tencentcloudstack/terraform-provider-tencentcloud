---
subcategory: "PrivateDNS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_private_dns_inbound_endpoint"
sidebar_current: "docs-tencentcloud-resource-private_dns_inbound_endpoint"
description: |-
  Provides a resource to create a Private Dns inbound endpoint
---

# tencentcloud_private_dns_inbound_endpoint

Provides a resource to create a Private Dns inbound endpoint

## Example Usage

```hcl
resource "tencentcloud_private_dns_inbound_endpoint" "example" {
  endpoint_name   = "tf-example"
  endpoint_region = "ap-guangzhou"
  endpoint_vpc    = "vpc-i5yyodl9"
  subnet_ip {
    subnet_id  = "subnet-hhi88a58"
    subnet_vip = "10.0.30.2"
  }

  subnet_ip {
    subnet_id  = "subnet-5rrirqyc"
    subnet_vip = "10.0.0.11"
  }

  subnet_ip {
    subnet_id = "subnet-60ut6n10"
  }
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_name` - (Required, String) Name.
* `endpoint_region` - (Required, String, ForceNew) Region.
* `endpoint_vpc` - (Required, String, ForceNew) VPC ID.
* `subnet_ip` - (Required, List, ForceNew) Subnet information.

The `subnet_ip` object supports the following:

* `subnet_id` - (Required, String, ForceNew) Subnet ID.
* `subnet_vip` - (Optional, String, ForceNew) IP address.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



