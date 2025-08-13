---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_native_spark_sessions"
sidebar_current: "docs-tencentcloud-datasource-dlc_native_spark_sessions"
description: |-
  Use this data source to query detailed information of DLC native spark sessions
---

# tencentcloud_dlc_native_spark_sessions

Use this data source to query detailed information of DLC native spark sessions

## Example Usage

```hcl
data "tencentcloud_dlc_native_spark_sessions" "example" {
  data_engine_id    = "DataEngine-5plqp7q7"
  resource_group_id = "rg-j3zolzg77b"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Optional, String) Data engine id.
* `resource_group_id` - (Optional, String) Resource group id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `spark_sessions_list` - Spark sessions list.


