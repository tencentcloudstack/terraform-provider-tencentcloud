---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_images"
sidebar_current: "docs-tencentcloud-datasource-images"
description: |-
  Use this data source to query images.
---

# tencentcloud_images

Use this data source to query images.

## Example Usage

```hcl
data "tencentcloud_images" "foo" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos 7.5"
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Optional) ID of the image to be queried.
* `image_name_regex` - (Optional) A regex string to apply to the image list returned by TencentCloud, conflict with 'os_name'. **NOTE**: it is not wildcard, should look like `image_name_regex = "^CentOS\s+6\.8\s+64\w*"`.
* `image_type` - (Optional) A list of the image type to be queried. Available values include: 'PUBLIC_IMAGE', 'PRIVATE_IMAGE', 'SHARED_IMAGE', 'MARKET_IMAGE'.
* `os_name` - (Optional) A string to apply with fuzzy match to the os_name atrribute on the image list returned by TencentCloud, conflict with 'image_name_regex'.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `images` - An information list of image. Each element contains the following attributes:
  * `architecture` - Architecture of the image.
  * `created_time` - Created time of the image.
  * `image_creator` - Image creator of the image.
  * `image_description` - Description of the image.
  * `image_id` - ID of the image.
  * `image_name` - Name of the image.
  * `image_size` - Size of the image.
  * `image_source` - Image source of the image.
  * `image_state` - State of the image.
  * `image_type` - Type of the image.
  * `os_name` - OS name of the image.
  * `platform` - Platform of the image.
  * `support_cloud_init` - Whether support cloud-init.
  * `sync_percent` - Sync percent of the image.


