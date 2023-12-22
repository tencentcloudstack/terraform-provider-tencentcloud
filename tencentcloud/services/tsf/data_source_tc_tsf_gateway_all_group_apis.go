package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfGatewayAllGroupApis() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfGatewayAllGroupApisRead,
		Schema: map[string]*schema.Schema{
			"gateway_deploy_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway group Id.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search keyword, supports api group name or API path.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Gateway group and API list information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway_deploy_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway deployment group ID.Note: This field may return null, which means no valid value was found.",
						},
						"gateway_deploy_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway deployment group name.Note: This field may return null, which means no valid value was found.",
						},
						"group_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Gateway deployment api group number.Note: This field may return null, which means no valid value was found.",
						},
						"groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Gateway deployment  api group list.Note: This field may return null, which means no valid value was found.",
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
									"group_api_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of APIs under the group. Note: This field may return null, which means no valid value was found.",
									},
									"group_apis": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of APIs under the group.Note: This field may return null, which means no valid value was found.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"api_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "API ID.",
												},
												"path": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "API path.",
												},
												"microservice_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "API service name.",
												},
												"method": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "API method.",
												},
												"namespace_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "namespace name.",
												},
											},
										},
									},
									"gateway_instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the gateway instance.Note: This field may return null, which means no valid value was found.",
									},
									"gateway_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "gateway instance id.Note: This field may return null, which means no valid value was found.",
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

func dataSourceTencentCloudTsfGatewayAllGroupApisRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_gateway_all_group_apis.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var gatewayDeployGroupId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_deploy_group_id"); ok {
		gatewayDeployGroupId = v.(string)
		paramMap["GatewayDeployGroupId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *tsf.GatewayVo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTsfGatewayAllGroupApisByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	gatewayVoMap := map[string]interface{}{}
	if result != nil {

		if result.GatewayDeployGroupId != nil {
			gatewayVoMap["gateway_deploy_group_id"] = result.GatewayDeployGroupId
		}

		if result.GatewayDeployGroupName != nil {
			gatewayVoMap["gateway_deploy_group_name"] = result.GatewayDeployGroupName
		}

		if result.GroupNum != nil {
			gatewayVoMap["group_num"] = result.GroupNum
		}

		if result.Groups != nil {
			groupsList := []interface{}{}
			for _, groups := range result.Groups {
				groupsMap := map[string]interface{}{}

				if groups.GroupId != nil {
					groupsMap["group_id"] = groups.GroupId
				}

				if groups.GroupName != nil {
					groupsMap["group_name"] = groups.GroupName
				}

				if groups.GroupApiCount != nil {
					groupsMap["group_api_count"] = groups.GroupApiCount
				}

				if groups.GroupApis != nil {
					groupApisList := []interface{}{}
					for _, groupApis := range groups.GroupApis {
						groupApisMap := map[string]interface{}{}

						if groupApis.ApiId != nil {
							groupApisMap["api_id"] = groupApis.ApiId
						}

						if groupApis.Path != nil {
							groupApisMap["path"] = groupApis.Path
						}

						if groupApis.MicroserviceName != nil {
							groupApisMap["microservice_name"] = groupApis.MicroserviceName
						}

						if groupApis.Method != nil {
							groupApisMap["method"] = groupApis.Method
						}

						if groupApis.NamespaceName != nil {
							groupApisMap["namespace_name"] = groupApis.NamespaceName
						}

						groupApisList = append(groupApisList, groupApisMap)
					}

					groupsMap["group_apis"] = groupApisList
				}

				if groups.GatewayInstanceType != nil {
					groupsMap["gateway_instance_type"] = groups.GatewayInstanceType
				}

				if groups.GatewayInstanceId != nil {
					groupsMap["gateway_instance_id"] = groups.GatewayInstanceId
				}

				groupsList = append(groupsList, groupsMap)
			}

			gatewayVoMap["groups"] = groupsList
		}

		_ = d.Set("result", []interface{}{gatewayVoMap})
	}

	d.SetId(helper.DataResourceIdsHash([]string{gatewayDeployGroupId}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), gatewayVoMap); e != nil {
			return e
		}
	}
	return nil
}
