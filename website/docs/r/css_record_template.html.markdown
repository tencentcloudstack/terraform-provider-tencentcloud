---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_record_template"
sidebar_current: "docs-tencentcloud-resource-css_record_template"
description: |-
  Provides a resource to create a css record_template
---

# tencentcloud_css_record_template

Provides a resource to create a css record_template

## Example Usage

```hcl
resource "tencentcloud_css_record_template" "record_template" {
  template_name = "demo"
  description   = "this is demo"
  flv_param {
    record_interval = 30
    storage_time    = 2
    enable          = 0
    vod_sub_app_id  = 123
    vod_file_name   = "demo"
    procedure       = ""
    storage_mode    = "Normal"
    class_id        = 123

  }
  hls_param {
    record_interval = 40
    storage_time    = 2
    enable          = 1
    vod_sub_app_id  = 123
    vod_file_name   = "test"
    procedure       = ""
    storage_mode    = "Cold"
    class_id        = 123

  }
  mp4_param {
    record_interval = 45
    storage_time    = 56
    enable          = 0
    vod_sub_app_id  = 234
    vod_file_name   = "test"
    procedure       = ""
    storage_mode    = "Cold"
    class_id        = 123

  }
  aac_param {
    record_interval = 5678
    storage_time    = 1234
    enable          = 1
    vod_sub_app_id  = 123
    vod_file_name   = "test"
    procedure       = ""
    storage_mode    = "Normal"
    class_id        = 123

  }
  is_delay_live = 1
  hls_special_param {
    flow_continue_duration = 1200

  }
  mp3_param {
    record_interval = 100
    storage_time    = 100
    enable          = 1
    vod_sub_app_id  = 123
    vod_file_name   = "test"
    procedure       = ""
    storage_mode    = "Normal"
    class_id        = 123

  }
  remove_watermark = true
  flv_special_param {
    upload_in_recording = true

  }
}
```

## Argument Reference

The following arguments are supported:

* `template_name` - (Required, String) Template name. Only `Chinese`, `English`, `numbers`, `_`, `-` are supported.
* `aac_param` - (Optional, List) Aac recording parameters are set when Aac recording is enabled.
* `description` - (Optional, String) Description.
* `flv_param` - (Optional, List) Flv recording parameters are set when Flv recording is enabled.
* `flv_special_param` - (Optional, List) FLV records special parameters.
* `hls_param` - (Optional, List) Hls recording parameters, which are set when hls recording is enabled.
* `hls_special_param` - (Optional, List) HLS specific recording parameters.
* `is_delay_live` - (Optional, Int) Live broadcast type, 0 by default. 0: Ordinary live broadcast, 1: Slow broadcast.
* `mp3_param` - (Optional, List) Mp3 recording parameters are set when Mp3 recording is turned on.
* `mp4_param` - (Optional, List) Mp4 recording parameters are set when Mp4 recording is enabled.
* `remove_watermark` - (Optional, Bool) Whether to remove the watermark. This parameter is invalid when the type is slow live broadcast.

The `aac_param` object supports the following:

* `class_id` - (Optional, Int) Classification of on-demand applications. This field may return null, indicating that no valid value can be obtained.
* `enable` - (Optional, Int) Whether to enable recording in the current format. The default value is 0. 0: No, 1: Yes.
* `procedure` - (Optional, String) Task flow. Note: This field may return null, indicating that no valid value can be obtained.
* `record_interval` - (Optional, Int) Recording interval. Unit: second, default: 1800. Value range: 30-7200. This parameter is invalid for HLS. When recording HLS, a file is generated from streaming to streaming.
* `storage_mode` - (Optional, String) Video storage strategy. Normal: standard storage. Cold: low frequency storage. Note: This field may return null, indicating that no valid value can be obtained.
* `storage_time` - (Optional, Int) Recording storage duration. Unit: second. Value range: 0 - 1500 days. 0: indicates permanent storage.
* `vod_file_name` - (Optional, String) Record file name.Special placeholders supported are: `StreamID`: Stream ID,`StartYear`: Start time - year,`StartMonth`: Start time - month,`StartDay`: Start time - day,`StartHour`: Start time - hour,`StartMinute`: Start time - minutes,`StartSecond`: Start time - seconds,`StartMillisecond`: Start time - milliseconds,`EndYear`: End time - year,`EndMonth`: End time - month,`EndDay`: End time - day,`EndHour`: End time - hour,`EndMinute`: End time - minutes,`EndSecond`: End time - seconds,`EndMillisecond`: End time - millisecondsIf not set, the default recording file name is: `StreamID`_`StartYear`-`StartMonth`-`StartDay`-`StartHour`-`StartMinute`-`StartSecond`_`EndYear`-`EndMonth`-`EndDay`-`EndHour`-`EndMinute`-`EndSecond`.
* `vod_sub_app_id` - (Optional, Int) The ID of the on-demand sub app.

The `flv_param` object supports the following:

* `class_id` - (Optional, Int) Classification of on-demand applications. Note: This field may return null, indicating that no valid value can be obtained.
* `enable` - (Optional, Int) Whether to enable recording in the current format. The default value is 0. 0: No, 1: Yes.
* `procedure` - (Optional, String) Task flow. This field may return null, indicating that no valid value can be obtained.
* `record_interval` - (Optional, Int) Recording interval.  Unit: second, default: 1800.  Value range: 30-7200.  This parameter is invalid for HLS. When recording HLS, a file is generated from streaming to streaming.
* `storage_mode` - (Optional, String) Video storage strategy. Normal: standard storage. Cold: low frequency storage. This field may return null, indicating that no valid value can be obtained.
* `storage_time` - (Optional, Int) Recording storage duration.  Unit: second. Value range: 0 - 1500 days.  0: indicates permanent storage.
* `vod_file_name` - (Optional, String) Record file name.Special placeholders supported are: `StreamID`: Stream ID,`StartYear`: Start time - year,`StartMonth`: Start time - month,`StartDay`: Start time - day,`StartHour`: Start time - hour,`StartMinute`: Start time - minutes,`StartSecond`: Start time - seconds,`StartMillisecond`: Start time - milliseconds,`EndYear`: End time - year,`EndMonth`: End time - month,`EndDay`: End time - day,`EndHour`: End time - hour,`EndMinute`: End time - minutes,`EndSecond`: End time - seconds,`EndMillisecond`: End time - millisecondsIf not set, the default recording file name is: `StreamID`_`StartYear`-`StartMonth`-`StartDay`-`StartHour`-`StartMinute`-`StartSecond`_`EndYear`-`EndMonth`-`EndDay`-`EndHour`-`EndMinute`-`EndSecond`.
* `vod_sub_app_id` - (Optional, Int) The ID of the vodSub app.

The `flv_special_param` object supports the following:

* `upload_in_recording` - (Optional, Bool) Whether to enable the transfer while recording is valid only in the flv format.

The `hls_param` object supports the following:

* `class_id` - (Optional, Int) Classification of on-demand applications. Note: This field may return null, indicating that no valid value can be obtained.
* `enable` - (Optional, Int) Whether to enable recording in the current format. The default value is 0. 0: No, 1: Yes.
* `procedure` - (Optional, String) Task flow. Note: This field may return null, indicating that no valid value can be obtained.
* `record_interval` - (Optional, Int) Recording interval. Unit: second, default: 1800. Value range: 30-7200. This parameter is invalid for HLS. When recording HLS, a file is generated from streaming to streaming.
* `storage_mode` - (Optional, String) Video storage strategy. Normal: standard storage. Cold: low frequency storage. Note: This field may return null, indicating that no valid value can be obtained.
* `storage_time` - (Optional, Int) Recording storage duration. Unit: second. Value range: 0 - 1500 days. 0: indicates permanent storage.
* `vod_file_name` - (Optional, String) Record file name.Special placeholders supported are: `StreamID`: Stream ID,`StartYear`: Start time - year,`StartMonth`: Start time - month,`StartDay`: Start time - day,`StartHour`: Start time - hour,`StartMinute`: Start time - minutes,`StartSecond`: Start time - seconds,`StartMillisecond`: Start time - milliseconds,`EndYear`: End time - year,`EndMonth`: End time - month,`EndDay`: End time - day,`EndHour`: End time - hour,`EndMinute`: End time - minutes,`EndSecond`: End time - seconds,`EndMillisecond`: End time - millisecondsIf not set, the default recording file name is: `StreamID`_`StartYear`-`StartMonth`-`StartDay`-`StartHour`-`StartMinute`-`StartSecond`_`EndYear`-`EndMonth`-`EndDay`-`EndHour`-`EndMinute`-`EndSecond`.
* `vod_sub_app_id` - (Optional, Int) The ID of the vodSub app.

The `hls_special_param` object supports the following:

* `flow_continue_duration` - (Optional, Int) HLS freewheeling timeout. Value range [0, 1800].

The `mp3_param` object supports the following:

* `class_id` - (Optional, Int) Classification of vod applications. Note: This field may return null, indicating that no valid value can be obtained.
* `enable` - (Optional, Int) Whether to enable recording in the current format. The default value is 0. 0: No, 1: Yes.
* `procedure` - (Optional, String) Task flow. Note: This field may return null, indicating that no valid value can be obtained.
* `record_interval` - (Optional, Int) Recording interval. Unit: second, default: 1800. Value range: 30-7200. This parameter is invalid for HLS. When recording HLS, a file is generated from streaming to streaming.
* `storage_mode` - (Optional, String) Video storage strategy. Normal: standard storage. Cold: low frequency storage. Note: This field may return null, indicating that no valid value can be obtained.
* `storage_time` - (Optional, Int) Recording storage duration. Unit: second. Value range: 0 - 1500 days. 0: indicates permanent storage.
* `vod_file_name` - (Optional, String) Record file name.Special placeholders supported are: `StreamID`: Stream ID,`StartYear`: Start time - year,`StartMonth`: Start time - month,`StartDay`: Start time - day,`StartHour`: Start time - hour,`StartMinute`: Start time - minutes,`StartSecond`: Start time - seconds,`StartMillisecond`: Start time - milliseconds,`EndYear`: End time - year,`EndMonth`: End time - month,`EndDay`: End time - day,`EndHour`: End time - hour,`EndMinute`: End time - minutes,`EndSecond`: End time - seconds,`EndMillisecond`: End time - millisecondsIf not set, the default recording file name is: `StreamID`_`StartYear`-`StartMonth`-`StartDay`-`StartHour`-`StartMinute`-`StartSecond`_`EndYear`-`EndMonth`-`EndDay`-`EndHour`-`EndMinute`-`EndSecond`.
* `vod_sub_app_id` - (Optional, Int) The ID of the on-demand sub app.

The `mp4_param` object supports the following:

* `class_id` - (Optional, Int) Classification of on-demand applications. Note: This field may return null, indicating that no valid value can be obtained.
* `enable` - (Optional, Int) Whether to enable recording in the current format. The default value is 0. 0: No, 1: Yes.
* `procedure` - (Optional, String) Task flow. Note: This field may return null, indicating that no valid value can be obtained.
* `record_interval` - (Optional, Int) Recording interval. Unit: second, default: 1800. Value range: 30-7200. This parameter is invalid for HLS. When recording HLS, a file is generated from streaming to streaming.
* `storage_mode` - (Optional, String) Video storage strategy. Normal: standard storage. Cold: low frequency storage. Note: This field may return null, indicating that no valid value can be obtained.
* `storage_time` - (Optional, Int) Recording storage duration. Unit: second. Value range: 0 - 1500 days. 0: indicates permanent storage.
* `vod_file_name` - (Optional, String) Record file name.Special placeholders supported are: `StreamID`: Stream ID,`StartYear`: Start time - year,`StartMonth`: Start time - month,`StartDay`: Start time - day,`StartHour`: Start time - hour,`StartMinute`: Start time - minutes,`StartSecond`: Start time - seconds,`StartMillisecond`: Start time - milliseconds,`EndYear`: End time - year,`EndMonth`: End time - month,`EndDay`: End time - day,`EndHour`: End time - hour,`EndMinute`: End time - minutes,`EndSecond`: End time - seconds,`EndMillisecond`: End time - millisecondsIf not set, the default recording file name is: `StreamID`_`StartYear`-`StartMonth`-`StartDay`-`StartHour`-`StartMinute`-`StartSecond`_`EndYear`-`EndMonth`-`EndDay`-`EndHour`-`EndMinute`-`EndSecond`.
* `vod_sub_app_id` - (Optional, Int) The ID of the on-demand sub app.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateLiveRecordTemplate](https://www.tencentcloud.com/document/api/267/30845)
* Read: [DescribeLiveRecordTemplate](https://www.tencentcloud.com/document/api/267/30840)
* Update: [ModifyLiveRecordTemplate](https://www.tencentcloud.com/document/api/267/30838)
* Delete: [DeleteLiveRecordTemplate](https://www.tencentcloud.com/document/api/267/30842)

## Import

css record_template can be imported using the id, e.g.

```
terraform import tencentcloud_css_record_template.record_template templateId
```

