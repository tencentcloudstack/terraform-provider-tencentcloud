---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api_doc"
sidebar_current: "docs-tencentcloud-resource-api_gateway_api_doc"
description: |-
  Provides a resource to create a APIGateway ApiDoc
---

# tencentcloud_api_gateway_api_doc

Provides a resource to create a APIGateway ApiDoc

## Example Usage

```hcl
resource "tencentcloud_api_gateway_api_doc" "my_api_doc" {
  api_doc_name = "doc_test1"
  service_id   = "service_test1"
  environment  = "release"
  api_ids      = ["api-test1", "api-test2"]
}
```

## Argument Reference

The following arguments are supported:

* `api_doc_name` - (Required, String) Api Document name.
* `api_ids` - (Required, Set: [`String`]) List of APIs for generating documents.
* `environment` - (Required, String) Env name.
* `service_id` - (Required, String) Service name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_count` - Api Document count.
* `api_doc_id` - Api Document ID.
* `api_doc_status` - API Document Build Status.
* `api_doc_uri` - API Document Access URI.
* `api_names` - List of names for generating documents.
* `release_count` - Number of API document releases.
* `service_name` - API Document service name.
* `share_password` - API Document Sharing Password.
* `updated_time` - API Document update time.
* `view_count` - API Document Viewing Times.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAPIDoc](https://www.tencentcloud.com/search/?wd=CreateAPIDoc&type=document&lang=en)
* Create: [DescribeAPIDocDetail](https://www.tencentcloud.com/search/?wd=DescribeAPIDocDetail&type=document&lang=en)
* Read: [DescribeAPIDocDetail](https://www.tencentcloud.com/search/?wd=DescribeAPIDocDetail&type=document&lang=en)
* Update: [DescribeAPIDocDetail](https://www.tencentcloud.com/search/?wd=DescribeAPIDocDetail&type=document&lang=en)
* Update: [ModifyAPIDoc](https://www.tencentcloud.com/search/?wd=ModifyAPIDoc&type=document&lang=en)
* Delete: [DeleteAPIDoc](https://www.tencentcloud.com/search/?wd=DeleteAPIDoc&type=document&lang=en)

