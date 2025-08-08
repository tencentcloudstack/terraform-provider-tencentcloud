package dlc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcDescribeWorkGroupInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeWorkGroupInfoRead,
		Schema: map[string]*schema.Schema{
			"work_group_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Working group ID.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Types of queried information. User: user information; DataAuth: data permissions; EngineAuth: engine permissions.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter criteria that are queriedWhen the type is User, the fuzzy search is supported as the key is user-name.When the type is DataAuth, the keys supported are:policy-type: types of permissions;policy-source: data sources;data-name: fuzzy search of the database and table.When the type is EngineAuth, the keys supported are:policy-type: types of permissions;policy-source: data sources;engine-name: fuzzy search of the database and table.",
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
				Description: "Sort fields.When the type is User, create-time and user-name are supported.When the type is DataAuth, create-time is supported.When the type is EngineAuth, create-time is supported.",
			},

			"sorting": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting methods: desc means in order; asc means in reverse order; it is asc by default.",
			},

			"work_group_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Details about working groupsNote: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"work_group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Working group IDNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"work_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Working group nameNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of information included. User: user information; DataAuth: data permissions; EngineAuth: engine permissionsNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"user_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of users bound to working groupsNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection of user informationNote: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "User Id which matches the sub-user UIN on the CAM side.",
												},
												"user_description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "User descriptionNote: The returned value of this field may be null, indicating that no valid value is obtained.",
												},
												"creator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The creator of the current user.",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The creation time of the current user, e.g. 16:19:32, July 28, 2021.",
												},
												"user_alias": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "User alias.",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total usersNote: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"data_policy_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of data permissionsNote: This field may return null, indicating that no valid values can be obtained.",
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
						"work_group_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Working group descriptionNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"row_filter_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of information about filtered rowsNote: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudDlcDescribeWorkGroupInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_describe_work_group_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOkExists("work_group_id"); v != nil {
		paramMap["WorkGroupId"] = helper.IntInt64(v.(int))
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

	var workGroupInfo *dlc.WorkGroupDetailInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDescribeWorkGroupInfoByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		workGroupInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	workGroupDetailInfoMap := map[string]interface{}{}

	if workGroupInfo != nil {

		if workGroupInfo.WorkGroupId != nil {
			workGroupDetailInfoMap["work_group_id"] = workGroupInfo.WorkGroupId
		}

		if workGroupInfo.WorkGroupName != nil {
			workGroupDetailInfoMap["work_group_name"] = workGroupInfo.WorkGroupName
		}

		if workGroupInfo.Type != nil {
			workGroupDetailInfoMap["type"] = workGroupInfo.Type
		}

		if workGroupInfo.UserInfo != nil {
			userInfoMap := map[string]interface{}{}

			if workGroupInfo.UserInfo.UserSet != nil {
				var userSetList []interface{}
				for _, userSet := range workGroupInfo.UserInfo.UserSet {
					userSetMap := map[string]interface{}{}

					if userSet.UserId != nil {
						userSetMap["user_id"] = userSet.UserId
					}

					if userSet.UserDescription != nil {
						userSetMap["user_description"] = userSet.UserDescription
					}

					if userSet.Creator != nil {
						userSetMap["creator"] = userSet.Creator
					}

					if userSet.CreateTime != nil {
						userSetMap["create_time"] = userSet.CreateTime
					}

					if userSet.UserAlias != nil {
						userSetMap["user_alias"] = userSet.UserAlias
					}

					userSetList = append(userSetList, userSetMap)
				}

				userInfoMap["user_set"] = userSetList
			}

			if workGroupInfo.UserInfo.TotalCount != nil {
				userInfoMap["total_count"] = workGroupInfo.UserInfo.TotalCount
			}

			workGroupDetailInfoMap["user_info"] = []interface{}{userInfoMap}
		}

		if workGroupInfo.DataPolicyInfo != nil {
			dataPolicyInfoMap := map[string]interface{}{}

			if workGroupInfo.DataPolicyInfo.PolicySet != nil {
				var policySetList []interface{}
				for _, policySet := range workGroupInfo.DataPolicyInfo.PolicySet {
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

			if workGroupInfo.DataPolicyInfo.TotalCount != nil {
				dataPolicyInfoMap["total_count"] = workGroupInfo.DataPolicyInfo.TotalCount
			}

			workGroupDetailInfoMap["data_policy_info"] = []interface{}{dataPolicyInfoMap}
		}

		if workGroupInfo.EnginePolicyInfo != nil {
			enginePolicyInfoMap := map[string]interface{}{}

			if workGroupInfo.EnginePolicyInfo.PolicySet != nil {
				var policySetList []interface{}
				for _, policySet := range workGroupInfo.EnginePolicyInfo.PolicySet {
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

			if workGroupInfo.EnginePolicyInfo.TotalCount != nil {
				enginePolicyInfoMap["total_count"] = workGroupInfo.EnginePolicyInfo.TotalCount
			}

			workGroupDetailInfoMap["engine_policy_info"] = []interface{}{enginePolicyInfoMap}
		}

		if workGroupInfo.WorkGroupDescription != nil {
			workGroupDetailInfoMap["work_group_description"] = workGroupInfo.WorkGroupDescription
		}

		if workGroupInfo.RowFilterInfo != nil {
			rowFilterInfoMap := map[string]interface{}{}

			if workGroupInfo.RowFilterInfo.PolicySet != nil {
				var policySetList []interface{}
				for _, policySet := range workGroupInfo.RowFilterInfo.PolicySet {
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

			if workGroupInfo.RowFilterInfo.TotalCount != nil {
				rowFilterInfoMap["total_count"] = workGroupInfo.RowFilterInfo.TotalCount
			}

			workGroupDetailInfoMap["row_filter_info"] = []interface{}{rowFilterInfoMap}
		}

		ids = append(ids, helper.Int64ToStr(*workGroupInfo.WorkGroupId))
		_ = d.Set("work_group_info", []interface{}{workGroupDetailInfoMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), workGroupDetailInfoMap); e != nil {
			return e
		}
	}
	return nil
}
