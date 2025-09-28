---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_data_engine_session_parameters"
sidebar_current: "docs-tencentcloud-datasource-dlc_data_engine_session_parameters"
description: |-
  Use this data source to query detailed information of DLC data engine session parameters
---

# tencentcloud_dlc_data_engine_session_parameters

Use this data source to query detailed information of DLC data engine session parameters

## Example Usage

```hcl
data "tencentcloud_dlc_data_engine_session_parameters" "example" {
  data_engine_id = "DataEngine-public-1308726196"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String) DataEngine Id.
* `data_engine_name` - (Optional, String) Engine name. When the engine name is specified, the name is used first to obtain the configuration.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data_engine_parameters` - Engine Session Configuration List.


