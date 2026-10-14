---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ipv6_address_bandwidth"
sidebar_current: "docs-tencentcloud-resource-ipv6_address_bandwidth"
description: |-
  Provides a resource to create a ipv6_address_bandwidth
---

# tencentcloud_ipv6_address_bandwidth

Provides a resource to create a ipv6_address_bandwidth

## Example Usage

```hcl
resource "tencentcloud_ipv6_address_bandwidth" "ipv6_address_bandwidth" {
  ipv6_address               = "2402:4e00:1019:9400:0:9905:a90b:2ef0"
  internet_max_bandwidth_out = 6
  internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
  #  bandwidth_package_id       = "bwp-34rfgt56"
}
```

## Argument Reference

The following arguments are supported:

* `ipv6_address` - (Required, String, ForceNew) IPV6 address that needs to be enabled for public network access.
* `bandwidth_package_id` - (Optional, String) The bandwidth package id, the Legacy account and the ipv6 address to apply for the bandwidth package charge type need to be passed in.
* `internet_charge_type` - (Optional, String) Network billing mode. IPV6 currently supports: `TRAFFIC_POSTPAID_BY_HOUR`, for standard account types; `BANDWIDTH_PACKAGE`, for traditional account types. The default network billing mode is: `TRAFFIC_POSTPAID_BY_HOUR`.
* `internet_max_bandwidth_out` - (Optional, Int) Bandwidth, in Mbps. The default is 1Mbps.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AllocateIp6AddressesBandwidth](https://cloud.tencent.com/document/api/215/40090) - This API is used to allocate IPv6 public network bandwidth for classic elastic public IPv6 addresses.
* Create: [DescribeTaskResult](https://cloud.tencent.com/document/api/215/36271) - This API is used to query the EIP async job execution results.
* Read: [DescribeIp6Addresses](https://cloud.tencent.com/document/api/215/40089) - This API is used to query the detailed information on one or multiple classic elastic public IPv6 instances.
* Update: [DescribeTaskResult](https://cloud.tencent.com/document/api/215/36271) - This API is used to query the EIP async job execution results.
* Update: [ModifyIp6AddressesBandwidth](https://cloud.tencent.com/document/api/215/40088) - This API is used to adjust the bandwidth limit of a classic elastic public IPv6 instance.
* Delete: [ReleaseIp6AddressesBandwidth](https://cloud.tencent.com/document/api/215/40087) - This API is used to release the IPv6 public network bandwidth of classic elastic public IPv6 instances.

