/*
Use this data source to query detailed information of dbbrain diag_event

Example Usage

```hcl
data "tencentcloud_dbbrain_diag_event" "diag_event" {
  instance_id = ""
  event_id =
  product = ""
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

func dataSourceTencentCloudDbbrainDiagEvent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainDiagEventRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Isntance id.",
			},

			"event_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Event ID. Obtain it through Get Instance Diagnosis History DescribeDBDiagHistory.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include： mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.",
			},

			"diag_item": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Diagnostic item.",
			},

			"diag_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Diagnostic type.",
			},

			"event_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Event ID.",
			},

			"explanation": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Diagnostic event details, output is empty if there is no additional explanatory information.",
			},

			"outline": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Diagnostic summary.",
			},

			"problem": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Diagnosed problem.",
			},

			"severity": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Severity. The severity is divided into 5 levels, according to the degree of impact from high to low: 1: Fatal, 2: Serious, 3: Warning, 4: Prompt, 5: Healthy.",
			},

			"start_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Starting time.",
			},

			"suggestions": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "A diagnostic suggestion, or empty if there is no suggestion.",
			},

			"metric": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Reserved text. Note: This field may return null, indicating that no valid value can be obtained.",
			},

			"end_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "End Time.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainDiagEventRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_diag_event.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("event_id"); v != nil {
		paramMap["EventId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainDiagEventByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		diagItem = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(diagItem))
	if diagItem != nil {
		_ = d.Set("diag_item", diagItem)
	}

	if diagType != nil {
		_ = d.Set("diag_type", diagType)
	}

	if eventId != nil {
		_ = d.Set("event_id", eventId)
	}

	if explanation != nil {
		_ = d.Set("explanation", explanation)
	}

	if outline != nil {
		_ = d.Set("outline", outline)
	}

	if problem != nil {
		_ = d.Set("problem", problem)
	}

	if severity != nil {
		_ = d.Set("severity", severity)
	}

	if startTime != nil {
		_ = d.Set("start_time", startTime)
	}

	if suggestions != nil {
		_ = d.Set("suggestions", suggestions)
	}

	if metric != nil {
		_ = d.Set("metric", metric)
	}

	if endTime != nil {
		_ = d.Set("end_time", endTime)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
