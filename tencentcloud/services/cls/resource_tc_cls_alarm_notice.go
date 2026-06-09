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
				Optional:    true,
				Computed:    true,
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

			"jump_domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Jump domain. Must start with http:// or https://, cannot end with /.",
			},

			"deliver_status": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Deliver log switch. Valid values: 1 (off, default), 2 (on). When set to 2, deliver_config is required.",
			},

			"alarm_shield_status": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Alarm shield status (no-login operation). Valid values: 1 (off), 2 (on, default).",
			},

			"callback_prioritize": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Callback prioritize. true: use custom callback params from notice content template; false: use params from alarm policy.",
			},

			"deliver_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Deliver log configuration. Required when deliver_status is 2.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region of the target log topic. e.g. ap-guangzhou.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Target log topic ID.",
						},
						"scope": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Deliver data scope. 0: all logs (default); 1: only alarm trigger and recovery logs.",
						},
					},
				},
			},

			"notice_rules": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Notice rules (advanced mode). Mutually exclusive with type/notice_receivers/web_callbacks (simple mode).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Matching rule JSON string.",
						},
						"notice_receivers": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Notice receivers for this rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"receiver_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Receiver type, Uin or Group.",
									},
									"receiver_ids": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Description: "Receiver id list.",
									},
									"receiver_channels": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
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
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Web callbacks for this rule.",
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
										Type:        schema.TypeSet,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
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
						"escalate": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Alarm escalate switch. true: enable; false: disable (default).",
						},
						"type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Alarm escalate condition. 1: unclaimed and unresolved (default); 2: unresolved.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Alarm escalate interval in minutes. Range: [1, 14400].",
						},
						"escalate_notices": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    5,
							Description: "Alarm escalate notice chain, ordered from level 1 to level 5 (max). Each element represents the next escalation level.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notice_receivers": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Notice receivers for this escalation level.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"receiver_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Receiver type, Uin or Group.",
												},
												"receiver_ids": {
													Type:     schema.TypeSet,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Description: "Receiver id list.",
												},
												"receiver_channels": {
													Type:     schema.TypeSet,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
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
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Web callbacks for this escalation level.",
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
													Type:        schema.TypeSet,
													Optional:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
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
									"escalate": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to continue escalating from this level. true: enable; false: disable.",
									},
									"type": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Escalate condition. 1: unclaimed and unresolved (default); 2: unresolved.",
									},
									"interval": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Escalate interval in minutes. Range: [1, 14400].",
									},
								},
							},
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

	if v, ok := d.GetOk("jump_domain"); ok {
		request.JumpDomain = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("deliver_status"); ok {
		request.DeliverStatus = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("alarm_shield_status"); ok {
		request.AlarmShieldStatus = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("callback_prioritize"); ok {
		request.CallbackPrioritize = helper.Bool(v.(bool))
	}

	if v, ok := helper.InterfacesHeadMap(d, "deliver_config"); ok {
		deliverConfig := cls.DeliverConfig{}
		if region, ok := v["region"].(string); ok && region != "" {
			deliverConfig.Region = helper.String(region)
		}
		if topicId, ok := v["topic_id"].(string); ok && topicId != "" {
			deliverConfig.TopicId = helper.String(topicId)
		}
		if scope, ok := v["scope"].(int); ok {
			deliverConfig.Scope = helper.IntUint64(scope)
		}
		request.DeliverConfig = &deliverConfig
	}

	if v, ok := d.GetOk("notice_rules"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			noticeRule := cls.NoticeRule{}
			if rule, ok := dMap["rule"].(string); ok && rule != "" {
				noticeRule.Rule = helper.String(rule)
			}
			if escalate, ok := dMap["escalate"].(bool); ok {
				noticeRule.Escalate = helper.Bool(escalate)
			}
			if t, ok := dMap["type"].(int); ok && t != 0 {
				noticeRule.Type = helper.IntUint64(t)
			}
			if interval, ok := dMap["interval"].(int); ok && interval != 0 {
				noticeRule.Interval = helper.IntUint64(interval)
			}
			if nrList, ok := dMap["notice_receivers"].([]interface{}); ok {
				for _, nrItem := range nrList {
					nrMap := nrItem.(map[string]interface{})
					nr := cls.NoticeReceiver{}
					if v, ok := nrMap["receiver_type"].(string); ok {
						nr.ReceiverType = helper.String(v)
					}
					if ids, ok := nrMap["receiver_ids"]; ok {
						for _, id := range ids.(*schema.Set).List() {
							nr.ReceiverIds = append(nr.ReceiverIds, helper.IntInt64(id.(int)))
						}
					}
					if channels, ok := nrMap["receiver_channels"]; ok {
						for _, ch := range channels.(*schema.Set).List() {
							c := ch.(string)
							nr.ReceiverChannels = append(nr.ReceiverChannels, &c)
						}
					}
					if v, ok := nrMap["notice_content_id"].(string); ok && v != "" {
						nr.NoticeContentId = helper.String(v)
					}
					if v, ok := nrMap["start_time"].(string); ok && v != "" {
						nr.StartTime = helper.String(v)
					}
					if v, ok := nrMap["end_time"].(string); ok && v != "" {
						nr.EndTime = helper.String(v)
					}
					noticeRule.NoticeReceivers = append(noticeRule.NoticeReceivers, &nr)
				}
			}
			if wcList, ok := dMap["web_callbacks"].([]interface{}); ok {
				for _, wcItem := range wcList {
					wcMap := wcItem.(map[string]interface{})
					wc := cls.WebCallback{}
					if v, ok := wcMap["callback_type"].(string); ok {
						wc.CallbackType = helper.String(v)
					}
					if v, ok := wcMap["url"].(string); ok {
						wc.Url = helper.String(v)
					}
					if v, ok := wcMap["web_callback_id"].(string); ok && v != "" {
						wc.WebCallbackId = helper.String(v)
					}
					if v, ok := wcMap["method"].(string); ok && v != "" {
						wc.Method = helper.String(v)
					}
					if v, ok := wcMap["notice_content_id"].(string); ok && v != "" {
						wc.NoticeContentId = helper.String(v)
					}
					if v, ok := wcMap["remind_type"].(int); ok {
						wc.RemindType = helper.IntUint64(v)
					}
					if mobiles, ok := wcMap["mobiles"]; ok {
						for _, m := range mobiles.(*schema.Set).List() {
							mobile := m.(string)
							wc.Mobiles = append(wc.Mobiles, &mobile)
						}
					}
					if userIds, ok := wcMap["user_ids"]; ok {
						for _, u := range userIds.(*schema.Set).List() {
							uid := u.(string)
							wc.UserIds = append(wc.UserIds, &uid)
						}
					}
					if headers, ok := wcMap["headers"]; ok {
						for _, h := range headers.(*schema.Set).List() {
							header := h.(string)
							wc.Headers = append(wc.Headers, &header)
						}
					}
					if v, ok := wcMap["body"].(string); ok && v != "" {
						wc.Body = helper.String(v)
					}
					noticeRule.WebCallbacks = append(noticeRule.WebCallbacks, &wc)
				}
			}
			if escalateList, ok := dMap["escalate_notices"].([]interface{}); ok && len(escalateList) > 0 {
				noticeRule.EscalateNotice = buildEscalateNoticeChain(escalateList)
			}
			request.NoticeRules = append(request.NoticeRules, &noticeRule)
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

	if alarmNotice.JumpDomain != nil {
		_ = d.Set("jump_domain", alarmNotice.JumpDomain)
	}

	if alarmNotice.DeliverStatus != nil {
		_ = d.Set("deliver_status", alarmNotice.DeliverStatus)
	}

	if alarmNotice.AlarmShieldStatus != nil {
		_ = d.Set("alarm_shield_status", alarmNotice.AlarmShieldStatus)
	}

	if alarmNotice.CallbackPrioritize != nil {
		_ = d.Set("callback_prioritize", alarmNotice.CallbackPrioritize)
	}

	if alarmNotice.AlarmNoticeDeliverConfig != nil && alarmNotice.AlarmNoticeDeliverConfig.DeliverConfig != nil {
		dc := alarmNotice.AlarmNoticeDeliverConfig.DeliverConfig
		deliverConfigMap := map[string]interface{}{}
		if dc.Region != nil {
			deliverConfigMap["region"] = dc.Region
		}
		if dc.TopicId != nil {
			deliverConfigMap["topic_id"] = dc.TopicId
		}
		if dc.Scope != nil {
			deliverConfigMap["scope"] = dc.Scope
		}
		_ = d.Set("deliver_config", []interface{}{deliverConfigMap})
	}

	if alarmNotice.NoticeRules != nil {
		noticeRulesList := []interface{}{}
		for _, nr := range alarmNotice.NoticeRules {
			nrMap := map[string]interface{}{}
			if nr.Rule != nil {
				nrMap["rule"] = nr.Rule
			}
			if nr.Escalate != nil {
				nrMap["escalate"] = nr.Escalate
			}
			if nr.Type != nil {
				nrMap["type"] = nr.Type
			}
			if nr.Interval != nil {
				nrMap["interval"] = nr.Interval
			}
			if nr.NoticeReceivers != nil {
				nrReceiverList := []interface{}{}
				for _, rec := range nr.NoticeReceivers {
					recMap := map[string]interface{}{}
					if rec.ReceiverType != nil {
						recMap["receiver_type"] = rec.ReceiverType
					}
					if rec.ReceiverIds != nil {
						recMap["receiver_ids"] = rec.ReceiverIds
					}
					if rec.ReceiverChannels != nil {
						recMap["receiver_channels"] = rec.ReceiverChannels
					}
					if rec.NoticeContentId != nil {
						recMap["notice_content_id"] = rec.NoticeContentId
					}
					if rec.StartTime != nil {
						recMap["start_time"] = rec.StartTime
					}
					if rec.EndTime != nil {
						recMap["end_time"] = rec.EndTime
					}
					if rec.Index != nil {
						recMap["index"] = rec.Index
					}
					nrReceiverList = append(nrReceiverList, recMap)
				}
				nrMap["notice_receivers"] = nrReceiverList
			}
			if nr.WebCallbacks != nil {
				nrWcList := []interface{}{}
				for _, wc := range nr.WebCallbacks {
					wcMap := map[string]interface{}{}
					if wc.CallbackType != nil {
						wcMap["callback_type"] = wc.CallbackType
					}
					if wc.Url != nil {
						wcMap["url"] = wc.Url
					}
					if wc.WebCallbackId != nil {
						wcMap["web_callback_id"] = wc.WebCallbackId
					}
					if wc.Method != nil {
						wcMap["method"] = wc.Method
					}
					if wc.NoticeContentId != nil {
						wcMap["notice_content_id"] = wc.NoticeContentId
					}
					if wc.RemindType != nil {
						wcMap["remind_type"] = wc.RemindType
					}
					if wc.Mobiles != nil {
						tmpList := make([]string, 0, len(wc.Mobiles))
						for _, m := range wc.Mobiles {
							tmpList = append(tmpList, *m)
						}
						wcMap["mobiles"] = tmpList
					}
					if wc.UserIds != nil {
						tmpList := make([]string, 0, len(wc.UserIds))
						for _, u := range wc.UserIds {
							tmpList = append(tmpList, *u)
						}
						wcMap["user_ids"] = tmpList
					}
					if wc.Headers != nil {
						tmpList := make([]string, 0, len(wc.Headers))
						for _, h := range wc.Headers {
							tmpList = append(tmpList, *h)
						}
						wcMap["headers"] = tmpList
					}
					if wc.Body != nil {
						wcMap["body"] = wc.Body
					}
					if wc.Index != nil {
						wcMap["index"] = wc.Index
					}
					nrWcList = append(nrWcList, wcMap)
				}
				nrMap["web_callbacks"] = nrWcList
			}
			if nr.EscalateNotice != nil {
				nrMap["escalate_notices"] = flattenEscalateNoticeChain(nr.EscalateNotice)
			}
			noticeRulesList = append(noticeRulesList, nrMap)
		}
		_ = d.Set("notice_rules", noticeRulesList)
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

	mutableArgs := []string{"name", "type", "notice_receivers", "web_callbacks", "jump_domain", "deliver_status", "deliver_config", "alarm_shield_status", "callback_prioritize", "notice_rules"}

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

		if v, ok := d.GetOk("jump_domain"); ok {
			request.JumpDomain = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("deliver_status"); ok {
			request.DeliverStatus = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("alarm_shield_status"); ok {
			request.AlarmShieldStatus = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("callback_prioritize"); ok {
			request.CallbackPrioritize = helper.Bool(v.(bool))
		}

		if v, ok := helper.InterfacesHeadMap(d, "deliver_config"); ok {
			deliverConfig := cls.DeliverConfig{}
			if region, ok := v["region"].(string); ok && region != "" {
				deliverConfig.Region = helper.String(region)
			}
			if topicId, ok := v["topic_id"].(string); ok && topicId != "" {
				deliverConfig.TopicId = helper.String(topicId)
			}
			if scope, ok := v["scope"].(int); ok {
				deliverConfig.Scope = helper.IntUint64(scope)
			}
			request.DeliverConfig = &deliverConfig
		}

		if v, ok := d.GetOk("notice_rules"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				noticeRule := cls.NoticeRule{}
				if rule, ok := dMap["rule"].(string); ok && rule != "" {
					noticeRule.Rule = helper.String(rule)
				}
				if escalate, ok := dMap["escalate"].(bool); ok {
					noticeRule.Escalate = helper.Bool(escalate)
				}
				if t, ok := dMap["type"].(int); ok && t != 0 {
					noticeRule.Type = helper.IntUint64(t)
				}
				if interval, ok := dMap["interval"].(int); ok && interval != 0 {
					noticeRule.Interval = helper.IntUint64(interval)
				}
				if nrList, ok := dMap["notice_receivers"].([]interface{}); ok {
					for _, nrItem := range nrList {
						nrMap := nrItem.(map[string]interface{})
						nr := cls.NoticeReceiver{}
						if v, ok := nrMap["receiver_type"].(string); ok {
							nr.ReceiverType = helper.String(v)
						}
						if ids, ok := nrMap["receiver_ids"]; ok {
							for _, id := range ids.(*schema.Set).List() {
								nr.ReceiverIds = append(nr.ReceiverIds, helper.IntInt64(id.(int)))
							}
						}
						if channels, ok := nrMap["receiver_channels"]; ok {
							for _, ch := range channels.(*schema.Set).List() {
								c := ch.(string)
								nr.ReceiverChannels = append(nr.ReceiverChannels, &c)
							}
						}
						if v, ok := nrMap["notice_content_id"].(string); ok && v != "" {
							nr.NoticeContentId = helper.String(v)
						}
						if v, ok := nrMap["start_time"].(string); ok && v != "" {
							nr.StartTime = helper.String(v)
						}
						if v, ok := nrMap["end_time"].(string); ok && v != "" {
							nr.EndTime = helper.String(v)
						}
						noticeRule.NoticeReceivers = append(noticeRule.NoticeReceivers, &nr)
					}
				}
				if wcList, ok := dMap["web_callbacks"].([]interface{}); ok {
					for _, wcItem := range wcList {
						wcMap := wcItem.(map[string]interface{})
						wc := cls.WebCallback{}
						if v, ok := wcMap["callback_type"].(string); ok {
							wc.CallbackType = helper.String(v)
						}
						if v, ok := wcMap["url"].(string); ok {
							wc.Url = helper.String(v)
						}
						if v, ok := wcMap["web_callback_id"].(string); ok && v != "" {
							wc.WebCallbackId = helper.String(v)
						}
						if v, ok := wcMap["method"].(string); ok && v != "" {
							wc.Method = helper.String(v)
						}
						if v, ok := wcMap["notice_content_id"].(string); ok && v != "" {
							wc.NoticeContentId = helper.String(v)
						}
						if v, ok := wcMap["remind_type"].(int); ok {
							wc.RemindType = helper.IntUint64(v)
						}
						if mobiles, ok := wcMap["mobiles"]; ok {
							for _, m := range mobiles.(*schema.Set).List() {
								mobile := m.(string)
								wc.Mobiles = append(wc.Mobiles, &mobile)
							}
						}
						if userIds, ok := wcMap["user_ids"]; ok {
							for _, u := range userIds.(*schema.Set).List() {
								uid := u.(string)
								wc.UserIds = append(wc.UserIds, &uid)
							}
						}
						if headers, ok := wcMap["headers"]; ok {
							for _, h := range headers.(*schema.Set).List() {
								header := h.(string)
								wc.Headers = append(wc.Headers, &header)
							}
						}
						if v, ok := wcMap["body"].(string); ok && v != "" {
							wc.Body = helper.String(v)
						}
						noticeRule.WebCallbacks = append(noticeRule.WebCallbacks, &wc)
					}
				}
				if escalateList, ok := dMap["escalate_notices"].([]interface{}); ok && len(escalateList) > 0 {
					noticeRule.EscalateNotice = buildEscalateNoticeChain(escalateList)
				}
				request.NoticeRules = append(request.NoticeRules, &noticeRule)
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

// buildEscalateNoticeChain converts an ordered list (plan D) into the SDK's recursive chain structure.
// escalateList[0] = level 1, escalateList[1] = level 2, ...
func buildEscalateNoticeChain(escalateList []interface{}) *cls.EscalateNoticeInfo {
	if len(escalateList) == 0 {
		return nil
	}
	m := escalateList[0].(map[string]interface{})
	info := &cls.EscalateNoticeInfo{}
	if v, ok := m["escalate"].(bool); ok {
		info.Escalate = helper.Bool(v)
	}
	if v, ok := m["type"].(int); ok && v != 0 {
		info.Type = helper.IntUint64(v)
	}
	if v, ok := m["interval"].(int); ok && v != 0 {
		info.Interval = helper.IntUint64(v)
	}
	if nrList, ok := m["notice_receivers"].([]interface{}); ok {
		for _, nrItem := range nrList {
			nrMap := nrItem.(map[string]interface{})
			nr := cls.NoticeReceiver{}
			if v, ok := nrMap["receiver_type"].(string); ok {
				nr.ReceiverType = helper.String(v)
			}
			if ids, ok := nrMap["receiver_ids"]; ok {
				for _, id := range ids.(*schema.Set).List() {
					nr.ReceiverIds = append(nr.ReceiverIds, helper.IntInt64(id.(int)))
				}
			}
			if channels, ok := nrMap["receiver_channels"]; ok {
				for _, ch := range channels.(*schema.Set).List() {
					c := ch.(string)
					nr.ReceiverChannels = append(nr.ReceiverChannels, &c)
				}
			}
			if v, ok := nrMap["notice_content_id"].(string); ok && v != "" {
				nr.NoticeContentId = helper.String(v)
			}
			if v, ok := nrMap["start_time"].(string); ok && v != "" {
				nr.StartTime = helper.String(v)
			}
			if v, ok := nrMap["end_time"].(string); ok && v != "" {
				nr.EndTime = helper.String(v)
			}
			info.NoticeReceivers = append(info.NoticeReceivers, &nr)
		}
	}
	if wcList, ok := m["web_callbacks"].([]interface{}); ok {
		for _, wcItem := range wcList {
			wcMap := wcItem.(map[string]interface{})
			wc := cls.WebCallback{}
			if v, ok := wcMap["callback_type"].(string); ok {
				wc.CallbackType = helper.String(v)
			}
			if v, ok := wcMap["url"].(string); ok {
				wc.Url = helper.String(v)
			}
			if v, ok := wcMap["web_callback_id"].(string); ok && v != "" {
				wc.WebCallbackId = helper.String(v)
			}
			if v, ok := wcMap["method"].(string); ok && v != "" {
				wc.Method = helper.String(v)
			}
			if v, ok := wcMap["notice_content_id"].(string); ok && v != "" {
				wc.NoticeContentId = helper.String(v)
			}
			if v, ok := wcMap["remind_type"].(int); ok {
				wc.RemindType = helper.IntUint64(v)
			}
			if mobiles, ok := wcMap["mobiles"]; ok {
				for _, mb := range mobiles.(*schema.Set).List() {
					mobile := mb.(string)
					wc.Mobiles = append(wc.Mobiles, &mobile)
				}
			}
			if userIds, ok := wcMap["user_ids"]; ok {
				for _, u := range userIds.(*schema.Set).List() {
					uid := u.(string)
					wc.UserIds = append(wc.UserIds, &uid)
				}
			}
			if headers, ok := wcMap["headers"]; ok {
				for _, h := range headers.(*schema.Set).List() {
					header := h.(string)
					wc.Headers = append(wc.Headers, &header)
				}
			}
			if v, ok := wcMap["body"].(string); ok && v != "" {
				wc.Body = helper.String(v)
			}
			info.WebCallbacks = append(info.WebCallbacks, &wc)
		}
	}
	// Recursively process the next level.
	info.EscalateNotice = buildEscalateNoticeChain(escalateList[1:])
	return info
}

// flattenEscalateNoticeChain flattens the SDK's recursive chain structure into an ordered list (plan D).
func flattenEscalateNoticeChain(info *cls.EscalateNoticeInfo) []interface{} {
	if info == nil {
		return nil
	}
	result := []interface{}{}
	cur := info
	for cur != nil {
		m := map[string]interface{}{}
		if cur.Escalate != nil {
			m["escalate"] = cur.Escalate
		}
		if cur.Type != nil {
			m["type"] = cur.Type
		}
		if cur.Interval != nil {
			m["interval"] = cur.Interval
		}
		if cur.NoticeReceivers != nil {
			nrList := []interface{}{}
			for _, rec := range cur.NoticeReceivers {
				recMap := map[string]interface{}{}
				if rec.ReceiverType != nil {
					recMap["receiver_type"] = rec.ReceiverType
				}
				if rec.ReceiverIds != nil {
					recMap["receiver_ids"] = rec.ReceiverIds
				}
				if rec.ReceiverChannels != nil {
					recMap["receiver_channels"] = rec.ReceiverChannels
				}
				if rec.NoticeContentId != nil {
					recMap["notice_content_id"] = rec.NoticeContentId
				}
				if rec.StartTime != nil {
					recMap["start_time"] = rec.StartTime
				}
				if rec.EndTime != nil {
					recMap["end_time"] = rec.EndTime
				}
				if rec.Index != nil {
					recMap["index"] = rec.Index
				}
				nrList = append(nrList, recMap)
			}
			m["notice_receivers"] = nrList
		}
		if cur.WebCallbacks != nil {
			wcList := []interface{}{}
			for _, wc := range cur.WebCallbacks {
				wcMap := map[string]interface{}{}
				if wc.CallbackType != nil {
					wcMap["callback_type"] = wc.CallbackType
				}
				if wc.Url != nil {
					wcMap["url"] = wc.Url
				}
				if wc.WebCallbackId != nil {
					wcMap["web_callback_id"] = wc.WebCallbackId
				}
				if wc.Method != nil {
					wcMap["method"] = wc.Method
				}
				if wc.NoticeContentId != nil {
					wcMap["notice_content_id"] = wc.NoticeContentId
				}
				if wc.RemindType != nil {
					wcMap["remind_type"] = wc.RemindType
				}
				if wc.Mobiles != nil {
					tmpList := make([]string, 0, len(wc.Mobiles))
					for _, mb := range wc.Mobiles {
						tmpList = append(tmpList, *mb)
					}
					wcMap["mobiles"] = tmpList
				}
				if wc.UserIds != nil {
					tmpList := make([]string, 0, len(wc.UserIds))
					for _, u := range wc.UserIds {
						tmpList = append(tmpList, *u)
					}
					wcMap["user_ids"] = tmpList
				}
				if wc.Headers != nil {
					tmpList := make([]string, 0, len(wc.Headers))
					for _, h := range wc.Headers {
						tmpList = append(tmpList, *h)
					}
					wcMap["headers"] = tmpList
				}
				if wc.Body != nil {
					wcMap["body"] = wc.Body
				}
				if wc.Index != nil {
					wcMap["index"] = wc.Index
				}
				wcList = append(wcList, wcMap)
			}
			m["web_callbacks"] = wcList
		}
		result = append(result, m)
		cur = cur.EscalateNotice
	}
	return result
}
