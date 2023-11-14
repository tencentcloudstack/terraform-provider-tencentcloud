/*
Use this data source to query detailed information of apigateway service

Example Usage

```hcl
data "tencentcloud_apigateway_service" "service" {
  service_id = ""
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

func dataSourceTencentCloudApigatewayService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApigatewayServiceRead,
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
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of service binding environments.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"environment_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of service binding environments.Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudApigatewayServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_apigateway_service.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
	}

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*apigateway.ServiceEnvironmentSet

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeApigatewayServiceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		serviceEnvironmentSetMap := map[string]interface{}{}

		if result.TotalCount != nil {
			serviceEnvironmentSetMap["total_count"] = result.TotalCount
		}

		if result.EnvironmentList != nil {
			environmentListList := []interface{}{}
			for _, environmentList := range result.EnvironmentList {
				environmentListMap := map[string]interface{}{}

				if environmentList.EnvironmentName != nil {
					environmentListMap["environment_name"] = environmentList.EnvironmentName
				}

				if environmentList.Url != nil {
					environmentListMap["url"] = environmentList.Url
				}

				if environmentList.Status != nil {
					environmentListMap["status"] = environmentList.Status
				}

				if environmentList.VersionName != nil {
					environmentListMap["version_name"] = environmentList.VersionName
				}

				environmentListList = append(environmentListList, environmentListMap)
			}

			serviceEnvironmentSetMap["environment_list"] = []interface{}{environmentListList}
		}

		ids = append(ids, *result.ServiceId)
		_ = d.Set("result", serviceEnvironmentSetMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), serviceEnvironmentSetMap); e != nil {
			return e
		}
	}
	return nil
}
