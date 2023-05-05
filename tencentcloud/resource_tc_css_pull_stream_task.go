/*
Provides a resource to create a css pull_stream_task

Example Usage

```hcl
resource "tencentcloud_css_pull_stream_task" "pull_stream_task" {
  source_type = "source_type"
  source_urls = ["source_urls"]
  domain_name = "domain_name"
  app_name = "app_name"
  stream_name = "stream_name"
  start_time = "2022-11-16T22:09:28Z"
  end_time = "2022-11-16T22:09:28Z"
  operator = "admin"
  comment = "comment."
  }

```
Import

css pull_stream_task can be imported using the id, e.g.
```
$ terraform import tencentcloud_css_pull_stream_task.pull_stream_task pullStreamTask_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssPullStreamTask() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudCssPullStreamTaskRead,
		Create: resourceTencentCloudCssPullStreamTaskCreate,
		Update: resourceTencentCloudCssPullStreamTaskUpdate,
		Delete: resourceTencentCloudCssPullStreamTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "&amp;#39;PullLivePushLive&amp;#39;: SourceUrls live type, &amp;#39;PullVodPushLive&amp;#39;: SourceUrls vod type.",
			},

			"source_urls": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Pull Source media, SourceType=PullLivePushLive only 1 value, SourceType=PullLivePushLive can input multi values.",
			},

			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "push domain name.",
			},

			"app_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "push app name.",
			},

			"stream_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "push stream name.",
			},

			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "task begin time.",
			},

			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "task end time.",
			},

			"operator": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "desc operator user name.",
			},

			"push_args": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "other pushing args.",
			},

			"callback_events": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Computed:    true,
				Description: "defind the callback event you need, null for all. TaskStart, TaskExit, VodSourceFileStart, VodSourceFileFinish, ResetTaskConfig, PullFileUnstable, PushStreamUnstable, PullFileFailed, PushStreamFailed, FileEndEarly.",
			},

			"vod_loop_times": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "loop time for vod.",
			},

			"vod_refresh_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "vod refresh method. &amp;#39;ImmediateNewSource&amp;#39;: switch to new source at once, &amp;#39;ContinueBreakPoint&amp;#39;: switch to new source while old source finish.",
			},

			"callback_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "task event callback url.",
			},

			"extra_cmd": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ignore_region for ignore the input region and reblance inside the server.",
			},

			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "desc for pull task.",
			},

			"to_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "full target push url, DomainName, AppName, StreamName field must be empty.",
			},

			"backup_source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "backup pull source type.",
			},

			"backup_source_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "backup pull source.",
			},

			"watermark_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "watermark list, max 4 setting.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"picture_url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "watermark picture url.",
						},
						"x_position": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "x position.",
						},
						"y_position": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "y position.",
						},
						"width": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "pic width.",
						},
						"height": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "pic height.",
						},
						"location": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "position type, 0:left top, 1:right top, 2:right bot, 3: left bot.",
						},
					},
				},
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "task enable or disable.",
			},

			"create_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "desc who create the task.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "create time.",
			},

			"update_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "desc who update the task.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "update time.",
			},

			"file_index": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "task enable or disable.",
			},

			"offset_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "task enable or disable.",
			},

			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "task run region.",
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
		response *css.CreateLivePullStreamTaskResponse
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

	if v, ok := d.GetOk("vod_loop_times"); ok {
		request.VodLoopTimes = helper.Int64ToStrPoint(v.(int64))
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().CreateLivePullStreamTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s create css pullStreamTask failed, reason:%+v", logId, err)
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

	taskId := d.Id()

	result, err := service.DescribeCssPullStreamTask(ctx, taskId)

	if err != nil {
		return err
	}
	log.Printf("[CRITICAL]##########%v ", len(result))
	if result == nil || len(result) < 1 {
		d.SetId("")
		return fmt.Errorf("resource `pullStreamTask` %s does not exist", taskId)
	}

	pullStreamTask := result[0]
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

	// if pullStreamTask.Operator != nil {
	// 	_ = d.Set("operator", pullStreamTask.Operator)
	// }

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

	// if pullStreamTask.ExtraCmd != nil {
	// 	_ = d.Set("extra_cmd", pullStreamTask.ExtraCmd)
	// }

	if pullStreamTask.Comment != nil {
		_ = d.Set("comment", pullStreamTask.Comment)
	}

	// if pullStreamTask.ToUrl != nil {
	// 	_ = d.Set("to_url", pullStreamTask.ToUrl)
	// }

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
			if watermarkList.PictureUrl != nil {
				watermarkListMap["picture_url"] = watermarkList.PictureUrl
			}
			if watermarkList.XPosition != nil {
				watermarkListMap["x_position"] = watermarkList.XPosition
			}
			if watermarkList.YPosition != nil {
				watermarkListMap["y_position"] = watermarkList.YPosition
			}
			if watermarkList.Width != nil {
				watermarkListMap["width"] = watermarkList.Width
			}
			if watermarkList.Height != nil {
				watermarkListMap["height"] = watermarkList.Height
			}
			if watermarkList.Location != nil {
				watermarkListMap["location"] = watermarkList.Location
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

	// if pullStreamTask.FileIndex != nil {
	// 	_ = d.Set("file_index", pullStreamTask.FileIndex)
	// }

	// if pullStreamTask.OffsetTime != nil {
	// 	_ = d.Set("offset_time", pullStreamTask.OffsetTime)
	// }

	if pullStreamTask.Region != nil {
		_ = d.Set("region", pullStreamTask.Region)
	}

	return nil
}

func resourceTencentCloudCssPullStreamTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_pull_stream_task.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	// ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := css.NewModifyLivePullStreamTaskRequest()

	taskId := d.Id()

	request.TaskId = &taskId

	if d.HasChange("source_type") {
		return fmt.Errorf("`source_type` do not support change now.")
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

	if d.HasChange("domain_name") {
		return fmt.Errorf("`domain_name` do not support change now.")
	}

	if d.HasChange("app_name") {
		return fmt.Errorf("`app_name` do not support change now.")
	}

	if d.HasChange("stream_name") {
		return fmt.Errorf("`stream_name` do not support change now.")
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

	if d.HasChange("push_args") {
		return fmt.Errorf("`push_args` do not support change now.")
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
		if v, ok := d.GetOk("vod_loop_times"); ok {
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

	if d.HasChange("extra_cmd") {
		return fmt.Errorf("`extra_cmd` do not support change now.")
	}

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}
	}

	if d.HasChange("to_url") {
		return fmt.Errorf("`to_url` do not support change now.")
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
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			request.Status = helper.String(v.(string))
		}
	}

	if d.HasChange("file_index") {
		if v, ok := d.GetOk("file_index"); ok {
			request.FileIndex = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("offset_time") {
		if v, ok := d.GetOk("offset_time"); ok {
			request.OffsetTime = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().ModifyLivePullStreamTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s create css pullStreamTask failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCssPullStreamTaskRead(d, meta)
}

func resourceTencentCloudCssPullStreamTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_pull_stream_task.delete")()
	defer inconsistentCheck(d, meta)()

	var operator *string
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	taskId := d.Id()
	if v, ok := d.GetOk("operator"); ok {
		operator = helper.String(v.(string))
	}

	if err := service.DeleteCssPullStreamTaskById(ctx, helper.String(taskId), operator); err != nil {
		return err
	}

	return nil
}
