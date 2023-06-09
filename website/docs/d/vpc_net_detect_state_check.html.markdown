---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_net_detect_state_check"
sidebar_current: "docs-tencentcloud-datasource-vpc_net_detect_state_check"
description: |-
  Use this data source to query detailed information of vpc net_detect_state_check
---

# tencentcloud_vpc_net_detect_state_check

Use this data source to query detailed information of vpc net_detect_state_check

## Example Usage

```hcl
data "tencentcloud_vpc_net_detect_state_check" "net_detect_state_check" {
  net_detect_id = "netd-12345678"
  detect_destination_ip = [
    "10.0.0.3",
    "10.0.0.2"
  ]
  next_hop_type        = "NORMAL_CVM"
  next_hop_destination = "10.0.0.4"
}
```

## Argument Reference

The following arguments are supported:

* `detect_destination_ip` - (Required, Set: [`String`]) The array of detection destination IPv4 addresses, which contains at most two IP addresses.
* `next_hop_destination` - (Required, String) The next-hop destination gateway. The value is related to NextHopType.If NextHopType is set to VPN, the value of this parameter is the VPN gateway ID, such as vpngw-12345678.If NextHopType is set to DIRECTCONNECT, the value of this parameter is the direct connect gateway ID, such as dcg-12345678.If NextHopType is set to PEERCONNECTION, the value of this parameter is the peering connection ID, such as pcx-12345678.If NextHopType is set to NAT, the value of this parameter is the NAT gateway ID, such as nat-12345678.If NextHopType is set to NORMAL_CVM, the value of this parameter is the IPv4 address of the CVM, such as 10.0.0.12.
* `next_hop_type` - (Required, String) The type of the next hop. Currently supported types are:VPN: VPN gateway;DIRECTCONNECT: direct connect gateway;PEERCONNECTION: peering connection;NAT: NAT gateway;NORMAL_CVM: normal CVM.
* `net_detect_id` - (Optional, String) ID of a network inspector instance, e.g. netd-12345678. Enter at least one of this parameter, VpcId, SubnetId, and NetDetectName. Use NetDetectId if it is present.
* `net_detect_name` - (Optional, String) The name of a network inspector, up to 60 bytes in length. It is used together with VpcId and NetDetectName. You should enter either this parameter or NetDetectId, or both. Use NetDetectId if it is present.
* `result_output_file` - (Optional, String) Used to save results.
* `subnet_id` - (Optional, String) ID of a subnet instance, e.g. `subnet-12345678`, which is used together with VpcId and NetDetectName. You should enter either this parameter or NetDetectId, or both. Use NetDetectId if it is present.
* `vpc_id` - (Optional, String) ID of a `VPC` instance, e.g. `vpc-12345678`, which is used together with SubnetId and NetDetectName. You should enter either this parameter or NetDetectId, or both. Use NetDetectId if it is present.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `net_detect_ip_state_set` - The array of network detection verification results.
  * `delay` - The latency. Unit: ms.
  * `detect_destination_ip` - The destination IPv4 address of network detection.
  * `packet_loss_rate` - The packet loss rate.
  * `state` - The detection result.0: successful;-1: no packet loss occurred during routing;-2: packet loss occurred when outbound traffic is blocked by the ACL;-3: packet loss occurred when inbound traffic is blocked by the ACL;-4: other errors.


