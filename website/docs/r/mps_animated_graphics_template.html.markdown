---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_animated_graphics_template"
sidebar_current: "docs-tencentcloud-resource-mps_animated_graphics_template"
description: |-
  Provides a resource to create a mps animated_graphics_template
---

# tencentcloud_mps_animated_graphics_template

Provides a resource to create a mps animated_graphics_template

## Example Usage

```hcl
resource "tencentcloud_mps_animated_graphics_template" "animated_graphics_template" {
  format              = "gif"
  fps                 = 20
  height              = 130
  name                = "terraform-test"
  quality             = 75
  resolution_adaptive = "open"
  width               = 140
}
```

## Argument Reference

The following arguments are supported:

* `fps` - (Required, Int) Frame rate, value range: [1, 30], unit: Hz.
* `comment` - (Optional, String) Template description information, length limit: 256 characters.
* `format` - (Optional, String) Animation format, the values are gif and webp. Default is gif.
* `height` - (Optional, Int) The maximum value of the animation height (or short side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.
* `name` - (Optional, String) Rotation diagram template name, length limit: 64 characters.
* `quality` - (Optional, Float64) Image quality, value range: [1, 100], default value is 75.
* `resolution_adaptive` - (Optional, String) Adaptive resolution, optional value:open: At this time, Width represents the long side of the video, Height represents the short side of the video.close: At this point, Width represents the width of the video, and Height represents the height of the video.Default value: open.
* `width` - (Optional, Int) The maximum value of the animation width (or long side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps animated_graphics_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_animated_graphics_template.animated_graphics_template animated_graphics_template_id
```

