package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTsfApiGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfApiGroupRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "search word.",
			},

			"group_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Group type. ms: Microservice group; external: External API group.",
			},

			"auth_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Authentication type. secret: Secret key authentication; none: No authentication.",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Publishing status. drafted: Not published. released: Published.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting field: created_time or group_context.",
			},

			"order_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Sorting type: 0 (ASC) or 1 (DESC).",
			},

			"gateway_instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Gateway Instance Id.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Pagination structure.Note: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "record count.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Api group info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Api Group Id.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Api Group Name.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"group_context": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Api Group Context.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"auth_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication type. secret: key authentication; none: no authentication.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Release status. drafted: not released. released: released.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group creation time.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"updated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group creation time, such as: 2019-06-20 15:51:28.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"binded_gateway_deploy_groups": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The gateway group bind with the api group list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"deploy_group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Gateway deployment group bound to the API group.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"deploy_group_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Deploy group name.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"application_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application ID.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"application_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application Name.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"application_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application Name.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"group_status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application category: V: virtual machine application, C: container application. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"cluster_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Cluster type, C: container, V: virtual machine.Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"api_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "api count.",
									},
									"acl_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Number of APIs.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"group_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group type.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"gateway_instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway Instance Type.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"gateway_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway Instance Id.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"namespace_name_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace name key.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"service_name_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Key value of microservice name parameter.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"namespace_name_key_position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace parameter location, path, header, or query, default is path. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"service_name_key_position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice name parameter location, path, header, or query, default is path.Note: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudTsfApiGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_api_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_type"); ok {
		paramMap["GroupType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auth_type"); ok {
		paramMap["AuthType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("order_type"); v != nil {
		paramMap["OrderType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("gateway_instance_id"); ok {
		paramMap["GatewayInstanceId"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var apiGroupInfo *tsf.TsfPageApiGroupInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfApiGroupByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		apiGroupInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(apiGroupInfo.Content))
	tsfPageApiGroupInfoMap := map[string]interface{}{}
	if apiGroupInfo != nil {

		if apiGroupInfo.TotalCount != nil {
			tsfPageApiGroupInfoMap["total_count"] = apiGroupInfo.TotalCount
		}

		if apiGroupInfo.Content != nil {
			contentList := []interface{}{}
			for _, content := range apiGroupInfo.Content {
				contentMap := map[string]interface{}{}

				if content.GroupId != nil {
					contentMap["group_id"] = content.GroupId
				}

				if content.GroupName != nil {
					contentMap["group_name"] = content.GroupName
				}

				if content.GroupContext != nil {
					contentMap["group_context"] = content.GroupContext
				}

				if content.AuthType != nil {
					contentMap["auth_type"] = content.AuthType
				}

				if content.Status != nil {
					contentMap["status"] = content.Status
				}

				if content.CreatedTime != nil {
					contentMap["created_time"] = content.CreatedTime
				}

				if content.UpdatedTime != nil {
					contentMap["updated_time"] = content.UpdatedTime
				}

				if content.BindedGatewayDeployGroups != nil {
					bindedGatewayDeployGroupsList := []interface{}{}
					for _, bindedGatewayDeployGroups := range content.BindedGatewayDeployGroups {
						bindedGatewayDeployGroupsMap := map[string]interface{}{}

						if bindedGatewayDeployGroups.DeployGroupId != nil {
							bindedGatewayDeployGroupsMap["deploy_group_id"] = bindedGatewayDeployGroups.DeployGroupId
						}

						if bindedGatewayDeployGroups.DeployGroupName != nil {
							bindedGatewayDeployGroupsMap["deploy_group_name"] = bindedGatewayDeployGroups.DeployGroupName
						}

						if bindedGatewayDeployGroups.ApplicationId != nil {
							bindedGatewayDeployGroupsMap["application_id"] = bindedGatewayDeployGroups.ApplicationId
						}

						if bindedGatewayDeployGroups.ApplicationName != nil {
							bindedGatewayDeployGroupsMap["application_name"] = bindedGatewayDeployGroups.ApplicationName
						}

						if bindedGatewayDeployGroups.ApplicationType != nil {
							bindedGatewayDeployGroupsMap["application_type"] = bindedGatewayDeployGroups.ApplicationType
						}

						if bindedGatewayDeployGroups.GroupStatus != nil {
							bindedGatewayDeployGroupsMap["group_status"] = bindedGatewayDeployGroups.GroupStatus
						}

						if bindedGatewayDeployGroups.ClusterType != nil {
							bindedGatewayDeployGroupsMap["cluster_type"] = bindedGatewayDeployGroups.ClusterType
						}

						bindedGatewayDeployGroupsList = append(bindedGatewayDeployGroupsList, bindedGatewayDeployGroupsMap)
					}

					contentMap["binded_gateway_deploy_groups"] = bindedGatewayDeployGroupsList
				}

				if content.ApiCount != nil {
					contentMap["api_count"] = content.ApiCount
				}

				if content.AclMode != nil {
					contentMap["acl_mode"] = content.AclMode
				}

				if content.Description != nil {
					contentMap["description"] = content.Description
				}

				if content.GroupType != nil {
					contentMap["group_type"] = content.GroupType
				}

				if content.GatewayInstanceType != nil {
					contentMap["gateway_instance_type"] = content.GatewayInstanceType
				}

				if content.GatewayInstanceId != nil {
					contentMap["gateway_instance_id"] = content.GatewayInstanceId
				}

				if content.NamespaceNameKey != nil {
					contentMap["namespace_name_key"] = content.NamespaceNameKey
				}

				if content.ServiceNameKey != nil {
					contentMap["service_name_key"] = content.ServiceNameKey
				}

				if content.NamespaceNameKeyPosition != nil {
					contentMap["namespace_name_key_position"] = content.NamespaceNameKeyPosition
				}

				if content.ServiceNameKeyPosition != nil {
					contentMap["service_name_key_position"] = content.ServiceNameKeyPosition
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.GroupId)
			}

			tsfPageApiGroupInfoMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageApiGroupInfoMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tsfPageApiGroupInfoMap); e != nil {
			return e
		}
	}
	return nil
}
