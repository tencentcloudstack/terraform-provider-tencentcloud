---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_media_pic_process_template"
sidebar_current: "docs-tencentcloud-resource-ci_media_pic_process_template"
description: |-
  Provides a resource to create a ci media_pic_process_template
---

# tencentcloud_ci_media_pic_process_template

Provides a resource to create a ci media_pic_process_template

## Example Usage

```hcl
resource "tencentcloud_ci_media_pic_process_template" "media_pic_process_template" {
  bucket = "terraform-ci-xxxxxx"
  name   = "pic_process_template"
  pic_process {
    is_pic_info  = "true"
    process_rule = "imageMogr2/rotate/90"

  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) bucket name.
* `name` - (Required, String) The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.
* `pic_process` - (Optional, List) container format.

The `pic_process` object supports the following:

* `process_rule` - (Required, String) Image processing rules, 1: basic image processing, please refer to the basic image processing document, 2: image compression, please refer to the image compression document, 3: blind watermark, please refer to the blind watermark document.
* `is_pic_info` - (Optional, String) Whether to return the original image information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ci media_pic_process_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_pic_process_template.media_pic_process_template terraform-ci-xxxxx#t184a8a26da4674c80bf260c1e34131a65
```

