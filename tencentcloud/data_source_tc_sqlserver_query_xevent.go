/*
Use this data source to query detailed information of sqlserver query_xevent

Example Usage

```hcl
data "tencentcloud_sqlserver_query_xevent" "query_xevent" {
  instance_id = ""
  event_type = ""
  start_time = ""
  end_time = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSqlserverQueryXevent() *schema.Resource {
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
	defer logElapsed("data_source.tencentcloud_sqlserver_query_xevent.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
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

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	var events []*sqlserver.Events

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverQueryXeventByFilter(ctx, paramMap)
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
		for _, events := range events {
			eventsMap := map[string]interface{}{}

			if events.Id != nil {
				eventsMap["id"] = events.Id
			}

			if events.FileName != nil {
				eventsMap["file_name"] = events.FileName
			}

			if events.Size != nil {
				eventsMap["size"] = events.Size
			}

			if events.EventType != nil {
				eventsMap["event_type"] = events.EventType
			}

			if events.Status != nil {
				eventsMap["status"] = events.Status
			}

			if events.StartTime != nil {
				eventsMap["start_time"] = events.StartTime
			}

			if events.EndTime != nil {
				eventsMap["end_time"] = events.EndTime
			}

			if events.InternalAddr != nil {
				eventsMap["internal_addr"] = events.InternalAddr
			}

			if events.ExternalAddr != nil {
				eventsMap["external_addr"] = events.ExternalAddr
			}

			ids = append(ids, *events.InstanceId)
			tmpList = append(tmpList, eventsMap)
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
