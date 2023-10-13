---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_financial_by_month"
sidebar_current: "docs-tencentcloud-datasource-organization_org_financial_by_month"
description: |-
  Use this data source to query detailed information of organization org_financial_by_month
---

# tencentcloud_organization_org_financial_by_month

Use this data source to query detailed information of organization org_financial_by_month

## Example Usage

```hcl
data "tencentcloud_organization_org_financial_by_month" "org_financial_by_month" {
  end_month   = "2023-05"
  member_uins = [100026517717]
}
```

## Argument Reference

The following arguments are supported:

* `end_month` - (Optional, String) Query for the end month. Format:yyyy-mm, for example:2021-01.
* `member_uins` - (Optional, Set: [`Int`]) Member uin list. Up to 100.
* `product_codes` - (Optional, Set: [`String`]) Product code list. Up to 100.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Organization financial info by month.
  * `growth_rate` - Growth rate compared to last month.
  * `id` - Record ID.
  * `month` - Month.
  * `total_cost` - Total cost of the month.


