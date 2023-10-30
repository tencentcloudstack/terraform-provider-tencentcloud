---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_update_data_engine_config_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_update_data_engine_config_operation"
description: |-
  Provides a resource to create a dlc update_data_engine_config_operation
---

# tencentcloud_dlc_update_data_engine_config_operation

Provides a resource to create a dlc update_data_engine_config_operation

## Example Usage

```hcl
resource "tencentcloud_dlc_update_data_engine_config_operation" "update_data_engine_config_operation" {
  data_engine_ids            =
  data_engine_config_command = "UpdateSparkSQLLakefsPath"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_config_command` - (Required, String, ForceNew) Engine configuration command, supports UpdateSparkSQLLakefsPath (update native table configuration), UpdateSparkSQLResultPath (update result path configuration).
* `data_engine_ids` - (Required, Set: [`String`], ForceNew) Engine unique id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dlc update_data_engine_config_operation can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_update_data_engine_config_operation.update_data_engine_config_operation update_data_engine_config_operation_id
```

