package tencentcloud

import (
	"context"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func dataSourceTencentCloudAPIGatewayPlugins() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayPluginRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The service ID to query.",
			},
			"plugin_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The plugin ID to query.",
			},
			"environment_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Environmental information.",
			},
			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of plugin related APIs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"api_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API type.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API path.",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API method.",
						},
						"attached_other_plugin": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the API is bound to other plugins.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"is_attached": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the API is bound to the current plugin.Note: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudAPIGatewayPluginRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_api_gateway_plugins.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		infos   []*apigateway.AvailableApiInfo
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plugin_id"); ok {
		paramMap["PluginId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("environment_name"); ok {
		paramMap["EnvironmentName"] = helper.String(v.(string))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAPIGatewayPluginByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		infos = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(infos))
	if infos != nil {
		apiSetList := []interface{}{}
		for _, apiSet := range infos {
			apiSetMap := map[string]interface{}{}
			if apiSet.ApiId != nil {
				apiSetMap["api_id"] = apiSet.ApiId
			}

			if apiSet.ApiName != nil {
				apiSetMap["api_name"] = apiSet.ApiName
			}

			if apiSet.ApiType != nil {
				apiSetMap["api_type"] = apiSet.ApiType
			}

			if apiSet.Path != nil {
				apiSetMap["path"] = apiSet.Path
			}

			if apiSet.Method != nil {
				apiSetMap["method"] = apiSet.Method
			}

			if apiSet.AttachedOtherPlugin != nil {
				apiSetMap["attached_other_plugin"] = apiSet.AttachedOtherPlugin
			}

			if apiSet.IsAttached != nil {
				apiSetMap["is_attached"] = apiSet.IsAttached
			}

			apiSetList = append(apiSetList, apiSetMap)
			ids = append(ids, *apiSet.ApiId)
		}

		_ = d.Set("result", apiSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
