---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_switch_data_engine_image_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_switch_data_engine_image_operation"
description: |-
  Provides a resource to create a dlc switch_data_engine_image_operation
---

# tencentcloud_dlc_switch_data_engine_image_operation

Provides a resource to create a dlc switch_data_engine_image_operation

## Example Usage

```hcl
resource "tencentcloud_dlc_switch_data_engine_image_operation" "switch_data_engine_image_operation" {
  data_engine_id       = "DataEngine-g5ds87d8"
  new_image_version_id = "344ba1c6-b7a9-403a-a255-422fffed6d38"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String, ForceNew) Engine unique id.
* `new_image_version_id` - (Required, String, ForceNew) New image version id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



