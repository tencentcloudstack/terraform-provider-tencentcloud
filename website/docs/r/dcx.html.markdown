---
subcategory: "Direct Connect(DC)"
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

### If network_type is VPC

```hcl
resource "tencentcloud_dcx" "example" {
  dc_id            = "dc-ink7y3qf"
  name             = "tf-example"
  dc_owner_account = "100017971194"
  network_type     = "VPC"
  network_region   = "ap-guangzhou"
  vpc_id           = "vpc-nzuu8dyj"
  dcg_id           = "dcg-ehr22qfb"
  bandwidth        = 100
  route_type       = "BGP"
  bgp_asn          = 64511
  vlan             = 60
  tencent_address  = "10.8.254.14/30"
  customer_address = "10.8.254.13/30"
}
```

### If network_type is CCN

```hcl
resource "tencentcloud_dcx" "example" {
  dc_id            = "dc-ink7y3qf"
  name             = "tf-example"
  dc_owner_account = "100017971194"
  network_type     = "CCN"
  network_region   = "ap-guangzhou"
  dcg_id           = "dcg-6d4uaubp"
  bandwidth        = 100
  route_type       = "BGP"
  bgp_asn          = 64511
  vlan             = 10
  tencent_address  = "10.8.254.10/30"
  customer_address = "10.8.254.9/30"
}
```

## Argument Reference

The following arguments are supported:

* `dc_id` - (Required, String, ForceNew) ID of the DC to be queried, application deployment offline.
* `dcg_id` - (Required, String, ForceNew) ID of the DC Gateway. Currently only new in the console.
* `name` - (Required, String) Name of the dedicated tunnel.
* `bandwidth` - (Optional, Int, ForceNew) Bandwidth of the DC.
* `bgp_asn` - (Optional, Int, ForceNew) BGP ASN of the user. A required field within BGP.
* `bgp_auth_key` - (Optional, String, ForceNew) BGP key of the user.
* `customer_address` - (Optional, String, ForceNew) Interconnect IP of the DC within client.
* `dc_owner_account` - (Optional, String, ForceNew) Connection owner, who is the current customer by default. The developer account ID should be entered for shared connections.
* `network_region` - (Optional, String, ForceNew) Network region.
* `network_type` - (Optional, String, ForceNew) Type of the network. Valid value: `VPC`, `BMVPC` and `CCN`. The default value is `VPC`.
* `route_filter_prefixes` - (Optional, Set: [`String`], ForceNew) Static route, the network address of the user IDC. It can be modified after setting but cannot be deleted. AN unable field within BGP.
* `route_type` - (Optional, String, ForceNew) Type of the route, and available values include BGP and STATIC. The default value is `BGP`.
* `tencent_address` - (Optional, String, ForceNew) Interconnect IP of the DC within Tencent.
* `vlan` - (Optional, Int, ForceNew) Vlan of the dedicated tunnels. Valid value ranges: (0~3000). `0` means that only one tunnel can be created for the physical connect.
* `vpc_id` - (Optional, String, ForceNew) ID of the VPC or BMVPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of resource.
* `state` - State of the dedicated tunnels. Valid value: `PENDING`, `ALLOCATING`, `ALLOCATED`, `ALTERING`, `DELETING`, `DELETED`, `COMFIRMING` and `REJECTED`.


## Import

DCX instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_dcx.example dcx-cbbr1gjk
```

