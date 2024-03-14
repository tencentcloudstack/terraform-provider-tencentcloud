---
subcategory: "Video on Demand(VOD)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_sample_snapshot_template"
sidebar_current: "docs-tencentcloud-resource-vod_sample_snapshot_template"
description: |-
  Provides a resource to create a vod snapshot template
---

# tencentcloud_vod_sample_snapshot_template

Provides a resource to create a vod snapshot template

## Example Usage

```hcl
resource "tencentcloud_vod_sub_application" "sub_application" {
  name        = "snapshotTemplateSubApplication"
  status      = "On"
  description = "this is sub application"
}

resource "tencentcloud_vod_sample_snapshot_template" "sample_snapshot_template" {
  sample_type         = "Percent"
  sample_interval     = 10
  sub_app_id          = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  name                = "testSampleSnapshot"
  width               = 500
  height              = 400
  resolution_adaptive = "open"
  format              = "jpg"
  comment             = "test sample snopshot"
  fill_type           = "black"
}
```

## Argument Reference

The following arguments are supported:

* `sample_interval` - (Required, Int) Sampling interval. If `SampleType` is `Percent`, sampling will be performed at an interval of the specified percentage. If `SampleType` is `Time`, sampling will be performed at the specified time interval in seconds.
* `sample_type` - (Required, String) Sampled screencapturing type. Valid values: Percent: by percent. Time: by time interval.
* `sub_app_id` - (Required, Int) The VOD [application](https://intl.cloud.tencent.com/document/product/266/14574) ID. For customers who activate VOD service from December 25, 2023, if they want to access resources in a VOD application (whether it's the default application or a newly created one), they must fill in this field with the application ID.
* `comment` - (Optional, String) Template description. Length limit: 256 characters.
* `fill_type` - (Optional, String) Fill type. Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported:  stretch: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot shorter or longer; black: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks. white: fill with white. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with white color blocks. gauss: fill with Gaussian blur. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with Gaussian blur.Default value: black.
* `format` - (Optional, String) Image format. Valid values: jpg, png. Default value: jpg.
* `height` - (Optional, Int) Maximum value of the height (or short side) of a screenshot in px. Value range: 0 and [128, 4,096]. If both `Width` and `Height` are 0, the resolution will be the same as that of the source video; If `Width` is 0, but `Height` is not 0, `Width` will be proportionally scaled; If `Width` is not 0, but `Height` is 0, `Height` will be proportionally scaled; If both `Width` and `Height` are not 0, the custom resolution will be used.Default value: 0.
* `name` - (Optional, String) Name of a sampled screencapturing template. Length limit: 64 characters.
* `resolution_adaptive` - (Optional, String) Resolution adaption. Valid values: open: enabled. In this case, `Width` represents the long side of a video, while `Height` the short side; close: disabled. In this case, `Width` represents the width of a video, while `Height` the height.Default value: open.
* `width` - (Optional, Int) Maximum value of the width (or long side) of a screenshot in px. Value range: 0 and [128, 4,096]. If both `Width` and `Height` are 0, the resolution will be the same as that of the source video; If `Width` is 0, but `Height` is not 0, `Width` will be proportionally scaled; If `Width` is not 0, but `Height` is 0, `Height` will be proportionally scaled; If both `Width` and `Height` are not 0, the custom resolution will be used.Default value: 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vod snapshot template can be imported using the id, e.g.

```
terraform import tencentcloud_vod_sample_snapshot_template.sample_snapshot_template $subAppId#$templateId
```

