/*
Use this data source to query detailed information of tsf usable_unit_namespaces

Example Usage

```hcl
data "tencentcloud_tsf_usable_unit_namespaces" "usable_unit_namespaces" {
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

func dataSourceTencentCloudTsfUsableUnitNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfUsableUnitNamespacesRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "search by namespace id or namespace Name.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "namespace object list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "total count.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "namespace list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "namespace id.",
									},
									"namespace_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "namespace name.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unit namespace ID. Note: This field may return null, indicating that no valid value was found.",
									},
									"gateway_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway instance id Note: This field may return null, indicating that no valid value was found.",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Create time. Note: This field may return null, indicating that no valid value was found.",
									},
									"updated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Update time. Note: This field may return null, indicating that no valid value was found.",
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

func dataSourceTencentCloudTsfUsableUnitNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_usable_unit_namespaces.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *tsf.TsfPageUnitNamespace
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTsfUsableUnitNamespacesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result.Content))
	tsfPageUnitNamespaceMap := map[string]interface{}{}
	if result != nil {

		if result.TotalCount != nil {
			tsfPageUnitNamespaceMap["total_count"] = result.TotalCount
		}

		if result.Content != nil {
			contentList := []interface{}{}
			for _, content := range result.Content {
				contentMap := map[string]interface{}{}

				if content.NamespaceId != nil {
					contentMap["namespace_id"] = content.NamespaceId
				}

				if content.NamespaceName != nil {
					contentMap["namespace_name"] = content.NamespaceName
				}

				if content.Id != nil {
					contentMap["id"] = content.Id
				}

				if content.GatewayInstanceId != nil {
					contentMap["gateway_instance_id"] = content.GatewayInstanceId
				}

				if content.CreatedTime != nil {
					contentMap["created_time"] = content.CreatedTime
				}

				if content.UpdatedTime != nil {
					contentMap["updated_time"] = content.UpdatedTime
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.NamespaceId)
			}

			tsfPageUnitNamespaceMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageUnitNamespaceMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tsfPageUnitNamespaceMap); e != nil {
			return e
		}
	}
	return nil
}
