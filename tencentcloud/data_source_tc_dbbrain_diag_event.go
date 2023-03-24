/*
Use this data source to query detailed information of dbbrain diag_event

Example Usage

```hcl
data "tencentcloud_dbbrain_diag_history" "diag_history" {
	instance_id = "%s"
	start_time = "%s"
	end_time = "%s"
	product = "mysql"
}

data "tencentcloud_dbbrain_diag_event" "diag_event" {
  instance_id = "%s"
  event_id = data.tencentcloud_dbbrain_diag_history.diag_history.events.0.event_id
  product = "mysql"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDbbrainDiagEvent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainDiagEventRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "isntance id.",
			},

			"event_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Event ID. Obtain it through `Get Instance Diagnosis History DescribeDBDiagHistory`.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include: `mysql` - cloud database MySQL, `cynosdb` - cloud database CynosDB for MySQL, the default is `mysql`.",
			},

			"diag_item": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "diagnostic item.",
			},

			"diag_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Diagnostic type.",
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
				Description: "severity. The severity is divided into 5 levels, according to the degree of impact from high to low: 1: Fatal, 2: Serious, 3: Warning, 4: Prompt, 5: Healthy.",
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
				Description: "reserved text. Note: This field may return null, indicating that no valid value can be obtained.",
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
	var id string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
		id = v.(string)
	}

	if v, _ := d.GetOk("event_id"); v != nil {
		paramMap["event_id"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["product"] = helper.String(v.(string))
	}

	var result *dbbrain.DescribeDBDiagEventResponseParams
	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		var e error
		result, e = service.DescribeDbbrainDiagEventByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if result != nil {
		if result.DiagItem != nil {
			_ = d.Set("diag_item", result.DiagItem)
		}

		if result.DiagType != nil {
			_ = d.Set("diag_type", result.DiagType)
		}

		if result.EventId != nil {
			_ = d.Set("event_id", result.EventId)
		}

		if result.Explanation != nil {
			_ = d.Set("explanation", result.Explanation)
		}

		if result.Outline != nil {
			_ = d.Set("outline", result.Outline)
		}

		if result.Problem != nil {
			_ = d.Set("problem", result.Problem)
		}

		if result.Severity != nil {
			_ = d.Set("severity", result.Severity)
		}

		if result.StartTime != nil {
			_ = d.Set("start_time", result.StartTime)
		}

		if result.Suggestions != nil {
			_ = d.Set("suggestions", result.Suggestions)
		}

		if result.Metric != nil {
			_ = d.Set("metric", result.Metric)
		}

		if result.EndTime != nil {
			_ = d.Set("end_time", result.EndTime)
		}

	}

	d.SetId(id)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
