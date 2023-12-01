---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_stream_monitor_list"
sidebar_current: "docs-tencentcloud-datasource-css_stream_monitor_list"
description: |-
  Use this data source to query detailed information of css stream_monitor_list
---

# tencentcloud_css_stream_monitor_list

Use this data source to query detailed information of css stream_monitor_list

## Example Usage

```hcl
data "tencentcloud_css_stream_monitor_list" "stream_monitor_list" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `live_stream_monitors` - The list of live stream monitoring tasks.Note: This field may return null, indicating that no valid value is available.
  * `ai_asr_input_index_list` - The list of input indices for enabling intelligent speech recognition.Note: This field may return null, indicating that no valid value is available.
  * `ai_format_diagnose` - Whether to enable format diagnosis. Note: This field may return null, indicating that no valid value is available.
  * `ai_ocr_input_index_list` - The list of input indices for enabling intelligent text recognition.Note: This field may return null, indicating that no valid value is available.
  * `allow_monitor_report` - Whether to store monitoring events in the monitoring report and allow querying of the monitoring report.Note: This field may return null, indicating that no valid value is available.
  * `asr_language` - The language for intelligent speech recognition:0: Disabled1: Chinese2: English3: Japanese4: KoreanNote: This field may return null, indicating that no valid value is available.
  * `audible_input_index_list` - The list of input indices for the output audio.Note: This field may return null, indicating that no valid value is available.
  * `check_stream_broken` - Whether to enable stream disconnection detection.Note: This field may return null, indicating that no valid value is available.
  * `check_stream_low_frame_rate` - Whether to enable low frame rate detection.Note: This field may return null, indicating that no valid value is available.
  * `create_time` - The creation time of the monitoring task, in Unix timestamp format.Note: This field may return null, indicating that no valid value is available.
  * `input_list` - The input stream information for the monitoring task.Note: This field may return null, indicating that no valid value is available.
    * `description` - Description of the monitoring task.It should be within 256 bytes.Note: This field may return null, indicating that no valid value is available.
    * `input_app` - The push path for the input stream to be monitored.It should be within 32 bytes and can only contain letters, numbers, `-`, `_`, and `.` characters.Note: This field may return null, indicating that no valid value is available.
    * `input_domain` - The push domain for the input stream to be monitored.It should be within 128 bytes and can only be filled with an enabled push domain.Note: This field may return null, indicating that no valid value is available.
    * `input_stream_name` - The name of the input stream for the monitoring task.It should be within 256 bytes and can only contain letters, numbers, `-`, `_`, and `.` characters.Note: This field may return null, indicating that no valid value is available.
    * `input_url` - The push URL for the input stream to be monitored. In most cases, this parameter is not required.Note: This field may return null, indicating that no valid value is available.
  * `monitor_id` - Monitoring task ID.Note: This field may return null, indicating that no valid value is available.
  * `monitor_name` - Monitoring task name. Up to 128 bytes.Note: This field may return null, indicating that no valid value is available.
  * `notify_policy` - The notification policy for monitoring events.Note: This field may return null, indicating that no valid value is available.
    * `callback_url` - The callback URL for notifications. It should be of length [0,512] and only support URLs with the http and https types.Note: This field may return null, indicating that no valid value is available.
    * `notify_policy_type` - The type of notification policy: Range [0,1]  0: Represents no notification policy is used.  1: Represents the use of a global callback policy, where all events are notified to the CallbackUrl.Note: This field may return null, indicating that no valid value is available.
  * `ocr_language` - The language for intelligent text recognition:0: Disabled1: Chinese and EnglishNote: This field may return null, indicating that no valid value is available.
  * `output_info` - Monitoring task output information.Note: This field may return null, indicating that no valid value is available.
    * `output_app` - The playback path for the monitoring task.It should be within 32 bytes and can only contain letters, numbers, `-`, `_`, and `.` characters.Note: This field may return null, indicating that no valid value is available.
    * `output_domain` - The playback domain for the monitoring task.It should be within 128 bytes and can only be filled with an enabled playback domain.Note: This field may return null, indicating that no valid value is available.
    * `output_stream_height` - The height of the output stream in pixels for the monitoring task. The range is [1, 1080]. It is recommended to be at least 100 pixels.Note: This field may return null, indicating that no valid value is available.
    * `output_stream_name` - The name of the output stream for the monitoring task.If not specified, the system will generate a name automatically.The name should be within 256 bytes and can only contain letters, numbers, `-`, `_`, and `.` characters.Note: This field may return null, indicating that no valid value is available.
    * `output_stream_width` - The width of the output stream in pixels for the monitoring task. The range is [1, 1920]. It is recommended to be at least 100 pixels.Note: This field may return null, indicating that no valid value is available.
  * `start_time` - The last start time of the monitoring task, in Unix timestamp format.Note: This field may return null, indicating that no valid value is available.
  * `status` - The status of the monitoring task.  0: Represents idle.  1: Represents monitoring in progress.Note: This field may return null, indicating that no valid value is available.
  * `stop_time` - The last stop time of the monitoring task, in Unix timestamp format.Note: This field may return null, indicating that no valid value is available.
  * `update_time` - The update time of the monitoring task, in Unix timestamp format.Note: This field may return null, indicating that no valid value is available.


