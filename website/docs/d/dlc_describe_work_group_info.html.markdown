---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_work_group_info"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_work_group_info"
description: |-
  Use this data source to query detailed information of dlc describe_work_group_info
---

# tencentcloud_dlc_describe_work_group_info

Use this data source to query detailed information of dlc describe_work_group_info

## Example Usage

```hcl
data "tencentcloud_dlc_describe_work_group_info" "describe_work_group_info" {
  work_group_id = 23181
  type          = "User"
  sort_by       = "create-time"
  sorting       = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Query filter conditions. when Type is User, fuzzy search with Key as user-name is supported; when Type is DataAuth, key is supported; policy-type: permission type; policy-source: data source; data-name: database table fuzzy search; when Type is EngineAuth, supports key; policy-type: permission type; policy-source: data source; engine-name: fuzzy search of library tables.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_by` - (Optional, String) Sorting fields, when Type is User, support create-time, user-name, when type is DataAuth, support create-time, when type is EngineAuth, support create-time.
* `sorting` - (Optional, String) Sorting method, desc means forward order, asc means reverse order, the default is asc.
* `type` - (Optional, String) Query information type, only support: User: user information/DataAuth: data permission/EngineAuth: engine permission.
* `work_group_id` - (Optional, Int) Work group id.

The `filters` object supports the following:

* `name` - (Required, String) Attribute name. If there are multiple Filters, the relationship between filters is a logical or (OR) relationship.
* `values` - (Required, Set) Attribute value, if there are multiple values in the same filter, the relationship between values under the same filter is a logical or relationship.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `work_group_info` - Workgroup details.
  * `data_policy_info` - Data permission collection.
    * `policy_set` - Policy set.
      * `catalog` - For the data source name that requires authorization, only * (representing all resources at this level) is supported under the administrator level; in the case of data source level and database level authentication, only COSDataCatalog or * is supported; in data table level authentication, it is possible Fill in the user-defined data source. If left blank, it defaults to DataLakeCatalog. note: If a user-defined data source is authenticated, the permissions that dlc can manage are a subset of the accounts provided by the user when accessing the data source.
      * `column` - For columns that require authorization, fill in * to represent all current columns. When the authorization type is administrator level, only * is allowed.
      * `create_time` - The time when the permission was created. Leave the input parameter blank.
      * `data_engine` - Data engines that require authorization, fill in * to represent all current engines. when the authorization type is administrator level, only * is allowed.
      * `database` - Database name that requires authorization, fill in * to represent all databases under the current catalog. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level, only blanks are allowed to be filled in. For other types, the database can be specified arbitrarily.
      * `function` - For the function name that requires authorization, fill in * to represent all functions under the current catalog. when the authorization type is administrator level, only * is allowed to be filled in. When the authorization type is data connection level, only blanks are allowed to be filled in. in other types, functions can be specified arbitrarily.
      * `id` - Policy id.
      * `mode` - Authorization mode, please leave this parameter blank. COMMON: normal mode; SENIOR: advanced mode.
      * `operation` - Authorized permission operations provide different operations for different levels of authentication. administrator permissions: ALL, default is ALL if left blank; data connection level authentication: CREATE; database level authentication: ALL, CREATE, ALTER, DROP; data table permissions: ALL, SELECT, INSERT, ALTER, DELETE, DROP, UPDATE. note: under data table permissions, only SELECT operations are supported when the specified data source is not COSDataCatalog.
      * `operator` - Operator, do not fill in the input parameters.
      * `policy_type` - Authorization type, currently supports eight authorization types: ADMIN: Administrator level authentication DATASOURCE: data connection level authentication DATABASE: database level authentication TABLE: Table level authentication VIEW: view level authentication FUNCTION: Function level authentication COLUMN: Column level authentication ENGINE: Data engine authentication. if left blank, the default is administrator level authentication.
      * `re_auth` - Whether the user can perform secondary authorization. when it is true, the authorized user can re-authorize the permissions obtained this time to other sub-users. default is false.
      * `source_id` - The id of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the Source field is WORKGROUP.
      * `source_name` - The name of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the source field is WORKGROUP.
      * `source` - Permission source, please leave it blank. USER: permissions come from the user itself; WORKGROUP: permissions come from the bound workgroup.
      * `table` - For the table name that requires authorization, fill in * to represent all tables under the current database. when the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. For other types, data tables can be specified arbitrarily.
      * `view` - For views that require authorization, fill in * to represent all views under the current database. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. for other types, the view can be specified arbitrarily.
    * `total_count` - Total count.
  * `engine_policy_info` - Engine permission collection.
    * `policy_set` - Policy set.
      * `catalog` - For the data source name that requires authorization, only * (representing all resources at this level) is supported under the administrator level; in the case of data source level and database level authentication, only COSDataCatalog or * is supported; in data table level authentication, it is possible Fill in the user-defined data source. If left blank, it defaults to DataLakeCatalog. note: If a user-defined data source is authenticated, the permissions that dlc can manage are a subset of the accounts provided by the user when accessing the data source.
      * `column` - For columns that require authorization, fill in * to represent all current columns. When the authorization type is administrator level, only * is allowed.
      * `create_time` - The time when the permission was created. Leave the input parameter blank.
      * `data_engine` - Data engines that require authorization, fill in * to represent all current engines. when the authorization type is administrator level, only * is allowed.
      * `database` - Database name that requires authorization, fill in * to represent all databases under the current catalog. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level, only blanks are allowed to be filled in. For other types, the database can be specified arbitrarily.
      * `function` - For the function name that requires authorization, fill in * to represent all functions under the current catalog. when the authorization type is administrator level, only * is allowed to be filled in. When the authorization type is data connection level, only blanks are allowed to be filled in. in other types, functions can be specified arbitrarily.
      * `id` - Policy id.
      * `mode` - Authorization mode, please leave this parameter blank. COMMON: normal mode; SENIOR: advanced mode.
      * `operation` - Authorized permission operations provide different operations for different levels of authentication. administrator permissions: ALL, default is ALL if left blank; data connection level authentication: CREATE; database level authentication: ALL, CREATE, ALTER, DROP; data table permissions: ALL, SELECT, INSERT, ALTER, DELETE, DROP, UPDATE. note: under data table permissions, only SELECT operations are supported when the specified data source is not COSDataCatalog.
      * `operator` - Operator, do not fill in the input parameters.
      * `policy_type` - Authorization type, currently supports eight authorization types: ADMIN: Administrator level authentication DATASOURCE: data connection level authentication DATABASE: database level authentication TABLE: Table level authentication VIEW: view level authentication FUNCTION: Function level authentication COLUMN: Column level authentication ENGINE: Data engine authentication. if left blank, the default is administrator level authentication.
      * `re_auth` - Whether the user can perform secondary authorization. when it is true, the authorized user can re-authorize the permissions obtained this time to other sub-users. default is false.
      * `source_id` - The id of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the Source field is WORKGROUP.
      * `source_name` - The name of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the source field is WORKGROUP.
      * `source` - Permission source, please leave it blank. USER: permissions come from the user itself; WORKGROUP: permissions come from the bound workgroup.
      * `table` - For the table name that requires authorization, fill in * to represent all tables under the current database. when the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. For other types, data tables can be specified arbitrarily.
      * `view` - For views that require authorization, fill in * to represent all views under the current database. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. for other types, the view can be specified arbitrarily.
    * `total_count` - Total count.
  * `row_filter_info` - Row filter information collection.
    * `policy_set` - Policy set.
      * `catalog` - For the data source name that requires authorization, only * (representing all resources at this level) is supported under the administrator level; in the case of data source level and database level authentication, only COSDataCatalog or * is supported; in data table level authentication, it is possible Fill in the user-defined data source. If left blank, it defaults to DataLakeCatalog. note: If a user-defined data source is authenticated, the permissions that dlc can manage are a subset of the accounts provided by the user when accessing the data source.
      * `column` - For columns that require authorization, fill in * to represent all current columns. When the authorization type is administrator level, only * is allowed.
      * `create_time` - The time when the permission was created. Leave the input parameter blank.
      * `data_engine` - Data engines that require authorization, fill in * to represent all current engines. when the authorization type is administrator level, only * is allowed.
      * `database` - Database name that requires authorization, fill in * to represent all databases under the current catalog. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level, only blanks are allowed to be filled in. For other types, the database can be specified arbitrarily.
      * `function` - For the function name that requires authorization, fill in * to represent all functions under the current catalog. when the authorization type is administrator level, only * is allowed to be filled in. When the authorization type is data connection level, only blanks are allowed to be filled in. in other types, functions can be specified arbitrarily.
      * `id` - Policy id.
      * `mode` - Authorization mode, please leave this parameter blank. COMMON: normal mode; SENIOR: advanced mode.
      * `operation` - Authorized permission operations provide different operations for different levels of authentication. administrator permissions: ALL, default is ALL if left blank; data connection level authentication: CREATE; database level authentication: ALL, CREATE, ALTER, DROP; data table permissions: ALL, SELECT, INSERT, ALTER, DELETE, DROP, UPDATE. note: under data table permissions, only SELECT operations are supported when the specified data source is not COSDataCatalog.
      * `operator` - Operator, do not fill in the input parameters.
      * `policy_type` - Authorization type, currently supports eight authorization types: ADMIN: Administrator level authentication DATASOURCE: data connection level authentication DATABASE: database level authentication TABLE: Table level authentication VIEW: view level authentication FUNCTION: Function level authentication COLUMN: Column level authentication ENGINE: Data engine authentication. if left blank, the default is administrator level authentication.
      * `re_auth` - Whether the user can perform secondary authorization. when it is true, the authorized user can re-authorize the permissions obtained this time to other sub-users. default is false.
      * `source_id` - The id of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the Source field is WORKGROUP.
      * `source_name` - The name of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the source field is WORKGROUP.
      * `source` - Permission source, please leave it blank. USER: permissions come from the user itself; WORKGROUP: permissions come from the bound workgroup.
      * `table` - For the table name that requires authorization, fill in * to represent all tables under the current database. when the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. For other types, data tables can be specified arbitrarily.
      * `view` - For views that require authorization, fill in * to represent all views under the current database. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. for other types, the view can be specified arbitrarily.
    * `total_count` - Total count.
  * `type` - The type of information contained. User: user information; DataAuth: data permissions; EngineAuth: engine permissions.
  * `user_info` - A collection of users bound to the workgroup.
    * `total_count` - Total count.
    * `user_set` - User information collection.
      * `create_time` - Create time.
      * `creator` - The creator of the current user.
      * `user_alias` - User alias.
      * `user_description` - User description.
      * `user_id` - User id, matches the CAM side sub-user uin.
  * `work_group_description` - Workgroup description information.
  * `work_group_id` - Work group id.
  * `work_group_name` - Work group name.


