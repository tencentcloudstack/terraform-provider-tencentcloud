---
subcategory: "Billing"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_billing_bill_detail"
sidebar_current: "docs-tencentcloud-datasource-billing_bill_detail"
description: |-
  Use this data source to query bill detail (L3 detail bill) of billing.
---

# tencentcloud_billing_bill_detail

Use this data source to query bill detail (L3 detail bill) of billing.

## Example Usage

```hcl
data "tencentcloud_billing_bill_detail" "example_by_month" {
  month           = "2024-01"
  need_record_num = 1
  pay_mode        = "postPay"
}
```



```hcl
data "tencentcloud_billing_bill_detail" "example_by_time_range" {
  begin_time      = "2024-01-01 00:00:00"
  end_time        = "2024-01-31 23:59:59"
  need_record_num = 1
  resource_id     = "ins-xxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `action_type` - (Optional, String) Transaction type name, such as prepay purchase, prepay renew, postpay deduction, etc.
* `begin_time` - (Optional, String) Period start time, format yyyy-mm-dd hh:ii:ss. Must be paired with EndTime, and must be in the same month. Cross-month query is not supported, the query result is the whole month data.
* `business_code` - (Optional, String) Product name code.
* `context` - (Optional, String) Context information returned by the last request, used to speed up pagination for Month>=2023-05.
* `end_time` - (Optional, String) Period end time, format yyyy-mm-dd hh:ii:ss. Must be paired with BeginTime, and must be in the same month. Cross-month query is not supported, the query result is the whole month data.
* `month` - (Optional, String) Month, format yyyy-mm. Either Month or BeginTime&EndTime must be provided. If BeginTime&EndTime is provided, Month is ignored. Data within the last 18 months can be pulled at most.
* `need_record_num` - (Optional, Int) Whether to return the total number of records, 1 means yes, 0 means no.
* `pay_mode` - (Optional, String) Payment mode, prePay (prepaid) / postPay (pay-as-you-go).
* `payer_uin` - (Optional, String) Payer account ID (the account ID is the user's unique account identifier on Tencent Cloud). Defaults to querying the current account's bill.
* `period_type` - (Optional, String) Period type, byUsedTime (by billing cycle) / byPayTime (by deduction cycle). It needs to be consistent with the cycle type of the bill for that month in the cost center. Deprecated but kept for compatibility.
* `product_code` - (Optional, String) Sub-product code. Deprecated but kept for compatibility.
* `project_id` - (Optional, Int) Project ID.
* `resource_id` - (Optional, String) Resource ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `detail_set` - Bill detail list.
  * `action_type_name` - Transaction type.
  * `action_type` - Transaction type code.
  * `associated_order` - Associated transaction order.
    * `new_order` - Adjusted order after discount.
    * `original` - Original order before discount adjustment.
    * `prepay_modify_up` - Upgrade order.
    * `prepay_purchase` - New purchase order.
    * `prepay_renew` - Renewal order.
    * `reverse_order` - Reverse order.
  * `bill_day` - Bill day.
  * `bill_id` - Transaction ID.
  * `bill_month` - Bill month.
  * `business_code_name` - Product name.
  * `business_code` - Product code.
  * `component_set` - Component list.
    * `blended_discount` - Blended discount rate.
    * `cash_pay_amount` - Cash account payment amount.
    * `component_code_name` - Component type.
    * `component_code` - Component name code.
    * `component_config` - Component config description list.
      * `name` - Config description name.
      * `value` - Config description value.
    * `contract_price` - Component unit price after discount.
    * `cost` - Component original price.
    * `deducted_measure` - Deducted usage/duration (including resource pack).
    * `discount` - Discount rate.
    * `incentive_pay_amount` - Incentive account payment amount.
    * `instance_type` - Instance type.
    * `item_code_name` - Component name.
    * `item_code` - Component type code.
    * `original_cost_with_ri` - Reserved instance deducted component original price.
    * `original_cost_with_sp` - Savings plan deducted component original price.
    * `price_unit` - Component price unit.
    * `real_cost` - Discounted total price.
    * `real_total_measure` - Original usage/duration before resource pack deduction.
    * `reduce_type` - Discount type.
    * `ri_time_span` - Reserved instance deducted usage duration.
    * `single_price` - Component list price.
    * `sp_deduction_rate` - Savings plan deduction rate.
    * `sp_deduction` - Savings plan deduction amount. Deprecated.
    * `specified_price` - Component specified price. Deprecated.
    * `time_span` - Usage duration.
    * `time_unit_name` - Duration unit.
    * `transfer_pay_amount` - Transfer account payment amount.
    * `used_amount_unit` - Component usage unit.
    * `used_amount` - Component usage.
    * `voucher_pay_amount` - Voucher payment amount.
  * `discount_content` - Discount content.
  * `discount_object` - Discount object.
  * `discount_type` - Discount type.
  * `fee_begin_time` - Start time of use.
  * `fee_end_time` - End time of use.
  * `formula_url` - Billing rule link.
  * `formula` - Calculation description.
  * `id` - Bill record ID.
  * `operate_uin` - Operator UIN.
  * `order_id` - Order ID.
  * `owner_uin` - Owner UIN.
  * `pay_mode_name` - Billing mode.
  * `pay_time` - Deduction time.
  * `payer_uin` - Payer UIN.
  * `price_info` - Price attributes.
  * `product_code_name` - Sub-product name.
  * `product_code` - Sub-product code.
  * `project_id` - Project ID.
  * `project_name` - Project name.
  * `region_id` - Region ID.
  * `region_name` - Region.
  * `region_type_name` - Domestic/International name.
  * `region_type` - Domestic/International code.
  * `reserve_detail` - Reserve detail (instance config).
  * `resource_id` - Resource ID.
  * `resource_name` - Resource alias.
  * `tags` - Tag information list.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
  * `zone_name` - Zone.
* `total` - Total record count.


