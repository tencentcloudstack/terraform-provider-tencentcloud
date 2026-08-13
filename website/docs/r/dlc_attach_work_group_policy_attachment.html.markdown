---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_attach_work_group_policy_attachment"
sidebar_current: "docs-tencentcloud-resource-dlc_attach_work_group_policy_attachment"
description: |-
  Provides a resource to create a DLC attach work group policy attachment
---

# tencentcloud_dlc_attach_work_group_policy_attachment

Provides a resource to create a DLC attach work group policy attachment

~> **NOTE:** `policy_id` format: `v1|{SubjectType}|{SubjectId}|{PolicyType}|{Mode}|{Catalog}|{Database}|{Table}|{View}|{Function}|{Column}|{DataEngine}|{Operation}`

## Example Usage

### If policy_type is ENGINE

```hcl
resource "tencentcloud_dlc_attach_work_group_policy_attachment" "example" {
  work_group_id = 21420
  policy_set {
    policy_type = "ENGINE"
    catalog     = ""
    database    = ""
    table       = ""
    data_engine = "test"
    operation   = "USE,MONITOR,MODIFY"
    source      = "WORKGROUP"
  }
}
```

### If policy_type is DATABASE

```hcl
resource "tencentcloud_dlc_attach_work_group_policy_attachment" "example1" {
  work_group_id = 21420
  policy_set {
    policy_type = "DATABASE"
    catalog     = "DataLakeCatalog"
    database    = "test_database"
    table       = ""
    operation   = "OWNER"
    source      = "WORKGROUP"
    mode        = "COMMON"
  }
}
```

### If policy_type is ROWFILTER

```hcl
resource "tencentcloud_dlc_attach_work_group_policy_attachment" "example2" {
  work_group_id = 21420
  policy_set {
    policy_type = "ROWFILTER"
    catalog     = "DataLakeCatalog"
    database    = "test_database"
    table       = "test_table"
    operation   = "year > 2026 and country == 'US'"
    source      = "WORKGROUP"
    mode        = "SENIOR"
  }
}
```

## Argument Reference

The following arguments are supported:

* `policy_set` - (Required, List, ForceNew) Collection of policies to be bound.
* `work_group_id` - (Required, Int, ForceNew) Work group ID.

The `policy_set` object supports the following:

* `catalog` - (Required, String) The name of the target data source. To grant admin permission, it must be `*` (all resources at this level); to grant data source and database permissions, it must be `COSDataCatalog` or `*`; to grant table permissions, it can be a custom data source; if it is left empty, `DataLakeCatalog` is used. Note: To grant permissions on a custom data source, the permissions that can be managed in the Data Lake Compute console are subsets of the account permissions granted when you connect the data source to the console.
* `database` - (Required, String) The name of the target database. `*` represents all databases in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any database.
* `operation` - (Required, String) The target permissions, which vary by permission level. Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`, `CREATE`, `ALTER`, and `DROP`; table: `ALL`, `SELECT`, `INSERT`, `ALTER`, `DELETE`, `DROP`, and `UPDATE`. Note: For table permissions, if a data source other than `COSDataCatalog` is specified, only the `SELECT` permission can be granted here.
* `table` - (Required, String) The name of the target table. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.
* `column` - (Optional, String) The name of the target column. `*` represents all columns. To grant admin permissions, it must be `*`.
* `data_engine` - (Optional, String) The name of the target data engine. `*` represents all engines. To grant admin permissions, it must be `*`.
* `engine_generation` - (Optional, String) The engine generation/type.
* `function` - (Optional, String) The name of the target function. `*` represents all functions in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any function.
* `mode` - (Optional, String) The grant mode, which is not required as an input parameter. Valid values: `COMMON` and `SENIOR`.
* `model` - (Optional, String) The name of the target Model. `*` represents all models in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any model.
* `policy_type` - (Optional, String) The permission type. Valid values: `ADMIN`, `DATASOURCE`, `DATABASE`, `TABLE`, `VIEW`, `FUNCTION`, `COLUMN`, and `ENGINE`. Note: If it is left empty, `ADMIN` is used.
* `source` - (Optional, String) The permission source, which is not required when input parameters are passed in. Valid values: `USER` (from the user) and `WORKGROUP` (from one or more associated work groups).
* `view` - (Optional, String) The name of the target view. `*` represents all views in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any view.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

DLC attach work group policy attachment can be imported using the composite id, e.g. The composite id is `WorkGroupId#PolicyId`.

```
terraform import tencentcloud_dlc_attach_work_group_policy_attachment.example 21420#v1|WORKGROUP|21420|DATABASE|COMMON|DataLakeCatalog|test_database||||||OWNER
```

