package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDlcDescribeWorkGroupInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeWorkGroupInfoRead,
		Schema: map[string]*schema.Schema{
			"work_group_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Work group id.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Query information type, only support: User: user information/DataAuth: data permission/EngineAuth: engine permission.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Query filter conditions. when Type is User, fuzzy search with Key as user-name is supported; when Type is DataAuth, key is supported; policy-type: permission type; policy-source: data source; data-name: database table fuzzy search; when Type is EngineAuth, supports key; policy-type: permission type; policy-source: data source; engine-name: fuzzy search of library tables.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attribute name. If there are multiple Filters, the relationship between filters is a logical or (OR) relationship.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Attribute value, if there are multiple values in the same filter, the relationship between values under the same filter is a logical or relationship.",
						},
					},
				},
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting fields, when Type is User, support create-time, user-name, when type is DataAuth, support create-time, when type is EngineAuth, support create-time.",
			},

			"sorting": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting method, desc means forward order, asc means reverse order, the default is asc.",
			},

			"work_group_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Workgroup details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"work_group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Work group id.",
						},
						"work_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Work group name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of information contained. User: user information; DataAuth: data permissions; EngineAuth: engine permissions.",
						},
						"user_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A collection of users bound to the workgroup.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "User information collection.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "User id, matches the CAM side sub-user uin.",
												},
												"user_description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "User description.",
												},
												"creator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The creator of the current user.",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Create time.",
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
										Description: "Total count.",
									},
								},
							},
						},
						"data_policy_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Data permission collection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Policy set.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Database name that requires authorization, fill in * to represent all databases under the current catalog. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level, only blanks are allowed to be filled in. For other types, the database can be specified arbitrarily.",
												},
												"catalog": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For the data source name that requires authorization, only * (representing all resources at this level) is supported under the administrator level; in the case of data source level and database level authentication, only COSDataCatalog or * is supported; in data table level authentication, it is possible Fill in the user-defined data source. If left blank, it defaults to DataLakeCatalog. note: If a user-defined data source is authenticated, the permissions that dlc can manage are a subset of the accounts provided by the user when accessing the data source.",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For the table name that requires authorization, fill in * to represent all tables under the current database. when the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. For other types, data tables can be specified arbitrarily.",
												},
												"operation": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Authorized permission operations provide different operations for different levels of authentication. administrator permissions: ALL, default is ALL if left blank; data connection level authentication: CREATE; database level authentication: ALL, CREATE, ALTER, DROP; data table permissions: ALL, SELECT, INSERT, ALTER, DELETE, DROP, UPDATE. note: under data table permissions, only SELECT operations are supported when the specified data source is not COSDataCatalog.",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Authorization type, currently supports eight authorization types: ADMIN: Administrator level authentication DATASOURCE: data connection level authentication DATABASE: database level authentication TABLE: Table level authentication VIEW: view level authentication FUNCTION: Function level authentication COLUMN: Column level authentication ENGINE: Data engine authentication. if left blank, the default is administrator level authentication.",
												},
												"function": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For the function name that requires authorization, fill in * to represent all functions under the current catalog. when the authorization type is administrator level, only * is allowed to be filled in. When the authorization type is data connection level, only blanks are allowed to be filled in. in other types, functions can be specified arbitrarily.",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For views that require authorization, fill in * to represent all views under the current database. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. for other types, the view can be specified arbitrarily.",
												},
												"column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For columns that require authorization, fill in * to represent all current columns. When the authorization type is administrator level, only * is allowed.",
												},
												"data_engine": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Data engines that require authorization, fill in * to represent all current engines. when the authorization type is administrator level, only * is allowed.",
												},
												"re_auth": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the user can perform secondary authorization. when it is true, the authorized user can re-authorize the permissions obtained this time to other sub-users. default is false.",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Permission source, please leave it blank. USER: permissions come from the user itself; WORKGROUP: permissions come from the bound workgroup.",
												},
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Authorization mode, please leave this parameter blank. COMMON: normal mode; SENIOR: advanced mode.",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Operator, do not fill in the input parameters.",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time when the permission was created. Leave the input parameter blank.",
												},
												"source_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The id of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the Source field is WORKGROUP.",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the source field is WORKGROUP.",
												},
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Policy id.",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total count.",
									},
								},
							},
						},
						"engine_policy_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Engine permission collection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Policy set.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Database name that requires authorization, fill in * to represent all databases under the current catalog. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level, only blanks are allowed to be filled in. For other types, the database can be specified arbitrarily.",
												},
												"catalog": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For the data source name that requires authorization, only * (representing all resources at this level) is supported under the administrator level; in the case of data source level and database level authentication, only COSDataCatalog or * is supported; in data table level authentication, it is possible Fill in the user-defined data source. If left blank, it defaults to DataLakeCatalog. note: If a user-defined data source is authenticated, the permissions that dlc can manage are a subset of the accounts provided by the user when accessing the data source.",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For the table name that requires authorization, fill in * to represent all tables under the current database. when the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. For other types, data tables can be specified arbitrarily.",
												},
												"operation": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Authorized permission operations provide different operations for different levels of authentication. administrator permissions: ALL, default is ALL if left blank; data connection level authentication: CREATE; database level authentication: ALL, CREATE, ALTER, DROP; data table permissions: ALL, SELECT, INSERT, ALTER, DELETE, DROP, UPDATE. note: under data table permissions, only SELECT operations are supported when the specified data source is not COSDataCatalog.",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Authorization type, currently supports eight authorization types: ADMIN: Administrator level authentication DATASOURCE: data connection level authentication DATABASE: database level authentication TABLE: Table level authentication VIEW: view level authentication FUNCTION: Function level authentication COLUMN: Column level authentication ENGINE: Data engine authentication. if left blank, the default is administrator level authentication.",
												},
												"function": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For the function name that requires authorization, fill in * to represent all functions under the current catalog. when the authorization type is administrator level, only * is allowed to be filled in. When the authorization type is data connection level, only blanks are allowed to be filled in. in other types, functions can be specified arbitrarily.",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For views that require authorization, fill in * to represent all views under the current database. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. for other types, the view can be specified arbitrarily.",
												},
												"column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For columns that require authorization, fill in * to represent all current columns. When the authorization type is administrator level, only * is allowed.",
												},
												"data_engine": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Data engines that require authorization, fill in * to represent all current engines. when the authorization type is administrator level, only * is allowed.",
												},
												"re_auth": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the user can perform secondary authorization. when it is true, the authorized user can re-authorize the permissions obtained this time to other sub-users. default is false.",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Permission source, please leave it blank. USER: permissions come from the user itself; WORKGROUP: permissions come from the bound workgroup.",
												},
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Authorization mode, please leave this parameter blank. COMMON: normal mode; SENIOR: advanced mode.",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Operator, do not fill in the input parameters.",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time when the permission was created. Leave the input parameter blank.",
												},
												"source_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The id of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the Source field is WORKGROUP.",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the source field is WORKGROUP.",
												},
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Policy id.",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total count.",
									},
								},
							},
						},
						"work_group_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workgroup description information.",
						},
						"row_filter_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Row filter information collection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Policy set.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Database name that requires authorization, fill in * to represent all databases under the current catalog. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level, only blanks are allowed to be filled in. For other types, the database can be specified arbitrarily.",
												},
												"catalog": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For the data source name that requires authorization, only * (representing all resources at this level) is supported under the administrator level; in the case of data source level and database level authentication, only COSDataCatalog or * is supported; in data table level authentication, it is possible Fill in the user-defined data source. If left blank, it defaults to DataLakeCatalog. note: If a user-defined data source is authenticated, the permissions that dlc can manage are a subset of the accounts provided by the user when accessing the data source.",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For the table name that requires authorization, fill in * to represent all tables under the current database. when the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. For other types, data tables can be specified arbitrarily.",
												},
												"operation": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Authorized permission operations provide different operations for different levels of authentication. administrator permissions: ALL, default is ALL if left blank; data connection level authentication: CREATE; database level authentication: ALL, CREATE, ALTER, DROP; data table permissions: ALL, SELECT, INSERT, ALTER, DELETE, DROP, UPDATE. note: under data table permissions, only SELECT operations are supported when the specified data source is not COSDataCatalog.",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Authorization type, currently supports eight authorization types: ADMIN: Administrator level authentication DATASOURCE: data connection level authentication DATABASE: database level authentication TABLE: Table level authentication VIEW: view level authentication FUNCTION: Function level authentication COLUMN: Column level authentication ENGINE: Data engine authentication. if left blank, the default is administrator level authentication.",
												},
												"function": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For the function name that requires authorization, fill in * to represent all functions under the current catalog. when the authorization type is administrator level, only * is allowed to be filled in. When the authorization type is data connection level, only blanks are allowed to be filled in. in other types, functions can be specified arbitrarily.",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For views that require authorization, fill in * to represent all views under the current database. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. for other types, the view can be specified arbitrarily.",
												},
												"column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "For columns that require authorization, fill in * to represent all current columns. When the authorization type is administrator level, only * is allowed.",
												},
												"data_engine": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Data engines that require authorization, fill in * to represent all current engines. when the authorization type is administrator level, only * is allowed.",
												},
												"re_auth": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the user can perform secondary authorization. when it is true, the authorized user can re-authorize the permissions obtained this time to other sub-users. default is false.",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Permission source, please leave it blank. USER: permissions come from the user itself; WORKGROUP: permissions come from the bound workgroup.",
												},
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Authorization mode, please leave this parameter blank. COMMON: normal mode; SENIOR: advanced mode.",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Operator, do not fill in the input parameters.",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time when the permission was created. Leave the input parameter blank.",
												},
												"source_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The id of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the Source field is WORKGROUP.",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the source field is WORKGROUP.",
												},
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Policy id.",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total count.",
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
	defer logElapsed("data_source.tencentcloud_dlc_describe_work_group_info.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
		paramMap["filters"] = tmpSet
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sorting"); ok {
		paramMap["Sorting"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var workGroupInfo *dlc.WorkGroupDetailInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDescribeWorkGroupInfoByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
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
		if e := writeToFile(output.(string), workGroupDetailInfoMap); e != nil {
			return e
		}
	}
	return nil
}
