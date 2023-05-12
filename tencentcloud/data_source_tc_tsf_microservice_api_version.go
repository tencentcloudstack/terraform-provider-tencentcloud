/*
Use this data source to query detailed information of tsf microservice_api_version

Example Usage

```hcl
data "tencentcloud_tsf_microservice_api_version" "microservice_api_version" {
  microservice_id = "ms-yq3jo6jd"
  path = ""
  method = "get"
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

func dataSourceTencentCloudTsfMicroserviceApiVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfMicroserviceApiVersionRead,
		Schema: map[string]*schema.Schema{
			"microservice_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Microservice ID.",
			},

			"path": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "api path.",
			},

			"method": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "request method.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "api version list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application ID.",
						},
						"application_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application Name.",
						},
						"pkg_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "application pkg version.",
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

func dataSourceTencentCloudTsfMicroserviceApiVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_microservice_api_version.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("microservice_id"); ok {
		paramMap["MicroserviceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("path"); ok {
		paramMap["Path"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("method"); ok {
		paramMap["Method"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var apiVersion []*tsf.ApiVersionArray
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfMicroserviceApiVersionByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		apiVersion = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(apiVersion))
	tmpList := make([]map[string]interface{}, 0, len(apiVersion))
	if apiVersion != nil {
		for _, apiVersionArray := range apiVersion {
			apiVersionArrayMap := map[string]interface{}{}

			if apiVersionArray.ApplicationId != nil {
				apiVersionArrayMap["application_id"] = apiVersionArray.ApplicationId
			}

			if apiVersionArray.ApplicationName != nil {
				apiVersionArrayMap["application_name"] = apiVersionArray.ApplicationName
			}

			if apiVersionArray.PkgVersion != nil {
				apiVersionArrayMap["pkg_version"] = apiVersionArray.PkgVersion
			}

			ids = append(ids, *apiVersionArray.ApplicationId)
			tmpList = append(tmpList, apiVersionArrayMap)
		}

		_ = d.Set("result", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
