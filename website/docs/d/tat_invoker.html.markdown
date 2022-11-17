---
subcategory: "TencentCloud Automation Tools(TAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tat_invoker"
sidebar_current: "docs-tencentcloud-datasource-tat_invoker"
description: |-
  Use this data source to query detailed information of tat invoker
---

# tencentcloud_tat_invoker

Use this data source to query detailed information of tat invoker

## Example Usage

```hcl
data "tencentcloud_tat_invoker" "invoker" {
  # invoker_id = ""
  # command_id = ""
  # type = ""
}
```

## Argument Reference

The following arguments are supported:

* `command_id` - (Optional, String) Command ID.
* `invoker_id` - (Optional, String) Invoker ID.
* `result_output_file` - (Optional, String) Used to save results.
* `type` - (Optional, String) Invoker type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `invoker_set` - Invoker information.
  * `command_id` - Command ID.
  * `created_time` - Creation time.
  * `enable` - Whether to enable the invoker.
  * `instance_ids` - Instance ID list.
  * `invoker_id` - Invoker ID.
  * `name` - Invoker name.
  * `parameters` - Custom parameters.
  * `schedule_settings` - Execution schedule of the invoker. This field is returned for recurring invokers.
    * `invoke_time` - The next execution time of the invoker. This field is required if Policy is ONCE.
    * `policy` - Execution policy: `ONCE`: Execute once; `RECURRENCE`: Execute repeatedly.
    * `recurrence` - Trigger the crontab expression. This field is required if `Policy` is `RECURRENCE`. The crontab expression is parsed in UTC+8.
  * `type` - Invoker type.
  * `updated_time` - Modification time.
  * `username` - Username.


