package monitor

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitorv20230616 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20230616"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorNoticeContentTmpl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorNoticeContentTmplCreate,
		Read:   resourceTencentCloudMonitorNoticeContentTmplRead,
		Update: resourceTencentCloudMonitorNoticeContentTmplUpdate,
		Delete: resourceTencentCloudMonitorNoticeContentTmplDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tmpl_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Template name.",
			},

			"monitor_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Monitor type, e.g. MT_QCE.",
			},

			"tmpl_language": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Template language, zh for Chinese, en for English.",
			},

			"tmpl_contents": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Template content configuration for different notification channels.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"qcloud_yehe": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "QCloud Yehe notification channel configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"matching_status": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Matching status list, e.g. Trigger, Recovery.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"template": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Template configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"email": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Email notification.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"content_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Content template.",
															},
															"title_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Title template.",
															},
														},
													},
												},
												"qywx": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Enterprise WeChat notification.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"content_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Content template.",
															},
															"title_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Title template.",
															},
														},
													},
												},
												"sms": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "SMS notification.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"content_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Content template.",
															},
															"title_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Title template.",
															},
														},
													},
												},
												"voice": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Voice notification.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"content_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Content template.",
															},
															"title_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Title template.",
															},
														},
													},
												},
												"wechat": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "WeChat notification.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"alarm_content_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Alarm content template.",
															},
															"alarm_object_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Alarm object template.",
															},
															"alarm_region_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Alarm region template.",
															},
															"alarm_time_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Alarm time template.",
															},
														},
													},
												},
												"site": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Site notification.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"content_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Content template.",
															},
															"title_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Title template.",
															},
														},
													},
												},
												"andon": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Andon notification.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"content_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Content template.",
															},
															"title_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Title template.",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"we_work_robot": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "WeWork Robot notification channel configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"matching_status": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Matching status list, e.g. Trigger, Recovery.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"template": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Template configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"content_tmpl": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Content template.",
												},
											},
										},
									},
								},
							},
						},
						"ding_ding_robot": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "DingDing Robot notification channel configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"matching_status": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Matching status list, e.g. Trigger, Recovery.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"template": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Template configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"content_tmpl": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Content template.",
												},
												"title_tmpl": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Title template.",
												},
											},
										},
									},
								},
							},
						},
						"fei_shu_robot": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "FeiShu Robot notification channel configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"matching_status": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Matching status list, e.g. Trigger, Recovery.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"template": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Template configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"content_tmpl": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Content template.",
												},
												"title_tmpl": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Title template.",
												},
											},
										},
									},
								},
							},
						},
						"webhook": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Webhook notification channel configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"matching_status": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Matching status list.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"template": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Webhook template.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"body": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Request body.",
												},
												"body_content_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Body content type.",
												},
												"headers": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Request headers.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Header key.",
															},
															"values": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Header values.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"teams_robot": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Teams Robot notification channel configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"matching_status": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Matching status list, e.g. Trigger, Recovery.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"template": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Template configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"content_tmpl": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Content template.",
												},
											},
										},
									},
								},
							},
						},
						"pager_duty_robot": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "PagerDuty Robot notification channel configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"matching_status": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Matching status list.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"template": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "PagerDuty template.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"body": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Request body template in JSON.",
												},
												"title_tmpl": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Title template.",
												},
												"headers": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Request headers.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Header key.",
															},
															"values": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Header values.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			// computed
			"tmpl_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template ID.",
			},
		},
	}
}

func resourceTencentCloudMonitorNoticeContentTmplCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_notice_content_tmpl.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = monitorv20230616.NewCreateNoticeContentTmplRequest()
		response = monitorv20230616.NewCreateNoticeContentTmplResponse()
		tmplID   string
	)

	if v, ok := d.GetOk("tmpl_name"); ok {
		request.TmplName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("monitor_type"); ok {
		request.MonitorType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tmpl_language"); ok {
		request.TmplLanguage = helper.String(v.(string))
	}

	if _, ok := d.GetOk("tmpl_contents"); ok {
		request.TmplContents = expandNoticeContentTmplItem(d)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorV20230616Client().CreateNoticeContentTmplWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create monitor notice_content_tmpl failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create monitor notice_content_tmpl failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.TmplID == nil {
		return fmt.Errorf("TmplID is nil.")
	}

	tmplID = *response.Response.TmplID
	d.SetId(tmplID)
	return resourceTencentCloudMonitorNoticeContentTmplRead(d, meta)
}

func resourceTencentCloudMonitorNoticeContentTmplRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_notice_content_tmpl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	tmplID := d.Id()

	respData, err := service.DescribeNoticeContentTmplById(ctx, tmplID)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_monitor_notice_content_tmpl` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.TmplID != nil {
		_ = d.Set("tmpl_id", respData.TmplID)
	}

	if respData.TmplName != nil {
		_ = d.Set("tmpl_name", respData.TmplName)
	}

	if respData.MonitorType != nil {
		_ = d.Set("monitor_type", respData.MonitorType)
	}

	if respData.TmplLanguage != nil {
		_ = d.Set("tmpl_language", respData.TmplLanguage)
	}

	if respData.TmplContents != nil {
		tmplContentsList := flattenNoticeContentTmplItem(respData.TmplContents)
		_ = d.Set("tmpl_contents", tmplContentsList)
	}

	return nil
}

func resourceTencentCloudMonitorNoticeContentTmplUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_notice_content_tmpl.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	tmplID := d.Id()

	needChange := false
	mutableArgs := []string{"tmpl_name", "tmpl_contents"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := monitorv20230616.NewModifyNoticeContentTmplRequest()

		request.TmplID = &tmplID
		if v, ok := d.GetOk("tmpl_name"); ok {
			request.TmplName = helper.String(v.(string))
		}
		if _, ok := d.GetOk("tmpl_contents"); ok {
			request.TmplContents = expandNoticeContentTmplItem(d)
		}
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorV20230616Client().ModifyNoticeContentTmplWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update monitor notice_content_tmpl failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudMonitorNoticeContentTmplRead(d, meta)
}

func resourceTencentCloudMonitorNoticeContentTmplDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_notice_content_tmpl.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = monitorv20230616.NewDeleteNoticeContentTmplsRequest()
	)

	tmplID := d.Id()

	request.TmplIDs = []*string{&tmplID}
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorV20230616Client().DeleteNoticeContentTmplsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete monitor notice_content_tmpl failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func expandNoticeContentTmplItem(d *schema.ResourceData) *monitorv20230616.NoticeContentTmplItem {
	tmplContentsRaw := d.Get("tmpl_contents").([]interface{})
	if len(tmplContentsRaw) == 0 {
		return nil
	}

	tmplContentsMap := tmplContentsRaw[0].(map[string]interface{})
	tmplContents := &monitorv20230616.NoticeContentTmplItem{}

	// QCloudYehe
	if v, ok := tmplContentsMap["qcloud_yehe"]; ok && len(v.([]interface{})) > 0 {
		qCloudYeheList := v.([]interface{})
		tmplContents.QCloudYehe = make([]*monitorv20230616.QCloudYeheNoticeTmplMatcher, 0, len(qCloudYeheList))
		for _, item := range qCloudYeheList {
			itemMap := item.(map[string]interface{})
			matcher := &monitorv20230616.QCloudYeheNoticeTmplMatcher{}

			if matchingStatus, ok := itemMap["matching_status"].([]interface{}); ok {
				for _, status := range matchingStatus {
					matcher.MatchingStatus = append(matcher.MatchingStatus, helper.String(status.(string)))
				}
			}

			if template, ok := itemMap["template"].([]interface{}); ok && len(template) > 0 {
				templateMap := template[0].(map[string]interface{})
				matcher.Template = &monitorv20230616.QCloudYeheNoticeTmpl{}

				// Email
				if email, ok := templateMap["email"].([]interface{}); ok && len(email) > 0 {
					emailMap := email[0].(map[string]interface{})
					matcher.Template.Email = &monitorv20230616.QCloudYeheNoticeTmplItem{}
					if contentTmpl, ok := emailMap["content_tmpl"].(string); ok && contentTmpl != "" {
						matcher.Template.Email.ContentTmpl = helper.String(contentTmpl)
					}
					if titleTmpl, ok := emailMap["title_tmpl"].(string); ok && titleTmpl != "" {
						matcher.Template.Email.TitleTmpl = helper.String(titleTmpl)
					}
				}

				// QYWX
				if qywx, ok := templateMap["qywx"].([]interface{}); ok && len(qywx) > 0 {
					qywxMap := qywx[0].(map[string]interface{})
					matcher.Template.QYWX = &monitorv20230616.QCloudYeheNoticeTmplItem{}
					if contentTmpl, ok := qywxMap["content_tmpl"].(string); ok && contentTmpl != "" {
						matcher.Template.QYWX.ContentTmpl = helper.String(contentTmpl)
					}
					if titleTmpl, ok := qywxMap["title_tmpl"].(string); ok && titleTmpl != "" {
						matcher.Template.QYWX.TitleTmpl = helper.String(titleTmpl)
					}
				}

				// SMS
				if sms, ok := templateMap["sms"].([]interface{}); ok && len(sms) > 0 {
					smsMap := sms[0].(map[string]interface{})
					matcher.Template.SMS = &monitorv20230616.QCloudYeheNoticeTmplItem{}
					if contentTmpl, ok := smsMap["content_tmpl"].(string); ok && contentTmpl != "" {
						matcher.Template.SMS.ContentTmpl = helper.String(contentTmpl)
					}
					if titleTmpl, ok := smsMap["title_tmpl"].(string); ok && titleTmpl != "" {
						matcher.Template.SMS.TitleTmpl = helper.String(titleTmpl)
					}
				}

				// Voice
				if voice, ok := templateMap["voice"].([]interface{}); ok && len(voice) > 0 {
					voiceMap := voice[0].(map[string]interface{})
					matcher.Template.Voice = &monitorv20230616.QCloudYeheNoticeTmplItem{}
					if contentTmpl, ok := voiceMap["content_tmpl"].(string); ok && contentTmpl != "" {
						matcher.Template.Voice.ContentTmpl = helper.String(contentTmpl)
					}
					if titleTmpl, ok := voiceMap["title_tmpl"].(string); ok && titleTmpl != "" {
						matcher.Template.Voice.TitleTmpl = helper.String(titleTmpl)
					}
				}

				// WeChat
				if wechat, ok := templateMap["wechat"].([]interface{}); ok && len(wechat) > 0 {
					wechatMap := wechat[0].(map[string]interface{})
					matcher.Template.WeChat = &monitorv20230616.QCloudYeheWeChatNoticeTmplItem{}
					if alarmContentTmpl, ok := wechatMap["alarm_content_tmpl"].(string); ok && alarmContentTmpl != "" {
						matcher.Template.WeChat.AlarmContentTmpl = helper.String(alarmContentTmpl)
					}
					if alarmObjectTmpl, ok := wechatMap["alarm_object_tmpl"].(string); ok && alarmObjectTmpl != "" {
						matcher.Template.WeChat.AlarmObjectTmpl = helper.String(alarmObjectTmpl)
					}
					if alarmRegionTmpl, ok := wechatMap["alarm_region_tmpl"].(string); ok && alarmRegionTmpl != "" {
						matcher.Template.WeChat.AlarmRegionTmpl = helper.String(alarmRegionTmpl)
					}
					if alarmTimeTmpl, ok := wechatMap["alarm_time_tmpl"].(string); ok && alarmTimeTmpl != "" {
						matcher.Template.WeChat.AlarmTimeTmpl = helper.String(alarmTimeTmpl)
					}
				}

				// Site
				if site, ok := templateMap["site"].([]interface{}); ok && len(site) > 0 {
					siteMap := site[0].(map[string]interface{})
					matcher.Template.Site = &monitorv20230616.QCloudYeheNoticeTmplItem{}
					if contentTmpl, ok := siteMap["content_tmpl"].(string); ok && contentTmpl != "" {
						matcher.Template.Site.ContentTmpl = helper.String(contentTmpl)
					}
					if titleTmpl, ok := siteMap["title_tmpl"].(string); ok && titleTmpl != "" {
						matcher.Template.Site.TitleTmpl = helper.String(titleTmpl)
					}
				}

				// Andon
				if andon, ok := templateMap["andon"].([]interface{}); ok && len(andon) > 0 {
					andonMap := andon[0].(map[string]interface{})
					matcher.Template.Andon = &monitorv20230616.QCloudYeheNoticeTmplItem{}
					if contentTmpl, ok := andonMap["content_tmpl"].(string); ok && contentTmpl != "" {
						matcher.Template.Andon.ContentTmpl = helper.String(contentTmpl)
					}
					if titleTmpl, ok := andonMap["title_tmpl"].(string); ok && titleTmpl != "" {
						matcher.Template.Andon.TitleTmpl = helper.String(titleTmpl)
					}
				}
			}

			tmplContents.QCloudYehe = append(tmplContents.QCloudYehe, matcher)
		}
	}

	// WeWorkRobot
	if v, ok := tmplContentsMap["we_work_robot"]; ok && len(v.([]interface{})) > 0 {
		weWorkRobotList := v.([]interface{})
		tmplContents.WeWorkRobot = make([]*monitorv20230616.WeWorkRobotNoticeTmplMatcher, 0, len(weWorkRobotList))
		for _, item := range weWorkRobotList {
			itemMap := item.(map[string]interface{})
			matcher := &monitorv20230616.WeWorkRobotNoticeTmplMatcher{}

			if matchingStatus, ok := itemMap["matching_status"].([]interface{}); ok {
				for _, status := range matchingStatus {
					matcher.MatchingStatus = append(matcher.MatchingStatus, helper.String(status.(string)))
				}
			}

			if template, ok := itemMap["template"].([]interface{}); ok && len(template) > 0 {
				templateMap := template[0].(map[string]interface{})
				matcher.Template = &monitorv20230616.WeWorkRobotNoticeTmpl{}
				if contentTmpl, ok := templateMap["content_tmpl"].(string); ok && contentTmpl != "" {
					matcher.Template.ContentTmpl = helper.String(contentTmpl)
				}
			}

			tmplContents.WeWorkRobot = append(tmplContents.WeWorkRobot, matcher)
		}
	}

	// DingDingRobot
	if v, ok := tmplContentsMap["ding_ding_robot"]; ok && len(v.([]interface{})) > 0 {
		dingDingRobotList := v.([]interface{})
		tmplContents.DingDingRobot = make([]*monitorv20230616.DingDingRobotNoticeTmplMatcher, 0, len(dingDingRobotList))
		for _, item := range dingDingRobotList {
			itemMap := item.(map[string]interface{})
			matcher := &monitorv20230616.DingDingRobotNoticeTmplMatcher{}

			if matchingStatus, ok := itemMap["matching_status"].([]interface{}); ok {
				for _, status := range matchingStatus {
					matcher.MatchingStatus = append(matcher.MatchingStatus, helper.String(status.(string)))
				}
			}

			if template, ok := itemMap["template"].([]interface{}); ok && len(template) > 0 {
				templateMap := template[0].(map[string]interface{})
				matcher.Template = &monitorv20230616.DingDingRobotNoticeTmpl{}
				if contentTmpl, ok := templateMap["content_tmpl"].(string); ok && contentTmpl != "" {
					matcher.Template.ContentTmpl = helper.String(contentTmpl)
				}
				if titleTmpl, ok := templateMap["title_tmpl"].(string); ok && titleTmpl != "" {
					matcher.Template.TitleTmpl = helper.String(titleTmpl)
				}
			}

			tmplContents.DingDingRobot = append(tmplContents.DingDingRobot, matcher)
		}
	}

	// FeiShuRobot
	if v, ok := tmplContentsMap["fei_shu_robot"]; ok && len(v.([]interface{})) > 0 {
		feiShuRobotList := v.([]interface{})
		tmplContents.FeiShuRobot = make([]*monitorv20230616.FeiShuRobotNoticeTmplMatcher, 0, len(feiShuRobotList))
		for _, item := range feiShuRobotList {
			itemMap := item.(map[string]interface{})
			matcher := &monitorv20230616.FeiShuRobotNoticeTmplMatcher{}

			if matchingStatus, ok := itemMap["matching_status"].([]interface{}); ok {
				for _, status := range matchingStatus {
					matcher.MatchingStatus = append(matcher.MatchingStatus, helper.String(status.(string)))
				}
			}

			if template, ok := itemMap["template"].([]interface{}); ok && len(template) > 0 {
				templateMap := template[0].(map[string]interface{})
				matcher.Template = &monitorv20230616.FeiShuRobotNoticeTmpl{}
				if contentTmpl, ok := templateMap["content_tmpl"].(string); ok && contentTmpl != "" {
					matcher.Template.ContentTmpl = helper.String(contentTmpl)
				}
				if titleTmpl, ok := templateMap["title_tmpl"].(string); ok && titleTmpl != "" {
					matcher.Template.TitleTmpl = helper.String(titleTmpl)
				}
			}

			tmplContents.FeiShuRobot = append(tmplContents.FeiShuRobot, matcher)
		}
	}

	// Webhook
	if v, ok := tmplContentsMap["webhook"]; ok && len(v.([]interface{})) > 0 {
		webhookList := v.([]interface{})
		tmplContents.Webhook = make([]*monitorv20230616.WebhookNoticeTmplMatcher, 0, len(webhookList))
		for _, item := range webhookList {
			itemMap := item.(map[string]interface{})
			matcher := &monitorv20230616.WebhookNoticeTmplMatcher{}

			if matchingStatus, ok := itemMap["matching_status"].([]interface{}); ok {
				for _, status := range matchingStatus {
					matcher.MatchingStatus = append(matcher.MatchingStatus, helper.String(status.(string)))
				}
			}

			if template, ok := itemMap["template"].([]interface{}); ok && len(template) > 0 {
				templateMap := template[0].(map[string]interface{})
				matcher.Template = &monitorv20230616.WebhookNoticeTmpl{}
				if body, ok := templateMap["body"].(string); ok && body != "" {
					matcher.Template.Body = helper.String(body)
				}
				if bodyContentType, ok := templateMap["body_content_type"].(string); ok && bodyContentType != "" {
					matcher.Template.BodyContentType = helper.String(bodyContentType)
				}

				if headers, ok := templateMap["headers"].([]interface{}); ok {
					for _, header := range headers {
						headerMap := header.(map[string]interface{})
						h := &monitorv20230616.WebhookNoticeTmplHeader{}
						if key, ok := headerMap["key"].(string); ok && key != "" {
							h.Key = helper.String(key)
						}
						if values, ok := headerMap["values"].([]interface{}); ok {
							for _, value := range values {
								h.Values = append(h.Values, helper.String(value.(string)))
							}
						}
						matcher.Template.Headers = append(matcher.Template.Headers, h)
					}
				}
			}

			tmplContents.Webhook = append(tmplContents.Webhook, matcher)
		}
	}

	// TeamsRobot
	if v, ok := tmplContentsMap["teams_robot"]; ok && len(v.([]interface{})) > 0 {
		teamsRobotList := v.([]interface{})
		tmplContents.TeamsRobot = make([]*monitorv20230616.TeamsRobotNoticeTmplMatcher, 0, len(teamsRobotList))
		for _, item := range teamsRobotList {
			itemMap := item.(map[string]interface{})
			matcher := &monitorv20230616.TeamsRobotNoticeTmplMatcher{}

			if matchingStatus, ok := itemMap["matching_status"].([]interface{}); ok {
				for _, status := range matchingStatus {
					matcher.MatchingStatus = append(matcher.MatchingStatus, helper.String(status.(string)))
				}
			}

			if template, ok := itemMap["template"].([]interface{}); ok && len(template) > 0 {
				templateMap := template[0].(map[string]interface{})
				matcher.Template = &monitorv20230616.TeamsRobotNoticeTmpl{}
				if contentTmpl, ok := templateMap["content_tmpl"].(string); ok && contentTmpl != "" {
					matcher.Template.ContentTmpl = helper.String(contentTmpl)
				}
			}

			tmplContents.TeamsRobot = append(tmplContents.TeamsRobot, matcher)
		}
	}

	// PagerDutyRobot
	if v, ok := tmplContentsMap["pager_duty_robot"]; ok && len(v.([]interface{})) > 0 {
		pagerDutyRobotList := v.([]interface{})
		tmplContents.PagerDutyRobot = make([]*monitorv20230616.PagerDutyRobotNoticeTmplMatcher, 0, len(pagerDutyRobotList))
		for _, item := range pagerDutyRobotList {
			itemMap := item.(map[string]interface{})
			matcher := &monitorv20230616.PagerDutyRobotNoticeTmplMatcher{}

			if matchingStatus, ok := itemMap["matching_status"].([]interface{}); ok {
				for _, status := range matchingStatus {
					matcher.MatchingStatus = append(matcher.MatchingStatus, helper.String(status.(string)))
				}
			}

			if template, ok := itemMap["template"].([]interface{}); ok && len(template) > 0 {
				templateMap := template[0].(map[string]interface{})
				matcher.Template = &monitorv20230616.PagerDutyRobotNoticeTmpl{}
				if body, ok := templateMap["body"].(string); ok && body != "" {
					matcher.Template.Body = helper.String(body)
				}
				if titleTmpl, ok := templateMap["title_tmpl"].(string); ok && titleTmpl != "" {
					matcher.Template.TitleTmpl = helper.String(titleTmpl)
				}

				if headers, ok := templateMap["headers"].([]interface{}); ok {
					for _, header := range headers {
						headerMap := header.(map[string]interface{})
						h := &monitorv20230616.PagerDutyRobotNoticeTmplHeader{}
						if key, ok := headerMap["key"].(string); ok && key != "" {
							h.Key = helper.String(key)
						}
						if values, ok := headerMap["values"].([]interface{}); ok {
							for _, value := range values {
								h.Values = append(h.Values, helper.String(value.(string)))
							}
						}
						matcher.Template.Headers = append(matcher.Template.Headers, h)
					}
				}
			}

			tmplContents.PagerDutyRobot = append(tmplContents.PagerDutyRobot, matcher)
		}
	}

	return tmplContents
}

func flattenNoticeContentTmplItem(item *monitorv20230616.NoticeContentTmplItem) []interface{} {
	if item == nil {
		return nil
	}

	result := make(map[string]interface{})

	// QCloudYehe
	if item.QCloudYehe != nil && len(item.QCloudYehe) > 0 {
		qcloudYeheMap := make(map[string]interface{})
		matcher := item.QCloudYehe[0]

		if matcher.MatchingStatus != nil {
			qcloudYeheMap["matching_status"] = matcher.MatchingStatus
		}

		if matcher.Template != nil {
			templateList := make([]interface{}, 0, 1)
			templateMap := make(map[string]interface{})

			if matcher.Template.Email != nil {
				emailMap := make(map[string]interface{})
				if matcher.Template.Email.TitleTmpl != nil {
					emailMap["title_tmpl"] = matcher.Template.Email.TitleTmpl
				}
				if matcher.Template.Email.ContentTmpl != nil {
					emailMap["content_tmpl"] = matcher.Template.Email.ContentTmpl
				}
				if len(emailMap) > 0 {
					emailList := make([]interface{}, 0, 1)
					emailList = append(emailList, emailMap)
					templateMap["email"] = emailList
				}
			}

			if matcher.Template.QYWX != nil {
				qywxMap := make(map[string]interface{})
				if matcher.Template.QYWX.TitleTmpl != nil {
					qywxMap["title_tmpl"] = matcher.Template.QYWX.TitleTmpl
				}
				if matcher.Template.QYWX.ContentTmpl != nil {
					qywxMap["content_tmpl"] = matcher.Template.QYWX.ContentTmpl
				}
				if len(qywxMap) > 0 {
					qywxList := make([]interface{}, 0, 1)
					qywxList = append(qywxList, qywxMap)
					templateMap["qywx"] = qywxList
				}
			}

			if matcher.Template.SMS != nil {
				smsMap := make(map[string]interface{})
				if matcher.Template.SMS.ContentTmpl != nil {
					smsMap["content_tmpl"] = matcher.Template.SMS.ContentTmpl
				}
				if len(smsMap) > 0 {
					smsList := make([]interface{}, 0, 1)
					smsList = append(smsList, smsMap)
					templateMap["sms"] = smsList
				}
			}

			if matcher.Template.Voice != nil {
				voiceMap := make(map[string]interface{})
				if matcher.Template.Voice.ContentTmpl != nil {
					voiceMap["content_tmpl"] = matcher.Template.Voice.ContentTmpl
				}
				if len(voiceMap) > 0 {
					voiceList := make([]interface{}, 0, 1)
					voiceList = append(voiceList, voiceMap)
					templateMap["voice"] = voiceList
				}
			}

			if matcher.Template.WeChat != nil {
				wechatMap := make(map[string]interface{})
				if matcher.Template.WeChat.AlarmContentTmpl != nil {
					wechatMap["alarm_content_tmpl"] = matcher.Template.WeChat.AlarmContentTmpl
				}
				if matcher.Template.WeChat.AlarmObjectTmpl != nil {
					wechatMap["alarm_object_tmpl"] = matcher.Template.WeChat.AlarmObjectTmpl
				}
				if matcher.Template.WeChat.AlarmRegionTmpl != nil {
					wechatMap["alarm_region_tmpl"] = matcher.Template.WeChat.AlarmRegionTmpl
				}
				if matcher.Template.WeChat.AlarmTimeTmpl != nil {
					wechatMap["alarm_time_tmpl"] = matcher.Template.WeChat.AlarmTimeTmpl
				}
				if len(wechatMap) > 0 {
					wechatList := make([]interface{}, 0, 1)
					wechatList = append(wechatList, wechatMap)
					templateMap["wechat"] = wechatList
				}
			}

			if matcher.Template.Site != nil {
				siteMap := make(map[string]interface{})
				if matcher.Template.Site.ContentTmpl != nil {
					siteMap["content_tmpl"] = matcher.Template.Site.ContentTmpl
				}
				if matcher.Template.Site.TitleTmpl != nil {
					siteMap["title_tmpl"] = matcher.Template.Site.TitleTmpl
				}
				if len(siteMap) > 0 {
					siteList := make([]interface{}, 0, 1)
					siteList = append(siteList, siteMap)
					templateMap["site"] = siteList
				}
			}

			if matcher.Template.Andon != nil {
				andonMap := make(map[string]interface{})
				if matcher.Template.Andon.ContentTmpl != nil {
					andonMap["content_tmpl"] = matcher.Template.Andon.ContentTmpl
				}
				if matcher.Template.Andon.TitleTmpl != nil {
					andonMap["title_tmpl"] = matcher.Template.Andon.TitleTmpl
				}
				if len(andonMap) > 0 {
					andonList := make([]interface{}, 0, 1)
					andonList = append(andonList, andonMap)
					templateMap["andon"] = andonList
				}
			}

			templateList = append(templateList, templateMap)
			qcloudYeheMap["template"] = templateList
		}

		result["qcloud_yehe"] = []interface{}{qcloudYeheMap}
	}

	// WeWorkRobot
	if item.WeWorkRobot != nil && len(item.WeWorkRobot) > 0 {
		weWorkRobotMap := make(map[string]interface{})
		matcher := item.WeWorkRobot[0]

		if matcher.MatchingStatus != nil {
			weWorkRobotMap["matching_status"] = matcher.MatchingStatus
		}

		if matcher.Template != nil {
			templateMap := make(map[string]interface{})
			if matcher.Template.ContentTmpl != nil {
				templateMap["content_tmpl"] = matcher.Template.ContentTmpl
			}
			if len(templateMap) > 0 {
				templateList := make([]interface{}, 0, 1)
				templateList = append(templateList, templateMap)
				weWorkRobotMap["template"] = templateList
			}
		}

		result["we_work_robot"] = []interface{}{weWorkRobotMap}
	}

	// DingDingRobot
	if item.DingDingRobot != nil && len(item.DingDingRobot) > 0 {
		dingDingRobotMap := make(map[string]interface{})
		matcher := item.DingDingRobot[0]

		if matcher.MatchingStatus != nil {
			dingDingRobotMap["matching_status"] = matcher.MatchingStatus
		}

		if matcher.Template != nil {
			templateMap := make(map[string]interface{})
			if matcher.Template.ContentTmpl != nil {
				templateMap["content_tmpl"] = matcher.Template.ContentTmpl
			}
			if matcher.Template.TitleTmpl != nil {
				templateMap["title_tmpl"] = matcher.Template.TitleTmpl
			}
			if len(templateMap) > 0 {
				templateList := make([]interface{}, 0, 1)
				templateList = append(templateList, templateMap)
				dingDingRobotMap["template"] = templateList
			}
		}

		result["ding_ding_robot"] = []interface{}{dingDingRobotMap}
	}

	// FeiShuRobot
	if item.FeiShuRobot != nil && len(item.FeiShuRobot) > 0 {
		feiShuRobotMap := make(map[string]interface{})
		matcher := item.FeiShuRobot[0]

		if matcher.MatchingStatus != nil {
			feiShuRobotMap["matching_status"] = matcher.MatchingStatus
		}

		if matcher.Template != nil {
			templateMap := make(map[string]interface{})
			if matcher.Template.ContentTmpl != nil {
				templateMap["content_tmpl"] = matcher.Template.ContentTmpl
			}
			if matcher.Template.TitleTmpl != nil {
				templateMap["title_tmpl"] = matcher.Template.TitleTmpl
			}
			if len(templateMap) > 0 {
				templateList := make([]interface{}, 0, 1)
				templateList = append(templateList, templateMap)
				feiShuRobotMap["template"] = templateList
			}
		}

		result["fei_shu_robot"] = []interface{}{feiShuRobotMap}
	}

	// Webhook
	if item.Webhook != nil && len(item.Webhook) > 0 {
		webhookMap := make(map[string]interface{})
		matcher := item.Webhook[0]

		if matcher.MatchingStatus != nil {
			webhookMap["matching_status"] = matcher.MatchingStatus
		}

		if matcher.Template != nil {
			templateMap := make(map[string]interface{})

			if matcher.Template.Body != nil {
				templateMap["body"] = matcher.Template.Body
			}

			if matcher.Template.BodyContentType != nil {
				templateMap["body_content_type"] = matcher.Template.BodyContentType
			}

			if matcher.Template.Headers != nil && len(matcher.Template.Headers) > 0 {
				headersList := make([]interface{}, 0, len(matcher.Template.Headers))
				for _, header := range matcher.Template.Headers {
					headerMap := make(map[string]interface{})
					if header.Key != nil {
						headerMap["key"] = header.Key
					}
					if header.Values != nil && len(header.Values) > 0 {
						headerMap["value"] = header.Values[0]
					}
					if len(headerMap) > 0 {
						headersList = append(headersList, headerMap)
					}
				}
				if len(headersList) > 0 {
					templateMap["headers"] = headersList
				}
			}

			if len(templateMap) > 0 {
				templateList := make([]interface{}, 0, 1)
				templateList = append(templateList, templateMap)
				webhookMap["template"] = templateList
			}
		}

		result["webhook"] = []interface{}{webhookMap}
	}

	// TeamsRobot
	if item.TeamsRobot != nil && len(item.TeamsRobot) > 0 {
		teamsRobotMap := make(map[string]interface{})
		matcher := item.TeamsRobot[0]

		if matcher.MatchingStatus != nil {
			teamsRobotMap["matching_status"] = matcher.MatchingStatus
		}

		if matcher.Template != nil {
			templateMap := make(map[string]interface{})
			if matcher.Template.ContentTmpl != nil {
				templateMap["content_tmpl"] = matcher.Template.ContentTmpl
			}
			if len(templateMap) > 0 {
				templateList := make([]interface{}, 0, 1)
				templateList = append(templateList, templateMap)
				teamsRobotMap["template"] = templateList
			}
		}

		result["teams_robot"] = []interface{}{teamsRobotMap}
	}

	// PagerDutyRobot
	if item.PagerDutyRobot != nil && len(item.PagerDutyRobot) > 0 {
		pagerDutyRobotMap := make(map[string]interface{})
		matcher := item.PagerDutyRobot[0]

		if matcher.MatchingStatus != nil {
			pagerDutyRobotMap["matching_status"] = matcher.MatchingStatus
		}

		if matcher.Template != nil {
			templateMap := make(map[string]interface{})

			if matcher.Template.Body != nil {
				templateMap["body"] = matcher.Template.Body
			}

			if matcher.Template.TitleTmpl != nil {
				templateMap["title_tmpl"] = matcher.Template.TitleTmpl
			}

			if matcher.Template.Headers != nil && len(matcher.Template.Headers) > 0 {
				headersList := make([]interface{}, 0, len(matcher.Template.Headers))
				for _, header := range matcher.Template.Headers {
					headerMap := make(map[string]interface{})
					if header.Key != nil {
						headerMap["key"] = header.Key
					}
					if header.Values != nil {
						headerMap["values"] = header.Values
					}
					if len(headerMap) > 0 {
						headersList = append(headersList, headerMap)
					}
				}
				if len(headersList) > 0 {
					templateMap["headers"] = headersList
				}
			}

			if len(templateMap) > 0 {
				templateList := make([]interface{}, 0, 1)
				templateList = append(templateList, templateMap)
				pagerDutyRobotMap["template"] = templateList
			}
		}

		result["pager_duty_robot"] = []interface{}{pagerDutyRobotMap}
	}

	return []interface{}{result}
}
