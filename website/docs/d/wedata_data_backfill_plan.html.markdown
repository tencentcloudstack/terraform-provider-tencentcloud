---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_data_backfill_plan"
sidebar_current: "docs-tencentcloud-datasource-wedata_data_backfill_plan"
description: |-
  Use this data source to query detailed information of wedata data backfill plan
---

# tencentcloud_wedata_data_backfill_plan

Use this data source to query detailed information of wedata data backfill plan

## Example Usage

```hcl
data "tencentcloud_wedata_data_backfill_plan" "wedata_data_backfill_plan" {
  project_id            = "1859317240494305280"
  data_backfill_plan_id = "deb71ea1-f708-47ab-8eb6-491ce5b9c011"
}
```

## Argument Reference

The following arguments are supported:

* `data_backfill_plan_id` - (Required, String) Backfill Plan ID.
* `project_id` - (Required, String) Project ID.
* `result_output_file` - (Optional, String) Used to save results.
* `time_zone` - (Optional, String) Display time zone, default UTC+8.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Backfill details.


