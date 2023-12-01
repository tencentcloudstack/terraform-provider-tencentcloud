/*
Use this data source to query detailed information of apiGateway service_environment_list

Example Usage

```hcl
data "tencentcloud_api_gateway_service_environment_list" "example" {
  service_id = "service-nxz6yync"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudApiGatewayServiceEnvironmentList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApiGatewayServiceEnvironmentListRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the service to be queried.",
			},
			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Service binding environment details.Note: This field may return null, indicating that no valid value can be obtained.",
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
							Description: "Access path.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Release status, 1 means released, 0 means not released.",
						},
						"version_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Running version.",
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

func dataSourceTencentCloudApiGatewayServiceEnvironmentListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_api_gateway_service_environment_list.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId           = getLogId(contextNil)
		ctx             = context.WithValue(context.TODO(), logIdKey, logId)
		service         = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		environmentList []*apigateway.Environment
		serviceId       string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
		serviceId = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeApiGatewayServiceEnvironmentListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		environmentList = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(environmentList))
	if environmentList != nil {
		for _, environment := range environmentList {
			environmentListMap := map[string]interface{}{}

			if environment.EnvironmentName != nil {
				environmentListMap["environment_name"] = environment.EnvironmentName
			}

			if environment.Url != nil {
				environmentListMap["url"] = environment.Url
			}

			if environment.Status != nil {
				environmentListMap["status"] = environment.Status
			}

			if environment.VersionName != nil {
				environmentListMap["version_name"] = environment.VersionName
			}

			tmpList = append(tmpList, environmentListMap)
		}

		_ = d.Set("result", tmpList)
	}

	d.SetId(serviceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
