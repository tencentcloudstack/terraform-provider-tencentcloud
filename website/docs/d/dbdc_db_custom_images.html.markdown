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

### Query dbdc db custom images by filters

```hcl
data "tencentcloud_dbdc_db_custom_images" "example" {
  filters {
    name   = "image-id"
    values = ["img-rm13akp3"]
  }

  filters {
    name   = "os-type"
    values = ["linux"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions of the DescribeDBCustomImages API. Each filter combines a name with one or more values, and multiple filters are combined with AND logic.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter name. Valid values: `image-id`, `os-type`, `image-type`, `architecture`.
* `values` - (Required, List) Filter values corresponding to the filter name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `image_set` - DB Custom available OS image list.
  * `architecture` - OS architecture. Values: x86_64, arm64.
  * `image_id` - Image ID.
  * `image_type` - Image type. Values: PUBLIC_IMAGE (TencentCloud official image), PRIVATE_IMAGE (customer dedicated image).
  * `os_name` - OS name.
  * `os_type` - OS type. Values: windows, linux.


