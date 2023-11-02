---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_monitor_report"
sidebar_current: "docs-tencentcloud-datasource-css_monitor_report"
description: |-
  Use this data source to query detailed information of css monitor_report
---

# tencentcloud_css_monitor_report

Use this data source to query detailed information of css monitor_report

## Example Usage

```hcl
data "tencentcloud_css_monitor_report" "monitor_report" {
  monitor_id = "0e8a12b5-df2a-4a1b-aa98-97d5610aa142"
}
```

## Argument Reference

The following arguments are supported:

* `monitor_id` - (Required, String) Monitor ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `diagnose_result` - The information about the media diagnostic result.Note: This field may return null, indicating that no valid value was found.
  * `low_frame_rate_results` - The information about low frame rate.Note: This field may return null, indicating that no valid value was found.
  * `stream_broken_results` - The information about the stream interruption.Note: This field may return null, indicating that no valid value was found.
  * `stream_format_results` - The information about the stream format diagnosis.Note: This field may return null, indicating that no valid value was found.
* `mps_result` - The information about the media processing result.Note: This field may return null, indicating that no valid value was found.
  * `ai_asr_results` - The result of intelligent speech recognition.Note: This field may return null, indicating that no valid value was found.
  * `ai_ocr_results` - The result of intelligent text recognition.Note: This field may return null, indicating that no valid value was found.


