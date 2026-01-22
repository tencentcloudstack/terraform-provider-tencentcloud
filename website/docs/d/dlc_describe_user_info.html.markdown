---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_user_info"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_user_info"
description: |-
  Use this data source to query detailed information of DLC describe user info
---

# tencentcloud_dlc_describe_user_info

Use this data source to query detailed information of DLC describe user info

## Example Usage

```hcl
data "tencentcloud_dlc_describe_user_info" "example" {
  user_id = "100021240189"
  type    = "Group"
  sort_by = "create-time"
  sorting = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter criteria that are queriedWhen the type is Group, the fuzzy search is supported as the key is workgroup-name.When the type is DataAuth, the keys supported are:policy-type: types of permissions;policy-source: data sources;data-name: fuzzy search of the database and table.When the type is EngineAuth, the keys supported are:policy-type: types of permissions;policy-source: data sources;engine-name: fuzzy search of the database and table.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_by` - (Optional, String) Sort fields.When the type is Group, the create-time and group-name are supported.When the type is DataAuth, create-time is supported.When the type is EngineAuth, create-time is supported.
* `sorting` - (Optional, String) Sorting methods: desc means in order; asc means in reverse order; it is asc by default.
* `type` - (Optional, String) Type of queried information. Group: working group; DataAuth: data permission; EngineAuth: engine permission.
* `user_id` - (Optional, String) User ID.

The `filters` object supports the following:

* `name` - (Required, String) Attribute name. If more than one filter exists, the logical relationship between these filters is `OR`.
* `values` - (Required, Set) Attribute value. If multiple values exist in one filter, the logical relationship between these values is `OR`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `user_info` - Detailed user informationNote: This field may return null, indicating that no valid values can be obtained.
  * `account_type` - Account type.
  * `catalog_policy_info` - Collection of catalog permissionsNote: This field may return null, indicating that no valid values can be obtained.
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
  * `data_policy_info` - Collection of data permission informationNote: This field may return null, indicating that no valid values can be obtained.
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
  * `row_filter_info` - Collection of filtered rowsNote: This field may return null, indicating that no valid values can be obtained.
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
  * `type` - Types of returned information. Group: returned information about the working group where the current user is; DataAuth: returned information about the current user&amp;#39;s data permission; EngineAuth: returned information about the current user&amp;#39;s engine permissionNote: This field may return null, indicating that no valid values can be obtained.
  * `user_alias` - User aliasNote: This field may return null, indicating that no valid values can be obtained.
  * `user_description` - User descriptionNote: This field may return null, indicating that no valid values can be obtained.
  * `user_id` - User IDNote: This field may return null, indicating that no valid values can be obtained.
  * `user_type` - Types of users. ADMIN: administrators; COMMON: general usersNote: This field may return null, indicating that no valid values can be obtained.
  * `work_group_info` - Information about collections of working groups bound to the userNote: This field may return null, indicating that no valid values can be obtained.
    * `total_count` - Total working groupsNote: This field may return null, indicating that no valid values can be obtained.
    * `work_group_set` - Collection of working group informationNote: This field may return null, indicating that no valid values can be obtained.
      * `create_time` - The creation time of the working group, e.g. at 16:19:32 on Jul 28, 2021.
      * `creator` - Creator.
      * `work_group_description` - Working group descriptionNote: This field may return null, indicating that no valid values can be obtained.
      * `work_group_id` - Unique ID of the working group.
      * `work_group_name` - Working group name.


