package apigateway

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudApiGatewayApiPlugins() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApiGatewayApiPluginsRead,
		Schema: map[string]*schema.Schema{
			"api_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "API ID to be queried.",
			},
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The service ID to be queried.",
			},
			"environment_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Environment information.",
			},
			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "API list information that the plug-in can bind.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plugin_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plugin ID.",
						},
						"environment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Environment information.",
						},
						"attached_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Binding time.",
						},
						"plugin_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plugin name.",
						},
						"plugin_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plugin type.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plugin description.",
						},
						"plugin_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plug-in definition statement.",
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

func dataSourceTencentCloudApiGatewayApiPluginsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_api_plugins.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		apiPlugins []*apigateway.AttachedPluginInfo
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("api_id"); ok {
		paramMap["APIId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("environment_name"); ok {
		paramMap["EnvironmentName"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeApiGatewayApiPluginsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		apiPlugins = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(apiPlugins))
	tmpList := make([]map[string]interface{}, 0, len(apiPlugins))
	if apiPlugins != nil {
		for _, pluginSummary := range apiPlugins {
			pluginSummaryMap := map[string]interface{}{}

			if pluginSummary.PluginId != nil {
				pluginSummaryMap["plugin_id"] = pluginSummary.PluginId
			}

			if pluginSummary.Environment != nil {
				pluginSummaryMap["environment"] = pluginSummary.Environment
			}

			if pluginSummary.AttachedTime != nil {
				pluginSummaryMap["attached_time"] = pluginSummary.AttachedTime
			}

			if pluginSummary.PluginName != nil {
				pluginSummaryMap["plugin_name"] = pluginSummary.PluginName
			}

			if pluginSummary.PluginType != nil {
				pluginSummaryMap["plugin_type"] = pluginSummary.PluginType
			}

			if pluginSummary.Description != nil {
				pluginSummaryMap["description"] = pluginSummary.Description
			}

			if pluginSummary.PluginData != nil {
				pluginSummaryMap["plugin_data"] = pluginSummary.PluginData
			}

			tmpList = append(tmpList, pluginSummaryMap)
			ids = append(ids, *pluginSummary.PluginId)
		}

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
