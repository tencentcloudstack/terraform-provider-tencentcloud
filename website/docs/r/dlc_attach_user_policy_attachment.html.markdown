---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_attach_user_policy_attachment"
sidebar_current: "docs-tencentcloud-resource-dlc_attach_user_policy_attachment"
description: |-
  Provides a resource to create a DLC attach user policy attachment
---

# tencentcloud_dlc_attach_user_policy_attachment

Provides a resource to create a DLC attach user policy attachment

~> **NOTE:** `policy_id` format: `v1|{SubjectType}|{SubjectId}|{PolicyType}|{Mode}|{Catalog}|{Database}|{Table}|{View}|{Function}|{Column}|{DataEngine}|{Operation}`

## Example Usage

### If policy_type is ENGINE

```hcl
resource "tencentcloud_dlc_attach_user_policy_attachment" "example" {
  user_id      = "100010109702"
  account_type = "TencentAccount"
  policy_set {
    policy_type = "ENGINE"
    data_engine = "test_engine"
    operation   = "USE,MONITOR"
    source      = "USER"
  }
}
```

### If policy_type is DATABASE

```hcl
resource "tencentcloud_dlc_attach_user_policy_attachment" "example1" {
  user_id      = "100010109702"
  account_type = "TencentAccount"
  policy_set {
    policy_type = "DATABASE"
    catalog     = "DataLakeCatalog"
    database    = "test_database"
    mode        = "COMMON"
    operation   = "ASSAYER"
    source      = "USER"
  }
}
```

### If policy_type is ROWFILTER

```hcl
resource "tencentcloud_dlc_attach_user_policy_attachment" "example2" {
  user_id      = "100010109702"
  account_type = "TencentAccount"
  policy_set {
    policy_type = "ROWFILTER"
    catalog     = "DataLakeCatalog"
    database    = "test_database"
    table       = "test_table"
    mode        = "SENIOR"
    operation   = "year > 2026 and country == 'US'"
    source      = "USER"
  }
}
```

## Argument Reference

The following arguments are supported:

* `policy_set` - (Required, List, ForceNew) Collection of authentication policies. Only one policy is allowed to be attached per resource.
* `user_id` - (Required, String, ForceNew) User ID, which is the same as the sub-user UIN. The CreateUser API is needed to create a user at first. The DescribeUsers API can be used for viewing.
* `account_type` - (Optional, String, ForceNew) User source type. Valid values: `TencentAccount` (common Tencent Cloud user) and `EntraAccount` (Microsoft user).

The `policy_set` object supports the following:

* `catalog` - (Optional, String) The name of the target data source. To grant admin permission, it must be `*` (all resources at this level); to grant data source and database permissions, it must be `COSDataCatalog` or `*`; to grant table permissions, it can be a custom data source; if it is left empty, `DataLakeCatalog` is used.
* `column` - (Optional, String) The name of the target column. `*` represents all columns. To grant admin permissions, it must be `*`.
* `data_engine` - (Optional, String) The name of the target data engine. `*` represents all engines. To grant admin permissions, it must be `*`.
* `database` - (Optional, String) The name of the target database. `*` represents all databases in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any database.
* `engine_generation` - (Optional, String) Engine type.
* `function` - (Optional, String) The name of the target function. `*` represents all functions in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any function.
* `mode` - (Optional, String) The grant mode, Valid values: `COMMON` and `SENIOR`.
* `model` - (Optional, String) The name of the target Model. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.
* `operation` - (Optional, String) The target permissions, which vary by permission level. Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`, `CREATE`, `ALTER`, and `DROP`; table: `ALL`, `SELECT`, `INSERT`, `ALTER`, `DELETE`, `DROP`, and `UPDATE`.
* `policy_type` - (Optional, String) The permission type. Valid values: `ADMIN`, `DATASOURCE`, `DATABASE`, `TABLE`, `VIEW`, `FUNCTION`, `COLUMN`, and `ENGINE`. Note: If it is left empty, `ADMIN` is used.
* `source` - (Optional, String) The permission source, Valid values: `USER` (from the user) and `WORKGROUP` (from one or more associated work groups).
* `table` - (Optional, String) The name of the target table. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.
* `view` - (Optional, String) The name of the target view. `*` represents all views in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any view.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

DLC attach user policy attachment can be imported using the composite id (`user_id#policy_id`), e.g.

```
terraform import tencentcloud_dlc_attach_user_policy_attachment.example 100010109702#v1|USER|100010109702|DATABASE|COMMON|DataLakeCatalog|test_database||||||ASSAYER
```

