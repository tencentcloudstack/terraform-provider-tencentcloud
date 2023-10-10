/*
Use this data source to query detailed information of eb plateform_event_template

Example Usage

```hcl
data "tencentcloud_eb_plateform_event_template" "plateform_event_template" {
  event_type = "eb_platform_test:TEST:ALL"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudEbPlateformEventTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEbPlateformEventTemplateRead,
		Schema: map[string]*schema.Schema{
			"event_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Platform product event type.",
			},

			"event_template": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Platform product event template.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudEbPlateformEventTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_eb_plateform_event_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var eventType string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("event_type"); ok {
		eventType = v.(string)
		paramMap["EventType"] = helper.String(v.(string))
	}

	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var eventTemplate *string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEbPlateformEventTemplateByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		eventTemplate = result
		return nil
	})
	if err != nil {
		return err
	}

	if eventTemplate != nil {
		_ = d.Set("event_template", eventTemplate)
	}

	d.SetId(helper.DataResourceIdsHash([]string{eventType, *eventTemplate}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), eventTemplate); e != nil {
			return e
		}
	}
	return nil
}
