---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_data_backfill_instances"
sidebar_current: "docs-tencentcloud-datasource-wedata_data_backfill_instances"
description: |-
  Use this data source to query detailed information of wedata data backfill instances
---

# tencentcloud_wedata_data_backfill_instances

Use this data source to query detailed information of wedata data backfill instances

## Example Usage

```hcl
data "tencentcloud_wedata_data_backfill_instances" "wedata_data_backfill_instances" {
  project_id            = "1859317240494305280"
  data_backfill_plan_id = "deb71ea1-f708-47ab-8eb6-491ce5b9c011"
  task_id               = "20231011152006462"
}
```

## Argument Reference

The following arguments are supported:

* `data_backfill_plan_id` - (Required, String) Backfill plan Id.
* `project_id` - (Required, String) Project ID.
* `task_id` - (Required, String) Task ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - All backfill  instances under one backfill  plan.


