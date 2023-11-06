---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_update_row_filter_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_update_row_filter_operation"
description: |-
  Provides a resource to create a dlc update_row_filter_operation
---

# tencentcloud_dlc_update_row_filter_operation

Provides a resource to create a dlc update_row_filter_operation

## Example Usage

```hcl
resource "tencentcloud_dlc_update_row_filter_operation" "update_row_filter_operation" {
  policy_id = 103704
  policy {
    database    = "test_iac_keep"
    catalog     = "DataLakeCatalog"
    table       = "test_table"
    operation   = "value!=\"0\""
    policy_type = "ROWFILTER"
    function    = ""
    view        = ""
    column      = ""
    source      = "USER"
    mode        = "SENIOR"
    re_auth     = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, Int, ForceNew) The id of the row filtering policy.
* `policy` - (Required, List, ForceNew) New filtering strategy.

The `policy` object supports the following:

* `catalog` - (Required, String) For the data source name that requires authorization, only * (representing all resources at this level) is supported under the administrator level; in the case of data source level and database level authentication, only COSDataCatalog or * is supported; in data table level authentication, it is possible Fill in the user-defined data source. If left blank, it defaults to DataLakeCatalog. note: If a user-defined data source is authenticated, the permissions that dlc can manage are a subset of the accounts provided by the user when accessing the data source.
* `database` - (Required, String) Database name that requires authorization, fill in * to represent all databases under the current catalog. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level, only blanks are allowed to be filled in. For other types, the database can be specified arbitrarily.
* `operation` - (Required, String) Authorized permission operations provide different operations for different levels of authentication. administrator permissions: ALL, default is ALL if left blank; data connection level authentication: CREATE; database level authentication: ALL, CREATE, ALTER, DROP; data table permissions: ALL, SELECT, INSERT, ALTER, DELETE, DROP, UPDATE. note: under data table permissions, only SELECT operations are supported when the specified data source is not COSDataCatalog.
* `table` - (Required, String) For the table name that requires authorization, fill in * to represent all tables under the current database. when the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. For other types, data tables can be specified arbitrarily.
* `column` - (Optional, String) For columns that require authorization, fill in * to represent all current columns. When the authorization type is administrator level, only * is allowed.
* `create_time` - (Optional, String) The time when the permission was created. Leave the input parameter blank.
* `data_engine` - (Optional, String) Data engines that require authorization, fill in * to represent all current engines. when the authorization type is administrator level, only * is allowed.
* `function` - (Optional, String) For the function name that requires authorization, fill in * to represent all functions under the current catalog. when the authorization type is administrator level, only * is allowed to be filled in. When the authorization type is data connection level, only blanks are allowed to be filled in. in other types, functions can be specified arbitrarily.
* `id` - (Optional, Int) Policy id.
* `mode` - (Optional, String) Authorization mode, please leave this parameter blank. COMMON: normal mode; SENIOR: advanced mode.
* `operator` - (Optional, String) Operator, do not fill in the input parameters.
* `policy_type` - (Optional, String) Authorization type, currently supports eight authorization types: ADMIN: Administrator level authentication DATASOURCE: data connection level authentication DATABASE: database level authentication TABLE: Table level authentication VIEW: view level authentication FUNCTION: Function level authentication COLUMN: Column level authentication ENGINE: Data engine authentication. if left blank, the default is administrator level authentication.
* `re_auth` - (Optional, Bool) Whether the user can perform secondary authorization. when it is true, the authorized user can re-authorize the permissions obtained this time to other sub-users. default is false.
* `source_id` - (Optional, Int) The id of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the Source field is WORKGROUP.
* `source_name` - (Optional, String) The name of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the source field is WORKGROUP.
* `source` - (Optional, String) Permission source, please leave it blank. USER: permissions come from the user itself; WORKGROUP: permissions come from the bound workgroup.
* `view` - (Optional, String) For views that require authorization, fill in * to represent all views under the current database. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. for other types, the view can be specified arbitrarily.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



