---
subcategory: "Video on Demand(VOD)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_procedure_template"
sidebar_current: "docs-tencentcloud-resource-vod_procedure_template"
description: |-
  Provide a resource to create a VOD procedure template.
---

# tencentcloud_vod_procedure_template

Provide a resource to create a VOD procedure template.

## Example Usage

```hcl
resource "tencentcloud_vod_sub_application" "sub_application" {
  name        = "procedure-subapplication"
  status      = "On"
  description = "this is sub application"
}

resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive"
  sub_app_id                      = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  drm_type                        = "SimpleAES"
  disable_higher_video_bitrate    = false
  disable_higher_video_resolution = false
  comment                         = "test"

  stream_info {
    video {
      codec   = "libx264"
      fps     = 3
      bitrate = 128
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 128
      sample_rate = 32000
    }
    remove_audio = true
  }
  stream_info {
    video {
      codec   = "libx264"
      fps     = 4
      bitrate = 256
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 256
      sample_rate = 44100
    }
    remove_audio = true
    tehd_config {
      type = "TEHD-100"
    }
  }
}

resource "tencentcloud_vod_snapshot_by_time_offset_template" "foo" {
  name                = "tf-snapshot"
  sub_app_id          = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  width               = 128
  height              = 128
  resolution_adaptive = false
  format              = "png"
  comment             = "test"
  fill_type           = "white"
}

resource "tencentcloud_vod_image_sprite_template" "foo" {
  sample_type         = "Percent"
  sub_app_id          = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  sample_interval     = 10
  row_count           = 3
  column_count        = 3
  name                = "tf-sprite"
  comment             = "test"
  fill_type           = "stretch"
  width               = 128
  height              = 128
  resolution_adaptive = false
}

resource "tencentcloud_vod_transcode_template" "transcode_template" {
  container    = "mp4"
  sub_app_id   = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  name         = "720pTranscodeTemplate"
  comment      = "test transcode mp4 720p update"
  remove_video = 0
  remove_audio = 0
  video_template {
    codec               = "libx264"
    fps                 = 26
    bitrate             = 1000
    resolution_adaptive = "open"
    width               = 0
    height              = 720
    fill_type           = "stretch"
    vcrf                = 1
    gop                 = 250
    preserve_hdr_switch = "OFF"
    codec_tag           = "hvc1"

  }
  audio_template {
    codec         = "libfdk_aac"
    bitrate       = 128
    sample_rate   = 44100
    audio_channel = 2

  }
  segment_type = "ts"
}

resource "tencentcloud_vod_procedure_template" "foo" {
  name       = "tf-procedure0"
  comment    = "test"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  media_process_task {
    adaptive_dynamic_streaming_task_list {
      definition = tonumber(split("#", tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id)[1])
    }
    snapshot_by_time_offset_task_list {
      definition = tonumber(split("#", tencentcloud_vod_snapshot_by_time_offset_template.foo.id)[1])
      ext_time_offset_list = [
        "3.5s"
      ]
    }
    image_sprite_task_list {
      definition = tonumber(split("#", tencentcloud_vod_image_sprite_template.foo.id)[1])
    }
    transcode_task_list {
      definition = tonumber(split("#", tencentcloud_vod_transcode_template.transcode_template.id)[1])
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Task flow name (up to 20 characters).
* `ai_analysis_task` - (Optional, List) Parameter of AI-based content analysis task.
* `ai_recognition_task` - (Optional, List) Type parameter of AI-based content recognition task.
* `comment` - (Optional, String) Template description. Length limit: 256 characters.
* `media_process_task` - (Optional, List) Parameter of video processing task.
* `review_audio_video_task` - (Optional, List) Type parameter of AI-based content recognition task.
* `sub_app_id` - (Optional, Int) The VOD [application](https://intl.cloud.tencent.com/document/product/266/14574) ID. For customers who activate VOD service from December 25, 2023, if they want to access resources in a VOD application (whether it's the default application or a newly created one), they must fill in this field with the application ID.

The `adaptive_dynamic_streaming_task_list` object of `media_process_task` supports the following:

* `definition` - (Required, String) Adaptive bitrate streaming template ID.
* `subtitle_list` - (Optional, List) Subtitle list, element is subtitle ID, support multiple subtitles, up to 16.
* `watermark_list` - (Optional, List) List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.

The `ai_analysis_task` object supports the following:

* `definition` - (Optional, String) Video content analysis template ID.

The `ai_recognition_task` object supports the following:

* `definition` - (Optional, String) Intelligent video recognition template ID.

The `animated_graphic_task_list` object of `media_process_task` supports the following:

* `definition` - (Required, String) Animated image generating template ID.
* `end_time_offset` - (Required, Float64) End time of animated image in video in seconds.
* `start_time_offset` - (Required, Float64) Start time of animated image in video in seconds.

The `copy_right_watermark` object of `transcode_task_list` supports the following:

* `text` - (Optional, String) Copyright information, maximum length is 200 characters.

The `cover_by_snapshot_task_list` object of `media_process_task` supports the following:

* `definition` - (Required, String) Time point screen capturing template ID.
* `position_type` - (Required, String) Screen capturing mode. Valid values: `Time`, `Percent`. `Time`: screen captures by time point, `Percent`: screen captures by percentage.
* `position_value` - (Required, Float64) Screenshot position: For time point screen capturing, this means to take a screenshot at a specified time point (in seconds) and use it as the cover. For percentage screen capturing, this value means to take a screenshot at a specified percentage of the video duration and use it as the cover.
* `watermark_list` - (Optional, List) List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.

The `head_tail_list` object of `transcode_task_list` supports the following:

* `definition` - (Optional, String) Video opening/closing credits configuration template ID.

The `image_sprite_task_list` object of `media_process_task` supports the following:

* `definition` - (Required, String) Image sprite generating template ID.

The `media_process_task` object supports the following:

* `adaptive_dynamic_streaming_task_list` - (Optional, List) List of adaptive bitrate streaming tasks. Note: this field may return null, indicating that no valid values can be obtained.
* `animated_graphic_task_list` - (Optional, List) List of animated image generating tasks. Note: this field may return null, indicating that no valid values can be obtained.
* `cover_by_snapshot_task_list` - (Optional, List) List of cover generating tasks. Note: this field may return null, indicating that no valid values can be obtained.
* `image_sprite_task_list` - (Optional, List) List of image sprite generating tasks. Note: this field may return null, indicating that no valid values can be obtained.
* `sample_snapshot_task_list` - (Optional, List) List of sampled screen capturing tasks. Note: this field may return null, indicating that no valid values can be obtained.
* `snapshot_by_time_offset_task_list` - (Optional, List) List of time point screen capturing tasks. Note: this field may return null, indicating that no valid values can be obtained.
* `transcode_task_list` - (Optional, List) List of transcoding tasks. Note: this field may return null, indicating that no valid values can be obtained.

The `mosaic_list` object of `transcode_task_list` supports the following:

* `coordinate_origin` - (Optional, String) Origin position, which currently can only be: `TopLeft`: the origin of coordinates is in the top-left corner of the video, and the origin of the blur is in the top-left corner of the image or text. Default value: TopLeft.
* `end_time_offset` - (Optional, Float64) End time offset of blur in seconds. If this parameter is left empty or `0` is entered, the blur will exist till the last video frame; If this value is greater than `0` (e.g., n), the blur will exist till second n; If this value is smaller than `0` (e.g., -n), the blur will exist till second n before the last video frame.
* `height` - (Optional, String) Blur height. `%` and `px` formats are supported: If the string ends in `%`, the `height` of the blur will be the specified percentage of the video height; for example, 10% means that Height is 10% of the video height; If the string ends in `px`, the `height` of the blur will be in px; for example, 100px means that Height is 100 px. Default value: `10%`.
* `start_time_offset` - (Optional, Float64) Start time offset of blur in seconds. If this parameter is left empty or `0` is entered, the blur will appear upon the first video frame. If this parameter is left empty or `0` is entered, the blur will appear upon the first video frame; If this value is greater than `0` (e.g., n), the blur will appear at second n after the first video frame; If this value is smaller than `0` (e.g., -n), the blur will appear at second n before the last video frame.
* `width` - (Optional, String) Blur width. `%` and `px` formats are supported: If the string ends in `%`, the `width` of the blur will be the specified percentage of the video width; for example, 10% means that `width` is 10% of the video width; If the string ends in `px`, the `width` of the blur will be in px; for example, 100px means that Width is 100 px. Default value: `10%`.
* `x_pos` - (Optional, String) The horizontal position of the origin of the blur relative to the origin of coordinates of the video. `%` and `px` formats are supported: If the string ends in `%`, the XPos of the blur will be the specified percentage of the video width; for example, 10% means that XPos is 10% of the video width; If the string ends in `px`, the XPos of the blur will be the specified px; for example, 100px means that XPos is 100 px. Default value: `0px`.
* `y_pos` - (Optional, String) Vertical position of the origin of blur relative to the origin of coordinates of video. `%` and `px` formats are supported: If the string ends in `%`, the YPos of the blur will be the specified percentage of the video height; for example, 10% means that YPos is 10% of the video height; If the string ends in `px`, the YPos of the blur will be the specified px; for example, 100px means that YPos is 100 px. Default value: `0px`.

The `review_audio_video_task` object supports the following:

* `definition` - (Optional, String) Review template.
* `review_contents` - (Optional, List) The type of moderated content. Valid values:
- `Media`: The original audio/video;
- `Cover`: Thumbnails.

The `sample_snapshot_task_list` object of `media_process_task` supports the following:

* `definition` - (Required, String) Sampled screen capturing template ID.
* `watermark_list` - (Optional, List) List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.

The `snapshot_by_time_offset_task_list` object of `media_process_task` supports the following:

* `definition` - (Required, String) Time point screen capturing template ID.
* `ext_time_offset_list` - (Optional, List) The list of screenshot time points. `s` and `%` formats are supported: When a time point string ends with `s`, its unit is second. For example, `3.5s` means the 3.5th second of the video; When a time point string ends with `%`, it is marked with corresponding percentage of the video duration. For example, `10%` means that the time point is at the 10% of the video entire duration.
* `time_offset_list` - (Optional, List) List of time points for screencapturing in milliseconds. Note: this field may return null, indicating that no valid values can be obtained.
* `watermark_list` - (Optional, List) List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.

The `trace_watermark` object of `transcode_task_list` supports the following:

* `switch` - (Optional, String) Whether to use digital watermarks. This parameter is required. Valid values: ON, OFF.

The `transcode_task_list` object of `media_process_task` supports the following:

* `definition` - (Required, String) Video transcoding template ID.
* `copy_right_watermark` - (Optional, List) opyright watermark.
* `end_time_offset` - (Optional, Float64) End time offset of blur in seconds. If this parameter is left empty or `0` is entered, the blur will exist till the last video frame; If this value is greater than `0` (e.g., n), the blur will exist till second n; If this value is smaller than `0` (e.g., -n), the blur will exist till second n before the last video frame.
* `head_tail_list` - (Optional, List) List of video opening/closing credits configuration template IDs. You can enter up to 10 IDs.
* `mosaic_list` - (Optional, List) List of blurs. Up to 10 ones can be supported.
* `start_time_offset` - (Optional, Float64) Start time offset of blur in seconds. If this parameter is left empty or `0` is entered, the blur will appear upon the first video frame. If this parameter is left empty or `0` is entered, the blur will appear upon the first video frame; If this value is greater than `0` (e.g., n), the blur will appear at second n after the first video frame; If this value is smaller than `0` (e.g., -n), the blur will appear at second n before the last video frame.
* `trace_watermark` - (Optional, List) Digital watermark.
* `watermark_list` - (Optional, List) List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.

The `watermark_list` object of `adaptive_dynamic_streaming_task_list` supports the following:

* `definition` - (Required, String) Watermarking template ID.
* `end_time_offset` - (Optional, Float64) End time offset of a watermark in seconds. If this parameter is left blank or `0` is entered, the watermark will exist till the last video frame; If this value is greater than `0` (e.g., n), the watermark will exist till second n; If this value is smaller than `0` (e.g., -n), the watermark will exist till second n before the last video frame.
* `start_time_offset` - (Optional, Float64) Start time offset of a watermark in seconds. If this parameter is left blank or `0` is entered, the watermark will appear upon the first video frame. If this parameter is left blank or `0` is entered, the watermark will appear upon the first video frame; If this value is greater than `0` (e.g., n), the watermark will appear at second n after the first video frame; If this value is smaller than `0` (e.g., -n), the watermark will appear at second n before the last video frame.
* `svg_content` - (Optional, String) SVG content of up to `2000000` characters. This needs to be entered only when the watermark type is `SVG`. Note: this field may return null, indicating that no valid values can be obtained.
* `text_content` - (Optional, String) Text content of up to `100` characters. This needs to be entered only when the watermark type is text. Note: this field may return null, indicating that no valid values can be obtained.

The `watermark_list` object of `cover_by_snapshot_task_list` supports the following:

* `definition` - (Required, String) Watermarking template ID.
* `end_time_offset` - (Optional, Float64) End time offset of a watermark in seconds. If this parameter is left blank or `0` is entered, the watermark will exist till the last video frame; If this value is greater than `0` (e.g., n), the watermark will exist till second n; If this value is smaller than `0` (e.g., -n), the watermark will exist till second n before the last video frame.
* `start_time_offset` - (Optional, Float64) Start time offset of a watermark in seconds. If this parameter is left blank or `0` is entered, the watermark will appear upon the first video frame. If this parameter is left blank or `0` is entered, the watermark will appear upon the first video frame; If this value is greater than `0` (e.g., n), the watermark will appear at second n after the first video frame; If this value is smaller than `0` (e.g., -n), the watermark will appear at second n before the last video frame.
* `svg_content` - (Optional, String) SVG content of up to `2000000` characters. This needs to be entered only when the watermark type is `SVG`. Note: this field may return null, indicating that no valid values can be obtained.
* `text_content` - (Optional, String) Text content of up to `100` characters. This needs to be entered only when the watermark type is text. Note: this field may return null, indicating that no valid values can be obtained.

The `watermark_list` object of `sample_snapshot_task_list` supports the following:

* `definition` - (Required, String) Watermarking template ID.
* `end_time_offset` - (Optional, Float64) End time offset of a watermark in seconds. If this parameter is left blank or `0` is entered, the watermark will exist till the last video frame; If this value is greater than `0` (e.g., n), the watermark will exist till second n; If this value is smaller than `0` (e.g., -n), the watermark will exist till second n before the last video frame.
* `start_time_offset` - (Optional, Float64) Start time offset of a watermark in seconds. If this parameter is left blank or `0` is entered, the watermark will appear upon the first video frame. If this parameter is left blank or `0` is entered, the watermark will appear upon the first video frame; If this value is greater than `0` (e.g., n), the watermark will appear at second n after the first video frame; If this value is smaller than `0` (e.g., -n), the watermark will appear at second n before the last video frame.
* `svg_content` - (Optional, String) SVG content of up to `2000000` characters. This needs to be entered only when the watermark type is `SVG`. Note: this field may return null, indicating that no valid values can be obtained.
* `text_content` - (Optional, String) Text content of up to `100` characters. This needs to be entered only when the watermark type is text. Note: this field may return null, indicating that no valid values can be obtained.

The `watermark_list` object of `snapshot_by_time_offset_task_list` supports the following:

* `definition` - (Required, String) Watermarking template ID.
* `end_time_offset` - (Optional, Float64) End time offset of a watermark in seconds. If this parameter is left blank or `0` is entered, the watermark will exist till the last video frame; If this value is greater than `0` (e.g., n), the watermark will exist till second n; If this value is smaller than `0` (e.g., -n), the watermark will exist till second n before the last video frame.
* `start_time_offset` - (Optional, Float64) Start time offset of a watermark in seconds. If this parameter is left blank or `0` is entered, the watermark will appear upon the first video frame. If this parameter is left blank or `0` is entered, the watermark will appear upon the first video frame; If this value is greater than `0` (e.g., n), the watermark will appear at second n after the first video frame; If this value is smaller than `0` (e.g., -n), the watermark will appear at second n before the last video frame.
* `svg_content` - (Optional, String) SVG content of up to `2000000` characters. This needs to be entered only when the watermark type is `SVG`. Note: this field may return null, indicating that no valid values can be obtained.
* `text_content` - (Optional, String) Text content of up to `100` characters. This needs to be entered only when the watermark type is text. Note: this field may return null, indicating that no valid values can be obtained.

The `watermark_list` object of `transcode_task_list` supports the following:

* `definition` - (Required, String) Watermarking template ID.
* `end_time_offset` - (Optional, Float64) End time offset of a watermark in seconds. If this parameter is left blank or `0` is entered, the watermark will exist till the last video frame; If this value is greater than `0` (e.g., n), the watermark will exist till second n; If this value is smaller than `0` (e.g., -n), the watermark will exist till second n before the last video frame.
* `start_time_offset` - (Optional, Float64) Start time offset of a watermark in seconds. If this parameter is left blank or `0` is entered, the watermark will appear upon the first video frame. If this parameter is left blank or `0` is entered, the watermark will appear upon the first video frame; If this value is greater than `0` (e.g., n), the watermark will appear at second n after the first video frame; If this value is smaller than `0` (e.g., -n), the watermark will appear at second n before the last video frame.
* `svg_content` - (Optional, String) SVG content of up to `2000000` characters. This needs to be entered only when the watermark type is `SVG`. Note: this field may return null, indicating that no valid values can be obtained.
* `text_content` - (Optional, String) Text content of up to `100` characters. This needs to be entered only when the watermark type is text. Note: this field may return null, indicating that no valid values can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of template in ISO date format.
* `type` - Template type, value range:
- Preset: system preset template;
- Custom: user-defined templates.
* `update_time` - Last modified time of template in ISO date format.


## Import

VOD procedure template can be imported using the name, e.g.

```
$ terraform import tencentcloud_vod_procedure_template.foo tf-procedure
```

