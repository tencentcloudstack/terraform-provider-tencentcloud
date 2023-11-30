package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudOceanusJobSubmissionLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusJobSubmissionLogRead,
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
				Default:     0,
				Description: "Job instance ID.",
			},
			"keyword": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Keyword, default empty.",
			},
			"cursor": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Cursor, default empty, first request does not need to pass in.",
			},
			"order_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "asc",
				Description: "Sorting method, default asc, asc: ascending, desc: descending.",
			},
			"list_over": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether the list is over.",
			},
			"job_request_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Request ID of starting job.",
			},
			"log_list": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Log list, deprecated.",
			},
			"job_instance_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Job instance list during the specified time period.",
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
			"log_content_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The list of log contents.",
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

func dataSourceTencentCloudOceanusJobSubmissionLogRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_oceanus_job_submission_log.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId            = getLogId(contextNil)
		ctx              = context.WithValue(context.TODO(), logIdKey, logId)
		service          = OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}
		jobSubmissionLog *oceanus.DescribeJobSubmissionLogResponseParams
		jobId            string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("job_id"); ok {
		paramMap["JobId"] = helper.String(v.(string))
		jobId = v.(string)
	}

	if v, ok := d.GetOkExists("start_time"); ok {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("running_order_id"); ok {
		paramMap["RunningOrderId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("keyword"); ok {
		paramMap["Keyword"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_type"); ok {
		paramMap["OrderType"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("cursor"); v != nil {
		paramMap["Cursor"] = helper.String(v.(string))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusJobSubmissionLogByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		jobSubmissionLog = result
		return nil
	})

	if err != nil {
		return err
	}

	if jobSubmissionLog.Cursor != nil {
		_ = d.Set("cursor", jobSubmissionLog.Cursor)
	}

	if jobSubmissionLog.ListOver != nil {
		_ = d.Set("list_over", jobSubmissionLog.ListOver)
	}

	if jobSubmissionLog.JobRequestId != nil {
		_ = d.Set("job_request_id", jobSubmissionLog.JobRequestId)
	}

	if jobSubmissionLog.LogList != nil {
		tmpList := make([]string, 0, len(jobSubmissionLog.LogList))
		for _, log := range jobSubmissionLog.LogList {
			tmpList = append(tmpList, *log)
		}

		_ = d.Set("log_list", tmpList)
	}

	if jobSubmissionLog.JobInstanceList != nil {
		tmpList := make([]map[string]interface{}, 0, len(jobSubmissionLog.JobInstanceList))
		for _, jobInstance := range jobSubmissionLog.JobInstanceList {
			jobInstanceMap := map[string]interface{}{}

			if jobInstance.RunningOrderId != nil {
				jobInstanceMap["running_order_id"] = jobInstance.RunningOrderId
			}

			if jobInstance.JobInstanceStartTime != nil {
				jobInstanceMap["job_instance_start_time"] = jobInstance.JobInstanceStartTime
			}

			if jobInstance.StartingMillis != nil {
				jobInstanceMap["starting_millis"] = jobInstance.StartingMillis
			}

			tmpList = append(tmpList, jobInstanceMap)
		}

		_ = d.Set("job_instance_list", tmpList)
	}

	if jobSubmissionLog.LogContentList != nil {
		tmpList := make([]map[string]interface{}, 0, len(jobSubmissionLog.LogContentList))
		for _, logContent := range jobSubmissionLog.LogContentList {
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

			tmpList = append(tmpList, logContentMap)
		}

		_ = d.Set("log_content_list", tmpList)
	}

	d.SetId(jobId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
