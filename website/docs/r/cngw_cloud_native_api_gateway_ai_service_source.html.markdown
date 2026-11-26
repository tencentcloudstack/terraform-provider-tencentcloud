---
subcategory: "Cloud Native API Gateway(CNGW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cngw_cloud_native_api_gateway_ai_service_source"
sidebar_current: "docs-tencentcloud-resource-cngw_cloud_native_api_gateway_ai_service_source"
description: |-
  Provides a resource to manage a cloud native API gateway AI service source.
---

# tencentcloud_cngw_cloud_native_api_gateway_ai_service_source

Provides a resource to manage a cloud native API gateway AI service source.

## Example Usage

```hcl
resource "tencentcloud_cngw_cloud_native_api_gateway_ai_service_source" "example" {
  gateway_id  = "gw-12345"
  source_type = "MCPRegistry"
  source_name = "tf-example-source"
  source_info = {
    instance_id = "ins-12345"
    addresses   = ["10.0.0.1:8848"]
  }
  description = "created by terraform"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String) Instance ID.
* `source_name` - (Required, String) Name of the service source.
* `source_type` - (Required, String) Source type. Valid values: `MCPRegistry`: MCP registry, `Registry`: registry, `DNS`: DNS service. It cannot be changed after creation.
* `description` - (Optional, String) Description.
* `source_info` - (Optional, Object) Source detail.
* `source_product` - (Optional, String) Source product. Valid values: `TSFNacos`: TSF Nacos.

The `source_info` object supports the following:

* `addresses` - (Optional, List) List of addresses.
* `auth` - (Optional, Object) Authorization information.
* `instance_id` - (Optional, String) Instance ID of the service source.

The `auth` object of `source_info` supports the following:

* `access_token` - (Optional, String) Access token.
* `password` - (Optional, String) Password.
* `username` - (Optional, String) Account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateCloudNativeAPIGatewayAIServiceSource](https://www.tencentcloud.com/search/?wd=CreateCloudNativeAPIGatewayAIServiceSource&type=document&lang=en)
  - [DescribeCloudNativeAPIGatewayAIServiceSourceList](https://www.tencentcloud.com/search/?wd=DescribeCloudNativeAPIGatewayAIServiceSourceList&type=document&lang=en)
* Read: [DescribeCloudNativeAPIGatewayAIServiceSourceList](https://www.tencentcloud.com/search/?wd=DescribeCloudNativeAPIGatewayAIServiceSourceList&type=document&lang=en)
* Update:
  - [DescribeCloudNativeAPIGatewayAIServiceSourceList](https://www.tencentcloud.com/search/?wd=DescribeCloudNativeAPIGatewayAIServiceSourceList&type=document&lang=en)
  - [ModifyCloudNativeAPIGatewayAIServiceSource](https://www.tencentcloud.com/search/?wd=ModifyCloudNativeAPIGatewayAIServiceSource&type=document&lang=en)
* Delete: [DeleteCloudNativeAPIGatewayAIServiceSource](https://www.tencentcloud.com/search/?wd=DeleteCloudNativeAPIGatewayAIServiceSource&type=document&lang=en)

## Import

Cloud native API gateway AI service source can be imported using the id, e.g.

```
terraform import tencentcloud_cngw_cloud_native_api_gateway_ai_service_source.example source-12345
```

