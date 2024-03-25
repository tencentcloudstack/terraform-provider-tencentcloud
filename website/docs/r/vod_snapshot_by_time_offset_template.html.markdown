---
subcategory: "Video on Demand(VOD)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_snapshot_by_time_offset_template"
sidebar_current: "docs-tencentcloud-resource-vod_snapshot_by_time_offset_template"
description: |-
  Provide a resource to create a VOD snapshot by time offset template.
---

# tencentcloud_vod_snapshot_by_time_offset_template

Provide a resource to create a VOD snapshot by time offset template.

## Example Usage

```hcl
resource "tencentcloud_vod_snapshot_by_time_offset_template" "foo" {
  name                = "tf-snapshot"
  width               = 130
  height              = 128
  resolution_adaptive = false
  format              = "png"
  comment             = "test"
  fill_type           = "white"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of a time point screen capturing template. Length limit: 64 characters.
* `comment` - (Optional, String) Template description. Length limit: 256 characters.
* `fill_type` - (Optional, String) Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported: `stretch`: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot `shorter` or `longer`; `black`: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks. `white`: fill with white. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with white color blocks. `gauss`: fill with Gaussian blur. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with Gaussian blur. Default value: `black`.
* `format` - (Optional, String) Image format. Valid values: `jpg`, `png`. Default value: `jpg`.
* `height` - (Optional, Int) Maximum value of the `height` (or short side) of a screenshot in px. Value range: 0 and [128, 4,096]. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, `width` will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used. Default value: `0`.
* `resolution_adaptive` - (Optional, Bool) Resolution adaption. Valid values: `true`,`false`. `true`: enabled. In this case, `width` represents the long side of a video, while `height` the short side; `false`: disabled. In this case, `width` represents the width of a video, while `height` the height. Default value: `true`.
* `sub_app_id` - (Optional, Int) The VOD [application](https://intl.cloud.tencent.com/document/product/266/14574) ID. For customers who activate VOD service from December 25, 2023, if they want to access resources in a VOD application (whether it's the default application or a newly created one), they must fill in this field with the application ID.
* `width` - (Optional, Int) Maximum value of the `width` (or long side) of a screenshot in px. Value range: 0 and [128, 4,096]. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, width will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used. Default value: `0`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of template in ISO date format.
* `type` - Template type, value range:
- Preset: system preset template;
- Custom: user-defined templates.
* `update_time` - Last modified time of template in ISO date format.


## Import

VOD snapshot by time offset template can be imported using the id($subAppId#$templateId), e.g.

```
$ terraform import tencentcloud_vod_snapshot_by_time_offset_template.foo $subAppId#$templateId
```

