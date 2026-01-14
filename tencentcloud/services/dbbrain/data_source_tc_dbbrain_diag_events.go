package dbbrain

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainDiagEvents() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainDiagEventsRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance ID list.",
			},

			"product": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Service product type; supported values include: `mysql` - Cloud Database MySQL, `redis` - Cloud Database Redis, `mariadb` - MariaDB database. The default is `mysql`.",
			},

			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Start time.",
			},

			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "End time.",
			},

			"severities": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Severity list, optional value is 1-fatal, 2-severity, 3-warning, 4-tips, 5-health.",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Diag event list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"diag_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Diag type.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End time.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time.",
						},
						"event_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Event ID.",
						},
						"severity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Severity.",
						},
						"outline": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Outline.",
						},
						"diag_item": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Diag item.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"metric": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Metric.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
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
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_diag_events.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["instance_ids"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["product"] = helper.String(v.(string))
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

	var items []*dbbrain.DiagHistoryEventItem
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainDiagEventsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		items = result
		return nil
	})

	if err != nil {
		return err
	}

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

			tmpList = append(tmpList, diagHistoryEventItemMap)
		}

		_ = d.Set("list", tmpList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
