---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_restart_data_engine_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_restart_data_engine_operation"
description: |-
  Provides a resource to create a dlc restart_data_engine
---

# tencentcloud_dlc_restart_data_engine_operation

Provides a resource to create a dlc restart_data_engine

## Example Usage

```hcl
resource "tencentcloud_dlc_restart_data_engine_operation" "restart_data_engine" {
  data_engine_id   = "DataEngine-g5ds87d8"
  forced_operation = false
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String, ForceNew) Engine unique id.
* `forced_operation` - (Optional, Bool, ForceNew) Whether to force restart and ignore tasks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



