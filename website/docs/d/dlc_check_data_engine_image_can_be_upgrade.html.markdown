---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_check_data_engine_image_can_be_upgrade"
sidebar_current: "docs-tencentcloud-datasource-dlc_check_data_engine_image_can_be_upgrade"
description: |-
  Use this data source to query detailed information of dlc check_data_engine_image_can_be_upgrade
---

# tencentcloud_dlc_check_data_engine_image_can_be_upgrade

Use this data source to query detailed information of dlc check_data_engine_image_can_be_upgrade

## Example Usage

```hcl
data "tencentcloud_dlc_check_data_engine_image_can_be_upgrade" "check_data_engine_image_can_be_upgrade" {
  data_engine_id = "DataEngine-cgkvbas6"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String) Engine unique id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `child_image_version_id` - The latest image version id that can be upgraded.
* `is_upgrade` - Is it possible to upgrade.


