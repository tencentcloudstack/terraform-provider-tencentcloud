---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_restart_db_instance_operation"
sidebar_current: "docs-tencentcloud-resource-postgresql_restart_db_instance_operation"
description: |-
  Provides a resource to create a postgresql restart_db_instance_operation
---

# tencentcloud_postgresql_restart_db_instance_operation

Provides a resource to create a postgresql restart_db_instance_operation

## Example Usage

```hcl
resource "tencentcloud_postgresql_restart_db_instance_operation" "restart_db_instance_operation" {
  db_instance_id = local.pgsql_id
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) dbInstance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeDBInstanceAttribute](https://cloud.tencent.com/document/api/409/16772)
* Create: [RestartDBInstance](https://cloud.tencent.com/document/api/409/18097)

