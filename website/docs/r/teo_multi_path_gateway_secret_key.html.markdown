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

* Create: [CreateMultiPathGatewaySecretKey](https://www.tencentcloud.com/search/?wd=CreateMultiPathGatewaySecretKey&type=document&lang=en)
* Create: [DescribeMultiPathGatewaySecretKey](https://www.tencentcloud.com/search/?wd=DescribeMultiPathGatewaySecretKey&type=document&lang=en)
* Create: [ModifyMultiPathGatewaySecretKey](https://www.tencentcloud.com/search/?wd=ModifyMultiPathGatewaySecretKey&type=document&lang=en)
* Read: [DescribeMultiPathGatewaySecretKey](https://www.tencentcloud.com/search/?wd=DescribeMultiPathGatewaySecretKey&type=document&lang=en)
* Update: [CreateMultiPathGatewaySecretKey](https://www.tencentcloud.com/search/?wd=CreateMultiPathGatewaySecretKey&type=document&lang=en)
* Update: [DescribeMultiPathGatewaySecretKey](https://www.tencentcloud.com/search/?wd=DescribeMultiPathGatewaySecretKey&type=document&lang=en)
* Update: [ModifyMultiPathGatewaySecretKey](https://www.tencentcloud.com/search/?wd=ModifyMultiPathGatewaySecretKey&type=document&lang=en)

## Import

TEO multi-path gateway secret key can be imported using the id, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway_secret_key.example zone-3edjdliiw3he
```

