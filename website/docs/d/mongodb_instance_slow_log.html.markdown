---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_slow_log"
sidebar_current: "docs-tencentcloud-datasource-mongodb_instance_slow_log"
description: |-
  Use this data source to query detailed information of mongodb instance_slow_log
---

# tencentcloud_mongodb_instance_slow_log

Use this data source to query detailed information of mongodb instance_slow_log

## Example Usage

```hcl
data "tencentcloud_mongodb_instance_slow_log" "instance_slow_log" {
  instance_id = "cmgo-9d0p6umb"
  start_time  = "2019-06-01 10:00:00"
  end_time    = "2019-06-02 12:00:00"
  slow_m_s    = 100
  format      = "json"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) Slow log termination time, format: yyyy-mm-dd hh:mm:ss, such as: 2019-06-02 12:00:00.The time interval between the start and end of the query cannot exceed 24 hours,and only slow logs within the last 7 days are allowed to be queried.
* `instance_id` - (Required, String) Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.
* `slow_ms` - (Required, Int) Slow log execution time threshold, return slow logs whose execution time exceeds this threshold,the unit is milliseconds (ms), and the minimum is 100 milliseconds.
* `start_time` - (Required, String) Slow log start time, format: yyyy-mm-dd hh:mm:ss, such as: 2019-06-01 10:00:00. The time intervalbetween the start and end of the query cannot exceed 24 hours,and only slow logs within the last 7 days are allowed to be queried.
* `format` - (Optional, String) Slow log return format. By default, the original slow log format is returned,and versions 4.4 and above can be set to json.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `slow_logs` - details of slow logs.


