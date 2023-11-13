/*
Provides a resource to create a cls alarm_notice

Example Usage

```hcl
resource "tencentcloud_cls_alarm_notice" "alarm_notice" {
  name = "notice"
  type = "Trigger"
  notice_receivers {
		receiver_type = "Uin"
		receiver_ids =
		receiver_channels =
		start_time = "00:00:00"
		end_time = "23:59:59"
		index = 1

  }
  web_callbacks {
		url = "http://www.testnotice.com/callback"
		callback_type = "WeCom"
		method = "POST"
		headers =
		body = "null"
		index = 10

  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

cls alarm_notice can be imported using the id, e.g.

```
terraform import tencentcloud_cls_alarm_notice.alarm_notice alarm_notice_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudClsAlarmNotice() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsAlarmNoticeCreate,
		Read:   resourceTencentCloudClsAlarmNoticeRead,
		Update: resourceTencentCloudClsAlarmNoticeUpdate,
		Delete: resourceTencentCloudClsAlarmNoticeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Alarm notice name.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Notice type.",
			},

			"notice_receivers": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Notice receivers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"receiver_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Receiver type, Uin or Group.",
						},
						"receiver_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Required:    true,
							Description: "Receiver id.",
						},
						"receiver_channels": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Receiver channels, Email,Sms,WeChat or Phone.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Start time allowed to receive messages.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "End time allowed to receive messages.",
						},
						"index": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Index.",
						},
					},
				},
			},

			"web_callbacks": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Callback info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Callback url.",
						},
						"callback_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Callback type, WeCom or Http.",
						},
						"method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Method, POST or PUT.",
						},
						"headers": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Abandoned.",
						},
						"body": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Abandoned.",
						},
						"index": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Index.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudClsAlarmNoticeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_alarm_notice.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = cls.NewCreateAlarmNoticeRequest()
		response      = cls.NewCreateAlarmNoticeResponse()
		alarmNoticeId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("notice_receivers"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			noticeReceiver := cls.NoticeReceiver{}
			if v, ok := dMap["receiver_type"]; ok {
				noticeReceiver.ReceiverType = helper.String(v.(string))
			}
			if v, ok := dMap["receiver_ids"]; ok {
				receiverIdsSet := v.(*schema.Set).List()
				for i := range receiverIdsSet {
					receiverIds := receiverIdsSet[i].(int)
					noticeReceiver.ReceiverIds = append(noticeReceiver.ReceiverIds, helper.IntInt64(receiverIds))
				}
			}
			if v, ok := dMap["receiver_channels"]; ok {
				receiverChannelsSet := v.(*schema.Set).List()
				for i := range receiverChannelsSet {
					receiverChannels := receiverChannelsSet[i].(string)
					noticeReceiver.ReceiverChannels = append(noticeReceiver.ReceiverChannels, &receiverChannels)
				}
			}
			if v, ok := dMap["start_time"]; ok {
				noticeReceiver.StartTime = helper.String(v.(string))
			}
			if v, ok := dMap["end_time"]; ok {
				noticeReceiver.EndTime = helper.String(v.(string))
			}
			if v, ok := dMap["index"]; ok {
				noticeReceiver.Index = helper.IntInt64(v.(int))
			}
			request.NoticeReceivers = append(request.NoticeReceivers, &noticeReceiver)
		}
	}

	if v, ok := d.GetOk("web_callbacks"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			webCallback := cls.WebCallback{}
			if v, ok := dMap["url"]; ok {
				webCallback.Url = helper.String(v.(string))
			}
			if v, ok := dMap["callback_type"]; ok {
				webCallback.CallbackType = helper.String(v.(string))
			}
			if v, ok := dMap["method"]; ok {
				webCallback.Method = helper.String(v.(string))
			}
			if v, ok := dMap["headers"]; ok {
				headersSet := v.(*schema.Set).List()
				for i := range headersSet {
					headers := headersSet[i].(string)
					webCallback.Headers = append(webCallback.Headers, &headers)
				}
			}
			if v, ok := dMap["body"]; ok {
				webCallback.Body = helper.String(v.(string))
			}
			if v, ok := dMap["index"]; ok {
				webCallback.Index = helper.IntInt64(v.(int))
			}
			request.WebCallbacks = append(request.WebCallbacks, &webCallback)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().CreateAlarmNotice(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cls alarmNotice failed, reason:%+v", logId, err)
		return err
	}

	alarmNoticeId = *response.Response.AlarmNoticeId
	d.SetId(alarmNoticeId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::cls:%s:uin/:alarmNotice/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudClsAlarmNoticeRead(d, meta)
}

func resourceTencentCloudClsAlarmNoticeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_alarm_notice.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}

	alarmNoticeId := d.Id()

	alarmNotice, err := service.DescribeClsAlarmNoticeById(ctx, alarmNoticeId)
	if err != nil {
		return err
	}

	if alarmNotice == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsAlarmNotice` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if alarmNotice.Name != nil {
		_ = d.Set("name", alarmNotice.Name)
	}

	if alarmNotice.Type != nil {
		_ = d.Set("type", alarmNotice.Type)
	}

	if alarmNotice.NoticeReceivers != nil {
		noticeReceiversList := []interface{}{}
		for _, noticeReceivers := range alarmNotice.NoticeReceivers {
			noticeReceiversMap := map[string]interface{}{}

			if alarmNotice.NoticeReceivers.ReceiverType != nil {
				noticeReceiversMap["receiver_type"] = alarmNotice.NoticeReceivers.ReceiverType
			}

			if alarmNotice.NoticeReceivers.ReceiverIds != nil {
				noticeReceiversMap["receiver_ids"] = alarmNotice.NoticeReceivers.ReceiverIds
			}

			if alarmNotice.NoticeReceivers.ReceiverChannels != nil {
				noticeReceiversMap["receiver_channels"] = alarmNotice.NoticeReceivers.ReceiverChannels
			}

			if alarmNotice.NoticeReceivers.StartTime != nil {
				noticeReceiversMap["start_time"] = alarmNotice.NoticeReceivers.StartTime
			}

			if alarmNotice.NoticeReceivers.EndTime != nil {
				noticeReceiversMap["end_time"] = alarmNotice.NoticeReceivers.EndTime
			}

			if alarmNotice.NoticeReceivers.Index != nil {
				noticeReceiversMap["index"] = alarmNotice.NoticeReceivers.Index
			}

			noticeReceiversList = append(noticeReceiversList, noticeReceiversMap)
		}

		_ = d.Set("notice_receivers", noticeReceiversList)

	}

	if alarmNotice.WebCallbacks != nil {
		webCallbacksList := []interface{}{}
		for _, webCallbacks := range alarmNotice.WebCallbacks {
			webCallbacksMap := map[string]interface{}{}

			if alarmNotice.WebCallbacks.Url != nil {
				webCallbacksMap["url"] = alarmNotice.WebCallbacks.Url
			}

			if alarmNotice.WebCallbacks.CallbackType != nil {
				webCallbacksMap["callback_type"] = alarmNotice.WebCallbacks.CallbackType
			}

			if alarmNotice.WebCallbacks.Method != nil {
				webCallbacksMap["method"] = alarmNotice.WebCallbacks.Method
			}

			if alarmNotice.WebCallbacks.Headers != nil {
				webCallbacksMap["headers"] = alarmNotice.WebCallbacks.Headers
			}

			if alarmNotice.WebCallbacks.Body != nil {
				webCallbacksMap["body"] = alarmNotice.WebCallbacks.Body
			}

			if alarmNotice.WebCallbacks.Index != nil {
				webCallbacksMap["index"] = alarmNotice.WebCallbacks.Index
			}

			webCallbacksList = append(webCallbacksList, webCallbacksMap)
		}

		_ = d.Set("web_callbacks", webCallbacksList)

	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "cls", "alarmNotice", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudClsAlarmNoticeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_alarm_notice.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cls.NewModifyAlarmNoticeRequest()

	alarmNoticeId := d.Id()

	request.AlarmNoticeId = &alarmNoticeId

	immutableArgs := []string{"name", "type", "notice_receivers", "web_callbacks"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}
	}

	if d.HasChange("notice_receivers") {
		if v, ok := d.GetOk("notice_receivers"); ok {
			for _, item := range v.([]interface{}) {
				noticeReceiver := cls.NoticeReceiver{}
				if v, ok := dMap["receiver_type"]; ok {
					noticeReceiver.ReceiverType = helper.String(v.(string))
				}
				if v, ok := dMap["receiver_ids"]; ok {
					receiverIdsSet := v.(*schema.Set).List()
					for i := range receiverIdsSet {
						receiverIds := receiverIdsSet[i].(int)
						noticeReceiver.ReceiverIds = append(noticeReceiver.ReceiverIds, helper.IntInt64(receiverIds))
					}
				}
				if v, ok := dMap["receiver_channels"]; ok {
					receiverChannelsSet := v.(*schema.Set).List()
					for i := range receiverChannelsSet {
						receiverChannels := receiverChannelsSet[i].(string)
						noticeReceiver.ReceiverChannels = append(noticeReceiver.ReceiverChannels, &receiverChannels)
					}
				}
				if v, ok := dMap["start_time"]; ok {
					noticeReceiver.StartTime = helper.String(v.(string))
				}
				if v, ok := dMap["end_time"]; ok {
					noticeReceiver.EndTime = helper.String(v.(string))
				}
				if v, ok := dMap["index"]; ok {
					noticeReceiver.Index = helper.IntInt64(v.(int))
				}
				request.NoticeReceivers = append(request.NoticeReceivers, &noticeReceiver)
			}
		}
	}

	if d.HasChange("web_callbacks") {
		if v, ok := d.GetOk("web_callbacks"); ok {
			for _, item := range v.([]interface{}) {
				webCallback := cls.WebCallback{}
				if v, ok := dMap["url"]; ok {
					webCallback.Url = helper.String(v.(string))
				}
				if v, ok := dMap["callback_type"]; ok {
					webCallback.CallbackType = helper.String(v.(string))
				}
				if v, ok := dMap["method"]; ok {
					webCallback.Method = helper.String(v.(string))
				}
				if v, ok := dMap["headers"]; ok {
					headersSet := v.(*schema.Set).List()
					for i := range headersSet {
						headers := headersSet[i].(string)
						webCallback.Headers = append(webCallback.Headers, &headers)
					}
				}
				if v, ok := dMap["body"]; ok {
					webCallback.Body = helper.String(v.(string))
				}
				if v, ok := dMap["index"]; ok {
					webCallback.Index = helper.IntInt64(v.(int))
				}
				request.WebCallbacks = append(request.WebCallbacks, &webCallback)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().ModifyAlarmNotice(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cls alarmNotice failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("cls", "alarmNotice", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudClsAlarmNoticeRead(d, meta)
}

func resourceTencentCloudClsAlarmNoticeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_alarm_notice.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	alarmNoticeId := d.Id()

	if err := service.DeleteClsAlarmNoticeById(ctx, alarmNoticeId); err != nil {
		return err
	}

	return nil
}
