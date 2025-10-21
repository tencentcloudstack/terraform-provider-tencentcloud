---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_check_data_engine_image_can_be_rollback"
sidebar_current: "docs-tencentcloud-datasource-dlc_check_data_engine_image_can_be_rollback"
description: |-
  Use this data source to query detailed information of DLC check data engine image can be rollback
---

# tencentcloud_dlc_check_data_engine_image_can_be_rollback

Use this data source to query detailed information of DLC check data engine image can be rollback

## Example Usage

```hcl
data "tencentcloud_dlc_check_data_engine_image_can_be_rollback" "example" {
  data_engine_id = "DataEngine-80ibn1cj"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String) Unique engine ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `from_record_id` - Log record ID before rolling back.
* `is_rollback` - Whether it can be rolled back.
* `to_record_id` - Log record ID after rolling back.


