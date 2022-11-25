---
subcategory: "DBbrain"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_sql_filters"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_sql_filters"
description: |-
  Use this data source to query detailed information of dbbrain sqlFilters
---

# tencentcloud_dbbrain_sql_filters

Use this data source to query detailed information of dbbrain sqlFilters

## Example Usage

```hcl
resource "tencentcloud_dbbrain_sql_filter" "sql_filter" {
  instance_id = "mysql_ins_id"
  session_token {
    user     = "user"
    password = "password"
  }
  sql_type        = "SELECT"
  filter_key      = "test"
  max_concurrency = 10
  duration        = 3600
}

data "tencentcloud_dbbrain_sql_filters" "sql_filters" {
  instance_id = "mysql_ins_id"
  filter_ids  = [tencentcloud_dbbrain_sql_filter.sql_filter.filter_id]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `filter_ids` - (Optional, Set: [`Int`]) filter id list.
* `result_output_file` - (Optional, String) Used to save results.
* `statuses` - (Optional, Set: [`String`]) status list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - sql filter list.
  * `create_time` - create time.
  * `current_concurrency` - current concurrency.
  * `current_time` - current time.
  * `expire_time` - expire time.
  * `id` - task id.
  * `max_concurrency` - maxmum concurrency.
  * `origin_keys` - origin keys.
  * `origin_rule` - origin rule.
  * `rejected_sql_count` - rejected sql count.
  * `sql_type` - sql type, optional value is SELECT, UPDATE, DELETE, INSERT, REPLACE.
  * `status` - task status, optional value is RUNNING, FINISHED, TERMINATED.


