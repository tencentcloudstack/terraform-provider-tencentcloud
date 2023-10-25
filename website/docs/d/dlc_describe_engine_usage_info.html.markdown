---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_engine_usage_info"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_engine_usage_info"
description: |-
  Use this data source to query detailed information of dlc describe_engine_usage_info
---

# tencentcloud_dlc_describe_engine_usage_info

Use this data source to query detailed information of dlc describe_engine_usage_info

## Example Usage

```hcl
data "tencentcloud_dlc_describe_engine_usage_info" "describe_engine_usage_info" {
  data_engine_id = "DataEngine-g5ds87d8"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String) Engine unique id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `available` - Remaining cluster specifications.
* `used` - Engine specifications occupied.


