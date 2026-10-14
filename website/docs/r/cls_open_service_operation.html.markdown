---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_open_service_operation"
sidebar_current: "docs-tencentcloud-resource-cls_open_service_operation"
description: |-
  Provides a resource to create a CLS open service operation
---

# tencentcloud_cls_open_service_operation

Provides a resource to create a CLS open service operation

## Example Usage

```hcl
resource "tencentcloud_cls_open_service_operation" "example" {}
```

## Argument Reference

The following arguments are supported:



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Account service status. `0`: service opened, `1`: service not opened.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [GetClsService](https://www.tencentcloud.com/document/product/614/80501) - This API is used to check whether CLS is enabled.
* Create: [OpenClsService](https://www.tencentcloud.com/document/product/614/80500) - This API is used to enable CLS in all regions by filling any region for Region, recommend using Guangzhou (ap-guangzhou).
* Read: [GetClsService](https://www.tencentcloud.com/document/product/614/80501) - This API is used to check whether CLS is enabled.

