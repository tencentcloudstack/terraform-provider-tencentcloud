/*
Provides a resource to create a live record_template

Example Usage

```hcl
resource "tencentcloud_live_record_template" "record_template" {
  template_name = ""
  description = ""
  flv_param {
		record_interval =
		storage_time =
		enable =
		vod_sub_app_id =
		vod_file_name = ""
		procedure = ""
		storage_mode = ""
		class_id =

  }
  hls_param {
		record_interval =
		storage_time =
		enable =
		vod_sub_app_id =
		vod_file_name = ""
		procedure = ""
		storage_mode = ""
		class_id =

  }
  mp4_param {
		record_interval =
		storage_time =
		enable =
		vod_sub_app_id =
		vod_file_name = ""
		procedure = ""
		storage_mode = ""
		class_id =

  }
  aac_param {
		record_interval =
		storage_time =
		enable =
		vod_sub_app_id =
		vod_file_name = ""
		procedure = ""
		storage_mode = ""
		class_id =

  }
  is_delay_live =
  hls_special_param {
		flow_continue_duration =

  }
  mp3_param {
		record_interval =
		storage_time =
		enable =
		vod_sub_app_id =
		vod_file_name = ""
		procedure = ""
		storage_mode = ""
		class_id =

  }
  remove_watermark =
  flv_special_param {
		upload_in_recording =

  }
}
```

Import

live record_template can be imported using the id, e.g.

```
terraform import tencentcloud_live_record_template.record_template record_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudLiveRecordTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveRecordTemplateCreate,
		Read:   resourceTencentCloudLiveRecordTemplateRead,
		Update: resourceTencentCloudLiveRecordTemplateUpdate,
		Delete: resourceTencentCloudLiveRecordTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template name. Only letters, digits, underscores, and hyphens can be contained.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Message description.",
			},

			"flv_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "FLV recording parameter, which is set when FLV recording is enabled.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Max recording time per fileDefault value: `1800` (seconds)Value range: 30-7200This parameter is invalid for HLS. Only one HLS file will be generated from push start to push end.",
						},
						"storage_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Storage duration of the recording fileValue range: 0-129600000 seconds (0-1500 days)`0`: permanent.",
						},
						"enable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether to enable recording in the current format. Default value: 0. 0: no, 1: yes.",
						},
						"vod_sub_app_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "VOD subapplication ID.",
						},
						"vod_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Recording filename.Supported special placeholders include:{StreamID}: stream ID{StartYear}: start time - year{StartMonth}: start time - month{StartDay}: start time - day{StartHour}: start time - hour{StartMinute}: start time - minute{StartSecond}: start time - second{StartMillisecond}: start time - millisecond{EndYear}: end time - year{EndMonth}: end time - month{EndDay}: end time - day{EndHour}: end time - hour{EndMinute}: end time - minute{EndSecond}: end time - second{EndMillisecond}: end time - millisecondIf this parameter is not set, the recording filename will be `{StreamID}_{StartYear}-{StartMonth}-{StartDay}-{StartHour}-{StartMinute}-{StartSecond}_{EndYear}-{EndMonth}-{EndDay}-{EndHour}-{EndMinute}-{EndSecond}` by default.",
						},
						"procedure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task flowNote: this field may return `null`, indicating that no valid value is obtained.",
						},
						"storage_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video storage class. Valid values:`normal`: STANDARD`cold`: STANDARD_IANote: this field may return `null`, indicating that no valid value is obtained.",
						},
						"class_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "VOD subapplication categoryNote: this field may return `null`, indicating that no valid value is obtained.",
						},
					},
				},
			},

			"hls_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "HLS recording parameter, which is set when HLS recording is enabled.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Max recording time per fileDefault value: `1800` (seconds)Value range: 30-7200This parameter is invalid for HLS. Only one HLS file will be generated from push start to push end.",
						},
						"storage_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Storage duration of the recording fileValue range: 0-129600000 seconds (0-1500 days)`0`: permanent.",
						},
						"enable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether to enable recording in the current format. Default value: 0. 0: no, 1: yes.",
						},
						"vod_sub_app_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "VOD subapplication ID.",
						},
						"vod_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Recording filename.Supported special placeholders include:{StreamID}: stream ID{StartYear}: start time - year{StartMonth}: start time - month{StartDay}: start time - day{StartHour}: start time - hour{StartMinute}: start time - minute{StartSecond}: start time - second{StartMillisecond}: start time - millisecond{EndYear}: end time - year{EndMonth}: end time - month{EndDay}: end time - day{EndHour}: end time - hour{EndMinute}: end time - minute{EndSecond}: end time - second{EndMillisecond}: end time - millisecondIf this parameter is not set, the recording filename will be `{StreamID}_{StartYear}-{StartMonth}-{StartDay}-{StartHour}-{StartMinute}-{StartSecond}_{EndYear}-{EndMonth}-{EndDay}-{EndHour}-{EndMinute}-{EndSecond}` by default.",
						},
						"procedure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task flowNote: this field may return `null`, indicating that no valid value is obtained.",
						},
						"storage_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video storage class. Valid values:`normal`: STANDARD`cold`: STANDARD_IANote: this field may return `null`, indicating that no valid value is obtained.",
						},
						"class_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "VOD subapplication categoryNote: this field may return `null`, indicating that no valid value is obtained.",
						},
					},
				},
			},

			"mp4_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Mp4 recording parameter, which is set when Mp4 recording is enabled.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Max recording time per fileDefault value: `1800` (seconds)Value range: 30-7200This parameter is invalid for HLS. Only one HLS file will be generated from push start to push end.",
						},
						"storage_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Storage duration of the recording fileValue range: 0-129600000 seconds (0-1500 days)`0`: permanent.",
						},
						"enable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether to enable recording in the current format. Default value: 0. 0: no, 1: yes.",
						},
						"vod_sub_app_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "VOD subapplication ID.",
						},
						"vod_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Recording filename.Supported special placeholders include:{StreamID}: stream ID{StartYear}: start time - year{StartMonth}: start time - month{StartDay}: start time - day{StartHour}: start time - hour{StartMinute}: start time - minute{StartSecond}: start time - second{StartMillisecond}: start time - millisecond{EndYear}: end time - year{EndMonth}: end time - month{EndDay}: end time - day{EndHour}: end time - hour{EndMinute}: end time - minute{EndSecond}: end time - second{EndMillisecond}: end time - millisecondIf this parameter is not set, the recording filename will be `{StreamID}_{StartYear}-{StartMonth}-{StartDay}-{StartHour}-{StartMinute}-{StartSecond}_{EndYear}-{EndMonth}-{EndDay}-{EndHour}-{EndMinute}-{EndSecond}` by default.",
						},
						"procedure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task flowNote: this field may return `null`, indicating that no valid value is obtained.",
						},
						"storage_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video storage class. Valid values:`normal`: STANDARD`cold`: STANDARD_IANote: this field may return `null`, indicating that no valid value is obtained.",
						},
						"class_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "VOD subapplication categoryNote: this field may return `null`, indicating that no valid value is obtained.",
						},
					},
				},
			},

			"aac_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "AAC recording parameter, which is set when AAC recording is enabled.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Max recording time per fileDefault value: `1800` (seconds)Value range: 30-7200This parameter is invalid for HLS. Only one HLS file will be generated from push start to push end.",
						},
						"storage_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Storage duration of the recording fileValue range: 0-129600000 seconds (0-1500 days)`0`: permanent.",
						},
						"enable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether to enable recording in the current format. Default value: 0. 0: no, 1: yes.",
						},
						"vod_sub_app_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "VOD subapplication ID.",
						},
						"vod_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Recording filename.Supported special placeholders include:{StreamID}: stream ID{StartYear}: start time - year{StartMonth}: start time - month{StartDay}: start time - day{StartHour}: start time - hour{StartMinute}: start time - minute{StartSecond}: start time - second{StartMillisecond}: start time - millisecond{EndYear}: end time - year{EndMonth}: end time - month{EndDay}: end time - day{EndHour}: end time - hour{EndMinute}: end time - minute{EndSecond}: end time - second{EndMillisecond}: end time - millisecondIf this parameter is not set, the recording filename will be `{StreamID}_{StartYear}-{StartMonth}-{StartDay}-{StartHour}-{StartMinute}-{StartSecond}_{EndYear}-{EndMonth}-{EndDay}-{EndHour}-{EndMinute}-{EndSecond}` by default.",
						},
						"procedure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task flowNote: this field may return `null`, indicating that no valid value is obtained.",
						},
						"storage_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video storage class. Valid values:`normal`: STANDARD`cold`: STANDARD_IANote: this field may return `null`, indicating that no valid value is obtained.",
						},
						"class_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "VOD subapplication categoryNote: this field may return `null`, indicating that no valid value is obtained.",
						},
					},
				},
			},

			"is_delay_live": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "LVB type. Default value: 0.0: LVB.1: LCB.",
			},

			"hls_special_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "HLS-specific recording parameter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flow_continue_duration": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Timeout period for restarting an interrupted HLS push.Value range: [0, 1,800].",
						},
					},
				},
			},

			"mp3_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Mp3 recording parameter, which is set when Mp3 recording is enabled.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Max recording time per fileDefault value: `1800` (seconds)Value range: 30-7200This parameter is invalid for HLS. Only one HLS file will be generated from push start to push end.",
						},
						"storage_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Storage duration of the recording fileValue range: 0-129600000 seconds (0-1500 days)`0`: permanent.",
						},
						"enable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether to enable recording in the current format. Default value: 0. 0: no, 1: yes.",
						},
						"vod_sub_app_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "VOD subapplication ID.",
						},
						"vod_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Recording filename.Supported special placeholders include:{StreamID}: stream ID{StartYear}: start time - year{StartMonth}: start time - month{StartDay}: start time - day{StartHour}: start time - hour{StartMinute}: start time - minute{StartSecond}: start time - second{StartMillisecond}: start time - millisecond{EndYear}: end time - year{EndMonth}: end time - month{EndDay}: end time - day{EndHour}: end time - hour{EndMinute}: end time - minute{EndSecond}: end time - second{EndMillisecond}: end time - millisecondIf this parameter is not set, the recording filename will be `{StreamID}_{StartYear}-{StartMonth}-{StartDay}-{StartHour}-{StartMinute}-{StartSecond}_{EndYear}-{EndMonth}-{EndDay}-{EndHour}-{EndMinute}-{EndSecond}` by default.",
						},
						"procedure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task flowNote: this field may return `null`, indicating that no valid value is obtained.",
						},
						"storage_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video storage class. Valid values:`normal`: STANDARD`cold`: STANDARD_IANote: this field may return `null`, indicating that no valid value is obtained.",
						},
						"class_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "VOD subapplication categoryNote: this field may return `null`, indicating that no valid value is obtained.",
						},
					},
				},
			},

			"remove_watermark": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to remove the watermark. This parameter is invalid if `IsDelayLive` is `1`.",
			},

			"flv_special_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "A special parameter for FLV recording.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"upload_in_recording": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable upload while recording. This parameter is only valid for FLV recording.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudLiveRecordTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_record_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = live.NewCreateLiveRecordTemplateRequest()
		response   = live.NewCreateLiveRecordTemplateResponse()
		templateId int
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "flv_param"); ok {
		recordParam := live.RecordParam{}
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
		recordParam := live.RecordParam{}
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
		recordParam := live.RecordParam{}
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
		recordParam := live.RecordParam{}
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
		hlsSpecialParam := live.HlsSpecialParam{}
		if v, ok := dMap["flow_continue_duration"]; ok {
			hlsSpecialParam.FlowContinueDuration = helper.IntUint64(v.(int))
		}
		request.HlsSpecialParam = &hlsSpecialParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "mp3_param"); ok {
		recordParam := live.RecordParam{}
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
		flvSpecialParam := live.FlvSpecialParam{}
		if v, ok := dMap["upload_in_recording"]; ok {
			flvSpecialParam.UploadInRecording = helper.Bool(v.(bool))
		}
		request.FlvSpecialParam = &flvSpecialParam
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().CreateLiveRecordTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create live recordTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(templateId))

	return resourceTencentCloudLiveRecordTemplateRead(d, meta)
}

func resourceTencentCloudLiveRecordTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_record_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	recordTemplateId := d.Id()

	recordTemplate, err := service.DescribeLiveRecordTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if recordTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LiveRecordTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

func resourceTencentCloudLiveRecordTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_record_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := live.NewModifyLiveRecordTemplateRequest()

	recordTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"template_name", "description", "flv_param", "hls_param", "mp4_param", "aac_param", "is_delay_live", "hls_special_param", "mp3_param", "remove_watermark", "flv_special_param"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

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
			recordParam := live.RecordParam{}
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
			recordParam := live.RecordParam{}
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
			recordParam := live.RecordParam{}
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
			recordParam := live.RecordParam{}
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
			hlsSpecialParam := live.HlsSpecialParam{}
			if v, ok := dMap["flow_continue_duration"]; ok {
				hlsSpecialParam.FlowContinueDuration = helper.IntUint64(v.(int))
			}
			request.HlsSpecialParam = &hlsSpecialParam
		}
	}

	if d.HasChange("mp3_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "mp3_param"); ok {
			recordParam := live.RecordParam{}
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
			flvSpecialParam := live.FlvSpecialParam{}
			if v, ok := dMap["upload_in_recording"]; ok {
				flvSpecialParam.UploadInRecording = helper.Bool(v.(bool))
			}
			request.FlvSpecialParam = &flvSpecialParam
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().ModifyLiveRecordTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update live recordTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLiveRecordTemplateRead(d, meta)
}

func resourceTencentCloudLiveRecordTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_record_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}
	recordTemplateId := d.Id()

	if err := service.DeleteLiveRecordTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
