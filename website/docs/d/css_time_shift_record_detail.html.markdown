---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_time_shift_record_detail"
sidebar_current: "docs-tencentcloud-datasource-css_time_shift_record_detail"
description: |-
  Use this data source to query detailed information of css time_shift_record_detail
---

# tencentcloud_css_time_shift_record_detail

Use this data source to query detailed information of css time_shift_record_detail

## Example Usage

```hcl
data "tencentcloud_css_time_shift_record_detail" "time_shift_record_detail" {
  domain        = "177154.push.tlivecloud.com"
  app_name      = "qqq"
  stream_name   = "live"
  start_time    = 1698768000
  end_time      = 1698820641
  domain_group  = "tf-test"
  trans_code_id = 0
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required, String) Push path.
* `domain` - (Required, String) Push domain.
* `end_time` - (Required, Int) The ending time of the query range is specified in Unix timestamp.
* `start_time` - (Required, Int) The starting time of the query range is specified in Unix timestamp.
* `stream_name` - (Required, String) Stream name.
* `domain_group` - (Optional, String) The streaming domain belongs to a group. If there is no domain group or the domain group is an empty string, it can be left blank.
* `result_output_file` - (Optional, String) Used to save results.
* `trans_code_id` - (Optional, Int) The transcoding template ID can be left blank if it is 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `record_list` - The array of time-shift recording sessions.Note: This field may return null, indicating that no valid value was found.
  * `end_time` - The end time of the recording session is specified in Unix timestamp.
  * `sid` - The identifier for the time-shift recording session.
  * `start_time` - The start time of the recording session is specified in Unix timestamp.


