package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAntiddosOverviewDdosEventList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAntiddosOverviewDdosEventListRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "StartTime.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "EndTime.",
			},

			"attack_status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "filter event by attack status, start: attacking; end: attack end.",
			},

			"event_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "EventList.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "event id.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ip.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "StartTime.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "EndTime.",
						},
						"attack_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AttackType.",
						},
						"attack_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Attack status, 0: Under attack; 1: End of attack.",
						},
						"mbps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Attack traffic, unit Mbps.",
						},
						"pps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "unit Mbps.",
						},
						"business": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dayu sub product code (bgpip represents advanced defense IP; net represents professional version of advanced defense IP).",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "InstanceId.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "InstanceId.",
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

func dataSourceTencentCloudAntiddosOverviewDdosEventListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_antiddos_overview_ddos_event_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("attack_status"); ok {
		paramMap["AttackStatus"] = helper.String(v.(string))
	}

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	var eventList []*antiddos.OverviewDDoSEvent

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAntiddosOverviewDdosEventListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		eventList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(eventList))
	tmpList := make([]map[string]interface{}, 0, len(eventList))

	if eventList != nil {
		for _, overviewDDoSEvent := range eventList {
			overviewDDoSEventMap := map[string]interface{}{}

			if overviewDDoSEvent.Id != nil {
				overviewDDoSEventMap["id"] = overviewDDoSEvent.Id
			}

			if overviewDDoSEvent.Vip != nil {
				overviewDDoSEventMap["vip"] = overviewDDoSEvent.Vip
			}

			if overviewDDoSEvent.StartTime != nil {
				overviewDDoSEventMap["start_time"] = overviewDDoSEvent.StartTime
			}

			if overviewDDoSEvent.EndTime != nil {
				overviewDDoSEventMap["end_time"] = overviewDDoSEvent.EndTime
			}

			if overviewDDoSEvent.AttackType != nil {
				overviewDDoSEventMap["attack_type"] = overviewDDoSEvent.AttackType
			}

			if overviewDDoSEvent.AttackStatus != nil {
				overviewDDoSEventMap["attack_status"] = overviewDDoSEvent.AttackStatus
			}

			if overviewDDoSEvent.Mbps != nil {
				overviewDDoSEventMap["mbps"] = overviewDDoSEvent.Mbps
			}

			if overviewDDoSEvent.Pps != nil {
				overviewDDoSEventMap["pps"] = overviewDDoSEvent.Pps
			}

			if overviewDDoSEvent.Business != nil {
				overviewDDoSEventMap["business"] = overviewDDoSEvent.Business
			}

			if overviewDDoSEvent.InstanceId != nil {
				overviewDDoSEventMap["instance_id"] = overviewDDoSEvent.InstanceId
			}

			if overviewDDoSEvent.InstanceName != nil {
				overviewDDoSEventMap["instance_name"] = overviewDDoSEvent.InstanceName
			}

			ids = append(ids, *overviewDDoSEvent.Id)
			tmpList = append(tmpList, overviewDDoSEventMap)
		}

		_ = d.Set("event_list", tmpList)
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
