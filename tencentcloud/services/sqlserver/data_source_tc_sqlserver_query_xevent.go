package sqlserver

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverQueryXevent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverQueryXeventRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"event_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Event type. Valid values: slow (Slow SQL event), blocked (blocking event), deadlock` (deadlock event).",
			},
			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Generation start time of an extended file.",
			},
			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Generation end time of an extended file.",
			},
			"events": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of extended events.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID.",
						},
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File name of an extended event.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "File size of an extended event.",
						},
						"event_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event type. Valid values: slow (Slow SQL event), blocked (blocking event), deadlock (deadlock event).",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Event record status. Valid values: 1 (succeeded), 2 (failed).",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Generation start time of an extended file.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Generation end time of an extended file.",
						},
						"internal_addr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Download address on the private network.",
						},
						"external_addr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Download address on the public network.",
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

func dataSourceTencentCloudSqlserverQueryXeventRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_query_xevent.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		events     []*sqlserver.Events
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("event_type"); ok {
		paramMap["EventType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverQueryXeventByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		events = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(events))

	if events != nil {
		for _, event := range events {
			eventsMap := map[string]interface{}{}

			if event.Id != nil {
				eventsMap["id"] = event.Id
			}

			if event.FileName != nil {
				eventsMap["file_name"] = event.FileName
			}

			if event.Size != nil {
				eventsMap["size"] = event.Size
			}

			if event.EventType != nil {
				eventsMap["event_type"] = event.EventType
			}

			if event.Status != nil {
				eventsMap["status"] = event.Status
			}

			if event.StartTime != nil {
				eventsMap["start_time"] = event.StartTime
			}

			if event.EndTime != nil {
				eventsMap["end_time"] = event.EndTime
			}

			if event.InternalAddr != nil {
				eventsMap["internal_addr"] = event.InternalAddr
			}

			if event.ExternalAddr != nil {
				eventsMap["external_addr"] = event.ExternalAddr
			}

			tmpList = append(tmpList, eventsMap)
		}

		_ = d.Set("events", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
