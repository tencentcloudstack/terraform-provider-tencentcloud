---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_engine_usage_info"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_engine_usage_info"
description: |-
  Use this data source to query detailed information of DLC describe engine usage info
---

# tencentcloud_dlc_describe_engine_usage_info

Use this data source to query detailed information of DLC describe engine usage info

## Example Usage

```hcl
data "tencentcloud_dlc_describe_engine_usage_info" "example" {
  data_engine_id = "DataEngine-80ibn1cj"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String) The data engine ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `available` - The available cluster spec.
* `used` - The used cluster spec.


