---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_multi_path_gateway"
sidebar_current: "docs-tencentcloud-datasource-teo_multi_path_gateway"
description: |-
  Use this data source to query detailed information of TEO multi-path gateways
---

# tencentcloud_teo_multi_path_gateway

Use this data source to query detailed information of TEO multi-path gateways

## Example Usage

### Query all gateways by zone_id

```hcl
data "tencentcloud_teo_multi_path_gateway" "example" {
  zone_id = "zone-2noq7st5t3t6"
}
```

### Query gateways by zone_id with filters

```hcl
data "tencentcloud_teo_multi_path_gateway" "example" {
  zone_id = "zone-2noq7st5t3t6"
  filters {
    name   = "gateway-type"
    values = ["cloud"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Site ID.
* `filters` - (Optional, List) Filter conditions for querying multi-path gateways. The detailed filtering conditions are as follows: <li>gateway-type: Filter by gateway type, supported values are cloud and private.</li><li>keyword: Filter by gateway name keyword.</li>.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter name.
* `values` - (Required, Set) Filter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `gateways` - Gateway details.
  * `gateway_id` - Gateway ID.
  * `gateway_ip` - Gateway IP, in IPv4 format.
  * `gateway_name` - Gateway name.
  * `gateway_port` - Gateway port, range 1-65535 (excluding 8888).
  * `gateway_type` - Gateway type. Valid values: <li>cloud: Cloud gateway, created and managed by Tencent Cloud.</li><li>private: Private gateway, deployed by the user.</li>.
  * `lines` - Line information, returned when querying gateway details via DescribeMultiPathGateway, but not returned when querying gateway list via DescribeMultiPathGateways.
    * `line_address` - Line address, in host:port format.
    * `line_id` - Line ID. line-0 and line-1 are built-in line IDs. Valid values: <li>line-0: Direct line, does not support adding, editing, or deletion;</li><li>line-1: EdgeOne Layer-4 proxy line, supports modifying instances and rules, does not support deletion;</li><li>line-2 and above: EdgeOne Layer-4 proxy line or custom line, supports modifying and deleting instances and rules.</li>.
    * `line_type` - Line type. Valid values: <li>direct: Direct line, does not support editing or deletion;</li><li>proxy: EdgeOne Layer-4 proxy line, supports editing instances and rules, does not support deletion;</li><li>custom: Custom line, supports editing and deletion.</li>.
    * `proxy_id` - Layer-4 proxy instance ID, returned when LineType is proxy (EdgeOne Layer-4 proxy).
    * `rule_id` - Forwarding rule ID, returned when LineType is proxy (EdgeOne Layer-4 proxy).
  * `need_confirm` - Whether reconfirmation is needed when the gateway origin IP list changes. Valid values: <li>true: The origin IP list has changed and needs confirmation;</li><li>false: The origin IP list has not changed and no confirmation is needed.</li>.
  * `region_id` - Gateway region ID, which can be obtained from the DescribeMultiPathGatewayRegions API.
  * `status` - Gateway status. Valid values: <li>creating: Creating;</li><li>online: Online;</li><li>offline: Offline;</li><li>disable: Disabled.</li>.


