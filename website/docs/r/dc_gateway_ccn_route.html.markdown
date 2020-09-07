---
subcategory: "Direct Connect Gateway(DCG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_gateway_ccn_route"
sidebar_current: "docs-tencentcloud-resource-dc_gateway_ccn_route"
description: |-
  Provides a resource to creating direct connect gateway route entry.
---

# tencentcloud_dc_gateway_ccn_route

Provides a resource to creating direct connect gateway route entry.

## Example Usage

```hcl
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_dc_gateway" "ccn_main" {
  name                = "ci-cdg-ccn-test"
  network_instance_id = tencentcloud_ccn.main.id
  network_type        = "CCN"
  gateway_type        = "NORMAL"
}

resource "tencentcloud_dc_gateway_ccn_route" "route1" {
  dcg_id     = tencentcloud_dc_gateway.ccn_main.id
  cidr_block = "10.1.1.0/32"
}

resource "tencentcloud_dc_gateway_ccn_route" "route2" {
  dcg_id     = tencentcloud_dc_gateway.ccn_main.id
  cidr_block = "192.1.1.0/32"
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required, ForceNew) A network address segment of IDC.
* `dcg_id` - (Required, ForceNew) ID of the DCG.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `as_path` - As_Path list of the BGP.


