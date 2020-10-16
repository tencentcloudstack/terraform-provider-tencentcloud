---
subcategory: "VOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_procedure_templates"
sidebar_current: "docs-tencentcloud-datasource-vod_procedure_templates"
description: |-
  Use this data source to query detailed information of VOD procedure templates.
---

# tencentcloud_vod_procedure_templates

Use this data source to query detailed information of VOD procedure templates.

## Example Usage

```hcl
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

data "tencentcloud_vod_procedure_templates" "foo" {
  type = "Custom"
  name = tencentcloud_vod_procedure_template.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of procedure template.
* `result_output_file` - (Optional) Used to save results.
* `sub_app_id` - (Optional) Subapplication ID in VOD. If you need to access a resource in a subapplication, enter the subapplication ID in this field; otherwise, leave it empty.
* `type` - (Optional) Template type filter. Valid values: `Preset`, `Custom`. `Preset`: preset template; `Custom`: custom template.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `template_list` - A list of adaptive dynamic streaming templates. Each element contains the following attributes:
  * `comment` - Template description.
  * `create_time` - Creation time of template in ISO date format.
  * `media_process_task` - Parameter of video processing task.
    * `adaptive_dynamic_streaming_task_list` - List of adaptive bitrate streaming tasks. Note: this field may return null, indicating that no valid values can be obtained.
      * `definition` - Adaptive bitrate streaming template ID.
      * `watermark_list` - List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.
    * `animated_graphic_task_list` - List of animated image generating tasks. Note: this field may return null, indicating that no valid values can be obtained.
      * `definition` - Animated image generating template ID.
      * `end_time_offset` - End time of animated image in video in seconds.
      * `start_time_offset` - Start time of animated image in video in seconds.
    * `cover_by_snapshot_task_list` - List of cover generating tasks. Note: this field may return null, indicating that no valid values can be obtained.
      * `definition` - Time point screen capturing template ID.
      * `position_type` - Screen capturing mode. Valid values: `Time`, `Percent`. `Time`: screen captures by time point, `Percent`: screen captures by percentage.
      * `position_value` - Screenshot position: For time point screen capturing, this means to take a screenshot at a specified time point (in seconds) and use it as the cover. For percentage screen capturing, this value means to take a screenshot at a specified percentage of the video duration and use it as the cover.
      * `watermark_list` - List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.
    * `image_sprite_task_list` - List of image sprite generating tasks. Note: this field may return null, indicating that no valid values can be obtained.
      * `definition` - Image sprite generating template ID.
    * `sample_snapshot_task_list` - List of sampled screen capturing tasks. Note: this field may return null, indicating that no valid values can be obtained.
      * `definition` - Sampled screen capturing template ID.
      * `watermark_list` - List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.
    * `snapshot_by_time_offset_task_list` - List of time point screen capturing tasks. Note: this field may return null, indicating that no valid values can be obtained.
      * `definition` - Time point screen capturing template ID.
      * `ext_time_offset_list` - The list of screenshot time points. `s` and `%` formats are supported: When a time point string ends with `s`, its unit is second. For example, `3.5s` means the 3.5th second of the video; When a time point string ends with `%`, it is marked with corresponding percentage of the video duration. For example, `10%` means that the time point is at the 10% of the video entire duration.
      * `watermark_list` - List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.
    * `transcode_task_list` - List of transcoding tasks. Note: this field may return null, indicating that no valid values can be obtained.
      * `definition` - Video transcoding template ID.
      * `mosaic_list` - List of blurs. Up to 10 ones can be supported.
        * `coordinate_origin` - Origin position, which currently can only be: `TopLeft`: the origin of coordinates is in the top-left corner of the video, and the origin of the blur is in the top-left corner of the image or text.
        * `end_time_offset` - End time offset of blur in seconds. If this parameter is left empty or `0` is entered, the blur will exist till the last video frame; If this value is greater than `0` (e.g., n), the blur will exist till second n; If this value is smaller than `0` (e.g., -n), the blur will exist till second n before the last video frame.
        * `height` - Blur height. `%` and `px` formats are supported: If the string ends in `%`, the `height` of the blur will be the specified percentage of the video height; for example, 10% means that Height is 10% of the video height; If the string ends in `px`, the `height` of the blur will be in px; for example, 100px means that Height is 100 px.
        * `start_time_offset` - Start time offset of blur in seconds. If this parameter is left empty or `0` is entered, the blur will appear upon the first video frame. If this parameter is left empty or `0` is entered, the blur will appear upon the first video frame; If this value is greater than `0` (e.g., n), the blur will appear at second n after the first video frame; If this value is smaller than `0` (e.g., -n), the blur will appear at second n before the last video frame.
        * `width` - Blur width. `%` and `px` formats are supported: If the string ends in `%`, the `width` of the blur will be the specified percentage of the video width; for example, 10% means that `width` is 10% of the video width; If the string ends in `px`, the `width` of the blur will be in px; for example, 100px means that Width is 100 px.
        * `x_pos` - The horizontal position of the origin of the blur relative to the origin of coordinates of the video. `%` and `px` formats are supported: If the string ends in `%`, the XPos of the blur will be the specified percentage of the video width; for example, 10% means that XPos is 10% of the video width; If the string ends in `px`, the XPos of the blur will be the specified px; for example, 100px means that XPos is 100 px.
        * `y_pos` - Vertical position of the origin of blur relative to the origin of coordinates of video. `%` and `px` formats are supported: If the string ends in `%`, the YPos of the blur will be the specified percentage of the video height; for example, 10% means that YPos is 10% of the video height; If the string ends in `px`, the YPos of the blur will be the specified px; for example, 100px means that YPos is 100 px.
      * `watermark_list` - List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.
  * `name` - Task flow name.
  * `type` - Template type filter. Valid values: `Preset`, `Custom`. `Preset`: preset template; `Custom`: custom template.
  * `update_time` - Last modified time of template in ISO date format.


