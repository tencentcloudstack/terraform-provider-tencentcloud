/*
Use this data source to query detailed information of dbbrain diag_history

Example Usage

```hcl
data "tencentcloud_dbbrain_diag_history" "diag_history" {
  instance_id = "%s"
  start_time = "%s"
  end_time = "%s"
  product = "mysql"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDbbrainDiagHistory() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainDiagHistoryRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time, such as `2019-09-10 12:13:14`.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time, such as `2019-09-11 12:13:14`, the interval between the end time and the start time can be up to 2 days.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include: `mysql` - cloud database MySQL, `cynosdb` - cloud database CynosDB for MySQL, the default is `mysql`.",
			},

			"events": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Event description.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"diag_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Diagnostic type.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End Time.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "start Time.",
						},
						"event_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Event unique ID.",
						},
						"severity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "severity. The severity is divided into 5 levels, according to the degree of impact from high to low: 1: Fatal, 2: Serious, 3: Warning, 4: Prompt, 5: Healthy.",
						},
						"outline": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Diagnostic summary.",
						},
						"diag_item": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the diagnostic item.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance id.",
						},
						"metric": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "reserved text. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region.",
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

func dataSourceTencentCloudDbbrainDiagHistoryRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_diag_history.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["start_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["end_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	var events []*dbbrain.DiagHistoryEventItem

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainDiagHistoryByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		events = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(events))
	tmpList := make([]map[string]interface{}, 0, len(events))

	if events != nil {
		for _, diagHistoryEventItem := range events {
			diagHistoryEventItemMap := map[string]interface{}{}

			if diagHistoryEventItem.DiagType != nil {
				diagHistoryEventItemMap["diag_type"] = diagHistoryEventItem.DiagType
			}

			if diagHistoryEventItem.EndTime != nil {
				diagHistoryEventItemMap["end_time"] = diagHistoryEventItem.EndTime
			}

			if diagHistoryEventItem.StartTime != nil {
				diagHistoryEventItemMap["start_time"] = diagHistoryEventItem.StartTime
			}

			if diagHistoryEventItem.EventId != nil {
				diagHistoryEventItemMap["event_id"] = diagHistoryEventItem.EventId
			}

			if diagHistoryEventItem.Severity != nil {
				diagHistoryEventItemMap["severity"] = diagHistoryEventItem.Severity
			}

			if diagHistoryEventItem.Outline != nil {
				diagHistoryEventItemMap["outline"] = diagHistoryEventItem.Outline
			}

			if diagHistoryEventItem.DiagItem != nil {
				diagHistoryEventItemMap["diag_item"] = diagHistoryEventItem.DiagItem
			}

			if diagHistoryEventItem.InstanceId != nil {
				diagHistoryEventItemMap["instance_id"] = diagHistoryEventItem.InstanceId
			}

			if diagHistoryEventItem.Metric != nil {
				diagHistoryEventItemMap["metric"] = diagHistoryEventItem.Metric
			}

			if diagHistoryEventItem.Region != nil {
				diagHistoryEventItemMap["region"] = diagHistoryEventItem.Region
			}

			ids = append(ids, *diagHistoryEventItem.InstanceId)
			tmpList = append(tmpList, diagHistoryEventItemMap)
		}

		_ = d.Set("events", tmpList)
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
