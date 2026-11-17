---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_modify_account_remark_operation"
sidebar_current: "docs-tencentcloud-resource-postgresql_modify_account_remark_operation"
description: |-
  Provides a resource to create a postgresql modify_account_remark_operation
---

# tencentcloud_postgresql_modify_account_remark_operation

Provides a resource to create a postgresql modify_account_remark_operation

## Example Usage

```hcl
resource "tencentcloud_postgresql_modify_account_remark_operation" "modify_account_remark_operation" {
  db_instance_id = local.pgsql_id
  user_name      = "root"
  remark         = "hello_world"
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Instance ID in the format of postgres-4wdeb0zv.
* `remark` - (Required, String, ForceNew) New remarks corresponding to user `UserName`.
* `user_name` - (Required, String, ForceNew) Instance username.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyAccountRemark](https://www.tencentcloud.com/document/api/409/18108)

