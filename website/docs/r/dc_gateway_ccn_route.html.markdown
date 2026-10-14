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

* `cidr_block` - (Required, String, ForceNew) A network address segment of IDC.
* `dcg_id` - (Required, String, ForceNew) ID of the DCG.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `as_path` - As path list of the BGP.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateDirectConnectGatewayCcnRoutes](https://www.tencentcloud.com/document/product/215/34777) - This API (CreateDirectConnectGatewayCcnRoutes) is used to create the CCN route (IDC IP range) of a Direct Connect gateway.
* Create: [DescribeDirectConnectGatewayCcnRoutes](https://www.tencentcloud.com/document/product/215/34775) - This API (DescribeDirectConnectGatewayCcnRoutes) is used to query the CCN routes (IDC IP range) of the Direct Connect gateway.
* Read: [DescribeDirectConnectGatewayCcnRoutes](https://www.tencentcloud.com/document/product/215/34775) - This API (DescribeDirectConnectGatewayCcnRoutes) is used to query the CCN routes (IDC IP range) of the Direct Connect gateway.
* Delete: [DeleteDirectConnectGatewayCcnRoutes](https://www.tencentcloud.com/document/product/215/34776) - This API (DeleteDirectConnectGatewayCcnRoutes) is used to delete the CCN routes (IDC IP range) of a Direct Connect gateway.
* Delete: [DescribeDirectConnectGatewayCcnRoutes](https://www.tencentcloud.com/document/product/215/34775) - This API (DescribeDirectConnectGatewayCcnRoutes) is used to query the CCN routes (IDC IP range) of the Direct Connect gateway.

## Import

dc_gateway_ccn_route can be imported using the id, e.g.

```
terraform import tencentcloud_dc_gateway_ccn_route.example dcg-dr1y0hu7#192.1.1.0/32
```

