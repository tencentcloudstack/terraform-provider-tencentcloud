/*
Use this data source to query detailed information of oceanus log

Example Usage

```hcl
data "tencentcloud_oceanus_log" "log" {
  job_id = "cql-6v1jkxrn"
  start_time = 1611754219108
  end_time = 1611754219108
  running_order_id = 1
  keyword = "xx"
    order_type = "asc"
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

func dataSourceTencentCloudOceanusLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusLogRead,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job ID.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Start time, unix timestamp, in milliseconds.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "End time, unix timestamp, in milliseconds.",
			},

			"running_order_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The ID of the instance, starting from 1 in the order of startup time. Default is 0, indicating that no instance is selected. Search for the logs of the most recent instance within the specified time range.",
			},

			"keyword": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The keyword for log search. Default is empty.",
			},

			"cursor": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The cursor for log search. It can be passed through the value returned last time. Default is empty.",
			},

			"order_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The sorting rule for timestamps. asc - ascending, desc - descending. Default is ascending.",
			},

			"list_over": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether all log records have been returned.",
			},

			"job_request_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The request ID for job startup.Note: This field may return null, indicating that no valid values can be obtained.",
			},

			"job_instance_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The list of all job instances that meet the keyword within the specified time range.Note: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"running_order_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the instance, starting from 1 in the order of startup time.",
						},
						"job_instance_start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The startup time of the instance.",
						},
						"starting_millis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The startup time of the instance in milliseconds.",
						},
					},
				},
			},

			"log_list": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Deprecated. Use LogContentList instead.Note: This field may return null, indicating that no valid values can be obtained.",
			},

			"log_content_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The list of logs.Note: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The content of the log.",
						},
						"time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timestamp in milliseconds.",
						},
						"pkg_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the log group.",
						},
						"pkg_log_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the log, which is unique within the log group.",
						},
						"container_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the container to which the log belongs.",
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

func dataSourceTencentCloudOceanusLogRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_oceanus_log.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("job_id"); ok {
		paramMap["JobId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("start_time"); v != nil {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("running_order_id"); v != nil {
		paramMap["RunningOrderId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("keyword"); ok {
		paramMap["Keyword"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_type"); ok {
		paramMap["OrderType"] = helper.String(v.(string))
	}

	service := OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusLogByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		cursor = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(cursor))
	if cursor != nil {
		_ = d.Set("cursor", cursor)
	}

	if listOver != nil {
		_ = d.Set("list_over", listOver)
	}

	if jobRequestId != nil {
		_ = d.Set("job_request_id", jobRequestId)
	}

	if jobInstanceList != nil {
		for _, jobInstanceForSubmissionLog := range jobInstanceList {
			jobInstanceForSubmissionLogMap := map[string]interface{}{}

			if jobInstanceForSubmissionLog.RunningOrderId != nil {
				jobInstanceForSubmissionLogMap["running_order_id"] = jobInstanceForSubmissionLog.RunningOrderId
			}

			if jobInstanceForSubmissionLog.JobInstanceStartTime != nil {
				jobInstanceForSubmissionLogMap["job_instance_start_time"] = jobInstanceForSubmissionLog.JobInstanceStartTime
			}

			if jobInstanceForSubmissionLog.StartingMillis != nil {
				jobInstanceForSubmissionLogMap["starting_millis"] = jobInstanceForSubmissionLog.StartingMillis
			}

			ids = append(ids, *jobInstanceForSubmissionLog.JobId)
			tmpList = append(tmpList, jobInstanceForSubmissionLogMap)
		}

		_ = d.Set("job_instance_list", tmpList)
	}

	if logList != nil {
		_ = d.Set("log_list", logList)
	}

	if logContentList != nil {
		for _, logContent := range logContentList {
			logContentMap := map[string]interface{}{}

			if logContent.Log != nil {
				logContentMap["log"] = logContent.Log
			}

			if logContent.Time != nil {
				logContentMap["time"] = logContent.Time
			}

			if logContent.PkgId != nil {
				logContentMap["pkg_id"] = logContent.PkgId
			}

			if logContent.PkgLogId != nil {
				logContentMap["pkg_log_id"] = logContent.PkgLogId
			}

			if logContent.ContainerName != nil {
				logContentMap["container_name"] = logContent.ContainerName
			}

			ids = append(ids, *logContent.JobId)
			tmpList = append(tmpList, logContentMap)
		}

		_ = d.Set("log_content_list", tmpList)
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
