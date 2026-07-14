---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_attach_user_policyr_attachment"
sidebar_current: "docs-tencentcloud-resource-dlc_attach_user_policyr_attachment"
description: |-
  Provides a resource to create a DLC attach user policyr attachment
---

# tencentcloud_dlc_attach_user_policyr_attachment

Provides a resource to create a DLC attach user policyr attachment

## Example Usage

```hcl
resource "tencentcloud_dlc_attach_user_policyr_attachment" "example" {
  user_id      = "100032676511"
  account_type = "TencentAccount"
  policy_set {
    database    = "tf_example_db"
    catalog     = "DataLakeCatalog"
    table       = "tf_example_table"
    operation   = "SELECT"
    policy_type = "TABLE"
  }
}
```

## Argument Reference

The following arguments are supported:

* `policy_set` - (Required, List, ForceNew) Collection of authentication policies.
* `user_id` - (Required, String, ForceNew) User ID, which is the same as the sub-user UIN. The CreateUser API is needed to create a user at first. The DescribeUsers API can be used for viewing.
* `account_type` - (Optional, String, ForceNew) User source type. Valid values: `TencentAccount` (common Tencent Cloud user) and `EntraAccount` (Microsoft user).

The `policy_set` object supports the following:

* `catalog` - (Optional, String) The name of the target data source. To grant admin permission, it must be `*` (all resources at this level); to grant data source and database permissions, it must be `COSDataCatalog` or `*`; to grant table permissions, it can be a custom data source; if it is left empty, `DataLakeCatalog` is used.
* `column` - (Optional, String) The name of the target column. `*` represents all columns. To grant admin permissions, it must be `*`.
* `data_engine` - (Optional, String) The name of the target data engine. `*` represents all engines. To grant admin permissions, it must be `*`.
* `database` - (Optional, String) The name of the target database. `*` represents all databases in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any database.
* `engine_generation` - (Optional, String) Engine type.
* `function` - (Optional, String) The name of the target function. `*` represents all functions in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any function.
* `model` - (Optional, String) The name of the target Model. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.
* `operation` - (Optional, String) The target permissions, which vary by permission level. Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`, `CREATE`, `ALTER`, and `DROP`; table: `ALL`, `SELECT`, `INSERT`, `ALTER`, `DELETE`, `DROP`, and `UPDATE`.
* `policy_id` - (Optional, String) The deterministic string PolicyId corresponding to the user and workgroup.
* `policy_type` - (Optional, String) The permission type. Valid values: `ADMIN`, `DATASOURCE`, `DATABASE`, `TABLE`, `VIEW`, `FUNCTION`, `COLUMN`, and `ENGINE`. Note: If it is left empty, `ADMIN` is used.
* `re_auth` - (Optional, Bool) Whether the grantee is allowed to further grant the permissions. Valid values: `false` (default) and `true` (the grantee can grant permissions gained here to other sub-users).
* `table` - (Optional, String) The name of the target table. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.
* `view` - (Optional, String) The name of the target view. `*` represents all views in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any view.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

DLC attach user policyr attachment can be imported using the composite id, e.g.

```
terraform import tencentcloud_dlc_attach_user_policyr_attachment.example 100032676511#TencentAccount
```

