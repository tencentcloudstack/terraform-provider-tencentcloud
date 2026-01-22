---
subcategory: "Billing"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_billing_budget_operation_log"
sidebar_current: "docs-tencentcloud-datasource-billing_budget_operation_log"
description: |-
  Use this data source to query detailed information of billing billing_budget_operation_log
---

# tencentcloud_billing_budget_operation_log

Use this data source to query detailed information of billing billing_budget_operation_log

## Example Usage

```hcl
data "tencentcloud_billing_budget_operation_log" "billing_budget_operation_log" {
  budget_id = "1971489821259956225"
}
```

## Argument Reference

The following arguments are supported:

* `budget_id` - (Required, String) Budget id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `records` - Query data list.
  * `action` - Modification type: ADD, UPDATE.
  * `bill_day` - Bill day.
  * `bill_month` - Bill month.
  * `budget_id` - Budget item id.
  * `create_time` - Create time.
  * `diff_value` - change information.
    * `after` - Content after change.
    * `before` - Content before change.
    * `property` - Change attributes.
  * `operate_uin` - Operate uin.
  * `operation_channel` - Operation channel.
  * `owner_uin` - Owner uin.
  * `payer_uin` - Payer uin.
  * `update_time` - Update time.


