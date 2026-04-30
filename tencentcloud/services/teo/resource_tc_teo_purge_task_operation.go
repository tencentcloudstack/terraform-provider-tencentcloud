package teo

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoPurgeTaskOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoPurgeTaskCreate,
		Read:   resourceTencentCloudTeoPurgeTaskRead,
		Delete: resourceTencentCloudTeoPurgeTaskDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Purge type. Valid values: purge_url, purge_prefix, purge_host, purge_all, purge_cache_tag.",
			},
			"method": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Purge method. Valid values: invalidate, delete. Default: invalidate.",
			},
			"targets": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of targets to purge.",
			},
			"cache_tag": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Cache tag configuration, required when type is purge_cache_tag.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domains": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Domain list for cache tag.",
						},
					},
				},
			},
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task ID returned by CreatePurgeTask.",
			},
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of purge task results.",
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
							Description: "Purge target.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Purge type.",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Purge method.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task status.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task update time.",
						},
						"fail_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Failure type.",
						},
						"fail_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Failure message.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoPurgeTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_purge_task_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := teov20220901.NewCreatePurgeTaskRequest()
	zoneId := d.Get("zone_id").(string)
	request.ZoneId = helper.String(zoneId)

	purgeType := d.Get("type").(string)
	request.Type = helper.String(purgeType)

	if v, ok := d.GetOk("method"); ok {
		request.Method = helper.String(v.(string))
	}

	if v, ok := d.GetOk("targets"); ok {
		targets := v.([]interface{})
		for _, target := range targets {
			request.Targets = append(request.Targets, helper.String(target.(string)))
		}
	}

	if v, ok := d.GetOk("cache_tag"); ok {
		cacheTagList := v.([]interface{})
		if len(cacheTagList) > 0 {
			cacheTagMap := cacheTagList[0].(map[string]interface{})
			cacheTag := &teov20220901.CacheTag{}
			if domains, ok := cacheTagMap["domains"].([]interface{}); ok && len(domains) > 0 {
				for _, domain := range domains {
					cacheTag.Domains = append(cacheTag.Domains, helper.String(domain.(string)))
				}
			}
			request.CacheTag = cacheTag
		}
	}

	var jobId string
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreatePurgeTask(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result.Response == nil || result.Response.JobId == nil {
			return resource.NonRetryableError(fmt.Errorf("CreatePurgeTask response is empty or JobId is nil"))
		}
		jobId = *result.Response.JobId
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(helper.BuildToken())
	_ = d.Set("job_id", jobId)

	// Poll DescribePurgeTasks to wait for task completion
	describeRequest := teov20220901.NewDescribePurgeTasksRequest()
	describeRequest.ZoneId = helper.String(zoneId)
	describeRequest.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("job-id"),
			Values: []*string{helper.String(jobId)},
		},
	}

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribePurgeTasks(describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result.Response == nil || result.Response.Tasks == nil || len(result.Response.Tasks) == 0 {
			return resource.RetryableError(fmt.Errorf("DescribePurgeTasks returned empty tasks, waiting for task to be created"))
		}

		allDone := true
		for _, task := range result.Response.Tasks {
			if task.Status != nil {
				status := *task.Status
				if status == "processing" {
					allDone = false
					break
				} else if status == "failed" {
					return resource.NonRetryableError(fmt.Errorf("purge task failed, fail_type: %s, fail_message: %s", *task.FailType, *task.FailMessage))
				} else if status == "timeout" {
					return resource.NonRetryableError(fmt.Errorf("purge task timed out"))
				} else if status == "canceled" {
					return resource.NonRetryableError(fmt.Errorf("purge task was canceled"))
				}
			}
		}

		if !allDone {
			return resource.RetryableError(fmt.Errorf("purge task is still processing"))
		}

		// Set tasks computed attribute
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

	return resourceTencentCloudTeoPurgeTaskRead(d, meta)
}

func resourceTencentCloudTeoPurgeTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_purge_task_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()
	return nil
}

func resourceTencentCloudTeoPurgeTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_purge_task_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()
	return nil
}
