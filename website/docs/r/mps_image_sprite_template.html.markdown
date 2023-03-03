---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_image_sprite_template"
sidebar_current: "docs-tencentcloud-resource-mps_image_sprite_template"
description: |-
  Provides a resource to create a mps image_sprite_template
---

# tencentcloud_mps_image_sprite_template

Provides a resource to create a mps image_sprite_template

## Example Usage

```hcl
resource "tencentcloud_mps_image_sprite_template" "image_sprite_template" {
  column_count        = 10
  fill_type           = "stretch"
  format              = "jpg"
  height              = 143
  name                = "terraform-test"
  resolution_adaptive = "open"
  row_count           = 10
  sample_interval     = 10
  sample_type         = "Time"
  width               = 182
}
```

## Argument Reference

The following arguments are supported:

* `column_count` - (Required, Int) The number of columns in the small image in the sprite.
* `row_count` - (Required, Int) The number of rows in the small image in the sprite.
* `sample_interval` - (Required, Int) Sampling interval.When SampleType is Percent, specify the percentage of the sampling interval.When SampleType is Time, specify the sampling interval time in seconds.
* `sample_type` - (Required, String) Sampling type, optional value:Percent/Time.
* `comment` - (Optional, String) Template description information, length limit: 256 characters.
* `fill_type` - (Optional, String) Filling type, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling type:stretch: Stretching, stretching each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched.black: Leave black, keep the video aspect ratio unchanged, and fill the rest of the edge with black.Default value: black.
* `format` - (Optional, String) Image format, the value can be jpg, png, webp. Default is jpg.
* `height` - (Optional, Int) The maximum value of the height (or short side) of the small image in the sprite image, value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.
* `name` - (Optional, String) Image sprite template name, length limit: 64 characters.
* `resolution_adaptive` - (Optional, String) Adaptive resolution, optional value:open: At this time, Width represents the long side of the video, Height represents the short side of the video.close: At this point, Width represents the width of the video, and Height represents the height of the video.Default value: open.
* `width` - (Optional, Int) The maximum value of the width (or long side) of the small image in the sprite image, value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps image_sprite_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_image_sprite_template.image_sprite_template image_sprite_template_id
```

