---
subcategory: "Video on Demand(VOD)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_watermark_template"
sidebar_current: "docs-tencentcloud-resource-vod_watermark_template"
description: |-
  Provides a resource to create a vod watermark template
---

# tencentcloud_vod_watermark_template

Provides a resource to create a vod watermark template

## Example Usage

```hcl
resource "tencentcloud_vod_sub_application" "sub_application" {
  name        = "watermarkTemplateSubApplication"
  status      = "On"
  description = "this is sub application"
}

resource "tencentcloud_vod_watermark_template" "watermark_template" {
  type              = "image"
  sub_app_id        = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  name              = "myImageWatermark"
  comment           = "a png watermark"
  coordinate_origin = "TopLeft"
  x_pos             = "10%"
  y_pos             = "10%"
  image_template {
    image_content = filebase64("xxx.png")
    width         = "10%"
    height        = "10px"
  }
}
```

## Argument Reference

The following arguments are supported:

* `sub_app_id` - (Required, Int) The VOD [application](https://intl.cloud.tencent.com/document/product/266/14574) ID. For customers who activate VOD service from December 25, 2023, if they want to access resources in a VOD application (whether it's the default application or a newly created one), they must fill in this field with the application ID.
* `type` - (Required, String) Watermarking type. Valid values: image: image watermark; text: text watermark; svg: SVG watermark.
* `comment` - (Optional, String) Template description. Length limit: 256 characters.
* `coordinate_origin` - (Optional, String) Origin position. Valid values: TopLeft: the origin of coordinates is in the top-left corner of the video, and the origin of the watermark is in the top-left corner of the image or text; TopRight: the origin of coordinates is in the top-right corner of the video, and the origin of the watermark is in the top-right corner of the image or text; BottomLeft: the origin of coordinates is in the bottom-left corner of the video, and the origin of the watermark is in the bottom-left corner of the image or text; BottomRight: the origin of coordinates is in the bottom-right corner of the video, and the origin of the watermark is in the bottom-right corner of the image or text.Default value: TopLeft.
* `image_template` - (Optional, List) Image watermarking template. This field is required when `Type` is `image` and is invalid when `Type` is `text`.
* `name` - (Optional, String) Watermarking template name. Length limit: 64 characters.
* `svg_template` - (Optional, List) SVG watermarking template. This field is required when `Type` is `svg` and is invalid when `Type` is `image` or `text`.
* `text_template` - (Optional, List) Text watermarking template. This field is required when `Type` is `text` and is invalid when `Type` is `image`.
* `x_pos` - (Optional, String) The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported: If the string ends in %, the `XPos` of the watermark will be the specified percentage of the video width; for example, `10%` means that `XPos` is 10% of the video width; If the string ends in px, the `XPos` of the watermark will be the specified px; for example, `100px` means that `XPos` is 100 px.Default value: 0 px.
* `y_pos` - (Optional, String) The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported: If the string ends in %, the `YPos` of the watermark will be the specified percentage of the video height; for example, `10%` means that `YPos` is 10% of the video height; If the string ends in px, the `YPos` of the watermark will be the specified px; for example, `100px` means that `YPos` is 100 px.Default value: 0 px.

The `image_template` object supports the following:

* `image_content` - (Required, String) The [Base64](https://tools.ietf.org/html/rfc4648) encoded string of a watermark image. Only JPEG, PNG, and GIF images are supported.
* `height` - (Optional, String) Watermark height. % and px formats are supported: If the string ends in %, the `Height` of the watermark will be the specified percentage of the video height; for example, `10%` means that `Height` is 10% of the video height;  If the string ends in px, the `Height` of the watermark will be in px; for example, `100px` means that `Height` is 100 px. Valid values: 0 or [8,4096]. Default value: 0 px, which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image.
* `repeat_type` - (Optional, String) Repeat type of an animated watermark. Valid values: once: no longer appears after watermark playback ends.  repeat_last_frame: stays on the last frame after watermark playback ends.  repeat (default): repeats the playback until the video ends.
* `transparency` - (Optional, Int) Image watermark transparency: 0: completely opaque  100: completely transparent Default value: 0.
* `width` - (Optional, String) Watermark width. % and px formats are supported: If the string ends in %, the `Width` of the watermark will be the specified percentage of the video width. For example, `10%` means that `Width` is 10% of the video width;  If the string ends in px, the `Width` of the watermark will be in pixels. For example, `100px` means that `Width` is 100 pixels. Value range: [8, 4096]. Default value: 10%.

The `svg_template` object supports the following:

* `height` - (Optional, String) Watermark height, which supports six formats of px, %, W%, H%, S%, and L%: If the string ends in px, the `Height` of the watermark will be in px; for example, `100px` means that `Height` is 100 px; if `0px` is entered and `Width` is not `0px`, the watermark height will be proportionally scaled based on the source SVG image; if `0px` is entered for both `Width` and `Height`, the watermark height will be the height of the source SVG image;  If the string ends in `W%`, the `Height` of the watermark will be the specified percentage of the video width; for example, `10W%` means that `Height` is 10% of the video width;  If the string ends in `H%`, the `Height` of the watermark will be the specified percentage of the video height; for example, `10H%` means that `Height` is 10% of the video height;  If the string ends in `S%`, the `Height` of the watermark will be the specified percentage of the short side of the video; for example, `10S%` means that `Height` is 10% of the short side of the video;  If the string ends in `L%`, the `Height` of the watermark will be the specified percentage of the long side of the video; for example, `10L%` means that `Height` is 10% of the long side of the video;  If the string ends in %, the meaning is the same as `H%`. Default value: 0 px.
* `width` - (Optional, String) Watermark width, which supports six formats of px, %, W%, H%, S%, and L%: If the string ends in px, the `Width` of the watermark will be in px; for example, `100px` means that `Width` is 100 px; if `0px` is entered and `Height` is not `0px`, the watermark width will be proportionally scaled based on the source SVG image; if `0px` is entered for both `Width` and `Height`, the watermark width will be the width of the source SVG image;  If the string ends in `W%`, the `Width` of the watermark will be the specified percentage of the video width; for example, `10W%` means that `Width` is 10% of the video width;  If the string ends in `H%`, the `Width` of the watermark will be the specified percentage of the video height; for example, `10H%` means that `Width` is 10% of the video height;  If the string ends in `S%`, the `Width` of the watermark will be the specified percentage of the short side of the video; for example, `10S%` means that `Width` is 10% of the short side of the video;  If the string ends in `L%`, the `Width` of the watermark will be the specified percentage of the long side of the video; for example, `10L%` means that `Width` is 10% of the long side of the video;  If the string ends in %, the meaning is the same as `W%`. Default value: 10W%.

The `text_template` object supports the following:

* `font_alpha` - (Required, Float64) Text transparency. Value range: (0, 1] 0: completely transparent  1: completely opaque Default value: 1.
* `font_color` - (Required, String) Font color in 0xRRGGBB format. Default value: 0xFFFFFF (white).
* `font_size` - (Required, String) Font size in Npx format where N is a numeric value.
* `font_type` - (Required, String) Font type. Currently, two types are supported: simkai.ttf: both Chinese and English are supported;  arial.ttf: only English is supported.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vod watermark template can be imported using the id, e.g.

```
terraform import tencentcloud_vod_watermark_template.watermark_template $subAppId#$templateId
```

