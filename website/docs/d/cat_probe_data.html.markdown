---
subcategory: "Cloud Automated Testing(CAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cat_probe_data"
sidebar_current: "docs-tencentcloud-datasource-cat_probe_data"
description: |-
  Use this data source to query detailed information of cat probe data
---

# tencentcloud_cat_probe_data

Use this data source to query detailed information of cat probe data

## Example Usage

```hcl
data "tencentcloud_cat_probe_data" "probe_data" {
  begin_time      = 1667923200000
  end_time        = 1667996208428
  task_type       = "AnalyzeTaskType_Network"
  sort_field      = "ProbeTime"
  ascending       = true
  selected_fields = ["terraform"]
  offset          = 0
  limit           = 20
  task_id         = ["task-knare1mk"]
}
```

## Argument Reference

The following arguments are supported:

* `ascending` - (Required, Bool) true is Ascending.
* `begin_time` - (Required, Int) Start timestamp (in milliseconds).
* `end_time` - (Required, Int) End timestamp (in milliseconds).
* `limit` - (Required, Int) Limit.
* `offset` - (Required, Int) Offset.
* `selected_fields` - (Required, Set: [`String`]) Selected Fields.
* `sort_field` - (Required, String) Fields to be sorted ProbeTime dial test time sorting can be filled in You can also fill in the selected fields in SelectedFields.
* `task_type` - (Required, String) Task Type in AnalyzeTaskType_Network,AnalyzeTaskType_Browse,AnalyzeTaskType_UploadDownload,AnalyzeTaskType_Transport,AnalyzeTaskType_MediaStream.
* `city` - (Optional, Set: [`String`]) City list.
* `code` - (Optional, Set: [`String`]) Code list.
* `districts` - (Optional, Set: [`String`]) Districts list.
* `error_types` - (Optional, Set: [`String`]) ErrorTypes list.
* `operators` - (Optional, Set: [`String`]) Operators list.
* `result_output_file` - (Optional, String) Used to save results.
* `task_id` - (Optional, Set: [`String`]) TaskID list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `detailed_single_data_define` - Probe node list.
  * `fields` - Fields.
    * `id` - ID.
    * `name` - Custom Field Name/Description.
    * `value` - Value.
  * `labels` - Labels.
    * `id` - ID.
    * `name` - Custom Field Name/Description.
    * `value` - Value.
  * `probe_time` - Probe time.


