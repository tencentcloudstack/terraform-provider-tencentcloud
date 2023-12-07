package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudApiGatewayServiceReleaseVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApiGatewayServiceReleaseVersionsRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the service to be queried.",
			},
			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of service releases.Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version number.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"version_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version description.Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudApiGatewayServiceReleaseVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_api_gateway_service_release_versions.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), logIdKey, logId)
		service     = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		versionList []*apigateway.DescribeServiceReleaseVersionResultVersionListInfo
		serviceId   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
		serviceId = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeApiGatewayServiceReleaseVersionsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		versionList = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(versionList))
	if versionList != nil {
		for _, version := range versionList {
			versionListMap := map[string]interface{}{}
			if version.VersionName != nil {
				versionListMap["version_name"] = version.VersionName
			}

			if version.VersionDesc != nil {
				versionListMap["version_desc"] = version.VersionDesc
			}

			tmpList = append(tmpList, versionListMap)
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
