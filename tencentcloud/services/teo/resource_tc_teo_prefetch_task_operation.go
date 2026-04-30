package teo

import (
	"fmt"
	"log"
	"time"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudTeoPrefetchTaskOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoPrefetchTaskOperationCreate,
		Read:   resourceTencentCloudTeoPrefetchTaskOperationRead,
		Delete: resourceTencentCloudTeoPrefetchTaskOperationDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},
			"targets": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of URLs to prefetch. Each element format is like: http://www.example.com/example.txt.",
			},
			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Prefetch mode. Valid values: `default` (prefetch to middle layer), `edge` (prefetch to edge and middle layer). Default: `default`.",
			},
			"headers": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "HTTP header name.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "HTTP header value.",
						},
					},
				},
				Description: "HTTP headers to carry during prefetch.",
			},
			"prefetch_media_segments": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Media segment prefetch control. Valid values: `on` (enable segment prefetch), `off` (only prefetch the submitted description file). Default: `off`.",
			},
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task job ID returned by CreatePrefetchTask.",
			},
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Task result list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task ID.",
						},
						"target": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource URL.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task type.",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cache purge method.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task status. Valid values: processing, success, failed, timeout, canceled.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task completion time.",
						},
						"fail_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Failure type.",
						},
						"fail_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Failure reason description.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoPrefetchTaskOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_prefetch_task_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Get("zone_id").(string)
	targetsInterface := d.Get("targets").([]interface{})

	request := teov20220901.NewCreatePrefetchTaskRequest()
	request.ZoneId = helper.String(zoneId)

	targets := make([]*string, 0, len(targetsInterface))
	for _, t := range targetsInterface {
		targets = append(targets, helper.String(t.(string)))
	}
	request.Targets = targets

	if v, ok := d.GetOk("mode"); ok {
		request.Mode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("headers"); ok {
		headersInterface := v.([]interface{})
		headers := make([]*teov20220901.Header, 0, len(headersInterface))
		for _, h := range headersInterface {
			hMap := h.(map[string]interface{})
			header := &teov20220901.Header{}
			if name, ok := hMap["name"]; ok {
				header.Name = helper.String(name.(string))
			}
			if value, ok := hMap["value"]; ok {
				header.Value = helper.String(value.(string))
			}
			headers = append(headers, header)
		}
		request.Headers = headers
	}

	if v, ok := d.GetOk("prefetch_media_segments"); ok {
		request.PrefetchMediaSegments = helper.String(v.(string))
	}

	var jobId string
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreatePrefetchTask(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("CreatePrefetchTask API returned empty response"))
		}
		if result.Response.JobId != nil {
			jobId = *result.Response.JobId
		}
		return nil
	})
	if err != nil {
		return err
	}

	if jobId == "" {
		return fmt.Errorf("CreatePrefetchTask API returned empty JobId")
	}

	_ = d.Set("job_id", jobId)

	// Poll DescribePrefetchTasks until the task is complete
	describeRequest := teov20220901.NewDescribePrefetchTasksRequest()
	describeRequest.ZoneId = helper.String(zoneId)
	describeRequest.Limit = helper.Int64(int64(1000))
	describeRequest.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("job-id"),
			Values: []*string{helper.String(jobId)},
		},
	}

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribePrefetchTasks(describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), describeRequest.ToJsonString(), result.ToJsonString())
		}

		if result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribePrefetchTasks API returned empty response"))
		}

		if result.Response.Tasks == nil || len(result.Response.Tasks) == 0 {
			return resource.RetryableError(fmt.Errorf("prefetch task is still processing, no tasks returned yet"))
		}

		allDone := true
		for _, task := range result.Response.Tasks {
			if task.Status != nil && *task.Status == "processing" {
				allDone = false
				break
			}
		}

		if !allDone {
			return resource.RetryableError(fmt.Errorf("prefetch task is still processing"))
		}

		// Check for failed/timeout statuses
		for _, task := range result.Response.Tasks {
			if task.Status != nil {
				status := *task.Status
				if status == "failed" {
					failType := ""
					failMessage := ""
					if task.FailType != nil {
						failType = *task.FailType
					}
					if task.FailMessage != nil {
						failMessage = *task.FailMessage
					}
					return resource.NonRetryableError(fmt.Errorf("prefetch task failed, fail_type: %s, fail_message: %s", failType, failMessage))
				}
				if status == "timeout" {
					return resource.NonRetryableError(fmt.Errorf("prefetch task timed out"))
				}
			}
		}

		// Set tasks computed fields
		tasksList := make([]map[string]interface{}, 0, len(result.Response.Tasks))
		for _, task := range result.Response.Tasks {
			taskMap := map[string]interface{}{}
			if task.JobId != nil {
				taskMap["job_id"] = *task.JobId
			}
			if task.Target != nil {
				taskMap["target"] = *task.Target
			}
			if task.Type != nil {
				taskMap["type"] = *task.Type
			}
			if task.Method != nil {
				taskMap["method"] = *task.Method
			}
			if task.Status != nil {
				taskMap["status"] = *task.Status
			}
			if task.CreateTime != nil {
				taskMap["create_time"] = *task.CreateTime
			}
			if task.UpdateTime != nil {
				taskMap["update_time"] = *task.UpdateTime
			}
			if task.FailType != nil {
				taskMap["fail_type"] = *task.FailType
			}
			if task.FailMessage != nil {
				taskMap["fail_message"] = *task.FailMessage
			}
			tasksList = append(tasksList, taskMap)
		}
		_ = d.Set("tasks", tasksList)

		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(zoneId + tccommon.FILED_SP + jobId)

	return resourceTencentCloudTeoPrefetchTaskOperationRead(d, meta)
}

func resourceTencentCloudTeoPrefetchTaskOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_prefetch_task_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoPrefetchTaskOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_prefetch_task_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
