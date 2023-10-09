/*
Use this data source to query detailed information of gaap custom header

Example Usage

```hcl
data "tencentcloud_gaap_custom_header" "custom_header" {
  rule_id = "rule-9sdhv655"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudGaapCustomHeader() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapCustomHeaderRead,
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Rule IdNote: This field may return null, indicating that a valid value cannot be obtained.",
			},

			"headers": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "HeadersNote: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Header Name.",
						},
						"header_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Header Value.",
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

func dataSourceTencentCloudGaapCustomHeaderRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_custom_header.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var headers []*gaap.HttpHeaderParam
	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}
	ruleId := d.Get("rule_id").(string)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapCustomHeader(ctx, ruleId)
		if e != nil {
			return retryError(e)
		}
		headers = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(ruleId))
	tmpList := make([]map[string]interface{}, 0)
	if headers != nil {
		for _, httpHeaderParam := range headers {
			httpHeaderParamMap := map[string]interface{}{}

			if httpHeaderParam.HeaderName != nil {
				httpHeaderParamMap["header_name"] = httpHeaderParam.HeaderName
			}

			if httpHeaderParam.HeaderValue != nil {
				httpHeaderParamMap["header_value"] = httpHeaderParam.HeaderValue
			}

			ids = append(ids, *httpHeaderParam.HeaderName)
			tmpList = append(tmpList, httpHeaderParamMap)
		}

		_ = d.Set("headers", tmpList)
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
