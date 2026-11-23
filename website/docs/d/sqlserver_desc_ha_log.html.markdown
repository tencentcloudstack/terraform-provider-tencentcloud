---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_desc_ha_log"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_desc_ha_log"
description: |-
  Use this data source to query detailed information of sqlserver desc_ha_log
---

# tencentcloud_sqlserver_desc_ha_log

Use this data source to query detailed information of sqlserver desc_ha_log

## Example Usage

```hcl
data "tencentcloud_sqlserver_desc_ha_log" "desc_ha_log" {
  instance_id = "mssql-jdk2pwld"
  start_time  = "2023-12-01 00:00:00"
  end_time    = "2023-12-15 00:00:00"
  switch_type = 1
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End time (yyyy-MM-dd HH:mm:ss).
* `instance_id` - (Required, String) Instance ID.
* `start_time` - (Required, String) Start time (yyyy-MM-dd HH:mm:ss).
* `result_output_file` - (Optional, String) Used to save results.
* `switch_type` - (Optional, Int) Switching mode 0-system automatically switches, 1-manual switch, if not filled in, all will be checked by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `switch_log` - Master/Slave switching log.
  * `end_time` - Switch end time Note: This field may return null, indicating that no valid value can be obtained.
  * `event_id` - Switch event ID Note: This field may return null, indicating that no valid value can be obtained.
  * `reason` - Machine failure causes automatic switching Note: This field may return null, indicating that no valid value can be obtained.
  * `start_time` - Switch start time Note: This field may return null, indicating that no valid value can be obtained.
  * `switch_type` - Switching mode 0-system automatic switching, 1-manual switching Note: This field may return null, indicating that no valid value can be obtained.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeHASwitchLog](https://www.tencentcloud.com/document/api/238/72752)

