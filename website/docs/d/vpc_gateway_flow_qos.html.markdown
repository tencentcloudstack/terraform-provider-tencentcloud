---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_gateway_flow_qos"
sidebar_current: "docs-tencentcloud-datasource-vpc_gateway_flow_qos"
description: |-
  Use this data source to query detailed information of vpc gateway_flow_qos
---

# tencentcloud_vpc_gateway_flow_qos

Use this data source to query detailed information of vpc gateway_flow_qos

## Example Usage

```hcl
data "tencentcloud_vpc_gateway_flow_qos" "gateway_flow_qos" {
  gateway_id = "vpngw-gt8bianl"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String) Network instance ID, the network instance types we currently support are:Private line gateway instance ID, in the form of `dcg-ltjahce6`;Nat gateway instance ID, in the form of `nat-ltjahce6`;VPN gateway instance ID, in the form of `vpn-ltjahce6`.
* `ip_addresses` - (Optional, Set: [`String`]) Intranet IP of the cloud server with traffic limitation.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `gateway_qos_set` - instance detail list.
  * `bandwidth` - bandwidth value.
  * `create_time` - create time.
  * `ip_address` - cvm ip address.
  * `vpc_id` - vpc id.


