---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_update_row_filter_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_update_row_filter_operation"
description: |-
  Provides a resource to create a DLC update row filter operation
---

# tencentcloud_dlc_update_row_filter_operation

Provides a resource to create a DLC update row filter operation

## Example Usage

```hcl
resource "tencentcloud_dlc_update_row_filter_operation" "example" {
  policy_id = 103704
  policy {
    database    = "tf_example_db"
    catalog     = "DataLakeCatalog"
    table       = "test_table"
    operation   = "value!=\"0\""
    policy_type = "ROWFILTER"
    source      = "USER"
    mode        = "SENIOR"
    re_auth     = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, Int, ForceNew) The ID of the row filter policy, which can be obtained using the `DescribeUserInfo` or `DescribeWorkGroupInfo` API.
* `policy` - (Required, List, ForceNew) The new filter policy.

The `policy` object supports the following:

* `catalog` - (Required, String) The name of the target data source. To grant admin permission, it must be `*` (all resources at this level); to grant data source and database permissions, it must be `COSDataCatalog` or `*`; to grant table permissions, it can be a custom data source; if it is left empty, `DataLakeCatalog` is used. Note: To grant permissions on a custom data source, the permissions that can be managed in the Data Lake Compute console are subsets of the account permissions granted when you connect the data source to the console.
* `database` - (Required, String) The name of the target database. `*` represents all databases in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any database.
* `operation` - (Required, String) The target permissions, which vary by permission level. Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`, `CREATE`, `ALTER`, and `DROP`; table: `ALL`, `SELECT`, `INSERT`, `ALTER`, `DELETE`, `DROP`, and `UPDATE`. Note: For table permissions, if a data source other than `COSDataCatalog` is specified, only the `SELECT` permission can be granted here.
* `table` - (Required, String) The name of the target table. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.
* `column` - (Optional, String) The name of the target column. `*` represents all columns. To grant admin permissions, it must be `*`.
* `create_time` - (Optional, String) The permission policy creation time, which is not required as an input parameter.
* `data_engine` - (Optional, String) The name of the target data engine. `*` represents all engines. To grant admin permissions, it must be `*`.
* `function` - (Optional, String) The name of the target function. `*` represents all functions in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any function.
* `id` - (Optional, Int) The policy ID.
* `mode` - (Optional, String) The grant mode, which is not required as an input parameter. Valid values: `COMMON` and `SENIOR`.
* `operator` - (Optional, String) The operator, which is not required as an input parameter.
* `policy_type` - (Optional, String) The permission type. Valid values: `ADMIN`, `DATASOURCE`, `DATABASE`, `TABLE`, `VIEW`, `FUNCTION`, `COLUMN`, and `ENGINE`. Note: If it is left empty, `ADMIN` is used.
* `re_auth` - (Optional, Bool) Whether the grantee is allowed to further grant the permissions. Valid values: `false` (default) and `true` (the grantee can grant permissions gained here to other sub-users).
* `source_id` - (Optional, Int) The ID of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.
* `source_name` - (Optional, String) The name of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.
* `source` - (Optional, String) The permission source, which is not required when input parameters are passed in. Valid values: `USER` (from the user) and `WORKGROUP` (from one or more associated work groups).
* `view` - (Optional, String) The name of the target view. `*` represents all views in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any view.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



