---
subcategory: "Video on Demand(VOD)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_image_sprite_templates"
sidebar_current: "docs-tencentcloud-datasource-vod_image_sprite_templates"
description: |-
  Use this data source to query detailed information of VOD image sprite templates.
---

# tencentcloud_vod_image_sprite_templates

Use this data source to query detailed information of VOD image sprite templates.

## Example Usage

```hcl
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

data "tencentcloud_vod_image_sprite_templates" "foo" {
  type       = "Custom"
  definition = tencentcloud_vod_image_sprite_template.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `definition` - (Optional) Unique ID filter of image sprite template.
* `result_output_file` - (Optional) Used to save results.
* `sub_app_id` - (Optional) Subapplication ID in VOD. If you need to access a resource in a subapplication, enter the subapplication ID in this field; otherwise, leave it empty.
* `type` - (Optional) Template type filter. Valid values: `Preset`, `Custom`. `Preset`: preset template; `Custom`: custom template.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `template_list` - A list of image sprite templates. Each element contains the following attributes:
  * `column_count` - Subimage column count of an image sprite.
  * `comment` - Template description.
  * `create_time` - Creation time of template in ISO date format.
  * `definition` - Unique ID of image sprite template.
  * `fill_type` - Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported: `stretch`: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot shorter or longer; `black`: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks.
  * `height` - Maximum value of the `height` (or short side) of a screenshot in px. Value range: 0 and [128, 4,096]. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, `width` will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used.
  * `name` - Name of a time point screen capturing template.
  * `resolution_adaptive` - Resolution adaption. Valid values: `true`: enabled. In this case, `width` represents the long side of a video, while `height` the short side; `false`: disabled. In this case, `width` represents the width of a video, while `height` the height.
  * `row_count` - Subimage row count of an image sprite.
  * `sample_interval` - Sampling interval. If `sample_type` is `Percent`, sampling will be performed at an interval of the specified percentage. If `sample_type` is `Time`, sampling will be performed at the specified time interval in seconds.
  * `sample_type` - Sampling type. Valid values: `Percent`, `Time`. `Percent`: by percent. `Time`: by time interval.
  * `type` - Template type filter. Valid values: `Preset`, `Custom`. `Preset`: preset template; `Custom`: custom template.
  * `update_time` - Last modified time of template in ISO date format.
  * `width` - Maximum value of the `width` (or long side) of a screenshot in px. Value range: 0 and [128, 4,096]. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, width will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used.


