---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elastic_public_ipv6"
sidebar_current: "docs-tencentcloud-resource-elastic_public_ipv6"
description: |-
  Provides a resource to create a vpc elastic_public_ipv6
---

# tencentcloud_elastic_public_ipv6

Provides a resource to create a vpc elastic_public_ipv6

## Example Usage

```hcl
resource "tencentcloud_elastic_public_ipv6" "elastic_public_ipv6" {
  address_name               = "test"
  internet_max_bandwidth_out = 1
  tags = {
    "test1key" = "test1value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `address_ip` - (Optional, String) External network IP address.
* `address_name` - (Optional, String) EIP name, used to customize the personalized name of the EIP when applying for EIP. Default value: unnamed.
* `address_type` - (Optional, String) Elastic IPv6 type, optional values:
	- EIPv6: Ordinary IPv6
	- HighQualityEIPv6: Premium IPv6
Note: You need to contact the product to open a premium IPv6 white list, and only some regions support premium IPv6
Default value: EIPv6.
* `bandwidth_package_id` - (Optional, String) Bandwidth packet unique ID parameter. If this parameter is set and the InternetChargeType is BANDWIDTH_PACKAGE, it means that the EIP created is added to the BGP bandwidth packet and the bandwidth packet is charged.
* `egress` - (Optional, String) Elastic IPv6 network exit, optional values:
	- CENTER_EGRESS_1: Center Exit 1
	- CENTER_EGRESS_2: Center Exit 2
	- CENTER_EGRESS_3: Center Exit 3
Note: Network exports corresponding to different operators or resource types need to contact the product for clarification
Default value: CENTER_EGRESS_1.
* `internet_charge_type` - (Optional, String) Elastic IPv6 charging method, optional values:
	- BANDWIDTH_PACKAGE: Payment for Shared Bandwidth Package
	- TRAFFIC_POSTPAID_BY_HOUR: Traffic is paid by the hour
Default value: TRAFFIC_POSTPAID_BY_HOUR.
* `internet_max_bandwidth_out` - (Optional, Int) Elastic IPv6 bandwidth limit in Mbps.
The range of selectable values depends on the EIP billing method:
	- BANDWIDTH_PACKAGE: 1 Mbps to 2000 Mbps
	- TRAFFIC_POSTPAID_BY_HOUR: 1 Mbps to 100 Mbps
Default value: 1 Mbps.
* `internet_service_provider` - (Optional, String) Elastic IPv6 line type, default value: BGP.
For users who have activated a static single-line IP whitelist, selectable values:
	- CMCC: China Mobile
	- CTCC: China Telecom
	- CUCC: China Unicom
Note: Static single-wire IP is only supported in some regions.
* `tags` - (Optional, Map) Tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc elastic_public_ipv6 can be imported using the id, e.g.

```
terraform import tencentcloud_elastic_public_ipv6.elastic_public_ipv6 elastic_public_ipv6_id
```

