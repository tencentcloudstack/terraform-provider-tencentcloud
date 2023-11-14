/*
Use this data source to query detailed information of mps parse_notification

Example Usage

```hcl
data "tencentcloud_mps_parse_notification" "parse_notification" {
  content = ""
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

func dataSourceTencentCloudMpsParseNotification() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsParseNotificationRead,
		Schema: map[string]*schema.Schema{
			"content": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Event notification obtained from CMQ.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMpsParseNotificationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_parse_notification.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("content"); ok {
		paramMap["Content"] = helper.String(v.(string))
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsParseNotificationByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		eventType = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(eventType))
	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
