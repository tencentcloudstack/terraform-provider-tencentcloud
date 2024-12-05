package cls

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsAlarmNotice() *schema.Resource {
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
				Description: "Notice type. Value: Trigger, Recovery, All.",
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
							Description: "Receiver id list.",
						},
						"receiver_channels": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Receiver channels, Value: Email, Sms, WeChat, Phone.",
						},
						"notice_content_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Notice content ID.",
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
							Computed:    true,
							Description: "Index. The input parameter is invalid, but the output parameter is valid.",
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
						"callback_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Callback type, Values: Http, WeCom, DingTalk, Lark.",
						},
						"url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Callback url.",
						},
						"web_callback_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Integration configuration ID.",
						},
						"method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Method, POST or PUT.",
						},
						"notice_content_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Notice content ID.",
						},
						"remind_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Remind type. 0: Do not remind; 1: Specified person; 2: Everyone.",
						},
						"mobiles": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Telephone list.",
						},
						"user_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "User ID list.",
						},
						"headers": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Deprecated:  "This parameter is deprecated. Please use `notice_content_id`.",
							Description: "Request headers.",
						},
						"body": {
							Type:        schema.TypeString,
							Optional:    true,
							Deprecated:  "This parameter is deprecated. Please use `notice_content_id`.",
							Description: "Request body.",
						},
						"index": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Index. The input parameter is invalid, but the output parameter is valid.",
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
	defer tccommon.LogElapsed("resource.tencentcloud_cls_alarm_notice.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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
			if v, ok := dMap["notice_content_id"].(string); ok && v != "" {
				noticeReceiver.NoticeContentId = helper.String(v)
			}
			if v, ok := dMap["start_time"].(string); ok && v != "" {
				noticeReceiver.StartTime = helper.String(v)
			}
			if v, ok := dMap["end_time"].(string); ok && v != "" {
				noticeReceiver.EndTime = helper.String(v)
			}
			if v, ok := dMap["index"].(int); ok && v != 0 {
				noticeReceiver.Index = helper.IntInt64(v)
			}
			request.NoticeReceivers = append(request.NoticeReceivers, &noticeReceiver)
		}
	}

	if v, ok := d.GetOk("web_callbacks"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			webCallback := cls.WebCallback{}
			if v, ok := dMap["callback_type"]; ok {
				webCallback.CallbackType = helper.String(v.(string))
			}
			if v, ok := dMap["url"]; ok {
				webCallback.Url = helper.String(v.(string))
			}
			if v, ok := dMap["web_callback_id"].(string); ok && v != "" {
				webCallback.WebCallbackId = helper.String(v)
			}
			if v, ok := dMap["method"].(string); ok && v != "" {
				webCallback.Method = helper.String(v)
			}
			if v, ok := dMap["notice_content_id"].(string); ok && v != "" {
				webCallback.NoticeContentId = helper.String(v)
			}
			if v, ok := dMap["remind_type"]; ok {
				webCallback.RemindType = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["mobiles"]; ok {
				mobilesSet := v.(*schema.Set).List()
				for i := range mobilesSet {
					mobile := mobilesSet[i].(string)
					webCallback.Mobiles = append(webCallback.Mobiles, &mobile)
				}
			}
			if v, ok := dMap["user_ids"]; ok {
				userIdsSet := v.(*schema.Set).List()
				for i := range userIdsSet {
					userId := userIdsSet[i].(string)
					webCallback.UserIds = append(webCallback.UserIds, &userId)
				}
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
			if v, ok := dMap["index"].(int); ok && v != 0 {
				webCallback.Index = helper.IntInt64(v)
			}
			request.WebCallbacks = append(request.WebCallbacks, &webCallback)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateAlarmNotice(request)
		if e != nil {
			return tccommon.RetryError(e)
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

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::cls:%s:uin/:alarmNotice/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudClsAlarmNoticeRead(d, meta)
}

func resourceTencentCloudClsAlarmNoticeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_alarm_notice.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

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
		for _, noticeReceiver := range alarmNotice.NoticeReceivers {
			noticeReceiversMap := map[string]interface{}{}

			if noticeReceiver.ReceiverType != nil {
				noticeReceiversMap["receiver_type"] = noticeReceiver.ReceiverType
			}

			if noticeReceiver.ReceiverIds != nil {
				noticeReceiversMap["receiver_ids"] = noticeReceiver.ReceiverIds
			}

			if noticeReceiver.ReceiverChannels != nil {
				noticeReceiversMap["receiver_channels"] = noticeReceiver.ReceiverChannels
			}

			if noticeReceiver.NoticeContentId != nil {
				noticeReceiversMap["notice_content_id"] = noticeReceiver.NoticeContentId
			}

			if noticeReceiver.StartTime != nil {
				noticeReceiversMap["start_time"] = noticeReceiver.StartTime
			}

			if noticeReceiver.EndTime != nil {
				noticeReceiversMap["end_time"] = noticeReceiver.EndTime
			}

			if noticeReceiver.Index != nil {
				noticeReceiversMap["index"] = noticeReceiver.Index
			}

			noticeReceiversList = append(noticeReceiversList, noticeReceiversMap)
		}

		_ = d.Set("notice_receivers", noticeReceiversList)

	}

	if alarmNotice.WebCallbacks != nil {
		webCallbacksList := []interface{}{}
		for _, webCallback := range alarmNotice.WebCallbacks {
			webCallbacksMap := map[string]interface{}{}

			if webCallback.CallbackType != nil {
				webCallbacksMap["callback_type"] = webCallback.CallbackType
			}

			if webCallback.Url != nil {
				webCallbacksMap["url"] = webCallback.Url
			}

			if webCallback.WebCallbackId != nil {
				webCallbacksMap["web_callback_id"] = webCallback.WebCallbackId
			}

			if webCallback.Method != nil {
				webCallbacksMap["method"] = webCallback.Method
			}

			if webCallback.NoticeContentId != nil {
				webCallbacksMap["notice_content_id"] = webCallback.NoticeContentId
			}

			if webCallback.RemindType != nil {
				webCallbacksMap["remind_type"] = webCallback.RemindType
			}

			if webCallback.Mobiles != nil {
				tmpList := make([]string, 0, len(webCallback.Mobiles))
				for _, item := range webCallback.Mobiles {
					tmpList = append(tmpList, *item)
				}

				webCallbacksMap["mobiles"] = tmpList
			}

			if webCallback.UserIds != nil {
				tmpList := make([]string, 0, len(webCallback.UserIds))
				for _, item := range webCallback.UserIds {
					tmpList = append(tmpList, *item)
				}

				webCallbacksMap["user_ids"] = tmpList
			}

			if webCallback.Headers != nil {
				tmpList := make([]string, 0, len(webCallback.Headers))
				for _, item := range webCallback.Headers {
					tmpList = append(tmpList, *item)
				}

				webCallbacksMap["headers"] = tmpList
			}

			if webCallback.Body != nil {
				webCallbacksMap["body"] = webCallback.Body
			}

			if webCallback.Index != nil {
				webCallbacksMap["index"] = webCallback.Index
			}

			webCallbacksList = append(webCallbacksList, webCallbacksMap)
		}

		_ = d.Set("web_callbacks", webCallbacksList)

	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "cls", "alarmNotice", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudClsAlarmNoticeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_alarm_notice.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cls.NewModifyAlarmNoticeRequest()

	alarmNoticeId := d.Id()

	needChange := false
	request.AlarmNoticeId = &alarmNoticeId

	mutableArgs := []string{"name", "type", "notice_receivers", "web_callbacks"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

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
				if v, ok := dMap["notice_content_id"].(string); ok && v != "" {
					noticeReceiver.NoticeContentId = helper.String(v)
				}
				if v, ok := dMap["start_time"].(string); ok && v != "" {
					noticeReceiver.StartTime = helper.String(v)
				}
				if v, ok := dMap["end_time"].(string); ok && v != "" {
					noticeReceiver.EndTime = helper.String(v)
				}
				if v, ok := dMap["index"].(int); ok && v != 0 {
					noticeReceiver.Index = helper.IntInt64(v)
				}
				request.NoticeReceivers = append(request.NoticeReceivers, &noticeReceiver)
			}
		}

		if v, ok := d.GetOk("web_callbacks"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				webCallback := cls.WebCallback{}
				if v, ok := dMap["callback_type"]; ok {
					webCallback.CallbackType = helper.String(v.(string))
				}
				if v, ok := dMap["url"]; ok {
					webCallback.Url = helper.String(v.(string))
				}
				if v, ok := dMap["web_callback_id"].(string); ok && v != "" {
					webCallback.WebCallbackId = helper.String(v)
				}
				if v, ok := dMap["method"].(string); ok && v != "" {
					webCallback.Method = helper.String(v)
				}
				if v, ok := dMap["notice_content_id"].(string); ok && v != "" {
					webCallback.NoticeContentId = helper.String(v)
				}
				if v, ok := dMap["remind_type"]; ok {
					webCallback.RemindType = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["mobiles"]; ok {
					mobilesSet := v.(*schema.Set).List()
					for i := range mobilesSet {
						mobile := mobilesSet[i].(string)
						webCallback.Mobiles = append(webCallback.Mobiles, &mobile)
					}
				}
				if v, ok := dMap["user_ids"]; ok {
					userIdsSet := v.(*schema.Set).List()
					for i := range userIdsSet {
						userId := userIdsSet[i].(string)
						webCallback.UserIds = append(webCallback.UserIds, &userId)
					}
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
				if v, ok := dMap["index"].(int); ok && v != 0 {
					webCallback.Index = helper.IntInt64(v)
				}
				request.WebCallbacks = append(request.WebCallbacks, &webCallback)
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyAlarmNotice(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cls alarmNotice failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("cls", "alarmNotice", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudClsAlarmNoticeRead(d, meta)
}

func resourceTencentCloudClsAlarmNoticeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_alarm_notice.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	alarmNoticeId := d.Id()

	if err := service.DeleteClsAlarmNoticeById(ctx, alarmNoticeId); err != nil {
		return err
	}

	return nil
}
