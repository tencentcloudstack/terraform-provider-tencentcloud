/*
Use this data source to query detailed information of dbbrain diag_events

Example Usage

```hcl
data "tencentcloud_dbbrain_diag_events" "diag_events" {
  instance_ids = ["%s"]
  start_time = "%s"
  end_time = "%s"
  severities = [1,4,5]
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

func dataSourceTencentCloudDbbrainDiagEvents() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainDiagEventsRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "instance id list.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "start time.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "end time.",
			},

			"severities": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "severity list, optional value is 1-fatal, 2-severity, 3-warning, 4-tips, 5-health.",
			},

			"list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "diag event list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"diag_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "diag type.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "end time.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "start time.",
						},
						"event_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "event id.",
						},
						"severity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "severity.",
						},
						"outline": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "outline.",
						},
						"diag_item": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "diag item.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance id.",
						},
						"metric": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "metric.",
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

func dataSourceTencentCloudDbbrainDiagEventsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_diag_events.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["instance_ids"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["start_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["end_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("severities"); ok {
		severitiesSet := v.(*schema.Set).List()
		tmpSet := make([]*int64, 0, len(severitiesSet))
		for i := range severitiesSet {
			severities := severitiesSet[i].(int)
			tmpSet = append(tmpSet, helper.IntInt64(severities))
		}
		paramMap["severities"] = tmpSet
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	var items []*dbbrain.DiagHistoryEventItem

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainDiagEventsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, diagHistoryEventItem := range items {
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

			ids = append(ids, helper.Int64ToStr(*diagHistoryEventItem.EventId))
			tmpList = append(tmpList, diagHistoryEventItemMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		_ = d.Set("list", tmpList)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
