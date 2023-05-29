---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_restart_db_instances_operation"
sidebar_current: "docs-tencentcloud-resource-mysql_restart_db_instances_operation"
description: |-
  Provides a resource to create a mysql restart_db_instances_operation
---

# tencentcloud_mysql_restart_db_instances_operation

Provides a resource to create a mysql restart_db_instances_operation

## Example Usage

```hcl
resource "tencentcloud_mysql_restart_db_instances_operation" "restart_db_instances_operation" {
  instance_id = "cdb-bohspx3j"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) An array of instance ID in the format: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Instance status.


## Import

mysql restart_db_instances_operation can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_restart_db_instances_operation.restart_db_instances_operation restart_db_instances_operation_id
```

