---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_image"
sidebar_current: "docs-tencentcloud-datasource-image"
description: |-
  Provides an available image for the user.
---

# tencentcloud_image

Provides an available image for the user.

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

The following arguments are supported:

* `filter` - (Optional) One or more name/value pairs to filter.
* `image_name_regex` - (Optional) A regex string to apply to the image list returned by TencentCloud. **NOTE**: it is not wildcard, should look like `image_name_regex = "^CentOS\s+6\.8\s+64\w*"`.
* `os_name` - (Optional) A string to apply with fuzzy match to the os_name atrribute on the image list returned by TencentCloud. **NOTE**: when os_name is provided, highest priority is applied in this field instead of `image_name_regex`.

The `filter` object supports the following:

* `name` - (Required) Key of the filter, valid keys: `image-id`, `image-type`, `image-name`.
* `values` - (Required) Values of the filter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `image_id` - An image id indicate the uniqueness of a certain image,  which can be used for instance creation or resetting.
* `image_name` - Name of this image.


