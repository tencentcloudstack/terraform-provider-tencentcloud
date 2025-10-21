---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_images"
sidebar_current: "docs-tencentcloud-datasource-images"
description: |-
  Use this data source to query images.
---

# tencentcloud_images

Use this data source to query images.

## Example Usage

### Query all images

```hcl
data "tencentcloud_images" "example" {}
```

### Query images by image ID

```hcl
data "tencentcloud_images" "example" {
  image_id = "img-9qrfy1xt"
}
```

### Query images by os name

```hcl
data "tencentcloud_images" "example" {
  os_name = "TencentOS Server 3.2 (Final)"
}
```

### Query images by image name regex

```hcl
data "tencentcloud_images" "example" {
  image_name_regex = "^TencentOS"
}
```

### Query images by image type

```hcl
data "tencentcloud_images" "example" {
  image_type = ["PUBLIC_IMAGE"]
}
```

### Query images by instance type

```hcl
data "tencentcloud_images" "example" {
  instance_type = "S1.SMALL1"
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Optional, String) ID of the image to be queried.
* `image_name_regex` - (Optional, String) A regex string to apply to the image list returned by TencentCloud, conflict with 'os_name'. **NOTE**: it is not wildcard, should look like `image_name_regex = "^CentOS\s+6\.8\s+64\w*"`.
* `image_type` - (Optional, List: [`String`]) A list of the image type to be queried. Valid values: 'PUBLIC_IMAGE', 'PRIVATE_IMAGE', 'SHARED_IMAGE', 'MARKET_IMAGE'.
* `instance_type` - (Optional, String) Instance type, such as `S1.SMALL1`.
* `os_name` - (Optional, String) A string to apply with fuzzy match to the os_name attribute on the image list returned by TencentCloud, conflict with 'image_name_regex'.
* `result_output_file` - (Optional, String) Used to save results.

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
  * `snapshots` - List of snapshot details.
    * `disk_size` - Size of the cloud disk used to create the snapshot; unit: GB.
    * `disk_usage` - Type of the cloud disk used to create the snapshot.
    * `snapshot_id` - Snapshot ID.
    * `snapshot_name` - Snapshot name, the user-defined snapshot alias.
  * `support_cloud_init` - Whether support cloud-init.
  * `sync_percent` - Sync percent of the image.


