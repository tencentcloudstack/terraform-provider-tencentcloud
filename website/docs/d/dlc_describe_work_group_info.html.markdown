---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_work_group_info"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_work_group_info"
description: |-
  Use this data source to query detailed information of DLC describe work group info
---

# tencentcloud_dlc_describe_work_group_info

Use this data source to query detailed information of DLC describe work group info

## Example Usage

```hcl
data "tencentcloud_dlc_describe_work_group_info" "example" {
  work_group_id = 70220
  type          = "User"
  sort_by       = "create-time"
  sorting       = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter criteria that are queriedWhen the type is User, the fuzzy search is supported as the key is user-name.When the type is DataAuth, the keys supported are:policy-type: types of permissions;policy-source: data sources;data-name: fuzzy search of the database and table.When the type is EngineAuth, the keys supported are:policy-type: types of permissions;policy-source: data sources;engine-name: fuzzy search of the database and table.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_by` - (Optional, String) Sort fields.When the type is User, create-time and user-name are supported.When the type is DataAuth, create-time is supported.When the type is EngineAuth, create-time is supported.
* `sorting` - (Optional, String) Sorting methods: desc means in order; asc means in reverse order; it is asc by default.
* `type` - (Optional, String) Types of queried information. User: user information; DataAuth: data permissions; EngineAuth: engine permissions.
* `work_group_id` - (Optional, Int) Working group ID.

The `filters` object supports the following:

* `name` - (Required, String) Attribute name. If more than one filter exists, the logical relationship between these filters is `OR`.
* `values` - (Required, Set) Attribute value. If multiple values exist in one filter, the logical relationship between these values is `OR`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `work_group_info` - Details about working groupsNote: This field may return null, indicating that no valid values can be obtained.
  * `data_policy_info` - Collection of data permissionsNote: This field may return null, indicating that no valid values can be obtained.
    * `policy_set` - Collection of policiesNote: This field may return null, indicating that no valid values can be obtained.
      * `catalog` - The name of the target data source. To grant admin permission, it must be `*` (all resources at this level); to grant data source and database permissions, it must be `COSDataCatalog` or `*`; to grant table permissions, it can be a custom data source; if it is left empty, `DataLakeCatalog` is used. Note: To grant permissions on a custom data source, the permissions that can be managed in the Data Lake Compute console are subsets of the account permissions granted when you connect the data source to the console.
      * `column` - The name of the target column. `*` represents all columns. To grant admin permissions, it must be `*`.Note: This field may return null, indicating that no valid values can be obtained.
      * `create_time` - The permission policy creation time, which is not required as an input parameter.Note: This field may return null, indicating that no valid values can be obtained.
      * `data_engine` - The name of the target data engine. `*` represents all engines. To grant admin permissions, it must be `*`.Note: This field may return null, indicating that no valid values can be obtained.
      * `database` - The name of the target database. `*` represents all databases in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any database.
      * `function` - The name of the target function. `*` represents all functions in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any function.Note: This field may return null, indicating that no valid values can be obtained.
      * `id` - The policy ID.Note: This field may return null, indicating that no valid values can be obtained.
      * `mode` - The grant mode, which is not required as an input parameter. Valid values: `COMMON` and `SENIOR`.Note: This field may return null, indicating that no valid values can be obtained.
      * `operation` - The target permissions, which vary by permission level. Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`, `CREATE`, `ALTER`, and `DROP`; table: `ALL`, `SELECT`, `INSERT`, `ALTER`, `DELETE`, `DROP`, and `UPDATE`. Note: For table permissions, if a data source other than `COSDataCatalog` is specified, only the `SELECT` permission can be granted here.
      * `operator` - The operator, which is not required as an input parameter.Note: This field may return null, indicating that no valid values can be obtained.
      * `policy_type` - The permission type. Valid values: `ADMIN`, `DATASOURCE`, `DATABASE`, `TABLE`, `VIEW`, `FUNCTION`, `COLUMN`, and `ENGINE`. Note: If it is left empty, `ADMIN` is used.
      * `re_auth` - Whether the grantee is allowed to further grant the permissions. Valid values: `false` (default) and `true` (the grantee can grant permissions gained here to other sub-users).Note: This field may return null, indicating that no valid values can be obtained.
      * `source_id` - The ID of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.Note: This field may return null, indicating that no valid values can be obtained.
      * `source_name` - The name of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.Note: This field may return null, indicating that no valid values can be obtained.
      * `source` - The permission source, which is not required when input parameters are passed in. Valid values: `USER` (from the user) and `WORKGROUP` (from one or more associated work groups).Note: This field may return null, indicating that no valid values can be obtained.
      * `table` - The name of the target table. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.
      * `view` - The name of the target view. `*` represents all views in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any view.Note: This field may return null, indicating that no valid values can be obtained.
    * `total_count` - Total policiesNote: This field may return null, indicating that no valid values can be obtained.
  * `engine_policy_info` - Collection of engine permissionsNote: This field may return null, indicating that no valid values can be obtained.
    * `policy_set` - Collection of policiesNote: This field may return null, indicating that no valid values can be obtained.
      * `catalog` - The name of the target data source. To grant admin permission, it must be `*` (all resources at this level); to grant data source and database permissions, it must be `COSDataCatalog` or `*`; to grant table permissions, it can be a custom data source; if it is left empty, `DataLakeCatalog` is used. Note: To grant permissions on a custom data source, the permissions that can be managed in the Data Lake Compute console are subsets of the account permissions granted when you connect the data source to the console.
      * `column` - The name of the target column. `*` represents all columns. To grant admin permissions, it must be `*`.Note: This field may return null, indicating that no valid values can be obtained.
      * `create_time` - The permission policy creation time, which is not required as an input parameter.Note: This field may return null, indicating that no valid values can be obtained.
      * `data_engine` - The name of the target data engine. `*` represents all engines. To grant admin permissions, it must be `*`.Note: This field may return null, indicating that no valid values can be obtained.
      * `database` - The name of the target database. `*` represents all databases in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any database.
      * `function` - The name of the target function. `*` represents all functions in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any function.Note: This field may return null, indicating that no valid values can be obtained.
      * `id` - The policy ID.Note: This field may return null, indicating that no valid values can be obtained.
      * `mode` - The grant mode, which is not required as an input parameter. Valid values: `COMMON` and `SENIOR`.Note: This field may return null, indicating that no valid values can be obtained.
      * `operation` - The target permissions, which vary by permission level. Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`, `CREATE`, `ALTER`, and `DROP`; table: `ALL`, `SELECT`, `INSERT`, `ALTER`, `DELETE`, `DROP`, and `UPDATE`. Note: For table permissions, if a data source other than `COSDataCatalog` is specified, only the `SELECT` permission can be granted here.
      * `operator` - The operator, which is not required as an input parameter.Note: This field may return null, indicating that no valid values can be obtained.
      * `policy_type` - The permission type. Valid values: `ADMIN`, `DATASOURCE`, `DATABASE`, `TABLE`, `VIEW`, `FUNCTION`, `COLUMN`, and `ENGINE`. Note: If it is left empty, `ADMIN` is used.
      * `re_auth` - Whether the grantee is allowed to further grant the permissions. Valid values: `false` (default) and `true` (the grantee can grant permissions gained here to other sub-users).Note: This field may return null, indicating that no valid values can be obtained.
      * `source_id` - The ID of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.Note: This field may return null, indicating that no valid values can be obtained.
      * `source_name` - The name of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.Note: This field may return null, indicating that no valid values can be obtained.
      * `source` - The permission source, which is not required when input parameters are passed in. Valid values: `USER` (from the user) and `WORKGROUP` (from one or more associated work groups).Note: This field may return null, indicating that no valid values can be obtained.
      * `table` - The name of the target table. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.
      * `view` - The name of the target view. `*` represents all views in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any view.Note: This field may return null, indicating that no valid values can be obtained.
    * `total_count` - Total policiesNote: This field may return null, indicating that no valid values can be obtained.
  * `row_filter_info` - Collection of information about filtered rowsNote: This field may return null, indicating that no valid values can be obtained.
    * `policy_set` - Collection of policiesNote: This field may return null, indicating that no valid values can be obtained.
      * `catalog` - The name of the target data source. To grant admin permission, it must be `*` (all resources at this level); to grant data source and database permissions, it must be `COSDataCatalog` or `*`; to grant table permissions, it can be a custom data source; if it is left empty, `DataLakeCatalog` is used. Note: To grant permissions on a custom data source, the permissions that can be managed in the Data Lake Compute console are subsets of the account permissions granted when you connect the data source to the console.
      * `column` - The name of the target column. `*` represents all columns. To grant admin permissions, it must be `*`.Note: This field may return null, indicating that no valid values can be obtained.
      * `create_time` - The permission policy creation time, which is not required as an input parameter.Note: This field may return null, indicating that no valid values can be obtained.
      * `data_engine` - The name of the target data engine. `*` represents all engines. To grant admin permissions, it must be `*`.Note: This field may return null, indicating that no valid values can be obtained.
      * `database` - The name of the target database. `*` represents all databases in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any database.
      * `function` - The name of the target function. `*` represents all functions in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any function.Note: This field may return null, indicating that no valid values can be obtained.
      * `id` - The policy ID.Note: This field may return null, indicating that no valid values can be obtained.
      * `mode` - The grant mode, which is not required as an input parameter. Valid values: `COMMON` and `SENIOR`.Note: This field may return null, indicating that no valid values can be obtained.
      * `operation` - The target permissions, which vary by permission level. Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`, `CREATE`, `ALTER`, and `DROP`; table: `ALL`, `SELECT`, `INSERT`, `ALTER`, `DELETE`, `DROP`, and `UPDATE`. Note: For table permissions, if a data source other than `COSDataCatalog` is specified, only the `SELECT` permission can be granted here.
      * `operator` - The operator, which is not required as an input parameter.Note: This field may return null, indicating that no valid values can be obtained.
      * `policy_type` - The permission type. Valid values: `ADMIN`, `DATASOURCE`, `DATABASE`, `TABLE`, `VIEW`, `FUNCTION`, `COLUMN`, and `ENGINE`. Note: If it is left empty, `ADMIN` is used.
      * `re_auth` - Whether the grantee is allowed to further grant the permissions. Valid values: `false` (default) and `true` (the grantee can grant permissions gained here to other sub-users).Note: This field may return null, indicating that no valid values can be obtained.
      * `source_id` - The ID of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.Note: This field may return null, indicating that no valid values can be obtained.
      * `source_name` - The name of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.Note: This field may return null, indicating that no valid values can be obtained.
      * `source` - The permission source, which is not required when input parameters are passed in. Valid values: `USER` (from the user) and `WORKGROUP` (from one or more associated work groups).Note: This field may return null, indicating that no valid values can be obtained.
      * `table` - The name of the target table. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.
      * `view` - The name of the target view. `*` represents all views in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any view.Note: This field may return null, indicating that no valid values can be obtained.
    * `total_count` - Total policiesNote: This field may return null, indicating that no valid values can be obtained.
  * `type` - Types of information included. User: user information; DataAuth: data permissions; EngineAuth: engine permissionsNote: This field may return null, indicating that no valid values can be obtained.
  * `user_info` - Collection of users bound to working groupsNote: This field may return null, indicating that no valid values can be obtained.
    * `total_count` - Total usersNote: This field may return null, indicating that no valid values can be obtained.
    * `user_set` - Collection of user informationNote: This field may return null, indicating that no valid values can be obtained.
      * `create_time` - The creation time of the current user, e.g. 16:19:32, July 28, 2021.
      * `creator` - The creator of the current user.
      * `user_alias` - User alias.
      * `user_description` - User descriptionNote: The returned value of this field may be null, indicating that no valid value is obtained.
      * `user_id` - User Id which matches the sub-user UIN on the CAM side.
  * `work_group_description` - Working group descriptionNote: This field may return null, indicating that no valid values can be obtained.
  * `work_group_id` - Working group IDNote: This field may return null, indicating that no valid values can be obtained.
  * `work_group_name` - Working group nameNote: This field may return null, indicating that no valid values can be obtained.


