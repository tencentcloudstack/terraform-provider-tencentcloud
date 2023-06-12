---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_orders"
sidebar_current: "docs-tencentcloud-datasource-mariadb_orders"
description: |-
  Use this data source to query detailed information of mariadb orders
---

# tencentcloud_mariadb_orders

Use this data source to query detailed information of mariadb orders

## Example Usage

```hcl
data "tencentcloud_mariadb_orders" "orders" {
  deal_name = "20230607164033835942781"
}
```

## Argument Reference

The following arguments are supported:

* `deal_name` - (Required, String) List of long order numbers to be queried, which are returned for the APIs for creating, renewing, or scaling instances.
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


