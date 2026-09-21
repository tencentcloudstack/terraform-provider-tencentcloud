package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoLogAnalysisDownloadTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoLogAnalysisDownloadTaskCreate,
		Read:   resourceTencentCloudTeoLogAnalysisDownloadTaskRead,
		Update: resourceTencentCloudTeoLogAnalysisDownloadTaskUpdate,
		Delete: resourceTencentCloudTeoLogAnalysisDownloadTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

			"area": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Data region, valid values: `mainland`, `overseas`.",
			},

			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Start time, example: `2020-04-29T00:00:00Z`.",
			},

			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "End time, example: `2020-04-30T00:00:00Z`. The maximum time span from start time to end time for a single query is 31 days.",
			},

			"log_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Log type, valid values: `l7-access-logs` (L7 access logs), `web-attack` (managed rules logs). Default: `l7-access-logs`.",
			},

			"condition": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Log match condition, max length 12KB.",
			},

			"format": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "File format, valid values: `csv`. Default: `csv`.",
			},

			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Time sort of raw logs, valid values: `asc`, `desc`. Default: `desc`.",
			},

			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Log analysis download task ID.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task status, valid values: `loading`, `failed`, `completed`.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task creation time. Tasks are retained for 3 days after creation.",
			},

			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Download URL, only returned when `status = completed`.",
			},

			"expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Download task expiration time. The download URL is unavailable after expiration.",
			},
		},
	}
}

func resourceTencentCloudTeoLogAnalysisDownloadTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_log_analysis_download_task.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = teov20220901.NewCreateLogAnalysisDownloadTaskRequest()
		response = teov20220901.NewCreateLogAnalysisDownloadTaskResponse()
		zoneId   string
		area     string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(zoneId)
	}

	if v, ok := d.GetOk("area"); ok {
		area = v.(string)
		request.Area = helper.String(area)
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_type"); ok {
		request.LogType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("condition"); ok {
		request.Condition = helper.String(v.(string))
	}

	if v, ok := d.GetOk("format"); ok {
		request.Format = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort"); ok {
		request.Sort = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateLogAnalysisDownloadTaskWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo log analysis download task failed, Response is nil."))
		}

		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo log analysis download task failed, reason:%+v", logId, err)
		return err
	}

	log.Printf("[DEBUG]%s teo_log_analysis_download_task d.Id()=%s", logId, d.Id())

	if response.Response.TaskId == nil || *response.Response.TaskId == "" {
		return fmt.Errorf("Create teo log analysis download task failed, TaskId is empty.")
	}
	taskId := *response.Response.TaskId
	d.SetId(strings.Join([]string{zoneId, area, taskId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoLogAnalysisDownloadTaskRead(d, meta)
}

func resourceTencentCloudTeoLogAnalysisDownloadTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_log_analysis_download_task.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	area := idSplit[1]
	taskId := idSplit[2]

	respData, err := service.DescribeTeoLogAnalysisDownloadTaskById(ctx, zoneId, area, taskId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] teo_log_analysis_download_task id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("area", area)

	if respData.TaskId != nil {
		_ = d.Set("task_id", respData.TaskId)
	}
	if respData.StartTime != nil {
		_ = d.Set("start_time", respData.StartTime)
	}
	if respData.EndTime != nil {
		_ = d.Set("end_time", respData.EndTime)
	}
	if respData.LogType != nil {
		_ = d.Set("log_type", respData.LogType)
	}
	if respData.Condition != nil {
		_ = d.Set("condition", respData.Condition)
	}
	if respData.Format != nil {
		_ = d.Set("format", respData.Format)
	}
	if respData.Sort != nil {
		_ = d.Set("sort", respData.Sort)
	}
	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}
	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}
	if respData.Url != nil {
		_ = d.Set("url", respData.Url)
	}
	if respData.ExpireTime != nil {
		_ = d.Set("expire_time", respData.ExpireTime)
	}

	return nil
}

func resourceTencentCloudTeoLogAnalysisDownloadTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_log_analysis_download_task.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"zone_id", "area", "start_time", "end_time", "log_type", "condition", "format", "sort"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("teo_log_analysis_download_task field `%s` is immutable, please recreate the resource", v)
		}
	}

	return resourceTencentCloudTeoLogAnalysisDownloadTaskRead(d, meta)
}

func resourceTencentCloudTeoLogAnalysisDownloadTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_log_analysis_download_task.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
