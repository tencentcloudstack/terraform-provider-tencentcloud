---
subcategory: "Cloud Automated Testing(CAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cat_probe_tasks"
sidebar_current: "docs-tencentcloud-datasource-cat_probe_tasks"
description: |-
  Use this data source to query the list of probe tasks of CAT (云拨测).
---

# tencentcloud_cat_probe_tasks

Use this data source to query the list of probe tasks of CAT (云拨测).

## Example Usage

```hcl
data "tencentcloud_cat_probe_tasks" "example" {
}

output "task_set" {
  value = data.tencentcloud_cat_probe_tasks.example.task_set
}

output "total" {
  value = data.tencentcloud_cat_probe_tasks.example.total
}
```

### Query probe tasks by name and target address

```hcl
data "tencentcloud_cat_probe_tasks" "by_name" {
  task_name      = "my-probe-task"
  target_address = "http://www.example.com"
}
```

### Query probe tasks by tag filters

```hcl
data "tencentcloud_cat_probe_tasks" "by_tags" {
  tag_filters {
    key   = "Environment"
    value = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `ascend` - (Optional, Bool) Whether to sort in ascending order.
* `order_by` - (Optional, String) Order by column.
* `order_state` - (Optional, Int) Order state.
* `pay_mode` - (Optional, Int) Pay mode.
* `result_output_file` - (Optional, String) Used to save results.
* `tag_filters` - (Optional, List) Tag filters.
* `target_address` - (Optional, String) Target address.
* `task_category` - (Optional, List: [`Int`]) Task category list.
* `task_i_ds` - (Optional, List: [`String`]) Task ID list.
* `task_name` - (Optional, String) Task name.
* `task_status` - (Optional, List: [`Int`]) Task status list.
* `task_type` - (Optional, List: [`Int`]) Task type list.

The `tag_filters` object supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `task_set` - Probe task list.
  * `created_at` - Created time.
  * `cron_state` - Scheduled task start status.
  * `cron` - Cron expression for scheduled task.
  * `interval` - Probe interval in minutes.
  * `name` - Task name.
  * `node_ip_type` - Probe node IP type.
  * `nodes` - Probe node list.
  * `order_state` - Order state.
  * `parameters` - Probe parameters.
  * `pay_mode` - Pay mode.
  * `status` - Task status.
  * `sub_sync_flag` - Whether it is a sync account.
  * `tag_info_list` - Tag info list.
    * `key` - Tag key.
    * `value` - Tag value.
  * `target_address` - Target address.
  * `task_category` - Task category.
  * `task_id` - Task ID.
  * `task_type` - Task type.
* `total` - Total number of probe tasks.


