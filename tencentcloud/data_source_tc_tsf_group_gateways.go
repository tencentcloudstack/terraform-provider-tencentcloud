/*
Use this data source to query detailed information of tsf group_gateways

Example Usage

```hcl
data "tencentcloud_tsf_group_gateways" "group_gateways" {
  gateway_deploy_group_id = "group-aeoej4qy"
  search_word = "test"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTsfGroupGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfGroupGatewaysRead,
		Schema: map[string]*schema.Schema{
			"gateway_deploy_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway group Id.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "search word.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "api group information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "total count.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "api group Info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "api group id.Note: This field may return null, which means no valid value was found.",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "api group name.Note: This field may return null, which means no valid value was found.",
									},
									"group_context": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "api group context.Note: This field may return null, which means no valid value was found.",
									},
									"auth_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication type. secret: key authentication; none: no authentication.Note: This field may return null, which means no valid value was found.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Release status. drafted: not released. released: released.Note: This field may return null, which means no valid value was found.",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group creation time, such as: 2019-06-20 15:51:28.Note: This field may return null, which means no valid value was found.",
									},
									"updated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group update time, such as: 2019-06-20 15:51:28.Note: This field may return null, which means no valid value was found.",
									},
									"binded_gateway_deploy_groups": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Gateway deployment group bound to the API group.Note: This field may return null, which means no valid value was found.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"deploy_group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Gateway deployment group ID.Note: This field may return null, which means no valid value was found.",
												},
												"deploy_group_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Gateway deployment group name.Note: This field may return null, which means no valid value was found.",
												},
												"application_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "application ID.Note: This field may return null, which means no valid value was found.",
												},
												"application_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "application name.Note: This field may return null, which means no valid value was found.",
												},
												"application_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application category: V: virtual machine application, C: container application.Note: This field may return null, which means no valid value was found.",
												},
												"group_status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application status of the deployment group, with possible values: Running, Waiting, Paused, Updating, RollingBack, Abnormal, Unknown.Note: This field may return null, which means no valid value was found.",
												},
												"cluster_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Cluster type, with possible values: C: container, V: virtual machine.Note: This field may return null, which means no valid value was found.",
												},
											},
										},
									},
									"api_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of APIs.Note: This field may return null, which means no valid value was found.",
									},
									"acl_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ACL type for accessing the group.Note: This field may return null, which means no valid value was found.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description.Note: This field may return null, which means no valid value was found.",
									},
									"group_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group type. ms: microservice group; external: external API group.This field may return null, which means no valid value was found.",
									},
									"gateway_instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway instance type.Note: This field may return null, which means no valid value was found.",
									},
									"gateway_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway instance ID.Note: This field may return null, which means no valid value was found.",
									},
									"namespace_name_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace parameter key.Note: This field may return null, which means no valid value was found.",
									},
									"service_name_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice name parameter key.Note: This field may return null, which means no valid value was found.",
									},
									"namespace_name_key_position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace parameter location, path, header, or query. The default is path.Note: This field may return null, which means no valid value was found.",
									},
									"service_name_key_position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice name parameter location, path, header, or query. The default is path.Note: This field may return null, which means no valid value was found.",
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

func dataSourceTencentCloudTsfGroupGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_group_gateways.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_deploy_group_id"); ok {
		paramMap["GatewayDeployGroupId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *tsf.TsfPageApiGroupInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTsfGroupGatewaysByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result.Content))
	tsfPageApiGroupInfoMap := map[string]interface{}{}
	if result != nil {

		if result.TotalCount != nil {
			tsfPageApiGroupInfoMap["total_count"] = result.TotalCount
		}

		if result.Content != nil {
			contentList := []interface{}{}
			for _, content := range result.Content {
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
