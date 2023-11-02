---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_time_shift_stream_list"
sidebar_current: "docs-tencentcloud-datasource-css_time_shift_stream_list"
description: |-
  Use this data source to query detailed information of css time_shift_stream_list
---

# tencentcloud_css_time_shift_stream_list

Use this data source to query detailed information of css time_shift_stream_list

## Example Usage

```hcl
data "tencentcloud_css_time_shift_stream_list" "time_shift_stream_list" {
  start_time   = 1698768000
  end_time     = 1698820641
  stream_name  = "live"
  domain       = "177154.push.tlivecloud.com"
  domain_group = "tf-test"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, Int) The end time, which must be a Unix timestamp.
* `start_time` - (Required, Int) The start time, which must be a Unix timestamp.
* `domain_group` - (Optional, String) The group the push domain belongs to.
* `domain` - (Optional, String) The push domain.
* `result_output_file` - (Optional, String) Used to save results.
* `stream_name` - (Optional, String) The stream name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `stream_list` - The information of the streams.Note: This field may return null, indicating that no valid values can be obtained.
  * `app_name` - The push path.
  * `domain_group` - The group the push domain belongs to.Note: This field may return null, indicating that no valid values can be obtained.
  * `domain` - The push domain.
  * `duration` - The storage duration (seconds) of the recording.Note: This field may return null, indicating that no valid values can be obtained.
  * `end_time` - The stream end time (for streams that ended before the time of query), which is a Unix timestamp.
  * `start_time` - The stream start time, which is a Unix timestamp.
  * `stream_name` - The stream name.
  * `stream_type` - The stream type. `0`: The original stream; `1`: The watermarked stream; `2`: The transcoded stream.
  * `trans_code_id` - The transcoding template ID.Note: This field may return null, indicating that no valid values can be obtained.
* `total_size` - The total number of records in the specified time period.


