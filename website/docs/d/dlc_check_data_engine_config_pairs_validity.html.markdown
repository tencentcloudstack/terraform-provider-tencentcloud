---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_check_data_engine_config_pairs_validity"
sidebar_current: "docs-tencentcloud-datasource-dlc_check_data_engine_config_pairs_validity"
description: |-
  Use this data source to query detailed information of dlc check_data_engine_config_pairs_validity
---

# tencentcloud_dlc_check_data_engine_config_pairs_validity

Use this data source to query detailed information of dlc check_data_engine_config_pairs_validity

## Example Usage

```hcl
data "tencentcloud_dlc_check_data_engine_config_pairs_validity" "check_data_engine_config_pairs_validity" {
  child_image_version_id = "d3ftghd4-9a7e-4f64-a3f4-f38507c69742"
}
```

## Argument Reference

The following arguments are supported:

* `child_image_version_id` - (Optional, String) Engine Image version id.
* `data_engine_config_pairs` - (Optional, List) User-defined parameters.
* `image_version_id` - (Optional, String) Engine major version id. If a minor version id exists, you only need to pass in the minor version id. If it does not exist, the latest minor version id under the current major version will be obtained.
* `result_output_file` - (Optional, String) Used to save results.

The `data_engine_config_pairs` object supports the following:

* `config_item` - (Required, String) Configuration item.
* `config_value` - (Required, String) Configuration value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `is_available` - Parameter validity: true: valid, false: at least one invalid parameter exists.
* `unavailable_config` - Invalid parameter set.


