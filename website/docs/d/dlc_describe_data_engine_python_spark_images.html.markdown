---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_data_engine_python_spark_images"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_data_engine_python_spark_images"
description: |-
  Use this data source to query detailed information of dlc describe_data_engine_python_spark_images
---

# tencentcloud_dlc_describe_data_engine_python_spark_images

Use this data source to query detailed information of dlc describe_data_engine_python_spark_images

## Example Usage

```hcl
data "tencentcloud_dlc_describe_data_engine_python_spark_images" "describe_data_engine_python_spark_images" {
  child_image_version_id = "d3ftghd4-9a7e-4f64-a3f4-f38507c69742"
}
```

## Argument Reference

The following arguments are supported:

* `child_image_version_id` - (Required, String) Engine Image version id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `python_spark_images` - Pyspark image list.
  * `child_image_version_id` - Engine Image version id.
  * `create_time` - Spark image create time.
  * `description` - Spark image description information.
  * `spark_image_id` - Spark image unique id.
  * `spark_image_version` - Spark image name.
  * `update_time` - Spark image update time.


