---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_db_space_status"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_db_space_status"
description: |-
  Use this data source to query detailed information of dbbrain db_space_status
---

# tencentcloud_dbbrain_db_space_status

Use this data source to query detailed information of dbbrain db_space_status

## Example Usage

```hcl
data "tencentcloud_dbbrain_db_space_status" "db_space_status" {
  instance_id = "%s"
  range_days  = 7
  product     = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `product` - (Optional, String) Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.
* `range_days` - (Optional, Int) The number of days in the time period, the deadline is the current day, and the default is 7 days.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `available_days` - Estimated number of days available.
* `growth` - Disk growth (MB).
* `remain` - Disk remaining (MB).
* `total` - Total disk size (MB).


