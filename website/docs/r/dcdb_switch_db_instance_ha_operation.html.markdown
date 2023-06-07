---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_switch_db_instance_ha_operation"
sidebar_current: "docs-tencentcloud-resource-dcdb_switch_db_instance_ha_operation"
description: |-
  Provides a resource to create a dcdb switch_db_instance_ha_operation
---

# tencentcloud_dcdb_switch_db_instance_ha_operation

Provides a resource to create a dcdb switch_db_instance_ha_operation

## Example Usage

```hcl
resource "tencentcloud_dcdb_switch_db_instance_ha_operation" "switch_db_instance_ha_operation" {
  instance_id = ""
  zone        = ""
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID in the format of tdsqlshard-ow728lmc.
* `zone` - (Required, String, ForceNew) Target AZ. The node with the lowest delay in the target AZ will be automatically promoted to primary node.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



