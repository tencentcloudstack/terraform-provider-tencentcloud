---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_slow_log_user_host_stats"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_slow_log_user_host_stats"
description: |-
  Use this data source to query detailed information of dbbrain slow_log_user_host_stats
---

# tencentcloud_dbbrain_slow_log_user_host_stats

Use this data source to query detailed information of dbbrain slow_log_user_host_stats

## Example Usage

```hcl
data "tencentcloud_dbbrain_slow_log_user_host_stats" "test" {
  instance_id = "%s"
  start_time  = "%s"
  end_time    = "%s"
  product     = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) EndTime time of the query range, time format such as: 2019-09-10 12:13:14.
* `instance_id` - (Required, String) instance id.
* `start_time` - (Required, String) Start time of the query range, time format such as: 2019-09-10 12:13:14.
* `md5` - (Optional, String) MD5 value of SOL template.
* `product` - (Optional, String) Types of service products, supported values:`mysql` - Cloud Database MySQL; `cynosdb` - Cloud Database TDSQL-C for MySQL, defaults to `mysql`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Detailed list of the slow log proportion for each source address.
  * `count` - The number of slow logs for this source address.
  * `ratio` - The ratio of the number of slow logs of the source address to the total, in %.
  * `user_host` - source address.


