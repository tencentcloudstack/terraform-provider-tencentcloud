package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAPIGatewayUpstreams() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayUpstreamRead,
		Schema: map[string]*schema.Schema{
			"upstream_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backend channel ID.",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "ServiceId and ApiId filtering queries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Fields that need to be filtered.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "The filtering value of the field.",
						},
					},
				},
			},
			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Query Results.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API Unique ID.",
						},
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service Unique ID.",
						},
						"api_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API nameNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service NameNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"bind_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "binding time.",
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

func dataSourceTencentCloudAPIGatewayUpstreamRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_api_gateway_upstreams.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		result  []*apigateway.BindApiInfo
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("upstream_id"); ok {
		paramMap["UpstreamId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*apigateway.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := apigateway.Filter{}
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

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeAPIGatewayUpstreamByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		result = response
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		bindApiSetList := []interface{}{}
		for _, bindApiSet := range result {
			bindApiSetMap := map[string]interface{}{}

			if bindApiSet.ApiId != nil {
				bindApiSetMap["api_id"] = bindApiSet.ApiId
			}

			if bindApiSet.ServiceId != nil {
				bindApiSetMap["service_id"] = bindApiSet.ServiceId
			}

			if bindApiSet.ApiName != nil {
				bindApiSetMap["api_name"] = bindApiSet.ApiName
			}

			if bindApiSet.ServiceName != nil {
				bindApiSetMap["service_name"] = bindApiSet.ServiceName
			}

			if bindApiSet.BindTime != nil {
				bindApiSetMap["bind_time"] = bindApiSet.BindTime
			}

			bindApiSetList = append(bindApiSetList, bindApiSetMap)
			ids = append(ids, *bindApiSet.ApiId)
		}

		_ = d.Set("result", bindApiSetList)
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
