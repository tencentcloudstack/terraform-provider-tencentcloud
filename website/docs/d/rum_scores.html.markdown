---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_scores"
sidebar_current: "docs-tencentcloud-datasource-rum_scores"
description: |-
  Use this data source to query detailed information of rum scores
---

# tencentcloud_rum_scores

Use this data source to query detailed information of rum scores

## Example Usage

```hcl
data "tencentcloud_rum_scores" "scores" {
  end_time   = "2023082215"
  start_time = "2023082214"
  project_id = 1
  is_demo    = 1
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End time.
* `start_time` - (Required, String) Start time.
* `is_demo` - (Optional, Int) Get data from demo. This parameter is deprecated.
* `project_id` - (Optional, Int) Project ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `score_set` - Score list.
  * `api_duration` - The mean duration of api request.
  * `api_fail` - The number of failed api.
  * `api_num` - The number of all request api.
  * `create_time` - Project record created time.
  * `page_duration` - The duration of page load.
  * `page_error` - The number of exception which happened on page.
  * `page_pv` - Pv.
  * `page_uv` - User view.
  * `project_id` - Project ID.
  * `record_num` - The number of record.
  * `score` - The score of project.
  * `static_duration` - Duration.
  * `static_fail` - The number of failed request static resource.
  * `static_num` - The number of static resource on page.


