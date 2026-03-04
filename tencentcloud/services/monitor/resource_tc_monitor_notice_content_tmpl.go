package monitor

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20230616"

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
				ForceNew:    true,
				Description: "Template name.",
			},

			"monitor_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Monitor type, e.g., `MT_QCE`.",
			},

			"tmpl_language": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Template language, valid values: `en`, `zh`.",
			},

			"tmpl_contents": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Template content configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching_status": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Matching status list, e.g., [\"Trigger\"].",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"template": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Template details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"q_cloud_yehe": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "QCloud Yehe configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sms": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "SMS template configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"title_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "SMS title template.",
															},
															"content_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "SMS content template.",
															},
														},
													},
												},
												"email": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Email template configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"title_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Email title template.",
															},
															"content_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Email content template.",
															},
														},
													},
												},
												"voice": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Voice template configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"content_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Voice content template.",
															},
														},
													},
												},
												"site": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Site message template configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"title_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Site message title template.",
															},
															"content_tmpl": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Site message content template.",
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
										MaxItems:    1,
										Description: "WeWork robot template configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"title_tmpl": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "WeWork robot title template.",
												},
												"content_tmpl": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "WeWork robot content template.",
												},
											},
										},
									},
									"ding_ding_robot": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "DingDing robot template configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"title_tmpl": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "DingDing robot title template.",
												},
												"content_tmpl": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "DingDing robot content template.",
												},
											},
										},
									},
									"fei_shu_robot": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "FeiShu robot template configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"title_tmpl": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "FeiShu robot title template.",
												},
												"content_tmpl": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "FeiShu robot content template.",
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
		request  = monitor.NewCreateNoticeContentTmplRequest()
		tmplId   string
		tmplName string
	)

	if v, ok := d.GetOk("tmpl_name"); ok {
		tmplName = v.(string)
		request.TmplName = helper.String(tmplName)
	}

	if v, ok := d.GetOk("monitor_type"); ok {
		request.MonitorType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tmpl_language"); ok {
		request.TmplLanguage = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tmpl_contents"); ok {
		tmplContents := make([]*monitor.TmplContents, 0)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			tmplContent := monitor.TmplContents{}

			if v, ok := m["matching_status"]; ok {
				matchingStatusSet := v.(*schema.Set).List()
				matchingStatusArr := make([]*string, 0, len(matchingStatusSet))
				for _, status := range matchingStatusSet {
					matchingStatusArr = append(matchingStatusArr, helper.String(status.(string)))
				}
				tmplContent.MatchingStatus = matchingStatusArr
			}

			if v, ok := m["template"]; ok && len(v.([]interface{})) > 0 {
				templateMap := v.([]interface{})[0].(map[string]interface{})
				template := &monitor.Template{}

				// QCloud Yehe
				if v, ok := templateMap["q_cloud_yehe"]; ok && len(v.([]interface{})) > 0 {
					qCloudYeheMap := v.([]interface{})[0].(map[string]interface{})
					qCloudYehe := &monitor.QCloudYehe{}

					if v, ok := qCloudYeheMap["sms"]; ok && len(v.([]interface{})) > 0 {
						smsMap := v.([]interface{})[0].(map[string]interface{})
						sms := &monitor.ChannelTemplate{}
						if title, ok := smsMap["title_tmpl"].(string); ok && title != "" {
							sms.TitleTmpl = helper.String(title)
						}
						if content, ok := smsMap["content_tmpl"].(string); ok && content != "" {
							sms.ContentTmpl = helper.String(content)
						}
						qCloudYehe.Sms = sms
					}

					if v, ok := qCloudYeheMap["email"]; ok && len(v.([]interface{})) > 0 {
						emailMap := v.([]interface{})[0].(map[string]interface{})
						email := &monitor.ChannelTemplate{}
						if title, ok := emailMap["title_tmpl"].(string); ok && title != "" {
							email.TitleTmpl = helper.String(title)
						}
						if content, ok := emailMap["content_tmpl"].(string); ok && content != "" {
							email.ContentTmpl = helper.String(content)
						}
						qCloudYehe.Email = email
					}

					if v, ok := qCloudYeheMap["voice"]; ok && len(v.([]interface{})) > 0 {
						voiceMap := v.([]interface{})[0].(map[string]interface{})
						voice := &monitor.ChannelTemplate{}
						if content, ok := voiceMap["content_tmpl"].(string); ok && content != "" {
							voice.ContentTmpl = helper.String(content)
						}
						qCloudYehe.Voice = voice
					}

					if v, ok := qCloudYeheMap["site"]; ok && len(v.([]interface{})) > 0 {
						siteMap := v.([]interface{})[0].(map[string]interface{})
						site := &monitor.ChannelTemplate{}
						if title, ok := siteMap["title_tmpl"].(string); ok && title != "" {
							site.TitleTmpl = helper.String(title)
						}
						if content, ok := siteMap["content_tmpl"].(string); ok && content != "" {
							site.ContentTmpl = helper.String(content)
						}
						qCloudYehe.Site = site
					}

					template.QCloudYehe = qCloudYehe
				}

				// WeWork Robot
				if v, ok := templateMap["we_work_robot"]; ok && len(v.([]interface{})) > 0 {
					robotMap := v.([]interface{})[0].(map[string]interface{})
					robot := &monitor.ChannelTemplate{}
					if title, ok := robotMap["title_tmpl"].(string); ok && title != "" {
						robot.TitleTmpl = helper.String(title)
					}
					if content, ok := robotMap["content_tmpl"].(string); ok && content != "" {
						robot.ContentTmpl = helper.String(content)
					}
					template.WeWorkRobot = robot
				}

				// DingDing Robot
				if v, ok := templateMap["ding_ding_robot"]; ok && len(v.([]interface{})) > 0 {
					robotMap := v.([]interface{})[0].(map[string]interface{})
					robot := &monitor.ChannelTemplate{}
					if title, ok := robotMap["title_tmpl"].(string); ok && title != "" {
						robot.TitleTmpl = helper.String(title)
					}
					if content, ok := robotMap["content_tmpl"].(string); ok && content != "" {
						robot.ContentTmpl = helper.String(content)
					}
					template.DingDingRobot = robot
				}

				// FeiShu Robot
				if v, ok := templateMap["fei_shu_robot"]; ok && len(v.([]interface{})) > 0 {
					robotMap := v.([]interface{})[0].(map[string]interface{})
					robot := &monitor.ChannelTemplate{}
					if title, ok := robotMap["title_tmpl"].(string); ok && title != "" {
						robot.TitleTmpl = helper.String(title)
					}
					if content, ok := robotMap["content_tmpl"].(string); ok && content != "" {
						robot.ContentTmpl = helper.String(content)
					}
					template.FeiShuRobot = robot
				}

				tmplContent.Template = template
			}

			tmplContents = append(tmplContents, &tmplContent)
		}
		request.TmplContents = tmplContents
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorV20230616Client().CreateNoticeContentTmpl(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result.Response.TmplID == nil {
			e = fmt.Errorf("monitor notice content tmpl TmplID is nil")
			return resource.NonRetryableError(e)
		}

		tmplId = *result.Response.TmplID
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create monitor notice content tmpl failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{tmplId, tmplName}, tccommon.FILED_SP))

	return resourceTencentCloudMonitorNoticeContentTmplRead(d, meta)
}

func resourceTencentCloudMonitorNoticeContentTmplRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_notice_content_tmpl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.Background()
		service = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	tmplId := idSplit[0]
	tmplName := idSplit[1]

	respData, err := service.DescribeMonitorNoticeContentTmplById(ctx, tmplId, tmplName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_monitor_notice_content_tmpl` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
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

	if respData.TmplID != nil {
		_ = d.Set("tmpl_id", respData.TmplID)
	}

	if respData.TmplContents != nil && len(respData.TmplContents) > 0 {
		tmplContentsList := make([]map[string]interface{}, 0, len(respData.TmplContents))
		for _, tmplContent := range respData.TmplContents {
			tmplContentMap := make(map[string]interface{})

			if tmplContent.MatchingStatus != nil {
				matchingStatusList := make([]string, 0, len(tmplContent.MatchingStatus))
				for _, status := range tmplContent.MatchingStatus {
					if status != nil {
						matchingStatusList = append(matchingStatusList, *status)
					}
				}
				tmplContentMap["matching_status"] = matchingStatusList
			}

			if tmplContent.Template != nil {
				templateList := make([]map[string]interface{}, 0, 1)
				templateMap := make(map[string]interface{})

				// QCloud Yehe
				if tmplContent.Template.QCloudYehe != nil {
					qCloudYeheList := make([]map[string]interface{}, 0, 1)
					qCloudYeheMap := make(map[string]interface{})

					if tmplContent.Template.QCloudYehe.Sms != nil {
						smsList := make([]map[string]interface{}, 0, 1)
						smsMap := make(map[string]interface{})
						if tmplContent.Template.QCloudYehe.Sms.TitleTmpl != nil {
							smsMap["title_tmpl"] = *tmplContent.Template.QCloudYehe.Sms.TitleTmpl
						}
						if tmplContent.Template.QCloudYehe.Sms.ContentTmpl != nil {
							smsMap["content_tmpl"] = *tmplContent.Template.QCloudYehe.Sms.ContentTmpl
						}
						smsList = append(smsList, smsMap)
						qCloudYeheMap["sms"] = smsList
					}

					if tmplContent.Template.QCloudYehe.Email != nil {
						emailList := make([]map[string]interface{}, 0, 1)
						emailMap := make(map[string]interface{})
						if tmplContent.Template.QCloudYehe.Email.TitleTmpl != nil {
							emailMap["title_tmpl"] = *tmplContent.Template.QCloudYehe.Email.TitleTmpl
						}
						if tmplContent.Template.QCloudYehe.Email.ContentTmpl != nil {
							emailMap["content_tmpl"] = *tmplContent.Template.QCloudYehe.Email.ContentTmpl
						}
						emailList = append(emailList, emailMap)
						qCloudYeheMap["email"] = emailList
					}

					if tmplContent.Template.QCloudYehe.Voice != nil {
						voiceList := make([]map[string]interface{}, 0, 1)
						voiceMap := make(map[string]interface{})
						if tmplContent.Template.QCloudYehe.Voice.ContentTmpl != nil {
							voiceMap["content_tmpl"] = *tmplContent.Template.QCloudYehe.Voice.ContentTmpl
						}
						voiceList = append(voiceList, voiceMap)
						qCloudYeheMap["voice"] = voiceList
					}

					if tmplContent.Template.QCloudYehe.Site != nil {
						siteList := make([]map[string]interface{}, 0, 1)
						siteMap := make(map[string]interface{})
						if tmplContent.Template.QCloudYehe.Site.TitleTmpl != nil {
							siteMap["title_tmpl"] = *tmplContent.Template.QCloudYehe.Site.TitleTmpl
						}
						if tmplContent.Template.QCloudYehe.Site.ContentTmpl != nil {
							siteMap["content_tmpl"] = *tmplContent.Template.QCloudYehe.Site.ContentTmpl
						}
						siteList = append(siteList, siteMap)
						qCloudYeheMap["site"] = siteList
					}

					qCloudYeheList = append(qCloudYeheList, qCloudYeheMap)
					templateMap["q_cloud_yehe"] = qCloudYeheList
				}

				// WeWork Robot
				if tmplContent.Template.WeWorkRobot != nil {
					robotList := make([]map[string]interface{}, 0, 1)
					robotMap := make(map[string]interface{})
					if tmplContent.Template.WeWorkRobot.TitleTmpl != nil {
						robotMap["title_tmpl"] = *tmplContent.Template.WeWorkRobot.TitleTmpl
					}
					if tmplContent.Template.WeWorkRobot.ContentTmpl != nil {
						robotMap["content_tmpl"] = *tmplContent.Template.WeWorkRobot.ContentTmpl
					}
					robotList = append(robotList, robotMap)
					templateMap["we_work_robot"] = robotList
				}

				// DingDing Robot
				if tmplContent.Template.DingDingRobot != nil {
					robotList := make([]map[string]interface{}, 0, 1)
					robotMap := make(map[string]interface{})
					if tmplContent.Template.DingDingRobot.TitleTmpl != nil {
						robotMap["title_tmpl"] = *tmplContent.Template.DingDingRobot.TitleTmpl
					}
					if tmplContent.Template.DingDingRobot.ContentTmpl != nil {
						robotMap["content_tmpl"] = *tmplContent.Template.DingDingRobot.ContentTmpl
					}
					robotList = append(robotList, robotMap)
					templateMap["ding_ding_robot"] = robotList
				}

				// FeiShu Robot
				if tmplContent.Template.FeiShuRobot != nil {
					robotList := make([]map[string]interface{}, 0, 1)
					robotMap := make(map[string]interface{})
					if tmplContent.Template.FeiShuRobot.TitleTmpl != nil {
						robotMap["title_tmpl"] = *tmplContent.Template.FeiShuRobot.TitleTmpl
					}
					if tmplContent.Template.FeiShuRobot.ContentTmpl != nil {
						robotMap["content_tmpl"] = *tmplContent.Template.FeiShuRobot.ContentTmpl
					}
					robotList = append(robotList, robotMap)
					templateMap["fei_shu_robot"] = robotList
				}

				templateList = append(templateList, templateMap)
				tmplContentMap["template"] = templateList
			}

			tmplContentsList = append(tmplContentsList, tmplContentMap)
		}
		_ = d.Set("tmpl_contents", tmplContentsList)
	}

	return nil
}

func resourceTencentCloudMonitorNoticeContentTmplUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_notice_content_tmpl.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = monitor.NewModifyNoticeContentTmplRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	tmplId := idSplit[0]

	request.TmplID = helper.String(tmplId)

	if d.HasChange("tmpl_contents") {
		if v, ok := d.GetOk("tmpl_contents"); ok {
			tmplContents := make([]*monitor.TmplContents, 0)
			for _, item := range v.([]interface{}) {
				m := item.(map[string]interface{})
				tmplContent := monitor.TmplContents{}

				if v, ok := m["matching_status"]; ok {
					matchingStatusSet := v.(*schema.Set).List()
					matchingStatusArr := make([]*string, 0, len(matchingStatusSet))
					for _, status := range matchingStatusSet {
						matchingStatusArr = append(matchingStatusArr, helper.String(status.(string)))
					}
					tmplContent.MatchingStatus = matchingStatusArr
				}

				if v, ok := m["template"]; ok && len(v.([]interface{})) > 0 {
					templateMap := v.([]interface{})[0].(map[string]interface{})
					template := &monitor.Template{}

					// QCloud Yehe
					if v, ok := templateMap["q_cloud_yehe"]; ok && len(v.([]interface{})) > 0 {
						qCloudYeheMap := v.([]interface{})[0].(map[string]interface{})
						qCloudYehe := &monitor.QCloudYehe{}

						if v, ok := qCloudYeheMap["sms"]; ok && len(v.([]interface{})) > 0 {
							smsMap := v.([]interface{})[0].(map[string]interface{})
							sms := &monitor.ChannelTemplate{}
							if title, ok := smsMap["title_tmpl"].(string); ok && title != "" {
								sms.TitleTmpl = helper.String(title)
							}
							if content, ok := smsMap["content_tmpl"].(string); ok && content != "" {
								sms.ContentTmpl = helper.String(content)
							}
							qCloudYehe.Sms = sms
						}

						if v, ok := qCloudYeheMap["email"]; ok && len(v.([]interface{})) > 0 {
							emailMap := v.([]interface{})[0].(map[string]interface{})
							email := &monitor.ChannelTemplate{}
							if title, ok := emailMap["title_tmpl"].(string); ok && title != "" {
								email.TitleTmpl = helper.String(title)
							}
							if content, ok := emailMap["content_tmpl"].(string); ok && content != "" {
								email.ContentTmpl = helper.String(content)
							}
							qCloudYehe.Email = email
						}

						if v, ok := qCloudYeheMap["voice"]; ok && len(v.([]interface{})) > 0 {
							voiceMap := v.([]interface{})[0].(map[string]interface{})
							voice := &monitor.ChannelTemplate{}
							if content, ok := voiceMap["content_tmpl"].(string); ok && content != "" {
								voice.ContentTmpl = helper.String(content)
							}
							qCloudYehe.Voice = voice
						}

						if v, ok := qCloudYeheMap["site"]; ok && len(v.([]interface{})) > 0 {
							siteMap := v.([]interface{})[0].(map[string]interface{})
							site := &monitor.ChannelTemplate{}
							if title, ok := siteMap["title_tmpl"].(string); ok && title != "" {
								site.TitleTmpl = helper.String(title)
							}
							if content, ok := siteMap["content_tmpl"].(string); ok && content != "" {
								site.ContentTmpl = helper.String(content)
							}
							qCloudYehe.Site = site
						}

						template.QCloudYehe = qCloudYehe
					}

					// WeWork Robot
					if v, ok := templateMap["we_work_robot"]; ok && len(v.([]interface{})) > 0 {
						robotMap := v.([]interface{})[0].(map[string]interface{})
						robot := &monitor.ChannelTemplate{}
						if title, ok := robotMap["title_tmpl"].(string); ok && title != "" {
							robot.TitleTmpl = helper.String(title)
						}
						if content, ok := robotMap["content_tmpl"].(string); ok && content != "" {
							robot.ContentTmpl = helper.String(content)
						}
						template.WeWorkRobot = robot
					}

					// DingDing Robot
					if v, ok := templateMap["ding_ding_robot"]; ok && len(v.([]interface{})) > 0 {
						robotMap := v.([]interface{})[0].(map[string]interface{})
						robot := &monitor.ChannelTemplate{}
						if title, ok := robotMap["title_tmpl"].(string); ok && title != "" {
							robot.TitleTmpl = helper.String(title)
						}
						if content, ok := robotMap["content_tmpl"].(string); ok && content != "" {
							robot.ContentTmpl = helper.String(content)
						}
						template.DingDingRobot = robot
					}

					// FeiShu Robot
					if v, ok := templateMap["fei_shu_robot"]; ok && len(v.([]interface{})) > 0 {
						robotMap := v.([]interface{})[0].(map[string]interface{})
						robot := &monitor.ChannelTemplate{}
						if title, ok := robotMap["title_tmpl"].(string); ok && title != "" {
							robot.TitleTmpl = helper.String(title)
						}
						if content, ok := robotMap["content_tmpl"].(string); ok && content != "" {
							robot.ContentTmpl = helper.String(content)
						}
						template.FeiShuRobot = robot
					}

					tmplContent.Template = template
				}

				tmplContents = append(tmplContents, &tmplContent)
			}
			request.TmplContents = tmplContents
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorV20230616Client().ModifyNoticeContentTmpl(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update monitor notice content tmpl failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorNoticeContentTmplRead(d, meta)
}

func resourceTencentCloudMonitorNoticeContentTmplDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_notice_content_tmpl.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = monitor.NewDeleteNoticeContentTmplsRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	tmplId := idSplit[0]

	request.TmplIDs = []*string{helper.String(tmplId)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorV20230616Client().DeleteNoticeContentTmpls(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete monitor notice content tmpl failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
