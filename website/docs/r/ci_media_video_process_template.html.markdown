---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_media_video_process_template"
sidebar_current: "docs-tencentcloud-resource-ci_media_video_process_template"
description: |-
  Provides a resource to create a ci media_video_process_template
---

# tencentcloud_ci_media_video_process_template

Provides a resource to create a ci media_video_process_template

## Example Usage

```hcl
resource "tencentcloud_ci_media_video_process_template" "media_video_process_template" {
  bucket = "terraform-ci-xxxxxx"
  name   = "video_process_template"
  color_enhance {
    enable     = "true"
    contrast   = ""
    correction = ""
    saturation = ""

  }
  ms_sharpen {
    enable        = "false"
    sharpen_level = ""

  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) bucket name.
* `name` - (Required, String) The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.
* `color_enhance` - (Optional, List) color enhancement.
* `ms_sharpen` - (Optional, List) detail enhancement, ColorEnhance and MsSharpen cannot both be empty.

The `color_enhance` object supports the following:

* `contrast` - (Optional, String) Contrast, value range: [0, 100], empty string (indicates automatic analysis).
* `correction` - (Optional, String) colorcorrection, value range: [0, 100], empty string (indicating automatic analysis).
* `enable` - (Optional, String) Whether color enhancement is turned on.
* `saturation` - (Optional, String) Saturation, value range: [0, 100], empty string (indicating automatic analysis).

The `ms_sharpen` object supports the following:

* `enable` - (Optional, String) Whether detail enhancement is enabled.
* `sharpen_level` - (Optional, String) Enhancement level, value range: [0, 10], empty string (indicates automatic analysis).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ci media_video_process_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_video_process_template.media_video_process_template terraform-ci-xxxxxx#t1d5694d87639a4593a9fd7e9025d26f52
```

