---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_bucket_pic_style"
sidebar_current: "docs-tencentcloud-resource-ci_bucket_pic_style"
description: |-
  Provides a resource to create a ci bucket_pic_style
---

# tencentcloud_ci_bucket_pic_style

Provides a resource to create a ci bucket_pic_style

## Example Usage

```hcl
resource "tencentcloud_ci_bucket_pic_style" "bucket_pic_style" {
  bucket     = "terraform-ci-xxxxxx"
  style_name = "rayscale_2"
  style_body = "imageMogr2/thumbnail/20x/crop/20x20/gravity/center/interlace/0/quality/100"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) bucket name.
* `style_body` - (Required, String, ForceNew) style details, example: mageMogr2/grayscale/1.
* `style_name` - (Required, String, ForceNew) style name, style names are case-sensitive, and a combination of uppercase and lowercase letters, numbers, and `$ + _ ( )` is supported.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ci bucket_pic_style can be imported using the bucket#styleName, e.g.

```
terraform import tencentcloud_ci_bucket_pic_style.bucket_pic_style terraform-ci-xxxxxx#rayscale_2
```

