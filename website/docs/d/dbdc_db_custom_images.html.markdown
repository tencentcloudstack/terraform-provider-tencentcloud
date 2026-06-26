---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_db_custom_images"
sidebar_current: "docs-tencentcloud-datasource-dbdc_db_custom_images"
description: |-
  Use this data source to query available DB Custom OS images from the TencentCloud DBDC product.
---

# tencentcloud_dbdc_db_custom_images

Use this data source to query available DB Custom OS images from the TencentCloud DBDC product.

## Example Usage

### Query all dbdc db custom images

```hcl
data "tencentcloud_dbdc_db_custom_images" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `image_set` - DB Custom available OS image list.
  * `architecture` - OS architecture. Values: x86_64, arm64.
  * `image_id` - Image ID.
  * `image_type` - Image type. Values: PUBLIC_IMAGE (TencentCloud official image), PRIVATE_IMAGE (customer dedicated image).
  * `os_name` - OS name.


