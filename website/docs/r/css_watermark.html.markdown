---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_watermark"
sidebar_current: "docs-tencentcloud-resource-css_watermark"
description: |-
  Provides a resource to create a css watermark
---

# tencentcloud_css_watermark

Provides a resource to create a css watermark

## Example Usage

```hcl
resource "tencentcloud_css_watermark" "watermark" {
  picture_url    = "picture_url"
  watermark_name = "watermark_name"
  x_position     = 0
  y_position     = 0
  width          = 0
  height         = 0
}
```

## Argument Reference

The following arguments are supported:

* `picture_url` - (Required, String) watermark url.
* `watermark_name` - (Required, String) watermark name.
* `height` - (Optional, Int) height of the picture.
* `width` - (Optional, Int) width of the picture.
* `x_position` - (Optional, Int) x position of the picture.
* `y_position` - (Optional, Int) y position of the picture.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - status. 0: not used, 1: used.


## Import

css watermark can be imported using the id, e.g.
```
$ terraform import tencentcloud_css_watermark.watermark watermark_id
```

