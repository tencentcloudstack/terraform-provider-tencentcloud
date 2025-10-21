---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_rollback_data_engine_image_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_rollback_data_engine_image_operation"
description: |-
  Provides a resource to create a DLC rollback data engine image
---

# tencentcloud_dlc_rollback_data_engine_image_operation

Provides a resource to create a DLC rollback data engine image

## Example Usage

```hcl
data "tencentcloud_dlc_check_data_engine_image_can_be_rollback" "example" {
  data_engine_id = "DataEngine-cgkvbas6"
}

resource "tencentcloud_dlc_rollback_data_engine_image_operation" "example" {
  data_engine_id = "DataEngine-cgkvbas6"
  from_record_id = data.tencentcloud_dlc_check_data_engine_image_can_be_rollback.example.from_record_id
  to_record_id   = data.tencentcloud_dlc_check_data_engine_image_can_be_rollback.example.to_record_id
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String, ForceNew) Engine ID.
* `from_record_id` - (Optional, String, ForceNew) FromRecordId parameters returned by the API for checking the availability of rolling back.
* `to_record_id` - (Optional, String, ForceNew) ToRecordId parameters returned by the API for checking the availability of rolling back.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



