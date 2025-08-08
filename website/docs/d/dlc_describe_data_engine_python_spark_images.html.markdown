---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_data_engine_python_spark_images"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_data_engine_python_spark_images"
description: |-
  Use this data source to query detailed information of DLC describe data engine python spark images
---

# tencentcloud_dlc_describe_data_engine_python_spark_images

Use this data source to query detailed information of DLC describe data engine python spark images

## Example Usage

```hcl
data "tencentcloud_dlc_describe_data_engine_python_spark_images" "example" {
  child_image_version_id = "d3ftghd4-9a7e-4f64-a3f4-f38507c69742"
}
```

## Argument Reference

The following arguments are supported:

* `child_image_version_id` - (Required, String) ID of the minor version of the cluster image.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `python_spark_images` - PYSPARK image information list.
  * `child_image_version_id` - ID of the cluster image of the minor version.
  * `create_time` - Spark image create time.
  * `description` - Description of the spark image.
  * `spark_image_id` - Unique ID of the spark image.
  * `spark_image_version` - Name of the spark image.
  * `update_time` - Spark image update time.


