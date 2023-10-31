---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_modify_data_engine_description_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_modify_data_engine_description_operation"
description: |-
  Provides a resource to create a dlc modify_data_engine_description_operation
---

# tencentcloud_dlc_modify_data_engine_description_operation

Provides a resource to create a dlc modify_data_engine_description_operation

## Example Usage

```hcl
resource "tencentcloud_dlc_modify_data_engine_description_operation" "modify_data_engine_description_operation" {
  data_engine_name = "testEngine"
  message          = "test"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_name` - (Required, String, ForceNew) The name of the engine to modify.
* `message` - (Required, String, ForceNew) Engine description information, the maximum length is 250.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dlc modify_data_engine_description_operation can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_modify_data_engine_description_operation.modify_data_engine_description_operation modify_data_engine_description_operation_id
```

