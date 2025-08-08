---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_check_data_engine_image_can_be_upgrade"
sidebar_current: "docs-tencentcloud-datasource-dlc_check_data_engine_image_can_be_upgrade"
description: |-
  Use this data source to query detailed information of DLC check data engine image can be upgrade
---

# tencentcloud_dlc_check_data_engine_image_can_be_upgrade

Use this data source to query detailed information of DLC check data engine image can be upgrade

## Example Usage

```hcl
data "tencentcloud_dlc_check_data_engine_image_can_be_upgrade" "example" {
  data_engine_id = "DataEngine-80ibn1cj"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String) Unique engine ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `child_image_version_id` - ID of the minor version of the cluster image that can be updated under the major version.
* `is_upgrade` - Whether it can be updated.


