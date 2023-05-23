---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_net_detect"
sidebar_current: "docs-tencentcloud-resource-vpc_net_detect"
description: |-
  Provides a resource to create a vpc net_detect
---

# tencentcloud_vpc_net_detect

Provides a resource to create a vpc net_detect

## Example Usage

```hcl
resource "tencentcloud_vpc_net_detect" "net_detect" {
  net_detect_name      = "terrform-test"
  vpc_id               = "vpc-4owdpnwr"
  subnet_id            = "subnet-c1l35990"
  next_hop_destination = "172.16.128.57"
  next_hop_type        = "NORMAL_CVM"
  subnet_id            = "subnet-c1l35990"
  detect_destination_ip = [
    "10.0.0.1",
    "10.0.0.2",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `detect_destination_ip` - (Required, Set: [`String`]) An array of probe destination IPv4 addresses. Up to two.
* `net_detect_name` - (Required, String) Network probe name, the maximum length cannot exceed 60 bytes.
* `subnet_id` - (Required, String, ForceNew) Subnet instance ID. Such as:subnet-12345678.
* `vpc_id` - (Required, String, ForceNew) `VPC` instance `ID`. Such as:`vpc-12345678`.
* `net_detect_description` - (Optional, String) Network probe description.
* `next_hop_destination` - (Optional, String) The destination gateway of the next hop, the value is related to the next hop type. If the next hop type is VPN, and the value is the VPN gateway ID, such as: vpngw-12345678; If the next hop type is DIRECTCONNECT, and the value is the private line gateway ID, such as: dcg-12345678; If the next hop type is PEERCONNECTION, which takes the value of the peer connection ID, such as: pcx-12345678; If the next hop type is NAT, and the value is Nat gateway, such as: nat-12345678; If the next hop type is NORMAL_CVM, which takes the IPv4 address of the cloud server, such as: 10.0.0.12; If the next hop type is CCN, and the value is the cloud network ID, such as: ccn-12345678; If the next hop type is NONEXTHOP, and the specified network probe is a network probe without a next hop.
* `next_hop_type` - (Optional, String) The next hop type, currently we support the following types: `VPN`: VPN gateway; `DIRECTCONNECT`: private line gateway; `PEERCONNECTION`: peer connection; `NAT`: NAT gateway; `NORMAL_CVM`: normal cloud server; `CCN`: cloud networking gateway; `NONEXTHOP`: no next hop.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc net_detect can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_net_detect.net_detect net_detect_id
```

