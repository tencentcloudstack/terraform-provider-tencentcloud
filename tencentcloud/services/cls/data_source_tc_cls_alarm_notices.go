package cls

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClsAlarmNotices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClsAlarmNoticesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions. Maximum 10 filters, each with up to 5 values. Multiple values within the same filter use OR logic, multiple filters use AND logic.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name. Supported values: `name` (alarm notice group name), `alarmNoticeId` (alarm notice ID), `uid` (receiver user ID), `groupId` (receiver user group ID), `deliverFlag` (delivery status: 1-not enabled, 2-enabled, 3-abnormal).",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Filter field values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"has_alarm_shield_count": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to query alarm shield count statistics. Default is false.",
			},

			"alarm_notices": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of alarm notice configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm notice name.",
						},
						"alarm_notice_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm notice ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last update time.",
						},
						"notice_receivers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of notice receivers.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"receiver_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Receiver type. Can be Uin or Group.",
									},
									"receiver_ids": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Receiver IDs.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"receiver_channels": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Notification channels (Email, Sms, WeChat, Phone).",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Allowed notification start time.",
									},
									"end_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Allowed notification end time.",
									},
									"index": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Index order.",
									},
								},
							},
						},
						"web_callbacks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of webhook callbacks.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Callback URL.",
									},
									"callback_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Callback type. WeCom or Http or DingTalk or Lark or Webhook.",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HTTP method. GET or POST.",
									},
									"headers": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Request headers.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"body": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Request body.",
									},
									"index": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Index order.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"jump_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Jump domain for alarm callback.",
						},
						"notice_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of notice rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notice_receivers": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Notice receivers for this rule.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"receiver_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Receiver type.",
												},
												"receiver_ids": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "Receiver IDs.",
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"receiver_channels": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "Notification channels.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"start_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Start time.",
												},
												"end_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "End time.",
												},
												"index": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Index.",
												},
											},
										},
									},
									"web_callbacks": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Webhook callbacks for this rule.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Callback URL.",
												},
												"callback_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Callback type.",
												},
												"method": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "HTTP method.",
												},
												"headers": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "Headers.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"body": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Body.",
												},
												"index": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Index.",
												},
											},
										},
									},
									"repeat_interval": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Repeat interval in minutes.",
									},
									"time_range_start": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Effective start time (24-hour format HH:mm:ss).",
									},
									"time_range_end": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Effective end time (24-hour format HH:mm:ss).",
									},
									"notify_way": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Notification ways.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"receiver_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Receiver type.",
									},
									"day_of_week": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Days of week (0-6, 0 is Sunday).",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"jump_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Jump domain.",
									},
								},
							},
						},
						"deliver_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Delivery status (0: delivered, 1: not delivered).",
						},
						"deliver_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Delivery flag (1: not enabled, 2: enabled, 3: abnormal).",
						},
						"alarm_shield_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Alarm shield status (0: not shielded, 1: shielded).",
						},
						"alarm_shield_count": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Alarm shield count statistics.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total count of shielded alarms.",
									},
								},
							},
						},

						"callback_prioritize": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether webhook callback takes priority.",
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

func dataSourceTencentCloudClsAlarmNoticesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cls_alarm_notices.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*cls.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := cls.Filter{}
			if v, ok := filtersMap["key"].(string); ok && v != "" {
				filter.Key = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valueSet := v.(*schema.Set).List()
				for i := range valueSet {
					value := valueSet[i].(string)
					filter.Values = append(filter.Values, helper.String(value))
				}
			}
			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOkExists("has_alarm_shield_count"); ok {
		paramMap["HasAlarmShieldCount"] = helper.Bool(v.(bool))
	}

	var respData []*cls.AlarmNotice

	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClsAlarmNoticesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	alarmNoticesList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, alarmNotice := range respData {
			alarmNoticeMap := map[string]interface{}{}

			if alarmNotice.Name != nil {
				alarmNoticeMap["name"] = alarmNotice.Name
			}

			if alarmNotice.AlarmNoticeId != nil {
				alarmNoticeMap["alarm_notice_id"] = alarmNotice.AlarmNoticeId
			}

			if alarmNotice.CreateTime != nil {
				alarmNoticeMap["create_time"] = alarmNotice.CreateTime
			}

			if alarmNotice.UpdateTime != nil {
				alarmNoticeMap["update_time"] = alarmNotice.UpdateTime
			}

			if alarmNotice.NoticeReceivers != nil {
				noticeReceiversList := []interface{}{}
				for _, noticeReceiver := range alarmNotice.NoticeReceivers {
					noticeReceiverMap := map[string]interface{}{}

					if noticeReceiver.ReceiverType != nil {
						noticeReceiverMap["receiver_type"] = noticeReceiver.ReceiverType
					}

					if noticeReceiver.ReceiverIds != nil {
						noticeReceiverMap["receiver_ids"] = noticeReceiver.ReceiverIds
					}

					if noticeReceiver.ReceiverChannels != nil {
						noticeReceiverMap["receiver_channels"] = noticeReceiver.ReceiverChannels
					}

					if noticeReceiver.StartTime != nil {
						noticeReceiverMap["start_time"] = noticeReceiver.StartTime
					}

					if noticeReceiver.EndTime != nil {
						noticeReceiverMap["end_time"] = noticeReceiver.EndTime
					}

					if noticeReceiver.Index != nil {
						noticeReceiverMap["index"] = noticeReceiver.Index
					}

					noticeReceiversList = append(noticeReceiversList, noticeReceiverMap)
				}

				alarmNoticeMap["notice_receivers"] = noticeReceiversList
			}

			if alarmNotice.WebCallbacks != nil {
				webCallbacksList := []interface{}{}
				for _, webCallback := range alarmNotice.WebCallbacks {
					webCallbackMap := map[string]interface{}{}

					if webCallback.Url != nil {
						webCallbackMap["url"] = webCallback.Url
					}

					if webCallback.CallbackType != nil {
						webCallbackMap["callback_type"] = webCallback.CallbackType
					}

					if webCallback.Method != nil {
						webCallbackMap["method"] = webCallback.Method
					}

					if webCallback.Headers != nil {
						webCallbackMap["headers"] = webCallback.Headers
					}

					if webCallback.Body != nil {
						webCallbackMap["body"] = webCallback.Body
					}

					if webCallback.Index != nil {
						webCallbackMap["index"] = webCallback.Index
					}

					webCallbacksList = append(webCallbacksList, webCallbackMap)
				}

				alarmNoticeMap["web_callbacks"] = webCallbacksList
			}

			if alarmNotice.Tags != nil {
				tagsList := []interface{}{}
				for _, tag := range alarmNotice.Tags {
					tagMap := map[string]interface{}{}

					if tag.Key != nil {
						tagMap["key"] = tag.Key
					}

					if tag.Value != nil {
						tagMap["value"] = tag.Value
					}

					tagsList = append(tagsList, tagMap)
				}

				alarmNoticeMap["tags"] = tagsList
			}

			if alarmNotice.JumpDomain != nil {
				alarmNoticeMap["jump_domain"] = alarmNotice.JumpDomain
			}

			if alarmNotice.NoticeRules != nil {
				noticeRulesList := []interface{}{}
				for _, noticeRule := range alarmNotice.NoticeRules {
					noticeRuleMap := map[string]interface{}{}

					if noticeRule.NoticeReceivers != nil {
						ruleNoticeReceiversList := []interface{}{}
						for _, receiver := range noticeRule.NoticeReceivers {
							receiverMap := map[string]interface{}{}

							if receiver.ReceiverType != nil {
								receiverMap["receiver_type"] = receiver.ReceiverType
							}

							if receiver.ReceiverIds != nil {
								receiverMap["receiver_ids"] = receiver.ReceiverIds
							}

							if receiver.ReceiverChannels != nil {
								receiverMap["receiver_channels"] = receiver.ReceiverChannels
							}

							if receiver.StartTime != nil {
								receiverMap["start_time"] = receiver.StartTime
							}

							if receiver.EndTime != nil {
								receiverMap["end_time"] = receiver.EndTime
							}

							if receiver.Index != nil {
								receiverMap["index"] = receiver.Index
							}

							ruleNoticeReceiversList = append(ruleNoticeReceiversList, receiverMap)
						}

						noticeRuleMap["notice_receivers"] = ruleNoticeReceiversList
					}

					if noticeRule.WebCallbacks != nil {
						ruleWebCallbacksList := []interface{}{}
						for _, callback := range noticeRule.WebCallbacks {
							callbackMap := map[string]interface{}{}

							if callback.Url != nil {
								callbackMap["url"] = callback.Url
							}

							if callback.CallbackType != nil {
								callbackMap["callback_type"] = callback.CallbackType
							}

							if callback.Method != nil {
								callbackMap["method"] = callback.Method
							}

							if callback.Headers != nil {
								callbackMap["headers"] = callback.Headers
							}

							if callback.Body != nil {
								callbackMap["body"] = callback.Body
							}

							if callback.Index != nil {
								callbackMap["index"] = callback.Index
							}

							ruleWebCallbacksList = append(ruleWebCallbacksList, callbackMap)
						}

						noticeRuleMap["web_callbacks"] = ruleWebCallbacksList
					}

					noticeRulesList = append(noticeRulesList, noticeRuleMap)
				}

				alarmNoticeMap["notice_rules"] = noticeRulesList
			}

			if alarmNotice.DeliverStatus != nil {
				alarmNoticeMap["deliver_status"] = alarmNotice.DeliverStatus
			}

			if alarmNotice.DeliverFlag != nil {
				alarmNoticeMap["deliver_flag"] = alarmNotice.DeliverFlag
			}

			if alarmNotice.AlarmShieldStatus != nil {
				alarmNoticeMap["alarm_shield_status"] = alarmNotice.AlarmShieldStatus
			}

			if alarmNotice.AlarmShieldCount != nil {
				alarmShieldCountList := []interface{}{}
				alarmShieldCountMap := map[string]interface{}{}

				if alarmNotice.AlarmShieldCount.TotalCount != nil {
					alarmShieldCountMap["total_count"] = alarmNotice.AlarmShieldCount.TotalCount
				}

				alarmShieldCountList = append(alarmShieldCountList, alarmShieldCountMap)
				alarmNoticeMap["alarm_shield_count"] = alarmShieldCountList
			}

			if alarmNotice.CallbackPrioritize != nil {
				alarmNoticeMap["callback_prioritize"] = alarmNotice.CallbackPrioritize
			}

			alarmNoticesList = append(alarmNoticesList, alarmNoticeMap)
		}

		_ = d.Set("alarm_notices", alarmNoticesList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
