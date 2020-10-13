---
subcategory: "VOD"
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
resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive"
  drm_type                        = "SimpleAES"
  disable_higher_video_bitrate    = false
  disable_higher_video_resolution = false
  comment                         = "test"

  stream_info {
    video {
      codec               = "libx265"
      fps                 = 4
      bitrate             = 129
      resolution_adaptive = false
      width               = 128
      height              = 128
      fill_type           = "stretch"
    }
    audio {
      codec         = "libmp3lame"
      bitrate       = 129
      sample_rate   = 44100
      audio_channel = "dual"
    }
    remove_audio = false
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
  }
}

resource "tencentcloud_vod_snapshot_by_time_offset_template" "foo" {
  name                = "tf-snapshot"
  width               = 130
  height              = 128
  resolution_adaptive = false
  format              = "png"
  comment             = "test"
  fill_type           = "white"
}

resource "tencentcloud_vod_image_sprite_template" "foo" {
  sample_type         = "Percent"
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

resource "tencentcloud_vod_procedure_template" "foo" {
  name    = "tf-procedure"
  comment = "test"
  media_process_task {
    adaptive_dynamic_streaming_task_list {
      definition = tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id
    }
    snapshot_by_time_offset_task_list {
      definition = tencentcloud_vod_snapshot_by_time_offset_template.foo.id
      ext_time_offset_list = [
        "3.5s"
      ]
    }
    image_sprite_task_list {
      definition = tencentcloud_vod_image_sprite_template.foo.id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Task flow name (up to 20 characters).
* `comment` - (Optional) Template description. Length limit: 256 characters.
* `media_process_task` - (Optional) Parameter of video processing task.
* `sub_app_id` - (Optional) Subapplication ID in VOD. If you need to access a resource in a subapplication, enter the subapplication ID in this field; otherwise, leave it empty.

The `adaptive_dynamic_streaming_task_list` object supports the following:

* `definition` - (Required) Adaptive bitrate streaming template ID.
* `watermark_list` - (Optional) List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.

The `animated_graphic_task_list` object supports the following:

* `definition` - (Required) Animated image generating template ID.
* `end_time_offset` - (Required) End time of animated image in video in seconds.
* `start_time_offset` - (Required) Start time of animated image in video in seconds.

The `cover_by_snapshot_task_list` object supports the following:

* `definition` - (Required) Time point screen capturing template ID.
* `position_type` - (Required) Screen capturing mode. Valid values: `Time`: screen captures by time point, `Percent`: screen captures by percentage.
* `position_value` - (Required) Screenshot position: For time point screen capturing, this means to take a screenshot at a specified time point (in seconds) and use it as the cover. For percentage screen capturing, this value means to take a screenshot at a specified percentage of the video duration and use it as the cover.
* `watermark_list` - (Optional) List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.

The `image_sprite_task_list` object supports the following:

* `definition` - (Required) Image sprite generating template ID.

The `media_process_task` object supports the following:

* `adaptive_dynamic_streaming_task_list` - (Optional) List of adaptive bitrate streaming tasks. Note: this field may return null, indicating that no valid values can be obtained.
* `animated_graphic_task_list` - (Optional) List of animated image generating tasks. Note: this field may return null, indicating that no valid values can be obtained.
* `cover_by_snapshot_task_list` - (Optional) List of cover generating tasks. Note: this field may return null, indicating that no valid values can be obtained.
* `image_sprite_task_list` - (Optional) List of image sprite generating tasks. Note: this field may return null, indicating that no valid values can be obtained.
* `sample_snapshot_task_list` - (Optional) List of sampled screen capturing tasks. Note: this field may return null, indicating that no valid values can be obtained.
* `snapshot_by_time_offset_task_list` - (Optional) List of time point screen capturing tasks. Note: this field may return null, indicating that no valid values can be obtained.
* `transcode_task_list` - (Optional) List of transcoding tasks. Note: this field may return null, indicating that no valid values can be obtained.

The `mosaic_list` object supports the following:

* `coordinate_origin` - (Optional) Origin position, which currently can only be: `TopLeft`: the origin of coordinates is in the top-left corner of the video, and the origin of the blur is in the top-left corner of the image or text. Default value: TopLeft.
* `end_time_offset` - (Optional) End time offset of blur in seconds. If this parameter is left empty or `0` is entered, the blur will exist till the last video frame; If this value is greater than `0` (e.g., n), the blur will exist till second n; If this value is smaller than `0` (e.g., -n), the blur will exist till second n before the last video frame.
* `height` - (Optional) Blur height. `%` and `px` formats are supported: If the string ends in `%`, the `height` of the blur will be the specified percentage of the video height; for example, 10% means that Height is 10% of the video height; If the string ends in `px`, the `height` of the blur will be in px; for example, 100px means that Height is 100 px. Default value: `10%`.
* `start_time_offset` - (Optional) Start time offset of blur in seconds. If this parameter is left empty or `0` is entered, the blur will appear upon the first video frame. If this parameter is left empty or `0` is entered, the blur will appear upon the first video frame; If this value is greater than `0` (e.g., n), the blur will appear at second n after the first video frame; If this value is smaller than `0` (e.g., -n), the blur will appear at second n before the last video frame.
* `width` - (Optional) Blur width. `%` and `px` formats are supported: If the string ends in `%`, the `width` of the blur will be the specified percentage of the video width; for example, 10% means that `width` is 10% of the video width; If the string ends in `px`, the `width` of the blur will be in px; for example, 100px means that Width is 100 px. Default value: `10%`.
* `x_pos` - (Optional) The horizontal position of the origin of the blur relative to the origin of coordinates of the video. `%` and `px` formats are supported: If the string ends in `%`, the XPos of the blur will be the specified percentage of the video width; for example, 10% means that XPos is 10% of the video width; If the string ends in `px`, the XPos of the blur will be the specified px; for example, 100px means that XPos is 100 px. Default value: `0px`.
* `y_pos` - (Optional) Vertical position of the origin of blur relative to the origin of coordinates of video. `%` and `px` formats are supported: If the string ends in `%`, the YPos of the blur will be the specified percentage of the video height; for example, 10% means that YPos is 10% of the video height; If the string ends in `px`, the YPos of the blur will be the specified px; for example, 100px means that YPos is 100 px. Default value: `0px`.

The `sample_snapshot_task_list` object supports the following:

* `definition` - (Required) Sampled screen capturing template ID.
* `watermark_list` - (Optional) List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.

The `snapshot_by_time_offset_task_list` object supports the following:

* `definition` - (Required) Time point screen capturing template ID.
* `ext_time_offset_list` - (Optional) The list of screenshot time points. `s` and `%` formats are supported: When a time point string ends with `s`, its unit is second. For example, `3.5s` means the 3.5th second of the video; When a time point string ends with `%`, it is marked with corresponding percentage of the video duration. For example, `10%` means that the time point is at the 10% of the video entire duration.
* `watermark_list` - (Optional) List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.

The `transcode_task_list` object supports the following:

* `definition` - (Required) Video transcoding template ID.
* `mosaic_list` - (Optional) List of blurs. Up to 10 ones can be supported.
* `watermark_list` - (Optional) List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.

The `watermark_list` object supports the following:

* `definition` - (Required) Watermarking template ID.
* `end_time_offset` - (Optional) End time offset of a watermark in seconds. If this parameter is left blank or `0` is entered, the watermark will exist till the last video frame; If this value is greater than `0` (e.g., n), the watermark will exist till second n; If this value is smaller than `0` (e.g., -n), the watermark will exist till second n before the last video frame.
* `start_time_offset` - (Optional) Start time offset of a watermark in seconds. If this parameter is left blank or `0` is entered, the watermark will appear upon the first video frame. If this parameter is left blank or `0` is entered, the watermark will appear upon the first video frame; If this value is greater than `0` (e.g., n), the watermark will appear at second n after the first video frame; If this value is smaller than `0` (e.g., -n), the watermark will appear at second n before the last video frame.
* `svg_content` - (Optional) SVG content of up to `2000000` characters. This needs to be entered only when the watermark type is `SVG`. Note: this field may return null, indicating that no valid values can be obtained.
* `text_content` - (Optional) Text content of up to `100` characters. This needs to be entered only when the watermark type is text. Note: this field may return null, indicating that no valid values can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of template in ISO date format.
* `update_time` - Last modified time of template in ISO date format.


## Import

VOD procedure template can be imported using the name, e.g.

```
$ terraform import tencentcloud_vod_procedure_template.foo tf-procedure
```

