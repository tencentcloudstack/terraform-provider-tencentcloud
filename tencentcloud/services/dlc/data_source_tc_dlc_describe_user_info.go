package dlc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcDescribeUserInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeUserInfoRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User ID.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Type of queried information. Group: working group; DataAuth: data permission; EngineAuth: engine permission.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter criteria that are queriedWhen the type is Group, the fuzzy search is supported as the key is workgroup-name.When the type is DataAuth, the keys supported are:policy-type: types of permissions;policy-source: data sources;data-name: fuzzy search of the database and table.When the type is EngineAuth, the keys supported are:policy-type: types of permissions;policy-source: data sources;engine-name: fuzzy search of the database and table.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attribute name. If more than one filter exists, the logical relationship between these filters is `OR`.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Attribute value. If multiple values exist in one filter, the logical relationship between these values is `OR`.",
						},
					},
				},
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort fields.When the type is Group, the create-time and group-name are supported.When the type is DataAuth, create-time is supported.When the type is EngineAuth, create-time is supported.",
			},

			"sorting": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting methods: desc means in order; asc means in reverse order; it is asc by default.",
			},

			"user_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Detailed user informationNote: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User IDNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of returned information. Group: returned information about the working group where the current user is; DataAuth: returned information about the current user&amp;#39;s data permission; EngineAuth: returned information about the current user&amp;#39;s engine permissionNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"user_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of users. ADMIN: administrators; COMMON: general usersNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"user_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User descriptionNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"data_policy_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of data permission informationNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection of policiesNote: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target database. `*` represents all databases in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any database.",
												},
												"catalog": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target data source. To grant admin permission, it must be `*` (all resources at this level); to grant data source and database permissions, it must be `COSDataCatalog` or `*`; to grant table permissions, it can be a custom data source; if it is left empty, `DataLakeCatalog` is used. Note: To grant permissions on a custom data source, the permissions that can be managed in the Data Lake Compute console are subsets of the account permissions granted when you connect the data source to the console.",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target table. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.",
												},
												"operation": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The target permissions, which vary by permission level. Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`, `CREATE`, `ALTER`, and `DROP`; table: `ALL`, `SELECT`, `INSERT`, `ALTER`, `DELETE`, `DROP`, and `UPDATE`. Note: For table permissions, if a data source other than `COSDataCatalog` is specified, only the `SELECT` permission can be granted here.",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission type. Valid values: `ADMIN`, `DATASOURCE`, `DATABASE`, `TABLE`, `VIEW`, `FUNCTION`, `COLUMN`, and `ENGINE`. Note: If it is left empty, `ADMIN` is used.",
												},
												"function": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target function. `*` represents all functions in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any function.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target view. `*` represents all views in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any view.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target column. `*` represents all columns. To grant admin permissions, it must be `*`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"data_engine": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target data engine. `*` represents all engines. To grant admin permissions, it must be `*`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"re_auth": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the grantee is allowed to further grant the permissions. Valid values: `false` (default) and `true` (the grantee can grant permissions gained here to other sub-users).Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission source, which is not required when input parameters are passed in. Valid values: `USER` (from the user) and `WORKGROUP` (from one or more associated work groups).Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The grant mode, which is not required as an input parameter. Valid values: `COMMON` and `SENIOR`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The operator, which is not required as an input parameter.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission policy creation time, which is not required as an input parameter.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"source_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The ID of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The policy ID.Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total policiesNote: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"engine_policy_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of engine permissionsNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection of policiesNote: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target database. `*` represents all databases in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any database.",
												},
												"catalog": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target data source. To grant admin permission, it must be `*` (all resources at this level); to grant data source and database permissions, it must be `COSDataCatalog` or `*`; to grant table permissions, it can be a custom data source; if it is left empty, `DataLakeCatalog` is used. Note: To grant permissions on a custom data source, the permissions that can be managed in the Data Lake Compute console are subsets of the account permissions granted when you connect the data source to the console.",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target table. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.",
												},
												"operation": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The target permissions, which vary by permission level. Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`, `CREATE`, `ALTER`, and `DROP`; table: `ALL`, `SELECT`, `INSERT`, `ALTER`, `DELETE`, `DROP`, and `UPDATE`. Note: For table permissions, if a data source other than `COSDataCatalog` is specified, only the `SELECT` permission can be granted here.",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission type. Valid values: `ADMIN`, `DATASOURCE`, `DATABASE`, `TABLE`, `VIEW`, `FUNCTION`, `COLUMN`, and `ENGINE`. Note: If it is left empty, `ADMIN` is used.",
												},
												"function": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target function. `*` represents all functions in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any function.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target view. `*` represents all views in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any view.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target column. `*` represents all columns. To grant admin permissions, it must be `*`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"data_engine": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target data engine. `*` represents all engines. To grant admin permissions, it must be `*`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"re_auth": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the grantee is allowed to further grant the permissions. Valid values: `false` (default) and `true` (the grantee can grant permissions gained here to other sub-users).Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission source, which is not required when input parameters are passed in. Valid values: `USER` (from the user) and `WORKGROUP` (from one or more associated work groups).Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The grant mode, which is not required as an input parameter. Valid values: `COMMON` and `SENIOR`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The operator, which is not required as an input parameter.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission policy creation time, which is not required as an input parameter.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"source_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The ID of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The policy ID.Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total policiesNote: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"work_group_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information about collections of working groups bound to the userNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"work_group_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection of working group informationNote: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"work_group_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Unique ID of the working group.",
												},
												"work_group_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Working group name.",
												},
												"work_group_description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Working group descriptionNote: This field may return null, indicating that no valid values can be obtained.",
												},
												"creator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Creator.",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The creation time of the working group, e.g. at 16:19:32 on Jul 28, 2021.",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total working groupsNote: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"user_alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User aliasNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"row_filter_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of filtered rowsNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection of policiesNote: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target database. `*` represents all databases in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any database.",
												},
												"catalog": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target data source. To grant admin permission, it must be `*` (all resources at this level); to grant data source and database permissions, it must be `COSDataCatalog` or `*`; to grant table permissions, it can be a custom data source; if it is left empty, `DataLakeCatalog` is used. Note: To grant permissions on a custom data source, the permissions that can be managed in the Data Lake Compute console are subsets of the account permissions granted when you connect the data source to the console.",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target table. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.",
												},
												"operation": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The target permissions, which vary by permission level. Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`, `CREATE`, `ALTER`, and `DROP`; table: `ALL`, `SELECT`, `INSERT`, `ALTER`, `DELETE`, `DROP`, and `UPDATE`. Note: For table permissions, if a data source other than `COSDataCatalog` is specified, only the `SELECT` permission can be granted here.",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission type. Valid values: `ADMIN`, `DATASOURCE`, `DATABASE`, `TABLE`, `VIEW`, `FUNCTION`, `COLUMN`, and `ENGINE`. Note: If it is left empty, `ADMIN` is used.",
												},
												"function": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target function. `*` represents all functions in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any function.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target view. `*` represents all views in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any view.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target column. `*` represents all columns. To grant admin permissions, it must be `*`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"data_engine": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target data engine. `*` represents all engines. To grant admin permissions, it must be `*`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"re_auth": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the grantee is allowed to further grant the permissions. Valid values: `false` (default) and `true` (the grantee can grant permissions gained here to other sub-users).Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission source, which is not required when input parameters are passed in. Valid values: `USER` (from the user) and `WORKGROUP` (from one or more associated work groups).Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The grant mode, which is not required as an input parameter. Valid values: `COMMON` and `SENIOR`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The operator, which is not required as an input parameter.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission policy creation time, which is not required as an input parameter.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"source_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The ID of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The policy ID.Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total policiesNote: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"account_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account type.",
						},
						"catalog_policy_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of catalog permissionsNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection of policiesNote: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target database. `*` represents all databases in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any database.",
												},
												"catalog": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target data source. To grant admin permission, it must be `*` (all resources at this level); to grant data source and database permissions, it must be `COSDataCatalog` or `*`; to grant table permissions, it can be a custom data source; if it is left empty, `DataLakeCatalog` is used. Note: To grant permissions on a custom data source, the permissions that can be managed in the Data Lake Compute console are subsets of the account permissions granted when you connect the data source to the console.",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target table. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.",
												},
												"operation": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The target permissions, which vary by permission level. Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`, `CREATE`, `ALTER`, and `DROP`; table: `ALL`, `SELECT`, `INSERT`, `ALTER`, `DELETE`, `DROP`, and `UPDATE`. Note: For table permissions, if a data source other than `COSDataCatalog` is specified, only the `SELECT` permission can be granted here.",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission type. Valid values: `ADMIN`, `DATASOURCE`, `DATABASE`, `TABLE`, `VIEW`, `FUNCTION`, `COLUMN`, and `ENGINE`. Note: If it is left empty, `ADMIN` is used.",
												},
												"function": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target function. `*` represents all functions in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any function.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target view. `*` represents all views in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any view.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target column. `*` represents all columns. To grant admin permissions, it must be `*`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"data_engine": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the target data engine. `*` represents all engines. To grant admin permissions, it must be `*`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"re_auth": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the grantee is allowed to further grant the permissions. Valid values: `false` (default) and `true` (the grantee can grant permissions gained here to other sub-users).Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission source, which is not required when input parameters are passed in. Valid values: `USER` (from the user) and `WORKGROUP` (from one or more associated work groups).Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The grant mode, which is not required as an input parameter. Valid values: `COMMON` and `SENIOR`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The operator, which is not required as an input parameter.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission policy creation time, which is not required as an input parameter.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"source_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The ID of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The policy ID.Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total policiesNote: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDlcDescribeUserInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_describe_user_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var userId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(string)
		paramMap["UserId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*dlc.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := dlc.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sorting"); ok {
		paramMap["Sorting"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var userInfo *dlc.UserDetailInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDescribeUserInfoByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		userInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	userDetailInfoMap := map[string]interface{}{}

	if userInfo != nil {

		if userInfo.UserId != nil {
			userDetailInfoMap["user_id"] = userInfo.UserId
		}

		if userInfo.Type != nil {
			userDetailInfoMap["type"] = userInfo.Type
		}

		if userInfo.UserType != nil {
			userDetailInfoMap["user_type"] = userInfo.UserType
		}

		if userInfo.UserDescription != nil {
			userDetailInfoMap["user_description"] = userInfo.UserDescription
		}

		if userInfo.DataPolicyInfo != nil {
			dataPolicyInfoMap := map[string]interface{}{}

			if userInfo.DataPolicyInfo.PolicySet != nil {
				var policySetList []interface{}
				for _, policySet := range userInfo.DataPolicyInfo.PolicySet {
					policySetMap := map[string]interface{}{}

					if policySet.Database != nil {
						policySetMap["database"] = policySet.Database
					}

					if policySet.Catalog != nil {
						policySetMap["catalog"] = policySet.Catalog
					}

					if policySet.Table != nil {
						policySetMap["table"] = policySet.Table
					}

					if policySet.Operation != nil {
						policySetMap["operation"] = policySet.Operation
					}

					if policySet.PolicyType != nil {
						policySetMap["policy_type"] = policySet.PolicyType
					}

					if policySet.Function != nil {
						policySetMap["function"] = policySet.Function
					}

					if policySet.View != nil {
						policySetMap["view"] = policySet.View
					}

					if policySet.Column != nil {
						policySetMap["column"] = policySet.Column
					}

					if policySet.DataEngine != nil {
						policySetMap["data_engine"] = policySet.DataEngine
					}

					if policySet.ReAuth != nil {
						policySetMap["re_auth"] = policySet.ReAuth
					}

					if policySet.Source != nil {
						policySetMap["source"] = policySet.Source
					}

					if policySet.Mode != nil {
						policySetMap["mode"] = policySet.Mode
					}

					if policySet.Operator != nil {
						policySetMap["operator"] = policySet.Operator
					}

					if policySet.CreateTime != nil {
						policySetMap["create_time"] = policySet.CreateTime
					}

					if policySet.SourceId != nil {
						policySetMap["source_id"] = policySet.SourceId
					}

					if policySet.SourceName != nil {
						policySetMap["source_name"] = policySet.SourceName
					}

					if policySet.Id != nil {
						policySetMap["id"] = policySet.Id
					}

					policySetList = append(policySetList, policySetMap)
				}

				dataPolicyInfoMap["policy_set"] = policySetList
			}

			if userInfo.DataPolicyInfo.TotalCount != nil {
				dataPolicyInfoMap["total_count"] = userInfo.DataPolicyInfo.TotalCount
			}

			userDetailInfoMap["data_policy_info"] = []interface{}{dataPolicyInfoMap}
		}

		if userInfo.EnginePolicyInfo != nil {
			enginePolicyInfoMap := map[string]interface{}{}

			if userInfo.EnginePolicyInfo.PolicySet != nil {
				var policySetList []interface{}
				for _, policySet := range userInfo.EnginePolicyInfo.PolicySet {
					policySetMap := map[string]interface{}{}

					if policySet.Database != nil {
						policySetMap["database"] = policySet.Database
					}

					if policySet.Catalog != nil {
						policySetMap["catalog"] = policySet.Catalog
					}

					if policySet.Table != nil {
						policySetMap["table"] = policySet.Table
					}

					if policySet.Operation != nil {
						policySetMap["operation"] = policySet.Operation
					}

					if policySet.PolicyType != nil {
						policySetMap["policy_type"] = policySet.PolicyType
					}

					if policySet.Function != nil {
						policySetMap["function"] = policySet.Function
					}

					if policySet.View != nil {
						policySetMap["view"] = policySet.View
					}

					if policySet.Column != nil {
						policySetMap["column"] = policySet.Column
					}

					if policySet.DataEngine != nil {
						policySetMap["data_engine"] = policySet.DataEngine
					}

					if policySet.ReAuth != nil {
						policySetMap["re_auth"] = policySet.ReAuth
					}

					if policySet.Source != nil {
						policySetMap["source"] = policySet.Source
					}

					if policySet.Mode != nil {
						policySetMap["mode"] = policySet.Mode
					}

					if policySet.Operator != nil {
						policySetMap["operator"] = policySet.Operator
					}

					if policySet.CreateTime != nil {
						policySetMap["create_time"] = policySet.CreateTime
					}

					if policySet.SourceId != nil {
						policySetMap["source_id"] = policySet.SourceId
					}

					if policySet.SourceName != nil {
						policySetMap["source_name"] = policySet.SourceName
					}

					if policySet.Id != nil {
						policySetMap["id"] = policySet.Id
					}

					policySetList = append(policySetList, policySetMap)
				}

				enginePolicyInfoMap["policy_set"] = policySetList
			}

			if userInfo.EnginePolicyInfo.TotalCount != nil {
				enginePolicyInfoMap["total_count"] = userInfo.EnginePolicyInfo.TotalCount
			}

			userDetailInfoMap["engine_policy_info"] = []interface{}{enginePolicyInfoMap}
		}

		if userInfo.WorkGroupInfo != nil {
			workGroupInfoMap := map[string]interface{}{}

			if userInfo.WorkGroupInfo.WorkGroupSet != nil {
				var workGroupSetList []interface{}
				for _, workGroupSet := range userInfo.WorkGroupInfo.WorkGroupSet {
					workGroupSetMap := map[string]interface{}{}

					if workGroupSet.WorkGroupId != nil {
						workGroupSetMap["work_group_id"] = workGroupSet.WorkGroupId
					}

					if workGroupSet.WorkGroupName != nil {
						workGroupSetMap["work_group_name"] = workGroupSet.WorkGroupName
					}

					if workGroupSet.WorkGroupDescription != nil {
						workGroupSetMap["work_group_description"] = workGroupSet.WorkGroupDescription
					}

					if workGroupSet.Creator != nil {
						workGroupSetMap["creator"] = workGroupSet.Creator
					}

					if workGroupSet.CreateTime != nil {
						workGroupSetMap["create_time"] = workGroupSet.CreateTime
					}

					workGroupSetList = append(workGroupSetList, workGroupSetMap)
				}

				workGroupInfoMap["work_group_set"] = workGroupSetList
			}

			if userInfo.WorkGroupInfo.TotalCount != nil {
				workGroupInfoMap["total_count"] = userInfo.WorkGroupInfo.TotalCount
			}

			userDetailInfoMap["work_group_info"] = []interface{}{workGroupInfoMap}
		}

		if userInfo.UserAlias != nil {
			userDetailInfoMap["user_alias"] = userInfo.UserAlias
		}

		if userInfo.RowFilterInfo != nil {
			rowFilterInfoMap := map[string]interface{}{}

			if userInfo.RowFilterInfo.PolicySet != nil {
				var policySetList []interface{}
				for _, policySet := range userInfo.RowFilterInfo.PolicySet {
					policySetMap := map[string]interface{}{}

					if policySet.Database != nil {
						policySetMap["database"] = policySet.Database
					}

					if policySet.Catalog != nil {
						policySetMap["catalog"] = policySet.Catalog
					}

					if policySet.Table != nil {
						policySetMap["table"] = policySet.Table
					}

					if policySet.Operation != nil {
						policySetMap["operation"] = policySet.Operation
					}

					if policySet.PolicyType != nil {
						policySetMap["policy_type"] = policySet.PolicyType
					}

					if policySet.Function != nil {
						policySetMap["function"] = policySet.Function
					}

					if policySet.View != nil {
						policySetMap["view"] = policySet.View
					}

					if policySet.Column != nil {
						policySetMap["column"] = policySet.Column
					}

					if policySet.DataEngine != nil {
						policySetMap["data_engine"] = policySet.DataEngine
					}

					if policySet.ReAuth != nil {
						policySetMap["re_auth"] = policySet.ReAuth
					}

					if policySet.Source != nil {
						policySetMap["source"] = policySet.Source
					}

					if policySet.Mode != nil {
						policySetMap["mode"] = policySet.Mode
					}

					if policySet.Operator != nil {
						policySetMap["operator"] = policySet.Operator
					}

					if policySet.CreateTime != nil {
						policySetMap["create_time"] = policySet.CreateTime
					}

					if policySet.SourceId != nil {
						policySetMap["source_id"] = policySet.SourceId
					}

					if policySet.SourceName != nil {
						policySetMap["source_name"] = policySet.SourceName
					}

					if policySet.Id != nil {
						policySetMap["id"] = policySet.Id
					}

					policySetList = append(policySetList, policySetMap)
				}

				rowFilterInfoMap["policy_set"] = policySetList
			}

			if userInfo.RowFilterInfo.TotalCount != nil {
				rowFilterInfoMap["total_count"] = userInfo.RowFilterInfo.TotalCount
			}

			userDetailInfoMap["row_filter_info"] = []interface{}{rowFilterInfoMap}
		}

		if userInfo.AccountType != nil {
			userDetailInfoMap["account_type"] = userInfo.AccountType
		}

		if userInfo.CatalogPolicyInfo != nil {
			catalogPolicyInfoMap := map[string]interface{}{}

			if userInfo.CatalogPolicyInfo.PolicySet != nil {
				policySetList := []interface{}{}
				for _, policySet := range userInfo.CatalogPolicyInfo.PolicySet {
					policySetMap := map[string]interface{}{}

					if policySet.Database != nil {
						policySetMap["database"] = policySet.Database
					}

					if policySet.Catalog != nil {
						policySetMap["catalog"] = policySet.Catalog
					}

					if policySet.Table != nil {
						policySetMap["table"] = policySet.Table
					}

					if policySet.Operation != nil {
						policySetMap["operation"] = policySet.Operation
					}

					if policySet.PolicyType != nil {
						policySetMap["policy_type"] = policySet.PolicyType
					}

					if policySet.Function != nil {
						policySetMap["function"] = policySet.Function
					}

					if policySet.View != nil {
						policySetMap["view"] = policySet.View
					}

					if policySet.Column != nil {
						policySetMap["column"] = policySet.Column
					}

					if policySet.DataEngine != nil {
						policySetMap["data_engine"] = policySet.DataEngine
					}

					if policySet.ReAuth != nil {
						policySetMap["re_auth"] = policySet.ReAuth
					}

					if policySet.Source != nil {
						policySetMap["source"] = policySet.Source
					}

					if policySet.Mode != nil {
						policySetMap["mode"] = policySet.Mode
					}

					if policySet.Operator != nil {
						policySetMap["operator"] = policySet.Operator
					}

					if policySet.CreateTime != nil {
						policySetMap["create_time"] = policySet.CreateTime
					}

					if policySet.SourceId != nil {
						policySetMap["source_id"] = policySet.SourceId
					}

					if policySet.SourceName != nil {
						policySetMap["source_name"] = policySet.SourceName
					}

					if policySet.Id != nil {
						policySetMap["id"] = policySet.Id
					}

					policySetList = append(policySetList, policySetMap)
				}

				catalogPolicyInfoMap["policy_set"] = policySetList
			}

			if userInfo.CatalogPolicyInfo.TotalCount != nil {
				catalogPolicyInfoMap["total_count"] = userInfo.CatalogPolicyInfo.TotalCount
			}

			userDetailInfoMap["catalog_policy_info"] = []interface{}{catalogPolicyInfoMap}
		}

		_ = d.Set("user_info", []interface{}{userDetailInfoMap})
	}

	d.SetId(userId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), userDetailInfoMap); e != nil {
			return e
		}
	}
	return nil
}
