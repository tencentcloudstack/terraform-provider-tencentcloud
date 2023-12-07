package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAPIGatewayThrottlingServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayThrottlingServicesRead,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Service ID for query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			//compute
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of Throttling policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service ID for query.",
						},
						"environments": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of Throttling policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"environment_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Environment name.",
									},
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Access service environment URL.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Release status.",
									},
									"version_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Published version number.",
									},
									"strategy": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Throttling value.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayThrottlingServicesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_api_gateway_throttling_services.read")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		infos             []*apigateway.Service
		serviceID         string
		err               error
		serviceIds        = make([]string, 0)
		resultLists       = make([]map[string]interface{}, 0)
		ids               = make([]string, 0)
	)
	if v, ok := d.GetOk("service_id"); ok {
		serviceID = v.(string)
	}

	if serviceID == "" {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			infos, err = apiGatewayService.DescribeServicesStatus(ctx, "", "")
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		})
		if err != nil {
			return err
		}

		for _, result := range infos {
			if result.ServiceId == nil {
				continue
			}
			serviceIds = append(serviceIds, *result.ServiceId)
		}
	} else {
		serviceIds = append(serviceIds, serviceID)
	}

	for _, serviceIdTmp := range serviceIds {
		environmentList, err := apiGatewayService.DescribeServiceEnvironmentStrategyList(ctx, serviceIdTmp)
		if err != nil {
			return err
		}

		environmentResults := make([]map[string]interface{}, 0, len(environmentList))
		for _, value := range environmentList {
			if value == nil {
				continue
			}
			item := map[string]interface{}{
				"environment_name": value.EnvironmentName,
				"url":              value.Url,
				"status":           value.Status,
				"version_name":     value.VersionName,
				"strategy":         value.Strategy,
			}
			environmentResults = append(environmentResults, item)
		}

		resultLists = append(resultLists, map[string]interface{}{
			"service_id":   serviceIdTmp,
			"environments": environmentResults,
		})
		ids = append(ids, serviceIdTmp)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if err = d.Set("list", resultLists); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), resultLists); err != nil {
			return err
		}
	}
	return nil
}
