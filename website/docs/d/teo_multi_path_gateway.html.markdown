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

* `gateways` - Multi-path gateway list.
  * `gateway_id` - Gateway ID.
  * `gateway_ip` - Gateway IP address in IPv4 format.
  * `gateway_name` - Gateway name.
  * `gateway_port` - Gateway port, range 1-65535 (excluding 8888).
  * `gateway_type` - Gateway type. Possible values are: cloud: cloud gateway; private: private gateway.
  * `need_confirm` - Whether the gateway origin IP list has changed and needs re-confirmation. Possible values are: true: changed, needs confirmation; false: not changed, no confirmation needed.
  * `region_id` - Gateway region ID.
  * `status` - Gateway status. Possible values are: creating: creating; online: online; offline: offline; disable: disabled.


