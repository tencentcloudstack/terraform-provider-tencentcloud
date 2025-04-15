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

* `clb_id` - (Optional, String) ID of the CLB to be queried.
* `clb_name` - (Optional, String) Name of the CLB to be queried.
* `master_zone` - (Optional, String) Master available zone id.
* `network_type` - (Optional, String) Type of CLB instance, and available values include `OPEN` and `INTERNAL`.
* `project_id` - (Optional, Int) Project ID of the CLB.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `clb_list` - A list of cloud load balancers. Each element contains the following attributes:
  * `address_ip_version` - IP version, only applicable to open CLB. Valid values are `IPV4`, `IPV6` and `IPv6FullChain`.
  * `clb_id` - ID of CLB.
  * `clb_name` - Name of CLB.
  * `clb_vips` - The virtual service address table of the CLB.
  * `cluster_id` - ID of the cluster.
  * `create_time` - Create time of the CLB.
  * `internet_bandwidth_max_out` - Max bandwidth out, only applicable to open CLB. Valid value ranges is [1, 2048]. Unit is MB.
  * `internet_charge_type` - Internet charge type, only applicable to open CLB. Valid values are `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.
  * `local_zone` - Whether this available zone is local zone, This field maybe null, means cannot get a valid value.
  * `network_type` - Types of CLB.
  * `numerical_vpc_id` - VPC ID in a numeric form. Note: This field may return null, indicating that no valid values can be obtained.
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
  * `zone_id` - Available zone unique id(numerical representation), This field maybe null, means cannot get a valid value.
  * `zone_name` - Available zone name, This field maybe null, means cannot get a valid value.
  * `zone_region` - Region that this available zone belong to, This field maybe null, means cannot get a valid value.
  * `zone` - Available zone unique id(string representation), This field maybe null, means cannot get a valid value.


