package apigateway

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiGateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudApiGatewayBindApiAppsStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApiGatewayBindApiAppsStatusRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Service ID.",
			},
			"api_ids": {
				Required:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Array of API IDs.",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions. Supports ApiAppId, Environment, KeyWord (can match name or ID).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Filter value of the field.",
						},
					},
				},
			},
			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of APIs bound by the application.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_app_api_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Application bound API information array.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_app_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application Name.",
									},
									"api_app_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application ID.",
									},
									"api_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API ID.",
									},
									"api_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API name.",
									},
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service ID.",
									},
									"authorized_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authorization binding time, expressed in accordance with the ISO8601 standard and using UTC time. The format is: YYYY-MM-DDThh:mm:ssZ.",
									},
									"api_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Apis region.",
									},
									"environment_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authorization binding environment.",
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

func dataSourceTencentCloudApiGatewayBindApiAppsStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_bind_api_apps_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service           = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		bindApiAppsStatus []*apiGateway.ApiAppApiInfo
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("api_ids"); ok {
		apiIdsSet := v.(*schema.Set).List()
		paramMap["APIIds"] = helper.InterfacesStringsPoint(apiIdsSet)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*apiGateway.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := apiGateway.Filter{}
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

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeApiGatewayBindApiAppsStatusByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		bindApiAppsStatus = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(bindApiAppsStatus))
	tmpList := make([]map[string]interface{}, 0, len(bindApiAppsStatus))
	if bindApiAppsStatus != nil {
		apiAppApiInfosMap := map[string]interface{}{}
		apiAppApiSetList := []interface{}{}
		for _, apiAppApiSet := range bindApiAppsStatus {
			apiAppApiSetMap := map[string]interface{}{}

			if apiAppApiSet.ApiAppName != nil {
				apiAppApiSetMap["api_app_name"] = apiAppApiSet.ApiAppName
			}

			if apiAppApiSet.ApiAppId != nil {
				apiAppApiSetMap["api_app_id"] = apiAppApiSet.ApiAppId
			}

			if apiAppApiSet.ApiId != nil {
				apiAppApiSetMap["api_id"] = apiAppApiSet.ApiId
			}

			if apiAppApiSet.ApiName != nil {
				apiAppApiSetMap["api_name"] = apiAppApiSet.ApiName
			}

			if apiAppApiSet.ServiceId != nil {
				apiAppApiSetMap["service_id"] = apiAppApiSet.ServiceId
			}

			if apiAppApiSet.AuthorizedTime != nil {
				apiAppApiSetMap["authorized_time"] = apiAppApiSet.AuthorizedTime
			}

			if apiAppApiSet.ApiRegion != nil {
				apiAppApiSetMap["api_region"] = apiAppApiSet.ApiRegion
			}

			if apiAppApiSet.EnvironmentName != nil {
				apiAppApiSetMap["environment_name"] = apiAppApiSet.EnvironmentName
			}

			apiAppApiSetList = append(apiAppApiSetList, apiAppApiSetMap)
			ids = append(ids, *apiAppApiSet.ApiAppId)
		}

		apiAppApiInfosMap["api_app_api_set"] = apiAppApiSetList
		tmpList = append(tmpList, apiAppApiInfosMap)
		_ = d.Set("result", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
