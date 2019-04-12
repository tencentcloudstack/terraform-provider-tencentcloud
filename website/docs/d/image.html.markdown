---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_image"
sidebar_current: "docs-tencentcloud-datasource-image"
description: |-
  Provides an available image for the user.
---

# tencentcloud_image

The Images data source fetch proper image, which could be one of the private images of the user and images of system resources provided by TencentCloud, as well as other public images and those available on the image market.

## Example Usage

```hcl
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"

  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}
```

## Argument Reference

 * `image_name_regex` - (Optional) A regex string to apply to the image list returned by TencentCloud. **NOTE**: it is not wildcard, should look like `image_name_regex = "^CentOS\\s+6\\.8\\s+64\\w*"`.
 * `os_name` - (Optional) A string to apply with fuzzy match to the os_name atrribute on the image list returned by TencentCloud. **NOTE**: when os_name is provided, highest priority is applied in this field instead of `image_name_regex`.
 * `filter` - (Optional) One or more name/value pairs to filter off of. There are several valid keys:  `image-id`,`image-type`,`image-name`. For a full reference, check out [DescribeImages in the TencentCloud API reference](https://intl.cloud.tencent.com/document/api/213/9451#filter).

## Attributes Reference

* `image_id` - An image id indicate the uniqueness of a certain image,  which can be used for instance creation or resetting.
* `image_name` - Name of this image.
