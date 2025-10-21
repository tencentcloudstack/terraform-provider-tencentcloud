---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_session_image_version"
sidebar_current: "docs-tencentcloud-datasource-dlc_session_image_version"
description: |-
  Use this data source to query detailed information of DLC session image version
---

# tencentcloud_dlc_session_image_version

Use this data source to query detailed information of DLC session image version

## Example Usage

```hcl
data "tencentcloud_dlc_session_image_version" "example" {
  data_engine_id = "DataEngine-e482ijv6"
  framework_type = "machine-learning"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String) Data engine ID.
* `framework_type` - (Required, String) Framework type: machine learning, Python, Spark ML.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `engine_session_images` - Engine session image information.


