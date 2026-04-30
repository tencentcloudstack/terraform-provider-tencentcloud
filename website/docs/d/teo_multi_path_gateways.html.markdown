---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_multi_path_gateways"
sidebar_current: "docs-tencentcloud-datasource-teo_multi_path_gateways"
description: |-
  Use this data source to query detailed information of TEO multi path gateways
---

# tencentcloud_teo_multi_path_gateways

Use this data source to query detailed information of TEO multi path gateways

## Example Usage

```hcl
data "tencentcloud_teo_multi_path_gateways" "example" {
  zone_id = "zone-2o1xvpmq7nn"
  filters {
    name   = "gateway-type"
    values = ["cloud"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Site ID.
* `filters` - (Optional, List) Filter conditions. The maximum value of Filters.Values is 20. If this parameter is not filled in, all gateway information under the current appid will be returned. Detailed filter conditions are as follows: gateway-type: filter by gateway type, supporting values cloud and private, representing filtering cloud gateways and private gateways respectively; keyword: filter by gateway name keyword.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name.
* `values` - (Required, Set) Filter field values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `gateways` - Gateway details.
  * `gateway_id` - Gateway ID.
  * `gateway_ip` - Gateway IP, in IPv4 format.
  * `gateway_name` - Gateway name.
  * `gateway_port` - Gateway port, range 1-65535 (excluding 8888).
  * `gateway_type` - Gateway type. Valid values: cloud (cloud gateway managed by Tencent Cloud), private (private gateway deployed by user).
  * `lines` - Line information. Returned when querying gateway details via DescribeMultiPathGateway, not returned when querying gateway list via DescribeMultiPathGateways.
    * `line_address` - Line address, in host:port format.
    * `line_id` - Line ID. line-0 and line-1 are built-in line IDs. line-0: direct line, does not support add/edit/delete; line-1: EdgeOne L4 proxy line, supports modifying instances and rules, does not support delete; line-2 and above: EdgeOne L4 proxy line or custom line, supports modify/delete instances and rules.
    * `line_type` - Line type. Valid values: direct (direct line, does not support edit/delete), proxy (EdgeOne L4 proxy line, supports editing instances and rules, does not support delete), custom (custom line, supports edit and delete).
    * `proxy_id` - L4 proxy instance ID, returned when LineType is proxy (EdgeOne L4 proxy).
    * `rule_id` - Forwarding rule ID, returned when LineType is proxy (EdgeOne L4 proxy).
  * `need_confirm` - Whether the gateway origin IP list has changed and needs re-confirmation. true: origin IP list changed, needs confirmation; false: origin IP list unchanged, no confirmation needed.
  * `region_id` - Gateway region ID, can be obtained from DescribeMultiPathGatewayRegions API.
  * `status` - Gateway status. Valid values: creating, online, offline, disable.


