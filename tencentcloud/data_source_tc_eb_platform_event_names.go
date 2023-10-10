/*
Use this data source to query detailed information of eb platform_event_names

Example Usage

```hcl
data "tencentcloud_eb_platform_event_names" "platform_event_names" {
  product_type = ""
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudEbPlatformEventNames() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEbPlateformRead,
		Schema: map[string]*schema.Schema{
			"product_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Platform product event type.",
			},

			"event_names": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Platform product list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event name.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"event_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event type.Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudEbPlateformRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_eb_platform_event_names.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var productType string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product_type"); ok {
		productType = v.(string)
		paramMap["ProductType"] = helper.String(v.(string))
	}

	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var eventNames []*eb.PlatformEventDetail
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEbPlateformByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		eventNames = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(eventNames))
	tmpList := make([]map[string]interface{}, 0, len(eventNames))
	if eventNames != nil {
		for _, platformEventDetail := range eventNames {
			platformEventDetailMap := map[string]interface{}{}

			if platformEventDetail.EventName != nil {
				platformEventDetailMap["event_name"] = platformEventDetail.EventName
			}

			if platformEventDetail.EventType != nil {
				platformEventDetailMap["event_type"] = platformEventDetail.EventType
				ids = append(ids, *platformEventDetail.EventType)
			}

			tmpList = append(tmpList, platformEventDetailMap)
		}

		_ = d.Set("event_names", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(append(ids, productType)))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
