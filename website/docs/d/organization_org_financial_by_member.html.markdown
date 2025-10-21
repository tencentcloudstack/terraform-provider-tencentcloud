---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_financial_by_member"
sidebar_current: "docs-tencentcloud-datasource-organization_org_financial_by_member"
description: |-
  Use this data source to query detailed information of organization org_financial_by_member
---

# tencentcloud_organization_org_financial_by_member

Use this data source to query detailed information of organization org_financial_by_member

## Example Usage

```hcl
data "tencentcloud_organization_org_financial_by_member" "org_financial_by_member" {
  month       = "2023-05"
  end_month   = "2023-10"
  member_uins = [100015591986, 100029796005]
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

* `items` - Member financial detail.
  * `member_name` - Member name.
  * `member_uin` - Member uin.
  * `ratio` - The percentage of the organization total cost that is accounted for by the member.
  * `total_cost` - Total cost of the member.
* `total_cost` - Total cost of the member.


