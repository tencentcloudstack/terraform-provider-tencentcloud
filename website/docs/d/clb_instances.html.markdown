---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_instances"
sidebar_current: "docs-tencentcloud-datasource-clb_instances"
description: |-
  Use this data source to query detailed information of CLB
---

# tencentcloud_clb_instances

Use this data source to query detailed information of CLB

## Example Usage

```hcl
data "tencentcloud_clb_instances" "foo" {
  clb_id             = "lb-k2zjp9lv"
  network_type       = "OPEN"
  clb_name           = "myclb"
  project_id         = 0
  result_output_file = "mytestpath"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Optional) ID of the CLB to be queried.
* `clb_name` - (Optional) Name of the CLB to be queried.
* `network_type` - (Optional) Type of CLB instance, and available values include `OPEN` and `INTERNAL`.
* `project_id` - (Optional) Project ID of the CLB.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `clb_list` - A list of cloud load balancers. Each element contains the following attributes:
  * `address_ip_version` - IP version, only applicable to open CLB. Valid values are `IPV4`, `IPV6` and `IPv6FullChain`.
  * `clb_id` - ID of CLB.
  * `clb_name` - Name of CLB.
  * `clb_vips` - The virtual service address table of the CLB.
  * `create_time` - Create time of the CLB.
  * `internet_bandwidth_max_out` - Max bandwidth out, only applicable to open CLB. Valid value ranges is [1, 2048]. Unit is MB.
  * `internet_charge_type` - Internet charge type, only applicable to open CLB. Valid values are `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.
  * `network_type` - Types of CLB.
  * `project_id` - ID of the project.
  * `security_groups` - ID set of the security groups.
  * `status_time` - Latest state transition time of CLB.
  * `status` - The status of CLB.
  * `subnet_id` - ID of the subnet.
  * `tags` - The available tags within this CLB.
  * `target_region_info_region` - Region information of backend service are attached the CLB.
  * `target_region_info_vpc_id` - VpcId information of backend service are attached the CLB.
  * `vip_isp` - Network operator, only applicable to open CLB. Valid values are `CMCC`(China Mobile), `CTCC`(Telecom), `CUCC`(China Unicom) and `BGP`. If this ISP is specified, network billing method can only use the bandwidth package billing (BANDWIDTH_PACKAGE).
  * `vpc_id` - ID of the VPC.


