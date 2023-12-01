---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_gateway_flow_monitor_detail"
sidebar_current: "docs-tencentcloud-datasource-vpc_gateway_flow_monitor_detail"
description: |-
  Use this data source to query detailed information of vpc gateway_flow_monitor_detail
---

# tencentcloud_vpc_gateway_flow_monitor_detail

Use this data source to query detailed information of vpc gateway_flow_monitor_detail

## Example Usage

```hcl
data "tencentcloud_vpc_gateway_flow_monitor_detail" "gateway_flow_monitor_detail" {
  time_point      = "2023-06-02 12:15:20"
  vpn_id          = "vpngw-gt8bianl"
  order_field     = "OutTraffic"
  order_direction = "DESC"
}
```

## Argument Reference

The following arguments are supported:

* `time_point` - (Required, String) The point in time. This indicates details of this minute will be queried. For example, in `2019-02-28 18:15:20`, details at `18:15` will be queried.
* `direct_connect_gateway_id` - (Optional, String) The instance ID of the Direct Connect gateway, such as `dcg-ltjahce6`.
* `nat_id` - (Optional, String) The instance ID of the NAT gateway, such as `nat-ltjahce6`.
* `order_direction` - (Optional, String) Order methods. Ascending: `ASC`, Descending: `DESC`.
* `order_field` - (Optional, String) The order field supports `InPkg`, `OutPkg`, `InTraffic`, and `OutTraffic`.
* `peering_connection_id` - (Optional, String) The instance ID of the peering connection, such as `pcx-ltjahce6`.
* `result_output_file` - (Optional, String) Used to save results.
* `vpn_id` - (Optional, String) The instance ID of the VPN gateway, such as `vpn-ltjahce6`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `gateway_flow_monitor_detail_set` - The gateway traffic monitoring details.
  * `in_pkg` - Inbound packets.
  * `in_traffic` - Inbound traffic, in Byte.
  * `out_pkg` - Outbound packets.
  * `out_traffic` - Outbound traffic, in Byte.
  * `private_ip_address` - Origin `IP`.


