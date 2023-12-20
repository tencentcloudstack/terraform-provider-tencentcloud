package css

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCssRecordTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssRecordTemplateCreate,
		Read:   resourceTencentCloudCssRecordTemplateRead,
		Update: resourceTencentCloudCssRecordTemplateUpdate,
		Delete: resourceTencentCloudCssRecordTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template name. Only `Chinese`, `English`, `numbers`, `_`, `-` are supported.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description.",
			},

			"flv_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Flv recording parameters are set when Flv recording is enabled.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Recording interval.  Unit: second, default: 1800.  Value range: 30-7200.  This parameter is invalid for HLS. When recording HLS, a file is generated from streaming to streaming.",
						},
						"storage_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Recording storage duration.  Unit: second. Value range: 0 - 1500 days.  0: indicates permanent storage.",
						},
						"enable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether to enable recording in the current format. The default value is 0. 0: No, 1: Yes.",
						},
						"vod_sub_app_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: " The ID of the vodSub app.",
						},
						"vod_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Record file name.Special placeholders supported are: `StreamID`: Stream ID,`StartYear`: Start time - year,`StartMonth`: Start time - month,`StartDay`: Start time - day,`StartHour`: Start time - hour,`StartMinute`: Start time - minutes,`StartSecond`: Start time - seconds,`StartMillisecond`: Start time - milliseconds,`EndYear`: End time - year,`EndMonth`: End time - month,`EndDay`: End time - day,`EndHour`: End time - hour,`EndMinute`: End time - minutes,`EndSecond`: End time - seconds,`EndMillisecond`: End time - millisecondsIf the default recording file name is not set as ,`StreamID`_ ,`StartYear`-,`StartMonth`-,`StartDay`-,`StartHour`-,`StartMinute`-,`StartSecond`_ ,`EndYear`-,`EndMonth`-,`EndDay`-,`EndHour`-,`EndMinute`-,`EndSecond`.",
						},
						"procedure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task flow. This field may return null, indicating that no valid value can be obtained.",
						},
						"storage_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video storage strategy. Normal: standard storage. Cold: low frequency storage. This field may return null, indicating that no valid value can be obtained.",
						},
						"class_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Classification of on-demand applications. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"hls_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Hls recording parameters, which are set when hls recording is enabled.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Recording interval. Unit: second, default: 1800. Value range: 30-7200. This parameter is invalid for HLS. When recording HLS, a file is generated from streaming to streaming.",
						},
						"storage_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Recording storage duration. Unit: second. Value range: 0 - 1500 days. 0: indicates permanent storage.",
						},
						"enable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether to enable recording in the current format. The default value is 0. 0: No, 1: Yes.",
						},
						"vod_sub_app_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The ID of the vodSub app.",
						},
						"vod_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Record file name.Special placeholders supported are: `StreamID`: Stream ID,`StartYear`: Start time - year,`StartMonth`: Start time - month,`StartDay`: Start time - day,`StartHour`: Start time - hour,`StartMinute`: Start time - minutes,`StartSecond`: Start time - seconds,`StartMillisecond`: Start time - milliseconds,`EndYear`: End time - year,`EndMonth`: End time - month,`EndDay`: End time - day,`EndHour`: End time - hour,`EndMinute`: End time - minutes,`EndSecond`: End time - seconds,`EndMillisecond`: End time - millisecondsIf the default recording file name is not set as ,`StreamID`,`StartYear`,`StartMonth`,`StartDay`,`StartHour`,`StartMinute`,`StartSecond`,`EndYear`,`EndMonth`,`EndDay`,`EndHour`,`EndMinute`,`EndSecond`.",
						},
						"procedure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task flow. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"storage_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video storage strategy. Normal: standard storage. Cold: low frequency storage. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"class_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Classification of on-demand applications. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"mp4_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Mp4 recording parameters are set when Mp4 recording is enabled.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Recording interval. Unit: second, default: 1800. Value range: 30-7200. This parameter is invalid for HLS. When recording HLS, a file is generated from streaming to streaming.",
						},
						"storage_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Recording storage duration. Unit: second. Value range: 0 - 1500 days. 0: indicates permanent storage.",
						},
						"enable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether to enable recording in the current format. The default value is 0. 0: No, 1: Yes.",
						},
						"vod_sub_app_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The ID of the on-demand sub app.",
						},
						"vod_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Record file name.Special placeholders supported are: `StreamID`: Stream ID,`StartYear`: Start time - year,`StartMonth`: Start time - month,`StartDay`: Start time - day,`StartHour`: Start time - hour,`StartMinute`: Start time - minutes,`StartSecond`: Start time - seconds,`StartMillisecond`: Start time - milliseconds,`EndYear`: End time - year,`EndMonth`: End time - month,`EndDay`: End time - day,`EndHour`: End time - hour,`EndMinute`: End time - minutes,`EndSecond`: End time - seconds,`EndMillisecond`: End time - millisecondsIf the default recording file name is not set as ,`StreamID`,`StartYear`,`StartMonth`,`StartDay`,`StartHour`,`StartMinute`,`StartSecond`,`EndYear`,`EndMonth`,`EndDay`,`EndHour`,`EndMinute`,`EndSecond`.",
						},
						"procedure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task flow. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"storage_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video storage strategy. Normal: standard storage. Cold: low frequency storage. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"class_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Classification of on-demand applications. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"aac_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Aac recording parameters are set when Aac recording is enabled.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Recording interval. Unit: second, default: 1800. Value range: 30-7200. This parameter is invalid for HLS. When recording HLS, a file is generated from streaming to streaming.",
						},
						"storage_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Recording storage duration. Unit: second. Value range: 0 - 1500 days. 0: indicates permanent storage.",
						},
						"enable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether to enable recording in the current format. The default value is 0. 0: No, 1: Yes.",
						},
						"vod_sub_app_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The ID of the on-demand sub app.",
						},
						"vod_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Record file name.Special placeholders supported are: `StreamID`: Stream ID,`StartYear`: Start time - year,`StartMonth`: Start time - month,`StartDay`: Start time - day,`StartHour`: Start time - hour,`StartMinute`: Start time - minutes,`StartSecond`: Start time - seconds,`StartMillisecond`: Start time - milliseconds,`EndYear`: End time - year,`EndMonth`: End time - month,`EndDay`: End time - day,`EndHour`: End time - hour,`EndMinute`: End time - minutes,`EndSecond`: End time - seconds,`EndMillisecond`: End time - millisecondsIf the default recording file name is not set as ,`StreamID`,`StartYear`,`StartMonth`,`StartDay`,`StartHour`,`StartMinute`,`StartSecond`,`EndYear`,`EndMonth`,`EndDay`,`EndHour`,`EndMinute`,`EndSecond`.",
						},
						"procedure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task flow. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"storage_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video storage strategy. Normal: standard storage. Cold: low frequency storage. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"class_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Classification of on-demand applications. This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"is_delay_live": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Live broadcast type, 0 by default. 0: Ordinary live broadcast, 1: Slow broadcast.",
			},

			"hls_special_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "HLS specific recording parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flow_continue_duration": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "HLS freewheeling timeout. Value range [0, 1800].",
						},
					},
				},
			},

			"mp3_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Mp3 recording parameters are set when Mp3 recording is turned on.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Recording interval. Unit: second, default: 1800. Value range: 30-7200. This parameter is invalid for HLS. When recording HLS, a file is generated from streaming to streaming.",
						},
						"storage_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Recording storage duration. Unit: second. Value range: 0 - 1500 days. 0: indicates permanent storage.",
						},
						"enable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether to enable recording in the current format. The default value is 0. 0: No, 1: Yes.",
						},
						"vod_sub_app_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The ID of the on-demand sub app.",
						},
						"vod_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Record file name.Special placeholders supported are: `StreamID`: Stream ID,`StartYear`: Start time - year,`StartMonth`: Start time - month,`StartDay`: Start time - day,`StartHour`: Start time - hour,`StartMinute`: Start time - minutes,`StartSecond`: Start time - seconds,`StartMillisecond`: Start time - milliseconds,`EndYear`: End time - year,`EndMonth`: End time - month,`EndDay`: End time - day,`EndHour`: End time - hour,`EndMinute`: End time - minutes,`EndSecond`: End time - seconds,`EndMillisecond`: End time - millisecondsIf the default recording file name is not set as ,`StreamID`,`StartYear`,`StartMonth`,`StartDay`,`StartHour`,`StartMinute`, `StartSecond`,`EndYear`,`EndMonth`,`EndDay`,`EndHour`,`EndMinute`,`EndSecond`.",
						},
						"procedure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task flow. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"storage_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video storage strategy. Normal: standard storage. Cold: low frequency storage. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"class_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Classification of vod applications. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"remove_watermark": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to remove the watermark. This parameter is invalid when the type is slow live broadcast.",
			},

			"flv_special_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "FLV records special parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"upload_in_recording": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable the transfer while recording is valid only in the flv format.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCssRecordTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_record_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = css.NewCreateLiveRecordTemplateRequest()
		response   = css.NewCreateLiveRecordTemplateResponse()
		templateId int64
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "flv_param"); ok {
		recordParam := css.RecordParam{}
		if v, ok := dMap["record_interval"]; ok {
			recordParam.RecordInterval = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["storage_time"]; ok {
			recordParam.StorageTime = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["enable"]; ok {
			recordParam.Enable = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["vod_sub_app_id"]; ok {
			recordParam.VodSubAppId = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["vod_file_name"]; ok {
			recordParam.VodFileName = helper.String(v.(string))
		}
		if v, ok := dMap["procedure"]; ok {
			recordParam.Procedure = helper.String(v.(string))
		}
		if v, ok := dMap["storage_mode"]; ok {
			recordParam.StorageMode = helper.String(v.(string))
		}
		if v, ok := dMap["class_id"]; ok {
			recordParam.ClassId = helper.IntInt64(v.(int))
		}
		request.FlvParam = &recordParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "hls_param"); ok {
		recordParam := css.RecordParam{}
		if v, ok := dMap["record_interval"]; ok {
			recordParam.RecordInterval = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["storage_time"]; ok {
			recordParam.StorageTime = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["enable"]; ok {
			recordParam.Enable = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["vod_sub_app_id"]; ok {
			recordParam.VodSubAppId = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["vod_file_name"]; ok {
			recordParam.VodFileName = helper.String(v.(string))
		}
		if v, ok := dMap["procedure"]; ok {
			recordParam.Procedure = helper.String(v.(string))
		}
		if v, ok := dMap["storage_mode"]; ok {
			recordParam.StorageMode = helper.String(v.(string))
		}
		if v, ok := dMap["class_id"]; ok {
			recordParam.ClassId = helper.IntInt64(v.(int))
		}
		request.HlsParam = &recordParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "mp4_param"); ok {
		recordParam := css.RecordParam{}
		if v, ok := dMap["record_interval"]; ok {
			recordParam.RecordInterval = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["storage_time"]; ok {
			recordParam.StorageTime = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["enable"]; ok {
			recordParam.Enable = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["vod_sub_app_id"]; ok {
			recordParam.VodSubAppId = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["vod_file_name"]; ok {
			recordParam.VodFileName = helper.String(v.(string))
		}
		if v, ok := dMap["procedure"]; ok {
			recordParam.Procedure = helper.String(v.(string))
		}
		if v, ok := dMap["storage_mode"]; ok {
			recordParam.StorageMode = helper.String(v.(string))
		}
		if v, ok := dMap["class_id"]; ok {
			recordParam.ClassId = helper.IntInt64(v.(int))
		}
		request.Mp4Param = &recordParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "aac_param"); ok {
		recordParam := css.RecordParam{}
		if v, ok := dMap["record_interval"]; ok {
			recordParam.RecordInterval = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["storage_time"]; ok {
			recordParam.StorageTime = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["enable"]; ok {
			recordParam.Enable = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["vod_sub_app_id"]; ok {
			recordParam.VodSubAppId = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["vod_file_name"]; ok {
			recordParam.VodFileName = helper.String(v.(string))
		}
		if v, ok := dMap["procedure"]; ok {
			recordParam.Procedure = helper.String(v.(string))
		}
		if v, ok := dMap["storage_mode"]; ok {
			recordParam.StorageMode = helper.String(v.(string))
		}
		if v, ok := dMap["class_id"]; ok {
			recordParam.ClassId = helper.IntInt64(v.(int))
		}
		request.AacParam = &recordParam
	}

	if v, ok := d.GetOkExists("is_delay_live"); ok {
		request.IsDelayLive = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "hls_special_param"); ok {
		hlsSpecialParam := css.HlsSpecialParam{}
		if v, ok := dMap["flow_continue_duration"]; ok {
			hlsSpecialParam.FlowContinueDuration = helper.IntUint64(v.(int))
		}
		request.HlsSpecialParam = &hlsSpecialParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "mp3_param"); ok {
		recordParam := css.RecordParam{}
		if v, ok := dMap["record_interval"]; ok {
			recordParam.RecordInterval = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["storage_time"]; ok {
			recordParam.StorageTime = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["enable"]; ok {
			recordParam.Enable = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["vod_sub_app_id"]; ok {
			recordParam.VodSubAppId = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["vod_file_name"]; ok {
			recordParam.VodFileName = helper.String(v.(string))
		}
		if v, ok := dMap["procedure"]; ok {
			recordParam.Procedure = helper.String(v.(string))
		}
		if v, ok := dMap["storage_mode"]; ok {
			recordParam.StorageMode = helper.String(v.(string))
		}
		if v, ok := dMap["class_id"]; ok {
			recordParam.ClassId = helper.IntInt64(v.(int))
		}
		request.Mp3Param = &recordParam
	}

	if v, ok := d.GetOkExists("remove_watermark"); ok {
		request.RemoveWatermark = helper.Bool(v.(bool))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "flv_special_param"); ok {
		flvSpecialParam := css.FlvSpecialParam{}
		if v, ok := dMap["upload_in_recording"]; ok {
			flvSpecialParam.UploadInRecording = helper.Bool(v.(bool))
		}
		request.FlvSpecialParam = &flvSpecialParam
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCssClient().CreateLiveRecordTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css recordTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(templateId))

	return resourceTencentCloudCssRecordTemplateRead(d, meta)
}

func resourceTencentCloudCssRecordTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_record_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	templateId := d.Id()
	templateIdInt64, err := strconv.ParseInt(templateId, 10, 64)
	if err != nil {
		return fmt.Errorf("TemplateId format type error: %s", err.Error())
	}

	recordTemplate, err := service.DescribeCssRecordTemplateById(ctx, templateIdInt64)
	if err != nil {
		return err
	}

	if recordTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssRecordTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if recordTemplate.TemplateName != nil {
		_ = d.Set("template_name", recordTemplate.TemplateName)
	}

	if recordTemplate.Description != nil {
		_ = d.Set("description", recordTemplate.Description)
	}

	if recordTemplate.FlvParam != nil {
		flvParamMap := map[string]interface{}{}

		if recordTemplate.FlvParam.RecordInterval != nil {
			flvParamMap["record_interval"] = recordTemplate.FlvParam.RecordInterval
		}

		if recordTemplate.FlvParam.StorageTime != nil {
			flvParamMap["storage_time"] = recordTemplate.FlvParam.StorageTime
		}

		if recordTemplate.FlvParam.Enable != nil {
			flvParamMap["enable"] = recordTemplate.FlvParam.Enable
		}

		if recordTemplate.FlvParam.VodSubAppId != nil {
			flvParamMap["vod_sub_app_id"] = recordTemplate.FlvParam.VodSubAppId
		}

		if recordTemplate.FlvParam.VodFileName != nil {
			flvParamMap["vod_file_name"] = recordTemplate.FlvParam.VodFileName
		}

		if recordTemplate.FlvParam.Procedure != nil {
			flvParamMap["procedure"] = recordTemplate.FlvParam.Procedure
		}

		if recordTemplate.FlvParam.StorageMode != nil {
			flvParamMap["storage_mode"] = recordTemplate.FlvParam.StorageMode
		}

		if recordTemplate.FlvParam.ClassId != nil {
			flvParamMap["class_id"] = recordTemplate.FlvParam.ClassId
		}

		_ = d.Set("flv_param", []interface{}{flvParamMap})
	}

	if recordTemplate.HlsParam != nil {
		hlsParamMap := map[string]interface{}{}

		if recordTemplate.HlsParam.RecordInterval != nil {
			hlsParamMap["record_interval"] = recordTemplate.HlsParam.RecordInterval
		}

		if recordTemplate.HlsParam.StorageTime != nil {
			hlsParamMap["storage_time"] = recordTemplate.HlsParam.StorageTime
		}

		if recordTemplate.HlsParam.Enable != nil {
			hlsParamMap["enable"] = recordTemplate.HlsParam.Enable
		}

		if recordTemplate.HlsParam.VodSubAppId != nil {
			hlsParamMap["vod_sub_app_id"] = recordTemplate.HlsParam.VodSubAppId
		}

		if recordTemplate.HlsParam.VodFileName != nil {
			hlsParamMap["vod_file_name"] = recordTemplate.HlsParam.VodFileName
		}

		if recordTemplate.HlsParam.Procedure != nil {
			hlsParamMap["procedure"] = recordTemplate.HlsParam.Procedure
		}

		if recordTemplate.HlsParam.StorageMode != nil {
			hlsParamMap["storage_mode"] = recordTemplate.HlsParam.StorageMode
		}

		if recordTemplate.HlsParam.ClassId != nil {
			hlsParamMap["class_id"] = recordTemplate.HlsParam.ClassId
		}

		_ = d.Set("hls_param", []interface{}{hlsParamMap})
	}

	if recordTemplate.Mp4Param != nil {
		mp4ParamMap := map[string]interface{}{}

		if recordTemplate.Mp4Param.RecordInterval != nil {
			mp4ParamMap["record_interval"] = recordTemplate.Mp4Param.RecordInterval
		}

		if recordTemplate.Mp4Param.StorageTime != nil {
			mp4ParamMap["storage_time"] = recordTemplate.Mp4Param.StorageTime
		}

		if recordTemplate.Mp4Param.Enable != nil {
			mp4ParamMap["enable"] = recordTemplate.Mp4Param.Enable
		}

		if recordTemplate.Mp4Param.VodSubAppId != nil {
			mp4ParamMap["vod_sub_app_id"] = recordTemplate.Mp4Param.VodSubAppId
		}

		if recordTemplate.Mp4Param.VodFileName != nil {
			mp4ParamMap["vod_file_name"] = recordTemplate.Mp4Param.VodFileName
		}

		if recordTemplate.Mp4Param.Procedure != nil {
			mp4ParamMap["procedure"] = recordTemplate.Mp4Param.Procedure
		}

		if recordTemplate.Mp4Param.StorageMode != nil {
			mp4ParamMap["storage_mode"] = recordTemplate.Mp4Param.StorageMode
		}

		if recordTemplate.Mp4Param.ClassId != nil {
			mp4ParamMap["class_id"] = recordTemplate.Mp4Param.ClassId
		}

		_ = d.Set("mp4_param", []interface{}{mp4ParamMap})
	}

	if recordTemplate.AacParam != nil {
		aacParamMap := map[string]interface{}{}

		if recordTemplate.AacParam.RecordInterval != nil {
			aacParamMap["record_interval"] = recordTemplate.AacParam.RecordInterval
		}

		if recordTemplate.AacParam.StorageTime != nil {
			aacParamMap["storage_time"] = recordTemplate.AacParam.StorageTime
		}

		if recordTemplate.AacParam.Enable != nil {
			aacParamMap["enable"] = recordTemplate.AacParam.Enable
		}

		if recordTemplate.AacParam.VodSubAppId != nil {
			aacParamMap["vod_sub_app_id"] = recordTemplate.AacParam.VodSubAppId
		}

		if recordTemplate.AacParam.VodFileName != nil {
			aacParamMap["vod_file_name"] = recordTemplate.AacParam.VodFileName
		}

		if recordTemplate.AacParam.Procedure != nil {
			aacParamMap["procedure"] = recordTemplate.AacParam.Procedure
		}

		if recordTemplate.AacParam.StorageMode != nil {
			aacParamMap["storage_mode"] = recordTemplate.AacParam.StorageMode
		}

		if recordTemplate.AacParam.ClassId != nil {
			aacParamMap["class_id"] = recordTemplate.AacParam.ClassId
		}

		_ = d.Set("aac_param", []interface{}{aacParamMap})
	}

	if recordTemplate.IsDelayLive != nil {
		_ = d.Set("is_delay_live", recordTemplate.IsDelayLive)
	}

	if recordTemplate.HlsSpecialParam != nil {
		hlsSpecialParamMap := map[string]interface{}{}

		if recordTemplate.HlsSpecialParam.FlowContinueDuration != nil {
			hlsSpecialParamMap["flow_continue_duration"] = recordTemplate.HlsSpecialParam.FlowContinueDuration
		}

		_ = d.Set("hls_special_param", []interface{}{hlsSpecialParamMap})
	}

	if recordTemplate.Mp3Param != nil {
		mp3ParamMap := map[string]interface{}{}

		if recordTemplate.Mp3Param.RecordInterval != nil {
			mp3ParamMap["record_interval"] = recordTemplate.Mp3Param.RecordInterval
		}

		if recordTemplate.Mp3Param.StorageTime != nil {
			mp3ParamMap["storage_time"] = recordTemplate.Mp3Param.StorageTime
		}

		if recordTemplate.Mp3Param.Enable != nil {
			mp3ParamMap["enable"] = recordTemplate.Mp3Param.Enable
		}

		if recordTemplate.Mp3Param.VodSubAppId != nil {
			mp3ParamMap["vod_sub_app_id"] = recordTemplate.Mp3Param.VodSubAppId
		}

		if recordTemplate.Mp3Param.VodFileName != nil {
			mp3ParamMap["vod_file_name"] = recordTemplate.Mp3Param.VodFileName
		}

		if recordTemplate.Mp3Param.Procedure != nil {
			mp3ParamMap["procedure"] = recordTemplate.Mp3Param.Procedure
		}

		if recordTemplate.Mp3Param.StorageMode != nil {
			mp3ParamMap["storage_mode"] = recordTemplate.Mp3Param.StorageMode
		}

		if recordTemplate.Mp3Param.ClassId != nil {
			mp3ParamMap["class_id"] = recordTemplate.Mp3Param.ClassId
		}

		_ = d.Set("mp3_param", []interface{}{mp3ParamMap})
	}

	if recordTemplate.RemoveWatermark != nil {
		_ = d.Set("remove_watermark", recordTemplate.RemoveWatermark)
	}

	if recordTemplate.FlvSpecialParam != nil {
		flvSpecialParamMap := map[string]interface{}{}

		if recordTemplate.FlvSpecialParam.UploadInRecording != nil {
			flvSpecialParamMap["upload_in_recording"] = recordTemplate.FlvSpecialParam.UploadInRecording
		}

		_ = d.Set("flv_special_param", []interface{}{flvSpecialParamMap})
	}

	return nil
}

func resourceTencentCloudCssRecordTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_record_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := css.NewModifyLiveRecordTemplateRequest()

	templateId := d.Id()
	templateIdInt64, _ := strconv.ParseInt(templateId, 10, 64)

	request.TemplateId = &templateIdInt64

	if d.HasChange("template_name") {
		if v, ok := d.GetOk("template_name"); ok {
			request.TemplateName = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("flv_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "flv_param"); ok {
			recordParam := css.RecordParam{}
			if v, ok := dMap["record_interval"]; ok {
				recordParam.RecordInterval = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["storage_time"]; ok {
				recordParam.StorageTime = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["enable"]; ok {
				recordParam.Enable = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["vod_sub_app_id"]; ok {
				recordParam.VodSubAppId = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["vod_file_name"]; ok {
				recordParam.VodFileName = helper.String(v.(string))
			}
			if v, ok := dMap["procedure"]; ok {
				recordParam.Procedure = helper.String(v.(string))
			}
			if v, ok := dMap["storage_mode"]; ok {
				recordParam.StorageMode = helper.String(v.(string))
			}
			if v, ok := dMap["class_id"]; ok {
				recordParam.ClassId = helper.IntInt64(v.(int))
			}
			request.FlvParam = &recordParam
		}
	}

	if d.HasChange("hls_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "hls_param"); ok {
			recordParam := css.RecordParam{}
			if v, ok := dMap["record_interval"]; ok {
				recordParam.RecordInterval = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["storage_time"]; ok {
				recordParam.StorageTime = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["enable"]; ok {
				recordParam.Enable = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["vod_sub_app_id"]; ok {
				recordParam.VodSubAppId = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["vod_file_name"]; ok {
				recordParam.VodFileName = helper.String(v.(string))
			}
			if v, ok := dMap["procedure"]; ok {
				recordParam.Procedure = helper.String(v.(string))
			}
			if v, ok := dMap["storage_mode"]; ok {
				recordParam.StorageMode = helper.String(v.(string))
			}
			if v, ok := dMap["class_id"]; ok {
				recordParam.ClassId = helper.IntInt64(v.(int))
			}
			request.HlsParam = &recordParam
		}
	}

	if d.HasChange("mp4_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "mp4_param"); ok {
			recordParam := css.RecordParam{}
			if v, ok := dMap["record_interval"]; ok {
				recordParam.RecordInterval = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["storage_time"]; ok {
				recordParam.StorageTime = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["enable"]; ok {
				recordParam.Enable = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["vod_sub_app_id"]; ok {
				recordParam.VodSubAppId = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["vod_file_name"]; ok {
				recordParam.VodFileName = helper.String(v.(string))
			}
			if v, ok := dMap["procedure"]; ok {
				recordParam.Procedure = helper.String(v.(string))
			}
			if v, ok := dMap["storage_mode"]; ok {
				recordParam.StorageMode = helper.String(v.(string))
			}
			if v, ok := dMap["class_id"]; ok {
				recordParam.ClassId = helper.IntInt64(v.(int))
			}
			request.Mp4Param = &recordParam
		}
	}

	if d.HasChange("aac_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "aac_param"); ok {
			recordParam := css.RecordParam{}
			if v, ok := dMap["record_interval"]; ok {
				recordParam.RecordInterval = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["storage_time"]; ok {
				recordParam.StorageTime = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["enable"]; ok {
				recordParam.Enable = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["vod_sub_app_id"]; ok {
				recordParam.VodSubAppId = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["vod_file_name"]; ok {
				recordParam.VodFileName = helper.String(v.(string))
			}
			if v, ok := dMap["procedure"]; ok {
				recordParam.Procedure = helper.String(v.(string))
			}
			if v, ok := dMap["storage_mode"]; ok {
				recordParam.StorageMode = helper.String(v.(string))
			}
			if v, ok := dMap["class_id"]; ok {
				recordParam.ClassId = helper.IntInt64(v.(int))
			}
			request.AacParam = &recordParam
		}
	}

	if d.HasChange("hls_special_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "hls_special_param"); ok {
			hlsSpecialParam := css.HlsSpecialParam{}
			if v, ok := dMap["flow_continue_duration"]; ok {
				hlsSpecialParam.FlowContinueDuration = helper.IntUint64(v.(int))
			}
			request.HlsSpecialParam = &hlsSpecialParam
		}
	}

	if d.HasChange("mp3_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "mp3_param"); ok {
			recordParam := css.RecordParam{}
			if v, ok := dMap["record_interval"]; ok {
				recordParam.RecordInterval = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["storage_time"]; ok {
				recordParam.StorageTime = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["enable"]; ok {
				recordParam.Enable = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["vod_sub_app_id"]; ok {
				recordParam.VodSubAppId = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["vod_file_name"]; ok {
				recordParam.VodFileName = helper.String(v.(string))
			}
			if v, ok := dMap["procedure"]; ok {
				recordParam.Procedure = helper.String(v.(string))
			}
			if v, ok := dMap["storage_mode"]; ok {
				recordParam.StorageMode = helper.String(v.(string))
			}
			if v, ok := dMap["class_id"]; ok {
				recordParam.ClassId = helper.IntInt64(v.(int))
			}
			request.Mp3Param = &recordParam
		}
	}

	if d.HasChange("remove_watermark") {
		if v, ok := d.GetOkExists("remove_watermark"); ok {
			request.RemoveWatermark = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("flv_special_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "flv_special_param"); ok {
			flvSpecialParam := css.FlvSpecialParam{}
			if v, ok := dMap["upload_in_recording"]; ok {
				flvSpecialParam.UploadInRecording = helper.Bool(v.(bool))
			}
			request.FlvSpecialParam = &flvSpecialParam
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCssClient().ModifyLiveRecordTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update css recordTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCssRecordTemplateRead(d, meta)
}

func resourceTencentCloudCssRecordTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_record_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	templateId := d.Id()
	templateIdInt64, _ := strconv.ParseInt(templateId, 10, 64)

	if err := service.DeleteCssRecordTemplateById(ctx, templateIdInt64); err != nil {
		return err
	}

	return nil
}
