---
subcategory: "Regional Management(region)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_products"
sidebar_current: "docs-tencentcloud-datasource-products"
description: |-
  Use this data source to query products that support region list queries.
---

# tencentcloud_products

Use this data source to query products that support region list queries.

## Example Usage

```hcl
data "tencentcloud_products" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `product_list` - Product list.
  * `name` - Product name, e.g. `cvm`.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeProducts](https://cloud.tencent.com/document/api/1596/77931) - This interface (DescribeProducts) is used for querying product information in each supported region list.

