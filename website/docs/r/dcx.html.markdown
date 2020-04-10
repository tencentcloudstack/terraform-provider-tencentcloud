---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcx"
sidebar_current: "docs-tencentcloud-resource-dcx"
description: |-
  Provides a resource to creating dedicated tunnels instances.
---

# tencentcloud_dcx

Provides a resource to creating dedicated tunnels instances.

~> **NOTE:** 1. ID of the DC is queried, can only apply for this resource offline.

## Example Usage

```hcl
variable "dc_id" {
  default = "dc-kax48sg7"
}

variable "dcg_id" {
  default = "dcg-dmbhf7jf"
}

variable "vpc_id" {
  default = "vpc-4h9v4mo3"
}

resource "tencentcloud_dcx" "bgp_main" {
  bandwidth    = 900
  dc_id        = var.dc_id
  dcg_id       = var.dcg_id
  name         = "bgp_main"
  network_type = "VPC"
  route_type   = "BGP"
  vlan         = 306
  vpc_id       = var.vpc_id
}

resource "tencentcloud_dcx" "static_main" {
  bandwidth             = 900
  dc_id                 = var.dc_id
  dcg_id                = var.dcg_id
  name                  = "static_main"
  network_type          = "VPC"
  route_type            = "STATIC"
  vlan                  = 301
  vpc_id                = var.vpc_id
  route_filter_prefixes = ["10.10.10.101/32"]
  tencent_address       = "100.93.46.1/30"
  customer_address      = "100.93.46.2/30"
}
```

## Argument Reference

The following arguments are supported:

* `dc_id` - (Required, ForceNew) ID of the DC to be queried, application deployment offline.
* `dcg_id` - (Required, ForceNew) ID of the DC Gateway. Currently only new in the console.
* `name` - (Required) Name of the dedicated tunnel.
* `vpc_id` - (Required, ForceNew) ID of the VPC or BMVPC.
* `bandwidth` - (Optional, ForceNew) Bandwidth of the DC.
* `bgp_asn` - (Optional, ForceNew) BGP ASN of the user. A required field within BGP.
* `bgp_auth_key` - (Optional, ForceNew) BGP key of the user.
* `customer_address` - (Optional, ForceNew) Interconnect IP of the DC within client.
* `network_type` - (Optional, ForceNew) Type of the network, and available values include VPC, BMVPC and CCN. The default value is VPC.
* `route_filter_prefixes` - (Optional, ForceNew) Static route, the network address of the user IDC. It can be modified after setting but cannot be deleted. AN unable field within BGP.
* `route_type` - (Optional, ForceNew) Type of the route, and available values include BGP and STATIC. The default value is BGP.
* `tencent_address` - (Optional, ForceNew) Interconnect IP of the DC within Tencent.
* `vlan` - (Optional, ForceNew) Vlan of the dedicated tunnels, and the range of values is [0-3000]. '0' means that only one tunnel can be created for the physical connect.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of resource.
* `state` - State of the dedicated tunnels, and available values include PENDING, ALLOCATING, ALLOCATED, ALTERING, DELETING, DELETED, COMFIRMING and REJECTED.


