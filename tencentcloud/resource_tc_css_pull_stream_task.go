/*
Provides a resource to create a css pull_stream_task

Example Usage

```hcl
resource "tencentcloud_css_pull_stream_task" "pull_stream_task" {
  source_type = "PullLivePushLive"
  source_urls = &lt;nil&gt;
  domain_name = &lt;nil&gt;
  app_name = &lt;nil&gt;
  stream_name = &lt;nil&gt;
  start_time = &lt;nil&gt;
  end_time = &lt;nil&gt;
  operator = &lt;nil&gt;
  push_args = &lt;nil&gt;
  callback_events = &lt;nil&gt;
  vod_loop_times = &lt;nil&gt;
  vod_refresh_type = &lt;nil&gt;
  callback_url = &lt;nil&gt;
  extra_cmd = &lt;nil&gt;
  comment = &lt;nil&gt;
  to_url = &lt;nil&gt;
  backup_source_type = &lt;nil&gt;
  backup_source_url = &lt;nil&gt;
  watermark_list {
		picture_url = &lt;nil&gt;
		x_position = &lt;nil&gt;
		y_position = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		location = &lt;nil&gt;

  }
  status = &lt;nil&gt;
          file_index = &lt;nil&gt;
  offset_time = &lt;nil&gt;
  }
```

Import

css pull_stream_task can be imported using the id, e.g.

```
terraform import tencentcloud_css_pull_stream_task.pull_stream_task pull_stream_task_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCssPullStreamTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssPullStreamTaskCreate,
		Read:   resourceTencentCloudCssPullStreamTaskRead,
		Update: resourceTencentCloudCssPullStreamTaskUpdate,
		Delete: resourceTencentCloudCssPullStreamTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"source_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "&amp;amp;#39;PullLivePushLive&amp;amp;#39;: SourceUrls live type, &amp;amp;#39;PullVodPushLive&amp;amp;#39;: SourceUrls vod type.",
			},

			"source_urls": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Pull Source media, SourceType=PullLivePushLive only 1 value, SourceType=PullLivePushLive can input multi values.",
			},

			"domain_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Push domain name.",
			},

			"app_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Push app name.",
			},

			"stream_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Push stream name.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task begin time.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task end time.",
			},

			"operator": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Desc operator user name.",
			},

			"push_args": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Other pushing args.",
			},

			"callback_events": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Defind the callback event you need, null for all. TaskStart, TaskExit, VodSourceFileStart, VodSourceFileFinish, ResetTaskConfig, PullFileUnstable, PushStreamUnstable, PullFileFailed, PushStreamFailed, FileEndEarly.",
			},

			"vod_loop_times": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Loop time for vod.",
			},

			"vod_refresh_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Vod refresh method. &amp;amp;#39;ImmediateNewSource&amp;amp;#39;: switch to new source at once, &amp;amp;#39;ContinueBreakPoint&amp;amp;#39;: switch to new source while old source finish.",
			},

			"callback_url": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Task event callback url.",
			},

			"extra_cmd": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Ignore_region for ignore the input region and reblance inside the server.",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Desc for pull task.",
			},

			"to_url": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Full target push url, DomainName, AppName, StreamName field must be empty.",
			},

			"backup_source_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Backup pull source type.",
			},

			"backup_source_url": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Backup pull source.",
			},

			"watermark_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Watermark list, max 4 setting.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"picture_url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Watermark picture url.",
						},
						"x_position": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "X position.",
						},
						"y_position": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Y position.",
						},
						"width": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Pic width.",
						},
						"height": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Pic height.",
						},
						"location": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Position type, 0:left top, 1:right top, 2:right bot, 3: left bot.",
						},
					},
				},
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Task enable or disable.",
			},

			"create_by": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Desc who create the task.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Create time.",
			},

			"update_by": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Desc who update the task.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Update time.",
			},

			"file_index": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Task enable or disable.",
			},

			"offset_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Task enable or disable.",
			},

			"region": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Task run region.",
			},
		},
	}
}

func resourceTencentCloudCssPullStreamTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_pull_stream_task.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = css.NewCreateLivePullStreamTaskRequest()
		response = css.NewCreateLivePullStreamTaskResponse()
		taskId   string
	)
	if v, ok := d.GetOk("source_type"); ok {
		request.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_urls"); ok {
		sourceUrlsSet := v.(*schema.Set).List()
		for i := range sourceUrlsSet {
			sourceUrls := sourceUrlsSet[i].(string)
			request.SourceUrls = append(request.SourceUrls, &sourceUrls)
		}
	}

	if v, ok := d.GetOk("domain_name"); ok {
		request.DomainName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("app_name"); ok {
		request.AppName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("stream_name"); ok {
		request.StreamName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("operator"); ok {
		request.Operator = helper.String(v.(string))
	}

	if v, ok := d.GetOk("push_args"); ok {
		request.PushArgs = helper.String(v.(string))
	}

	if v, ok := d.GetOk("callback_events"); ok {
		callbackEventsSet := v.(*schema.Set).List()
		for i := range callbackEventsSet {
			callbackEvents := callbackEventsSet[i].(string)
			request.CallbackEvents = append(request.CallbackEvents, &callbackEvents)
		}
	}

	if v, ok := d.GetOkExists("vod_loop_times"); ok {
		request.VodLoopTimes = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vod_refresh_type"); ok {
		request.VodRefreshType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("callback_url"); ok {
		request.CallbackUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("extra_cmd"); ok {
		request.ExtraCmd = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if v, ok := d.GetOk("to_url"); ok {
		request.ToUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_source_type"); ok {
		request.BackupSourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_source_url"); ok {
		request.BackupSourceUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("watermark_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			pullPushWatermarkInfo := css.PullPushWatermarkInfo{}
			if v, ok := dMap["picture_url"]; ok {
				pullPushWatermarkInfo.PictureUrl = helper.String(v.(string))
			}
			if v, ok := dMap["x_position"]; ok {
				pullPushWatermarkInfo.XPosition = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["y_position"]; ok {
				pullPushWatermarkInfo.YPosition = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["width"]; ok {
				pullPushWatermarkInfo.Width = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["height"]; ok {
				pullPushWatermarkInfo.Height = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["location"]; ok {
				pullPushWatermarkInfo.Location = helper.IntInt64(v.(int))
			}
			request.WatermarkList = append(request.WatermarkList, &pullPushWatermarkInfo)
		}
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("file_index"); ok {
		request.FileIndex = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("offset_time"); ok {
		request.OffsetTime = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().CreateLivePullStreamTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css pullStreamTask failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudCssPullStreamTaskRead(d, meta)
}

func resourceTencentCloudCssPullStreamTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_pull_stream_task.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	pullStreamTaskId := d.Id()

	pullStreamTask, err := service.DescribeCssPullStreamTaskById(ctx, taskId)
	if err != nil {
		return err
	}

	if pullStreamTask == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssPullStreamTask` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if pullStreamTask.SourceType != nil {
		_ = d.Set("source_type", pullStreamTask.SourceType)
	}

	if pullStreamTask.SourceUrls != nil {
		_ = d.Set("source_urls", pullStreamTask.SourceUrls)
	}

	if pullStreamTask.DomainName != nil {
		_ = d.Set("domain_name", pullStreamTask.DomainName)
	}

	if pullStreamTask.AppName != nil {
		_ = d.Set("app_name", pullStreamTask.AppName)
	}

	if pullStreamTask.StreamName != nil {
		_ = d.Set("stream_name", pullStreamTask.StreamName)
	}

	if pullStreamTask.StartTime != nil {
		_ = d.Set("start_time", pullStreamTask.StartTime)
	}

	if pullStreamTask.EndTime != nil {
		_ = d.Set("end_time", pullStreamTask.EndTime)
	}

	if pullStreamTask.Operator != nil {
		_ = d.Set("operator", pullStreamTask.Operator)
	}

	if pullStreamTask.PushArgs != nil {
		_ = d.Set("push_args", pullStreamTask.PushArgs)
	}

	if pullStreamTask.CallbackEvents != nil {
		_ = d.Set("callback_events", pullStreamTask.CallbackEvents)
	}

	if pullStreamTask.VodLoopTimes != nil {
		_ = d.Set("vod_loop_times", pullStreamTask.VodLoopTimes)
	}

	if pullStreamTask.VodRefreshType != nil {
		_ = d.Set("vod_refresh_type", pullStreamTask.VodRefreshType)
	}

	if pullStreamTask.CallbackUrl != nil {
		_ = d.Set("callback_url", pullStreamTask.CallbackUrl)
	}

	if pullStreamTask.ExtraCmd != nil {
		_ = d.Set("extra_cmd", pullStreamTask.ExtraCmd)
	}

	if pullStreamTask.Comment != nil {
		_ = d.Set("comment", pullStreamTask.Comment)
	}

	if pullStreamTask.ToUrl != nil {
		_ = d.Set("to_url", pullStreamTask.ToUrl)
	}

	if pullStreamTask.BackupSourceType != nil {
		_ = d.Set("backup_source_type", pullStreamTask.BackupSourceType)
	}

	if pullStreamTask.BackupSourceUrl != nil {
		_ = d.Set("backup_source_url", pullStreamTask.BackupSourceUrl)
	}

	if pullStreamTask.WatermarkList != nil {
		watermarkListList := []interface{}{}
		for _, watermarkList := range pullStreamTask.WatermarkList {
			watermarkListMap := map[string]interface{}{}

			if pullStreamTask.WatermarkList.PictureUrl != nil {
				watermarkListMap["picture_url"] = pullStreamTask.WatermarkList.PictureUrl
			}

			if pullStreamTask.WatermarkList.XPosition != nil {
				watermarkListMap["x_position"] = pullStreamTask.WatermarkList.XPosition
			}

			if pullStreamTask.WatermarkList.YPosition != nil {
				watermarkListMap["y_position"] = pullStreamTask.WatermarkList.YPosition
			}

			if pullStreamTask.WatermarkList.Width != nil {
				watermarkListMap["width"] = pullStreamTask.WatermarkList.Width
			}

			if pullStreamTask.WatermarkList.Height != nil {
				watermarkListMap["height"] = pullStreamTask.WatermarkList.Height
			}

			if pullStreamTask.WatermarkList.Location != nil {
				watermarkListMap["location"] = pullStreamTask.WatermarkList.Location
			}

			watermarkListList = append(watermarkListList, watermarkListMap)
		}

		_ = d.Set("watermark_list", watermarkListList)

	}

	if pullStreamTask.Status != nil {
		_ = d.Set("status", pullStreamTask.Status)
	}

	if pullStreamTask.CreateBy != nil {
		_ = d.Set("create_by", pullStreamTask.CreateBy)
	}

	if pullStreamTask.CreateTime != nil {
		_ = d.Set("create_time", pullStreamTask.CreateTime)
	}

	if pullStreamTask.UpdateBy != nil {
		_ = d.Set("update_by", pullStreamTask.UpdateBy)
	}

	if pullStreamTask.UpdateTime != nil {
		_ = d.Set("update_time", pullStreamTask.UpdateTime)
	}

	if pullStreamTask.FileIndex != nil {
		_ = d.Set("file_index", pullStreamTask.FileIndex)
	}

	if pullStreamTask.OffsetTime != nil {
		_ = d.Set("offset_time", pullStreamTask.OffsetTime)
	}

	if pullStreamTask.Region != nil {
		_ = d.Set("region", pullStreamTask.Region)
	}

	return nil
}

func resourceTencentCloudCssPullStreamTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_pull_stream_task.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := css.NewModifyLivePullStreamTaskRequest()

	pullStreamTaskId := d.Id()

	request.TaskId = &taskId

	immutableArgs := []string{"source_type", "source_urls", "domain_name", "app_name", "stream_name", "start_time", "end_time", "operator", "push_args", "callback_events", "vod_loop_times", "vod_refresh_type", "callback_url", "extra_cmd", "comment", "to_url", "backup_source_type", "backup_source_url", "watermark_list", "status", "create_by", "create_time", "update_by", "update_time", "file_index", "offset_time", "region"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("source_urls") {
		if v, ok := d.GetOk("source_urls"); ok {
			sourceUrlsSet := v.(*schema.Set).List()
			for i := range sourceUrlsSet {
				sourceUrls := sourceUrlsSet[i].(string)
				request.SourceUrls = append(request.SourceUrls, &sourceUrls)
			}
		}
	}

	if d.HasChange("start_time") {
		if v, ok := d.GetOk("start_time"); ok {
			request.StartTime = helper.String(v.(string))
		}
	}

	if d.HasChange("end_time") {
		if v, ok := d.GetOk("end_time"); ok {
			request.EndTime = helper.String(v.(string))
		}
	}

	if d.HasChange("operator") {
		if v, ok := d.GetOk("operator"); ok {
			request.Operator = helper.String(v.(string))
		}
	}

	if d.HasChange("callback_events") {
		if v, ok := d.GetOk("callback_events"); ok {
			callbackEventsSet := v.(*schema.Set).List()
			for i := range callbackEventsSet {
				callbackEvents := callbackEventsSet[i].(string)
				request.CallbackEvents = append(request.CallbackEvents, &callbackEvents)
			}
		}
	}

	if d.HasChange("vod_loop_times") {
		if v, ok := d.GetOkExists("vod_loop_times"); ok {
			request.VodLoopTimes = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("vod_refresh_type") {
		if v, ok := d.GetOk("vod_refresh_type"); ok {
			request.VodRefreshType = helper.String(v.(string))
		}
	}

	if d.HasChange("callback_url") {
		if v, ok := d.GetOk("callback_url"); ok {
			request.CallbackUrl = helper.String(v.(string))
		}
	}

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}
	}

	if d.HasChange("backup_source_type") {
		if v, ok := d.GetOk("backup_source_type"); ok {
			request.BackupSourceType = helper.String(v.(string))
		}
	}

	if d.HasChange("backup_source_url") {
		if v, ok := d.GetOk("backup_source_url"); ok {
			request.BackupSourceUrl = helper.String(v.(string))
		}
	}

	if d.HasChange("watermark_list") {
		if v, ok := d.GetOk("watermark_list"); ok {
			for _, item := range v.([]interface{}) {
				pullPushWatermarkInfo := css.PullPushWatermarkInfo{}
				if v, ok := dMap["picture_url"]; ok {
					pullPushWatermarkInfo.PictureUrl = helper.String(v.(string))
				}
				if v, ok := dMap["x_position"]; ok {
					pullPushWatermarkInfo.XPosition = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["y_position"]; ok {
					pullPushWatermarkInfo.YPosition = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["width"]; ok {
					pullPushWatermarkInfo.Width = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["height"]; ok {
					pullPushWatermarkInfo.Height = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["location"]; ok {
					pullPushWatermarkInfo.Location = helper.IntInt64(v.(int))
				}
				request.WatermarkList = append(request.WatermarkList, &pullPushWatermarkInfo)
			}
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			request.Status = helper.String(v.(string))
		}
	}

	if d.HasChange("file_index") {
		if v, ok := d.GetOkExists("file_index"); ok {
			request.FileIndex = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("offset_time") {
		if v, ok := d.GetOkExists("offset_time"); ok {
			request.OffsetTime = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().ModifyLivePullStreamTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update css pullStreamTask failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCssPullStreamTaskRead(d, meta)
}

func resourceTencentCloudCssPullStreamTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_pull_stream_task.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}
	pullStreamTaskId := d.Id()

	if err := service.DeleteCssPullStreamTaskById(ctx, taskId); err != nil {
		return err
	}

	return nil
}
