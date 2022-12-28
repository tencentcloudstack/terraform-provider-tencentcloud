---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_media_watermark_template"
sidebar_current: "docs-tencentcloud-resource-ci_media_watermark_template"
description: |-
  Provides a resource to create a ci media_watermark_template
---

# tencentcloud_ci_media_watermark_template

Provides a resource to create a ci media_watermark_template

## Example Usage

```hcl
resource "tencentcloud_ci_media_watermark_template" "media_watermark_template" {
  name = ""
  watermark {
    type       = ""
    pos        = ""
    loc_mode   = ""
    dx         = ""
    dy         = ""
    start_time = ""
    end_time   = ""
    image {
      url          = ""
      mode         = ""
      width        = ""
      height       = ""
      transparency = ""
      background   = ""
    }
    text {
      font_size    = ""
      font_type    = ""
      font_color   = ""
      transparency = ""
      text         = ""
    }

  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) bucket name.
* `name` - (Required, String) The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.
* `watermark` - (Required, List) container format.

The `image` object supports the following:

* `background` - (Required, String) Whether the background image.
* `mode` - (Required, String) Size mode, Original: original size, Proportion: proportional, Fixed: fixed size.
* `transparency` - (Required, String) Transparency, value range: [1 100], unit %.
* `url` - (Required, String) Address of watermark map (pass in after Urlencode is required).
* `height` - (Optional, String) High, 1: When the Mode is Original, it does not support setting the width of the watermark image, 2: When the Mode is Proportion, the unit is %, the value range of the background image: [100 300]; the value range of the foreground image: [1 100], relative to Video width, up to 4096px, 3: When Mode is Fixed, the unit is px, value range: [8, 4096], 4: If only Width is set, Height is calculated according to the proportion of the watermark image.
* `width` - (Optional, String) Width, 1: When the Mode is Original, it does not support setting the width of the watermark image, 2: When the Mode is Proportion, the unit is %, the value range of the background image: [100 300]; the value range of the foreground image: [1 100], relative to Video width, up to 4096px, 3: When Mode is Fixed, the unit is px, value range: [8, 4096], 4: If only Width is set, Height is calculated according to the proportion of the watermark image.

The `text` object supports the following:

* `font_color` - (Required, String) Font color, format: 0xRRGGBB.
* `font_size` - (Required, String) Font size, value range: [5 100], unit px.
* `font_type` - (Required, String) font type.
* `text` - (Required, String) Watermark content, the length does not exceed 64 characters, only supports Chinese, English, numbers, _, - and *.
* `transparency` - (Required, String) Transparency, value range: [1 100], unit %.

The `watermark` object supports the following:

* `dx` - (Required, String) Horizontal offset, 1: In the picture watermark, if Background is true, when locMode is Relativity, it is %, value range: [-300 0]; when locMode is Absolute, it is px, value range: [-4096 0] ], 2: In the picture watermark, if Background is false, when locMode is Relativity, it is %, value range: [0 100]; when locMode is Absolute, it is px, value range: [0 4096], 3: In text watermark, when locMode is Relativity, it is %, value range: [0 100]; when locMode is Absolute, it is px, value range: [0 4096], 4: When Pos is Top, Bottom and Center, the parameter is invalid.
* `dy` - (Required, String) Vertical offset, 1: In the picture watermark, if Background is true, when locMode is Relativity, it is %, value range: [-300 0]; when locMode is Absolute, it is px, value range: [-4096 0] ],2: In the picture watermark, if Background is false, when locMode is Relativity, it is %, value range: [0 100]; when locMode is Absolute, it is px, value range: [0 4096],3: In text watermark, when locMode is Relativity, it is %, value range: [0 100]; when locMode is Absolute, it is px, value range: [0 4096], 4: When Pos is Left, Right and Center, the parameter is invalid.
* `loc_mode` - (Required, String) Offset method, Relativity: proportional, Absolute: fixed position.
* `pos` - (Required, String) Reference position, TopRight, TopLeft, BottomRight, BottomLeft, Left, Right, Top, Bottom, Center.
* `type` - (Required, String) Watermark type, Text: text watermark, Image: image watermark.
* `end_time` - (Optional, String) Watermark end time, 1: [0 video duration], 2: unit is second, 3: support float format, execution accuracy is accurate to milliseconds.
* `image` - (Optional, List) Image watermark node.
* `start_time` - (Optional, String) Watermark start time, 1: [0 video duration], 2: unit is second, 3: support float format, execution accuracy is accurate to milliseconds.
* `text` - (Optional, List) Text Watermark Node.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ci media_watermark_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_watermark_template.media_watermark_template media_watermark_template_id
```

