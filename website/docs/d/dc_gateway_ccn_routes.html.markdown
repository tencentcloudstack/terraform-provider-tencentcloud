---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_gateway_ccn_routes"
sidebar_current: "docs-tencentcloud-datasource-dc_gateway_ccn_routes"
description: |-
  Use this data source to query detailed information of direct connect gateway route entries.
---

# tencentcloud_dc_gateway_ccn_routes

Use this data source to query detailed information of direct connect gateway route entries.

## Example Usage

```hcl
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_dc_gateway" "ccn_main" {
  name                = "ci-cdg-ccn-test"
  network_instance_id = "${tencentcloud_ccn.main.id}"
  network_type        = "CCN"
  gateway_type        = "NORMAL"
}

resource "tencentcloud_dc_gateway_ccn_route" "route1" {
  dcg_id     = "${tencentcloud_dc_gateway.ccn_main.id}"
  cidr_block = "10.1.1.0/32"
}

resource "tencentcloud_dc_gateway_ccn_route" "route2" {
  dcg_id     = "${tencentcloud_dc_gateway.ccn_main.id}"
  cidr_block = "192.1.1.0/32"
}

#You need to sleep for a few seconds because there is a cache on the server
data "tencentcloud_dc_gateway_ccn_routes"  "test" {
  dcg_id = "${tencentcloud_dc_gateway.ccn_main.id}"
}
```

## Argument Reference

The following arguments are supported:

* `dcg_id` - (Required) ID of the DCG to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - Information list of the DCG route entries.
  * `as_path` - As_Path list of the BGP.
  * `cidr_block` - A network address segment of IDC.
  * `dcg_id` - ID of the DCG.
  * `route_id` - ID of the DCG route.


