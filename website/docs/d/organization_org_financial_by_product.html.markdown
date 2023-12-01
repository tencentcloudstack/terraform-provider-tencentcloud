---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_financial_by_product"
sidebar_current: "docs-tencentcloud-datasource-organization_org_financial_by_product"
description: |-
  Use this data source to query detailed information of organization org_financial_by_product
---

# tencentcloud_organization_org_financial_by_product

Use this data source to query detailed information of organization org_financial_by_product

## Example Usage

```hcl
data "tencentcloud_organization_org_financial_by_product" "org_financial_by_product" {
  month         = "2023-05"
  end_month     = "2023-09"
  product_codes = ["p_eip"]
}
```

## Argument Reference

The following arguments are supported:

* `month` - (Required, String) Query for the start month. Format:yyyy-mm, for example:2021-01.
* `end_month` - (Optional, String) Query for the end month. Format:yyyy-mm, for example:2021-01.The default value is the `Month`.
* `member_uins` - (Optional, Set: [`Int`]) Member uin list. Up to 100.
* `product_codes` - (Optional, Set: [`String`]) Product code list. Up to 100.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Organization financial info by products.
  * `product_code` - Product code.
  * `product_name` - Product name.
  * `ratio` - The percentage of the organization total cost that is accounted for by the product.
  * `total_cost` - Total cost of the product.
* `total_cost` - Total cost of the product.


