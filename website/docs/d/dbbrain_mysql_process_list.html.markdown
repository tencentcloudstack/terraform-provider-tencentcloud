---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_mysql_process_list"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_mysql_process_list"
description: |-
  Use this data source to query detailed information of dbbrain mysql_process_list
---

# tencentcloud_dbbrain_mysql_process_list

Use this data source to query detailed information of dbbrain mysql_process_list

## Example Usage

```hcl
data "tencentcloud_dbbrain_mysql_process_list" "mysql_process_list" {
  instance_id = local.mysql_id
  product     = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `command` - (Optional, String) The execution type of the thread, used to filter the thread list.
* `db` - (Optional, String) The threads operations database, used to filter the thread list.
* `host` - (Optional, String) The operating host address of the thread, used to filter the thread list.
* `id` - (Optional, Int) thread ID, used to filter the thread list.
* `info` - (Optional, String) The threads operation statement is used to filter the thread list.
* `product` - (Optional, String) Service product type, supported values: `mysql` - cloud database MySQL; `cynosdb` - cloud database TDSQL-C for MySQL, the default is `mysql`.
* `result_output_file` - (Optional, String) Used to save results.
* `state` - (Optional, String) The operational state of the thread, used to filter the thread list.
* `time` - (Optional, Int) The minimum value of the operation duration of a thread, in seconds, used to filter the list of threads whose operation duration is longer than this value.
* `user` - (Optional, String) The operating account name of the thread, used to filter the thread list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `process_list` - Live thread list.
  * `command` - The execution type of the thread.
  * `db` - The thread that operates the database.
  * `host` - The operating host address of the thread.
  * `id` - thread ID.
  * `info` - The operation statement for the thread.
  * `state` - The operational state of the thread.
  * `time` - The operation duration of the thread, in seconds.
  * `user` - The operating account name of the thread.


