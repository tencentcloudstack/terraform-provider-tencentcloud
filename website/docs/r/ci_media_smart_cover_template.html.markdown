---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_media_smart_cover_template"
sidebar_current: "docs-tencentcloud-resource-ci_media_smart_cover_template"
description: |-
  Provides a resource to create a ci media_smart_cover_template
---

# tencentcloud_ci_media_smart_cover_template

Provides a resource to create a ci media_smart_cover_template

## Example Usage

```hcl
resource "tencentcloud_ci_media_smart_cover_template" "media_smart_cover_template" {
  bucket = "terraform-ci-xxxxxx"
  name   = "smart_cover_template"
  smart_cover {
    format            = "jpg"
    width             = "1280"
    height            = "960"
    count             = "10"
    delete_duplicates = "true"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) bucket name.
* `name` - (Required, String) The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.
* `smart_cover` - (Required, List) Smart Cover Parameters.

The `smart_cover` object supports the following:

* `format` - (Required, String) Image Format, value jpg, png, webp.
* `count` - (Optional, String) Number of screenshots, [1,10].
* `delete_duplicates` - (Optional, String) cover deduplication, true/false.
* `height` - (Optional, String) Height, value range: [128, 4096], unit: px, if only Height is set, Width is calculated according to the original video ratio.
* `width` - (Optional, String) Width, value range: [128, 4096], unit: px, if only Width is set, Height is calculated according to the original ratio of the video.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ci media_smart_cover_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_smart_cover_template.media_smart_cover_template terraform-ci-xxxxxx#t1ede83acc305e423799d638044d859fb7
```

