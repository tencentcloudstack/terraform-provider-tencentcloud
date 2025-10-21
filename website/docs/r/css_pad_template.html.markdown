---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_pad_template"
sidebar_current: "docs-tencentcloud-resource-css_pad_template"
description: |-
  Provides a resource to create a css pad_template
---

# tencentcloud_css_pad_template

Provides a resource to create a css pad_template

## Example Usage

```hcl
resource "tencentcloud_css_pad_template" "pad_template" {
  description   = "pad template"
  max_duration  = 120000
  template_name = "tf-pad"
  type          = 1
  url           = "https://livewatermark-1251132611.cos.ap-guangzhou.myqcloud.com/1308919341/watermark_img_1698736540399_1441698123618_.pic.jpg"
  wait_duration = 2000
}
```

## Argument Reference

The following arguments are supported:

* `template_name` - (Required, String) Template namelimit 255 bytes.
* `url` - (Required, String) Pad content.
* `description` - (Optional, String) Description content.limit length 1024 bytes.
* `max_duration` - (Optional, Int) Max pad duration.limit: 0 - 9999999 ms.
* `type` - (Optional, Int) Pad content type.1: picture.2: video.default: 1.
* `wait_duration` - (Optional, Int) Stop stream wait time.limit: 0 - 30000 ms.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

css pad_template can be imported using the id, e.g.

```
terraform import tencentcloud_css_pad_template.pad_template templateId
```

