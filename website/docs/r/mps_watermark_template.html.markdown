---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_watermark_template"
sidebar_current: "docs-tencentcloud-resource-mps_watermark_template"
description: |-
  Provides a resource to create a mps watermark_template
---

# tencentcloud_mps_watermark_template

Provides a resource to create a mps watermark_template

## Example Usage

```hcl
resource "tencentcloud_mps_watermark_template" "watermark_template" {
  coordinate_origin = "TopLeft"
  name              = "xZxasd"
  type              = "image"
  x_pos             = "12%"
  y_pos             = "21%"

  image_template {
    height        = "17px"
    image_content = filebase64("./logo.png")
    repeat_type   = "repeat"
    width         = "12px"
  }
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required, String, ForceNew) Watermark type, optional value:image, text, svg.
* `comment` - (Optional, String) Template description information, length limit: 256 characters.
* `coordinate_origin` - (Optional, String) Origin position, optional value:TopLeft: Indicates that the origin of the coordinates is at the upper left corner of the video image, and the origin of the watermark is the upper left corner of the picture or text.TopRight: Indicates that the origin of the coordinates is at the upper right corner of the video image, and the origin of the watermark is at the upper right corner of the picture or text.BottomLeft: Indicates that the origin of the coordinates is at the lower left corner of the video image, and the origin of the watermark is the lower left corner of the picture or text.BottomRight: Indicates that the origin of the coordinates is at the lower right corner of the video image, and the origin of the watermark is at the lower right corner of the picture or text.Default value: TopLeft.
* `image_template` - (Optional, List) Image watermark template, only when Type is image, this field is required and valid.
* `name` - (Optional, String) Watermark template name, length limit: 64 characters.
* `svg_template` - (Optional, List) SVG watermark template, only when Type is svg, this field is required and valid.
* `text_template` - (Optional, List) Text watermark template, only when Type is text, this field is required and valid.
* `x_pos` - (Optional, String) The horizontal position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats.When the string ends with %, it means that the watermark XPos specifies a percentage for the video width, such as 10% means that XPos is 10% of the video width.When the string ends with px, it means that the watermark XPos is the specified pixel, such as 100px means that the XPos is 100 pixels.Default value: 0px.
* `y_pos` - (Optional, String) The vertical position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats.When the string ends with %, it means that the watermark YPos specifies a percentage for the video height, such as 10% means that YPos is 10% of the video height.When the string ends with px, it means that the watermark YPos is the specified pixel, such as 100px means that the YPos is 100 pixels.Default value: 0px.

The `image_template` object supports the following:

* `image_content` - (Required, String) Watermark image[Base64](https://tools.ietf.org/html/rfc4648) encoded string. Support jpeg, png image format.
* `height` - (Optional, String) The height of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.When the string ends with px, it means that the watermark Height unit is pixel, such as 100px means that the Height is 100 pixels. The value range is 0 or [8, 4096].Default value: 0px. Indicates that Height is scaled according to the aspect ratio of the original watermark image.
* `repeat_type` - (Optional, String) Watermark repeat type. Usage scenario: The watermark is a dynamic image. Ranges:once: After the dynamic watermark is played, it will no longer appear.repeat_last_frame: After the watermark is played, stay on the last frame.repeat: the watermark loops until the end of the video (default).
* `width` - (Optional, String) The width of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Width is a percentage of the video width, such as 10% means that the Width is 10% of the video width.When the string ends with px, it means that the watermark Width unit is pixel, such as 100px means that the Width is 100 pixels. The value range is [8, 4096].Default value: 10%.

The `svg_template` object supports the following:

* `height` - (Optional, String) The height of the watermark, supports px, W%, H%, S%, L% six formats:When the string ends with px, it means that the watermark Height unit is pixels, such as 100px means that the Height is 100 pixels; when filling 0px and Width is not 0px, it means that the height of the watermark is proportionally scaled according to the original SVG image; when both Width and Height are filled When 0px, it means that the height of the watermark takes the height of the original SVG image.When the string ends with W%, it means that the watermark Height is a percentage of the video width, such as 10W% means that the Height is 10% of the video width.When the string ends with H%, it means that the watermark Height is the percentage size of the video height, such as 10H% means that the Height is 10% of the video height.When the string ends with S%, it means that the watermark Height is the percentage size of the short side of the video, such as 10S% means that the Height is 10% of the short side of the video.When the string ends with L%, it means that the watermark Height is the percentage size of the long side of the video, such as 10L% means that the Height is 10% of the long side of the video.When the string ends with %, the meaning is the same as H%.Default value: 0px.
* `width` - (Optional, String) The width of the watermark, supports px, %, W%, H%, S%, L% six formats.When the string ends with px, it means that the watermark Width unit is pixels, such as 100px means that the Width is 100 pixels; when filling 0px and the Height is not 0px, it means that the width of the watermark is proportionally scaled according to the original SVG image; when both Width and Height are filled When 0px, it means that the width of the watermark takes the width of the original SVG image.When the string ends with W%, it means that the watermark Width is a percentage of the video width, such as 10W% means that the Width is 10% of the video width.When the string ends with H%, it means that the watermark Width is a percentage of the video height, such as 10H% means that the Width is 10% of the video height.When the string ends with S%, it means that the watermark Width is the percentage size of the short side of the video, such as 10S% means that the Width is 10% of the short side of the video.When the string ends with L%, it means that the watermark Width is the percentage size of the long side of the video, such as 10L% means that the Width is 10% of the long side of the video.When the string ends with %, it has the same meaning as W%.Default value: 10W%.

The `text_template` object supports the following:

* `font_alpha` - (Required, Float64) Text transparency, value range: (0, 1].0: fully transparent.1: fully opaque.Default value: 1.
* `font_color` - (Required, String) Font color, format: 0xRRGGBB, default value: 0xFFFFFF (white).
* `font_size` - (Required, String) Font size, format: Npx, N is a number.
* `font_type` - (Required, String) Font type, currently supports two:simkai.ttf: can support Chinese and English.arial.ttf: English only.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps watermark_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_watermark_template.watermark_template watermark_template_id
```

