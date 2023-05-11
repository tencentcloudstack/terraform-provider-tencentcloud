---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_general_clone"
sidebar_current: "docs-tencentcloud-resource-sqlserver_general_clone"
description: |-
  Provides a resource to create a sqlserver general_communication
---

# tencentcloud_sqlserver_general_clone

Provides a resource to create a sqlserver general_communication

## Example Usage

```hcl
resource "tencentcloud_sqlserver_general_clone" "general_clone" {
  instance_id = "Instance ID in the format of mssql-j8kv137v"
  old_name    = "old_db_name"
  new_name    = "new_db_name"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `new_name` - (Required, String) New database name. In offline migration, OldName will be used if NewName is left empty (OldName and NewName cannot be both empty). In database cloning, OldName and NewName must be both specified and cannot have the same value.
* `old_name` - (Required, String) Database name. If the OldName database does not exist, a failure will be returned. It can be left empty in offline migration tasks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `db_detail` - Sqlserver db Clone detail.
  * `is_fulltext_enabled` - Whether to enable full text 0: No 1: Yes.


## Import

sqlserver general_communication can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_general_communication.general_communication general_communication_id
```

