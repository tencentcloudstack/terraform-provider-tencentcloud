---
subcategory: "EventBridge(EB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eb_platform_products"
sidebar_current: "docs-tencentcloud-datasource-eb_platform_products"
description: |-
  Use this data source to query detailed information of eb platform_products
---

# tencentcloud_eb_platform_products

Use this data source to query detailed information of eb platform_products

## Example Usage

```hcl
data "tencentcloud_eb_platform_products" "platform_products" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `platform_products` - Platform product list.
  * `product_name` - Platform product name.
  * `product_type` - Platform product type.


