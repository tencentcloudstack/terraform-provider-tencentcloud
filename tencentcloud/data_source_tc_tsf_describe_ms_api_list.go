/*
Use this data source to query detailed information of tsf describe_ms_api_list

Example Usage

```hcl
data "tencentcloud_tsf_describe_ms_api_list" "describe_ms_api_list" {
  microservice_id = "ms-xxxxxxxx"
  search_word = ""
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

func dataSourceTencentCloudTsfDescribeMsApiList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfDescribeMsApiListRead,
		Schema: map[string]*schema.Schema{
			"microservice_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Microservice Id.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search word, support  service name.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Result list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Quantity.Note: This field may return null, indicating that no valid value was found.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Api list .Note: This field may return null, indicating that no valid value was found.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Api path.",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Api method.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Method description. Note: This field may return null, indicating that no valid value was found.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "API status. 0: offline, 1: online.Note: This field may return null, indicating that no valid value was found.",
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

func dataSourceTencentCloudTsfDescribeMsApiListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_describe_ms_api_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("microservice_id"); ok {
		paramMap["MicroserviceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*tsf.TsfApiListResponse

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfDescribeMsApiListByFilter(ctx, paramMap)
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
		tsfApiListResponseMap := map[string]interface{}{}

		if result.TotalCount != nil {
			tsfApiListResponseMap["total_count"] = result.TotalCount
		}

		if result.Content != nil {
			contentList := []interface{}{}
			for _, content := range result.Content {
				contentMap := map[string]interface{}{}

				if content.Path != nil {
					contentMap["path"] = content.Path
				}

				if content.Method != nil {
					contentMap["method"] = content.Method
				}

				if content.Description != nil {
					contentMap["description"] = content.Description
				}

				if content.Status != nil {
					contentMap["status"] = content.Status
				}

				contentList = append(contentList, contentMap)
			}

			tsfApiListResponseMap["content"] = []interface{}{contentList}
		}

		ids = append(ids, *result.MicroserviceId)
		_ = d.Set("result", tsfApiListResponseMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tsfApiListResponseMap); e != nil {
			return e
		}
	}
	return nil
}
