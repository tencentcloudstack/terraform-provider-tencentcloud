/*
Use this data source to query detailed information of live describe_live_pull_stream_task_status

Example Usage

```hcl
data "tencentcloud_live_describe_live_pull_stream_task_status" "describe_live_pull_stream_task_status" {
  task_id = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLiveDescribeLivePullStreamTaskStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLiveDescribeLivePullStreamTaskStatusRead,
		Schema: map[string]*schema.Schema{
			"task_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task ID。.",
			},

			"task_status_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Task status info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current use source url.",
						},
						"looped_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of times a VOD source task is played in a loop.",
						},
						"offset_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The playback offset of the VOD source, in seconds.",
						},
						"report_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The latest heartbeat reporting time in UTC format, for example: 2022-02-11T10:00:00Z.Note: UTC time is 8 hours ahead of Beijing time.",
						},
						"run_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Real run status:active,inactive.",
						},
						"file_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The duration of the VOD source file, in seconds.",
						},
						"next_file_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the next progress VOD file.",
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

func dataSourceTencentCloudLiveDescribeLivePullStreamTaskStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_live_describe_live_pull_stream_task_status.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("task_id"); ok {
		paramMap["TaskId"] = helper.String(v.(string))
	}

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	var taskStatusInfo []*live.TaskStatusInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLiveDescribeLivePullStreamTaskStatusByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		taskStatusInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(taskStatusInfo))
	if taskStatusInfo != nil {
		taskStatusInfoMap := map[string]interface{}{}

		if taskStatusInfo.FileUrl != nil {
			taskStatusInfoMap["file_url"] = taskStatusInfo.FileUrl
		}

		if taskStatusInfo.LoopedTimes != nil {
			taskStatusInfoMap["looped_times"] = taskStatusInfo.LoopedTimes
		}

		if taskStatusInfo.OffsetTime != nil {
			taskStatusInfoMap["offset_time"] = taskStatusInfo.OffsetTime
		}

		if taskStatusInfo.ReportTime != nil {
			taskStatusInfoMap["report_time"] = taskStatusInfo.ReportTime
		}

		if taskStatusInfo.RunStatus != nil {
			taskStatusInfoMap["run_status"] = taskStatusInfo.RunStatus
		}

		if taskStatusInfo.FileDuration != nil {
			taskStatusInfoMap["file_duration"] = taskStatusInfo.FileDuration
		}

		if taskStatusInfo.NextFileUrl != nil {
			taskStatusInfoMap["next_file_url"] = taskStatusInfo.NextFileUrl
		}

		ids = append(ids, *taskStatusInfo.TaskId)
		_ = d.Set("task_status_info", taskStatusInfoMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), taskStatusInfoMap); e != nil {
			return e
		}
	}
	return nil
}
