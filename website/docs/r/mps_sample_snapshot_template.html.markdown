---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_sample_snapshot_template"
sidebar_current: "docs-tencentcloud-resource-mps_sample_snapshot_template"
description: |-
  Provides a resource to create a mps sample_snapshot_template
---

# tencentcloud_mps_sample_snapshot_template

Provides a resource to create a mps sample_snapshot_template

## Example Usage

```hcl
resource "tencentcloud_mps_sample_snapshot_template" "sample_snapshot_template" {
  fill_type           = "stretch"
  format              = "jpg"
  height              = 128
  name                = "terraform-test-for"
  resolution_adaptive = "open"
  sample_interval     = 10
  sample_type         = "Percent"
  width               = 140
}
```

## Argument Reference

The following arguments are supported:

* `sample_interval` - (Required, Int) Sampling interval.When SampleType is Percent, specify the percentage of the sampling interval.When SampleType is Time, specify the sampling interval time in seconds.
* `sample_type` - (Required, String) Sampling snapshot type, optional value:Percent/Time.
* `comment` - (Optional, String) Template description information, length limit: 256 characters.
* `fill_type` - (Optional, String) Filling type, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling type:stretch: Stretching, stretching each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched.black: Leave black, keep the video aspect ratio unchanged, and fill the rest of the edge with black.white: Leave blank, keep the aspect ratio of the video, and fill the rest of the edge with white.gauss: Gaussian blur, keep the aspect ratio of the video unchanged, and use Gaussian blur for the rest of the edge.Default value: black.
* `format` - (Optional, String) Image format, the value can be jpg, png, webp. Default is jpg.
* `height` - (Optional, Int) The maximum value of the snapshot height (or short side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.
* `name` - (Optional, String) Sample snapshot template name, length limit: 64 characters.
* `resolution_adaptive` - (Optional, String) Adaptive resolution, optional value:open: At this time, Width represents the long side of the video, Height represents the short side of the video.close: At this point, Width represents the width of the video, and Height represents the height of the video.Default value: open.
* `width` - (Optional, Int) The maximum value of the snapshot width (or long side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps sample_snapshot_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_sample_snapshot_template.sample_snapshot_template sample_snapshot_template_id
```

