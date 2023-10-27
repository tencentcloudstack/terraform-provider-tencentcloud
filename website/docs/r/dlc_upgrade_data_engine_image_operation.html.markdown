---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_upgrade_data_engine_image_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_upgrade_data_engine_image_operation"
description: |-
  Provides a resource to create a dlc upgrade_data_engine_image_operation
---

# tencentcloud_dlc_upgrade_data_engine_image_operation

Provides a resource to create a dlc upgrade_data_engine_image_operation

## Example Usage

```hcl
resource "tencentcloud_dlc_upgrade_data_engine_image_operation" "upgrade_data_engine_image_operation" {
  data_engine_id = "DataEngine-g5ds87d8"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String, ForceNew) Engine unique id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



