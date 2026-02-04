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

### If network_type is VPC

```hcl
// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create dc gateway
resource "tencentcloud_dc_gateway" "example" {
  name                = "tf-example"
  network_instance_id = tencentcloud_vpc.vpc.id
  network_type        = "VPC"
  gateway_type        = "NORMAL"

  tags = {
    Environment = "production"
    Owner       = "ops-team"
  }
}
```

### If network_type is CCN

```hcl
// create ccn
resource "tencentcloud_ccn" "ccn" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
  tags = {
    createBy = "terraform"
  }
}

// create dc gateway
resource "tencentcloud_dc_gateway" "example" {
  name                = "tf-example"
  network_instance_id = tencentcloud_ccn.ccn.id
  network_type        = "CCN"
  gateway_type        = "NORMAL"

  tags = {
    Team    = "networking"
    Purpose = "production"
  }
}
```

### Update tags

```hcl
resource "tencentcloud_dc_gateway" "example" {
  name                = "tf-example"
  network_instance_id = tencentcloud_ccn.ccn.id
  network_type        = "CCN"
  gateway_type        = "NORMAL"

  # Tags can be updated without recreating the gateway
  tags = {
    Environment = "staging"
    Team        = "devops"
    CostCenter  = "IT-001"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the DCG.
* `network_instance_id` - (Required, String, ForceNew) If the `network_type` value is `VPC`, the available value is VPC ID. But when the `network_type` value is `CCN`, the available value is CCN instance ID.
* `network_type` - (Required, String, ForceNew) Type of associated network. Valid value: `VPC` and `CCN`.
* `gateway_type` - (Optional, String, ForceNew) Type of the gateway. Valid value: `NORMAL` and `NAT`. Default is `NORMAL`. NOTES: CCN only supports `NORMAL` and a VPC can create two DCGs, the one is NAT type and the other is non-NAT type.
* `tags` - (Optional, Map) Tag key-value pairs for the DC gateway. Multiple tags can be set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cnn_route_type` - Type of CCN route. Valid value: `BGP` and `STATIC`. The property is available when the DCG type is CCN gateway and BGP enabled.
* `create_time` - Creation time of resource.
* `enable_bgp` - Indicates whether the BGP is enabled.


## Import

Direct connect gateway instance can be imported, e.g. Tags will be imported automatically.

```
$ terraform import tencentcloud_dc_gateway.example dcg-dr1y0hu7
```

