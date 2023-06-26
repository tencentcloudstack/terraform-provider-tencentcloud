---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_orders"
sidebar_current: "docs-tencentcloud-datasource-dcdb_orders"
description: |-
  Use this data source to query detailed information of dcdb orders
---

# tencentcloud_dcdb_orders

Use this data source to query detailed information of dcdb orders

## Example Usage

```hcl
data "tencentcloud_dcdb_orders" "orders" {
  deal_names = ["2023061224903413767xxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `deal_names` - (Required, Set: [`String`]) List of long order numbers to be queried, which are returned for the APIs for creating, renewing, or scaling instances.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `deals` - Order information list.
  * `count` - Number of items.
  * `deal_name` - Order number.
  * `flow_id` - ID of the associated process, which can be used to query the process execution status.
  * `instance_ids` - The ID of the created instance, which is required only for the order that creates an instance.Note: This field may return null, indicating that no valid values can be obtained.
  * `owner_uin` - Account.
  * `pay_mode` - Payment mode. Valid values: 0 (postpaid), 1 (prepaid).


