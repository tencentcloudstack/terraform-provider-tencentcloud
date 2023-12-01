---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_check_data_engine_image_can_be_rollback"
sidebar_current: "docs-tencentcloud-datasource-dlc_check_data_engine_image_can_be_rollback"
description: |-
  Use this data source to query detailed information of dlc check_data_engine_image_can_be_rollback
---

# tencentcloud_dlc_check_data_engine_image_can_be_rollback

Use this data source to query detailed information of dlc check_data_engine_image_can_be_rollback

## Example Usage

```hcl
data "tencentcloud_dlc_check_data_engine_image_can_be_rollback" "check_data_engine_image_can_be_rollback" {
  data_engine_id = "DataEngine-public-1308919341"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String) Engine unique id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `from_record_id` - Log record id before rollback.
* `is_rollback` - Is it possible to roll back.
* `to_record_id` - Log record id after rollback.


