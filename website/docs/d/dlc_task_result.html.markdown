---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_task_result"
sidebar_current: "docs-tencentcloud-datasource-dlc_task_result"
description: |-
  Use this data source to query detailed information of DLC task result
---

# tencentcloud_dlc_task_result

Use this data source to query detailed information of DLC task result

## Example Usage

```hcl
data "tencentcloud_dlc_task_result" "example" {
  task_id = "fdd9c5fa21ca11eca6fb5254006c64af"
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String) Unique task ID.
* `is_transform_data_type` - (Optional, Bool) Whether to convert the data type.
* `max_results` - (Optional, Int) Maximum number of returned rows. Value range: 0-1,000. Default value: 1,000.
* `next_token` - (Optional, String) The pagination information returned by the last response. This parameter can be omitted for the first response, where the data will be returned from the beginning. The data with a volume set by the `MaxResults` field is returned each time.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `task_info` - The queried task information. If the returned value is empty, the task with the entered task ID does not exist. The task result will be returned only if the task status is `2` (succeeded).
Note: This field may return null, indicating that no valid values can be obtained.


