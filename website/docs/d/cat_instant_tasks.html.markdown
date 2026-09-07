---
subcategory: "Cloud Automated Testing(CAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cat_instant_tasks"
sidebar_current: "docs-tencentcloud-datasource-cat_instant_tasks"
description: |-
  Use this data source to query historical instant tasks of CAT (Cat).
---

# tencentcloud_cat_instant_tasks

Use this data source to query historical instant tasks of CAT (Cat).

## Example Usage

```hcl
data "tencentcloud_cat_instant_tasks" "example" {
}

output "tasks" {
  value = data.tencentcloud_cat_instant_tasks.example.tasks
}

output "total" {
  value = data.tencentcloud_cat_instant_tasks.example.total
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tasks` - History instant tasks list.
  * `node_count` - Node count.
  * `probe_time` - Probe time.
  * `status` - Task status.
  * `success_rate` - Success rate.
  * `target_address` - Target address.
  * `task_category` - Task category.
  * `task_id` - Task ID.
  * `task_type` - Task type.
* `total` - Total number of instant tasks.


