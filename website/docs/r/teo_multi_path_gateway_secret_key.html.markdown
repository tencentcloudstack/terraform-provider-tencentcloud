---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_multi_path_gateway_secret_key"
sidebar_current: "docs-tencentcloud-resource-teo_multi_path_gateway_secret_key"
description: |-
  Provides a resource to manage a TEO multi-path gateway secret key config.
---

# tencentcloud_teo_multi_path_gateway_secret_key

Provides a resource to manage a TEO multi-path gateway secret key config.

## Example Usage

```hcl
resource "tencentcloud_teo_multi_path_gateway_secret_key" "example" {
  zone_id    = "zone-359h725djt7h"
  secret_key = base64encode("123123123")
}
```

## Argument Reference

The following arguments are supported:

* `secret_key` - (Required, String) Multi-path gateway secret key, base64 string, the string length before encoding is 32-48 characters.
* `zone_id` - (Required, String, ForceNew) Zone ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateMultiPathGatewaySecretKey](https://cloud.tencent.com/document/api/1552/121341) - This API creates an access key for the multi-channel security acceleration gateway. Customers use the access key to sign requests for accessing the gateway. Each site can have only one key, which is applicable to all gateways under that site. Query the key via the DescribeMultiPathGatewaySecretKey API.
* Create: [DescribeMultiPathGatewaySecretKey](https://cloud.tencent.com/document/api/1552/121335) - This API is used to query keys for integrating multi-channel security acceleration gateways. Customers access multi-channel security acceleration gateways based on key signature.
* Create: [ModifyMultiPathGatewaySecretKey](https://cloud.tencent.com/document/api/1552/121331) - This API is used to modify the access key for the multi-channel security acceleration gateway.The access key is used by customers to sign requests for gateway access. The original key becomes invalid after modification.
* Read: [DescribeMultiPathGatewaySecretKey](https://cloud.tencent.com/document/api/1552/121335) - This API is used to query keys for integrating multi-channel security acceleration gateways. Customers access multi-channel security acceleration gateways based on key signature.
* Update: [CreateMultiPathGatewaySecretKey](https://cloud.tencent.com/document/api/1552/121341) - This API creates an access key for the multi-channel security acceleration gateway. Customers use the access key to sign requests for accessing the gateway. Each site can have only one key, which is applicable to all gateways under that site. Query the key via the DescribeMultiPathGatewaySecretKey API.
* Update: [DescribeMultiPathGatewaySecretKey](https://cloud.tencent.com/document/api/1552/121335) - This API is used to query keys for integrating multi-channel security acceleration gateways. Customers access multi-channel security acceleration gateways based on key signature.
* Update: [ModifyMultiPathGatewaySecretKey](https://cloud.tencent.com/document/api/1552/121331) - This API is used to modify the access key for the multi-channel security acceleration gateway.The access key is used by customers to sign requests for gateway access. The original key becomes invalid after modification.

## Import

TEO multi-path gateway secret key can be imported using the id, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway_secret_key.example zone-3edjdliiw3he
```

