---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_cloud_native_api_gateway_ip_restriction"
sidebar_current: "docs-tencentcloud-resource-tse_cloud_native_api_gateway_ip_restriction"
description: |-
  Provides a resource to manage a TSE cloud native API gateway IP restriction.
---

# tencentcloud_tse_cloud_native_api_gateway_ip_restriction

Provides a resource to manage a TSE cloud native API gateway IP restriction.

## Example Usage

```hcl
resource "tencentcloud_tse_cloud_native_api_gateway_ip_restriction" "example" {
  gateway_id       = "gateway-ddbb709b"
  source_type      = "route"
  source_id        = "route-xxxx"
  enabled          = true
  restriction_type = "whiteList"
  address_list     = ["10.0.0.0/8", "192.168.1.1"]
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Gateway ID.
* `source_id` - (Required, String, ForceNew) Route or service ID.
* `source_type` - (Required, String, ForceNew) Resource type bound to the IP restriction plugin: route|service.
* `address_list` - (Optional, Set: [`String`]) IP/CIDR address list.
* `enabled` - (Optional, Bool) Whether to enable the plugin.
* `restriction_type` - (Optional, String) IP restriction type: whiteList|blackList.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateOrModifyCloudNativeAPIGatewayIPRestriction](https://www.tencentcloud.com/search/?wd=CreateOrModifyCloudNativeAPIGatewayIPRestriction&type=document&lang=en)
* Create: [DescribeCloudNativeAPIGatewayIPRestriction](https://www.tencentcloud.com/search/?wd=DescribeCloudNativeAPIGatewayIPRestriction&type=document&lang=en)
* Read: [DescribeCloudNativeAPIGatewayIPRestriction](https://www.tencentcloud.com/search/?wd=DescribeCloudNativeAPIGatewayIPRestriction&type=document&lang=en)
* Update: [CreateOrModifyCloudNativeAPIGatewayIPRestriction](https://www.tencentcloud.com/search/?wd=CreateOrModifyCloudNativeAPIGatewayIPRestriction&type=document&lang=en)
* Update: [DescribeCloudNativeAPIGatewayIPRestriction](https://www.tencentcloud.com/search/?wd=DescribeCloudNativeAPIGatewayIPRestriction&type=document&lang=en)
* Delete: [DeleteCloudNativeAPIGatewayIPRestriction](https://www.tencentcloud.com/search/?wd=DeleteCloudNativeAPIGatewayIPRestriction&type=document&lang=en)

## Import

tse cloud native api gateway ip restriction can be imported using the composite id, e.g.

```
terraform import tencentcloud_tse_cloud_native_api_gateway_ip_restriction.example gatewayId#sourceType#sourceId
```

