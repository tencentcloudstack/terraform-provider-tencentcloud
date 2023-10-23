/*
Use this data source to query detailed information of dlc describe_user_info

Example Usage

```hcl
data "tencentcloud_dlc_describe_user_info" "describe_user_info" {
  user_id = "100032772113"
  type = "Group"
  sort_by = "create-time"
  sorting = "desc"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDlcDescribeUserInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeUserInfoRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User id, the same as the sub-user uin.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Query information type, Group: work group DataAuth: data permission EngineAuth: engine permission.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Query filter conditions. when type is Group, fuzzy search with Key as workgroup-name is supported. when type is DataAuth, key is supported. policy-type: permission type, policy-source: data source, data-name: database table. Fuzzy search, when type is EngineAuth, supports fuzzy search of key, policy-type: permission type, policy-source: data source, engine-name: library table.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attribute name. If there are multiple Filters, the relationship between Filters is a logical OR (OR) relationship.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Attribute value, if there are multiple Values in the same filter, the relationship between values under the same filter is a logical OR relationship.",
						},
					},
				},
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting field, when type is Group, support create-time, group-name, when type is DataAuth, support create-time, when type is EngineAuth, support create-time.",
			},

			"sorting": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting method, desc means forward order, asc means reverse order, the default is asc.",
			},

			"user_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "User details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User id, the same as the sub-user uin.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of information returned, Group: the returned workgroup information of the current user; DataAuth: the returned data permission information of the current user; EngineAuth: the returned engine permission information of the current user.",
						},
						"user_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User type: ADMIN: Administrator COMMON: General user.",
						},
						"user_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User description.",
						},
						"data_policy_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Data permission information collection.",
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
						"work_group_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Workgroup collection information bound to this user.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"work_group_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Work group set.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"work_group_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Work group unique id.",
												},
												"work_group_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Work group name.",
												},
												"work_group_description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Work group description.",
												},
												"creator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Creator.",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time the workgroup was created.",
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
						"user_alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User alias.",
						},
						"row_filter_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Row filter collection.",
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

func dataSourceTencentCloudDlcDescribeUserInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dlc_describe_user_info.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
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
		paramMap["filters"] = tmpSet
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sorting"); ok {
		paramMap["Sorting"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var userInfo *dlc.UserDetailInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDescribeUserInfoByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
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

				dataPolicyInfoMap["policy_set"] = []interface{}{policySetList}
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

				rowFilterInfoMap["policy_set"] = []interface{}{policySetList}
			}

			if userInfo.RowFilterInfo.TotalCount != nil {
				rowFilterInfoMap["total_count"] = userInfo.RowFilterInfo.TotalCount
			}

			userDetailInfoMap["row_filter_info"] = []interface{}{rowFilterInfoMap}
		}

		_ = d.Set("user_info", []interface{}{userDetailInfoMap})
	}

	d.SetId(userId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), userDetailInfoMap); e != nil {
			return e
		}
	}
	return nil
}
