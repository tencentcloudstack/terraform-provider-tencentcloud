---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_modify_switch_time_period_operation"
sidebar_current: "docs-tencentcloud-resource-postgresql_modify_switch_time_period_operation"
description: |-
  Provides a resource to create a postgresql modify_switch_time_period_operation
---

# tencentcloud_postgresql_modify_switch_time_period_operation

Provides a resource to create a postgresql modify_switch_time_period_operation

## Example Usage

```hcl
resource "tencentcloud_postgresql_modify_switch_time_period_operation" "modify_switch_time_period_operation" {
  db_instance_id = local.pgsql_id
  switch_tag     = 0
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) The ID of the instance waiting for a switch.
* `switch_tag` - (Required, Int, ForceNew) Valid value: `0` (switch immediately).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



