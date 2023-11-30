package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudEbPlatformEventPatterns() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEbPlatformEventPatternsRead,
		Schema: map[string]*schema.Schema{
			"product_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Platform product type.",
			},

			"event_patterns": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Platform product event matching rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Platform event name.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"event_pattern": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Platform event matching rules.Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudEbPlatformEventPatternsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_eb_platform_event_patterns.read")()
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

	var eventPatterns []*eb.PlatformEventSummary

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEbPlatformEventPatternsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		eventPatterns = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(eventPatterns))
	tmpList := make([]map[string]interface{}, 0, len(eventPatterns))

	if eventPatterns != nil {
		for _, platformEventSummary := range eventPatterns {
			platformEventSummaryMap := map[string]interface{}{}

			if platformEventSummary.EventName != nil {
				platformEventSummaryMap["event_name"] = platformEventSummary.EventName
			}

			if platformEventSummary.EventPattern != nil {
				platformEventSummaryMap["event_pattern"] = platformEventSummary.EventPattern
				ids = append(ids, *platformEventSummary.EventPattern)
			}
			tmpList = append(tmpList, platformEventSummaryMap)
		}

		_ = d.Set("event_patterns", tmpList)
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
