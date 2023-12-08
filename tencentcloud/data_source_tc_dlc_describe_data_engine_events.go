package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDlcDescribeDataEngineEvents() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeDataEngineEventsRead,
		Schema: map[string]*schema.Schema{
			"data_engine_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Data engine name.",
			},

			"events": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Event details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Event time.",
						},
						"events_action": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Event action.",
						},
						"cluster_info": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Cluster information.",
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

func dataSourceTencentCloudDlcDescribeDataEngineEventsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dlc_describe_data_engine_events.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("data_engine_name"); ok {
		paramMap["DataEngineName"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var events []*dlc.HouseEventsInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDescribeDataEngineEventsByFilter(ctx, paramMap)
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
		for _, houseEventsInfo := range events {
			houseEventsInfoMap := map[string]interface{}{}

			if houseEventsInfo.Time != nil {
				houseEventsInfoMap["time"] = houseEventsInfo.Time
			}

			if houseEventsInfo.EventsAction != nil {
				houseEventsInfoMap["events_action"] = houseEventsInfo.EventsAction
			}

			if houseEventsInfo.ClusterInfo != nil {
				houseEventsInfoMap["cluster_info"] = houseEventsInfo.ClusterInfo
			}

			tmpList = append(tmpList, houseEventsInfoMap)
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
