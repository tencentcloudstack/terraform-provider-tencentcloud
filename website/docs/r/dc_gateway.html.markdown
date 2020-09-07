---
subcategory: "Direct Connect Gateway(DCG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_gateway"
sidebar_current: "docs-tencentcloud-resource-dc_gateway"
description: |-
  Provides a resource to creating direct connect gateway instance.
---

# tencentcloud_dc_gateway

Provides a resource to creating direct connect gateway instance.

## Example Usage

```hcl
resource "tencentcloud_vpc" "main" {
  name       = "ci-vpc-instance-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_dc_gateway" "vpc_main" {
  name                = "ci-cdg-vpc-test"
  network_instance_id = tencentcloud_vpc.main.id
  network_type        = "VPC"
  gateway_type        = "NAT"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the DCG.
* `network_instance_id` - (Required, ForceNew) If the 'network_type' value is 'VPC', the available value is VPC ID. But when the 'network_type' value is 'CCN', the available value is CCN instance ID.
* `network_type` - (Required, ForceNew) Type of associated network, the available value include 'VPC' and 'CCN'.
* `gateway_type` - (Optional, ForceNew) Type of the gateway, the available value include 'NORMAL' and 'NAT'. Default is 'NORMAL'. NOTES: CCN only supports 'NORMAL' and a vpc can create two DCGs, the one is NAT type and the other is non-NAT type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cnn_route_type` - Type of CCN route, the available value include 'BGP' and 'STATIC'. The property is available when the DCG type is CCN gateway and BGP enabled.
* `create_time` - Creation time of resource.
* `enable_bgp` - Indicates whether the BGP is enabled.


## Import

Direct connect gateway instance can be imported, e.g.

```
$ terraform import tencentcloud_dc_gateway.instance dcg-id
```

