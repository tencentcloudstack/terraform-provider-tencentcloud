/*
Provides a resource to create a live stream_monitor

Example Usage

```hcl
resource "tencentcloud_live_stream_monitor" "stream_monitor" {
  output_info {
		output_stream_width =
		output_stream_height =
		output_stream_name = ""
		output_domain = ""
		output_app = ""

  }
  input_list {
		input_stream_name = ""
		input_domain = ""
		input_app = ""
		input_url = ""
		description = ""

  }
  monitor_name = ""
  notify_policy {
		notify_policy_type =
		callback_url = ""

  }
  asr_language =
  ocr_language =
  ai_asr_input_index_list =
  ai_ocr_input_index_list =
  check_stream_broken =
  check_stream_low_frame_rate =
  allow_monitor_report =
  ai_format_diagnose =
}
```

Import

live stream_monitor can be imported using the id, e.g.

```
terraform import tencentcloud_live_stream_monitor.stream_monitor stream_monitor_id
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

func resourceTencentCloudLiveStreamMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveStreamMonitorCreate,
		Read:   resourceTencentCloudLiveStreamMonitorRead,
		Update: resourceTencentCloudLiveStreamMonitorUpdate,
		Delete: resourceTencentCloudLiveStreamMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"output_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Monitor task output info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"output_stream_width": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Output stream width, limit[1, 1920].",
						},
						"output_stream_height": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Monitor task output height, limit[1, 1080].",
						},
						"output_stream_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Monitor task output stream name.limit 256 bytes.",
						},
						"output_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Monitor task output play domain.limit 128 bytes.",
						},
						"output_app": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Monitor task play path.limit 32 bytes.",
						},
					},
				},
			},

			"input_list": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Wait monitor input info list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"input_stream_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Wait monitor input stream name.limit 256 bytes.",
						},
						"input_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Wait monitor input push domain.limit 128 bytes.",
						},
						"input_app": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Wait monitor input push path.limit 32 bytes.",
						},
						"input_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Wait monitor input stream push url.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description content.limit 256 bytes.",
						},
					},
				},
			},

			"monitor_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Monitor task name.",
			},

			"notify_policy": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Monitor event notify policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notify_policy_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Notify policy type.0: not notify.1: use global policy.",
						},
						"callback_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Callback url.limit [0,512].only http or https.",
						},
					},
				},
			},

			"asr_language": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Asr language.0: close.1: Chinese 2: English3: Japanese4: Korean.",
			},

			"ocr_language": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Ocr language.0: close.1. Chinese,English.",
			},

			"ai_asr_input_index_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "AI asr input index list.(first input index is 1.).",
			},

			"ai_ocr_input_index_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Ai ocr input index list(first input index is 1.).",
			},

			"check_stream_broken": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "If enable stream broken check.",
			},

			"check_stream_low_frame_rate": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "If enable low frame rate check.",
			},

			"allow_monitor_report": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "If store monitor event.",
			},

			"ai_format_diagnose": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "If enable format diagnose.",
			},
		},
	}
}

func resourceTencentCloudLiveStreamMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_stream_monitor.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = live.NewCreateLiveStreamMonitorRequest()
		response  = live.NewCreateLiveStreamMonitorResponse()
		monitorId string
	)
	if dMap, ok := helper.InterfacesHeadMap(d, "output_info"); ok {
		liveStreamMonitorOutputInfo := live.LiveStreamMonitorOutputInfo{}
		if v, ok := dMap["output_stream_width"]; ok {
			liveStreamMonitorOutputInfo.OutputStreamWidth = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["output_stream_height"]; ok {
			liveStreamMonitorOutputInfo.OutputStreamHeight = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["output_stream_name"]; ok {
			liveStreamMonitorOutputInfo.OutputStreamName = helper.String(v.(string))
		}
		if v, ok := dMap["output_domain"]; ok {
			liveStreamMonitorOutputInfo.OutputDomain = helper.String(v.(string))
		}
		if v, ok := dMap["output_app"]; ok {
			liveStreamMonitorOutputInfo.OutputApp = helper.String(v.(string))
		}
		request.OutputInfo = &liveStreamMonitorOutputInfo
	}

	if v, ok := d.GetOk("input_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			liveStreamMonitorInputInfo := live.LiveStreamMonitorInputInfo{}
			if v, ok := dMap["input_stream_name"]; ok {
				liveStreamMonitorInputInfo.InputStreamName = helper.String(v.(string))
			}
			if v, ok := dMap["input_domain"]; ok {
				liveStreamMonitorInputInfo.InputDomain = helper.String(v.(string))
			}
			if v, ok := dMap["input_app"]; ok {
				liveStreamMonitorInputInfo.InputApp = helper.String(v.(string))
			}
			if v, ok := dMap["input_url"]; ok {
				liveStreamMonitorInputInfo.InputUrl = helper.String(v.(string))
			}
			if v, ok := dMap["description"]; ok {
				liveStreamMonitorInputInfo.Description = helper.String(v.(string))
			}
			request.InputList = append(request.InputList, &liveStreamMonitorInputInfo)
		}
	}

	if v, ok := d.GetOk("monitor_name"); ok {
		request.MonitorName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "notify_policy"); ok {
		liveStreamMonitorNotifyPolicy := live.LiveStreamMonitorNotifyPolicy{}
		if v, ok := dMap["notify_policy_type"]; ok {
			liveStreamMonitorNotifyPolicy.NotifyPolicyType = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["callback_url"]; ok {
			liveStreamMonitorNotifyPolicy.CallbackUrl = helper.String(v.(string))
		}
		request.NotifyPolicy = &liveStreamMonitorNotifyPolicy
	}

	if v, ok := d.GetOkExists("asr_language"); ok {
		request.AsrLanguage = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("ocr_language"); ok {
		request.OcrLanguage = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("ai_asr_input_index_list"); ok {
		aiAsrInputIndexListSet := v.(*schema.Set).List()
		for i := range aiAsrInputIndexListSet {
			aiAsrInputIndexList := aiAsrInputIndexListSet[i].(int)
			request.AiAsrInputIndexList = append(request.AiAsrInputIndexList, helper.IntUint64(aiAsrInputIndexList))
		}
	}

	if v, ok := d.GetOk("ai_ocr_input_index_list"); ok {
		aiOcrInputIndexListSet := v.(*schema.Set).List()
		for i := range aiOcrInputIndexListSet {
			aiOcrInputIndexList := aiOcrInputIndexListSet[i].(int)
			request.AiOcrInputIndexList = append(request.AiOcrInputIndexList, helper.IntUint64(aiOcrInputIndexList))
		}
	}

	if v, ok := d.GetOkExists("check_stream_broken"); ok {
		request.CheckStreamBroken = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("check_stream_low_frame_rate"); ok {
		request.CheckStreamLowFrameRate = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("allow_monitor_report"); ok {
		request.AllowMonitorReport = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("ai_format_diagnose"); ok {
		request.AiFormatDiagnose = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().CreateLiveStreamMonitor(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create live streamMonitor failed, reason:%+v", logId, err)
		return err
	}

	monitorId = *response.Response.MonitorId
	d.SetId(monitorId)

	return resourceTencentCloudLiveStreamMonitorRead(d, meta)
}

func resourceTencentCloudLiveStreamMonitorRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_stream_monitor.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	streamMonitorId := d.Id()

	streamMonitor, err := service.DescribeLiveStreamMonitorById(ctx, monitorId)
	if err != nil {
		return err
	}

	if streamMonitor == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LiveStreamMonitor` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if streamMonitor.OutputInfo != nil {
		outputInfoMap := map[string]interface{}{}

		if streamMonitor.OutputInfo.OutputStreamWidth != nil {
			outputInfoMap["output_stream_width"] = streamMonitor.OutputInfo.OutputStreamWidth
		}

		if streamMonitor.OutputInfo.OutputStreamHeight != nil {
			outputInfoMap["output_stream_height"] = streamMonitor.OutputInfo.OutputStreamHeight
		}

		if streamMonitor.OutputInfo.OutputStreamName != nil {
			outputInfoMap["output_stream_name"] = streamMonitor.OutputInfo.OutputStreamName
		}

		if streamMonitor.OutputInfo.OutputDomain != nil {
			outputInfoMap["output_domain"] = streamMonitor.OutputInfo.OutputDomain
		}

		if streamMonitor.OutputInfo.OutputApp != nil {
			outputInfoMap["output_app"] = streamMonitor.OutputInfo.OutputApp
		}

		_ = d.Set("output_info", []interface{}{outputInfoMap})
	}

	if streamMonitor.InputList != nil {
		inputListList := []interface{}{}
		for _, inputList := range streamMonitor.InputList {
			inputListMap := map[string]interface{}{}

			if streamMonitor.InputList.InputStreamName != nil {
				inputListMap["input_stream_name"] = streamMonitor.InputList.InputStreamName
			}

			if streamMonitor.InputList.InputDomain != nil {
				inputListMap["input_domain"] = streamMonitor.InputList.InputDomain
			}

			if streamMonitor.InputList.InputApp != nil {
				inputListMap["input_app"] = streamMonitor.InputList.InputApp
			}

			if streamMonitor.InputList.InputUrl != nil {
				inputListMap["input_url"] = streamMonitor.InputList.InputUrl
			}

			if streamMonitor.InputList.Description != nil {
				inputListMap["description"] = streamMonitor.InputList.Description
			}

			inputListList = append(inputListList, inputListMap)
		}

		_ = d.Set("input_list", inputListList)

	}

	if streamMonitor.MonitorName != nil {
		_ = d.Set("monitor_name", streamMonitor.MonitorName)
	}

	if streamMonitor.NotifyPolicy != nil {
		notifyPolicyMap := map[string]interface{}{}

		if streamMonitor.NotifyPolicy.NotifyPolicyType != nil {
			notifyPolicyMap["notify_policy_type"] = streamMonitor.NotifyPolicy.NotifyPolicyType
		}

		if streamMonitor.NotifyPolicy.CallbackUrl != nil {
			notifyPolicyMap["callback_url"] = streamMonitor.NotifyPolicy.CallbackUrl
		}

		_ = d.Set("notify_policy", []interface{}{notifyPolicyMap})
	}

	if streamMonitor.AsrLanguage != nil {
		_ = d.Set("asr_language", streamMonitor.AsrLanguage)
	}

	if streamMonitor.OcrLanguage != nil {
		_ = d.Set("ocr_language", streamMonitor.OcrLanguage)
	}

	if streamMonitor.AiAsrInputIndexList != nil {
		_ = d.Set("ai_asr_input_index_list", streamMonitor.AiAsrInputIndexList)
	}

	if streamMonitor.AiOcrInputIndexList != nil {
		_ = d.Set("ai_ocr_input_index_list", streamMonitor.AiOcrInputIndexList)
	}

	if streamMonitor.CheckStreamBroken != nil {
		_ = d.Set("check_stream_broken", streamMonitor.CheckStreamBroken)
	}

	if streamMonitor.CheckStreamLowFrameRate != nil {
		_ = d.Set("check_stream_low_frame_rate", streamMonitor.CheckStreamLowFrameRate)
	}

	if streamMonitor.AllowMonitorReport != nil {
		_ = d.Set("allow_monitor_report", streamMonitor.AllowMonitorReport)
	}

	if streamMonitor.AiFormatDiagnose != nil {
		_ = d.Set("ai_format_diagnose", streamMonitor.AiFormatDiagnose)
	}

	return nil
}

func resourceTencentCloudLiveStreamMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_stream_monitor.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := live.NewModifyLiveStreamMonitorRequest()

	streamMonitorId := d.Id()

	request.MonitorId = &monitorId

	immutableArgs := []string{"output_info", "input_list", "monitor_name", "notify_policy", "asr_language", "ocr_language", "ai_asr_input_index_list", "ai_ocr_input_index_list", "check_stream_broken", "check_stream_low_frame_rate", "allow_monitor_report", "ai_format_diagnose"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("output_info") {
		if dMap, ok := helper.InterfacesHeadMap(d, "output_info"); ok {
			liveStreamMonitorOutputInfo := live.LiveStreamMonitorOutputInfo{}
			if v, ok := dMap["output_stream_width"]; ok {
				liveStreamMonitorOutputInfo.OutputStreamWidth = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["output_stream_height"]; ok {
				liveStreamMonitorOutputInfo.OutputStreamHeight = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["output_stream_name"]; ok {
				liveStreamMonitorOutputInfo.OutputStreamName = helper.String(v.(string))
			}
			if v, ok := dMap["output_domain"]; ok {
				liveStreamMonitorOutputInfo.OutputDomain = helper.String(v.(string))
			}
			if v, ok := dMap["output_app"]; ok {
				liveStreamMonitorOutputInfo.OutputApp = helper.String(v.(string))
			}
			request.OutputInfo = &liveStreamMonitorOutputInfo
		}
	}

	if d.HasChange("input_list") {
		if v, ok := d.GetOk("input_list"); ok {
			for _, item := range v.([]interface{}) {
				liveStreamMonitorInputInfo := live.LiveStreamMonitorInputInfo{}
				if v, ok := dMap["input_stream_name"]; ok {
					liveStreamMonitorInputInfo.InputStreamName = helper.String(v.(string))
				}
				if v, ok := dMap["input_domain"]; ok {
					liveStreamMonitorInputInfo.InputDomain = helper.String(v.(string))
				}
				if v, ok := dMap["input_app"]; ok {
					liveStreamMonitorInputInfo.InputApp = helper.String(v.(string))
				}
				if v, ok := dMap["input_url"]; ok {
					liveStreamMonitorInputInfo.InputUrl = helper.String(v.(string))
				}
				if v, ok := dMap["description"]; ok {
					liveStreamMonitorInputInfo.Description = helper.String(v.(string))
				}
				request.InputList = append(request.InputList, &liveStreamMonitorInputInfo)
			}
		}
	}

	if d.HasChange("monitor_name") {
		if v, ok := d.GetOk("monitor_name"); ok {
			request.MonitorName = helper.String(v.(string))
		}
	}

	if d.HasChange("notify_policy") {
		if dMap, ok := helper.InterfacesHeadMap(d, "notify_policy"); ok {
			liveStreamMonitorNotifyPolicy := live.LiveStreamMonitorNotifyPolicy{}
			if v, ok := dMap["notify_policy_type"]; ok {
				liveStreamMonitorNotifyPolicy.NotifyPolicyType = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["callback_url"]; ok {
				liveStreamMonitorNotifyPolicy.CallbackUrl = helper.String(v.(string))
			}
			request.NotifyPolicy = &liveStreamMonitorNotifyPolicy
		}
	}

	if d.HasChange("asr_language") {
		if v, ok := d.GetOkExists("asr_language"); ok {
			request.AsrLanguage = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("ocr_language") {
		if v, ok := d.GetOkExists("ocr_language"); ok {
			request.OcrLanguage = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("ai_asr_input_index_list") {
		if v, ok := d.GetOk("ai_asr_input_index_list"); ok {
			aiAsrInputIndexListSet := v.(*schema.Set).List()
			for i := range aiAsrInputIndexListSet {
				aiAsrInputIndexList := aiAsrInputIndexListSet[i].(int)
				request.AiAsrInputIndexList = append(request.AiAsrInputIndexList, helper.IntUint64(aiAsrInputIndexList))
			}
		}
	}

	if d.HasChange("ai_ocr_input_index_list") {
		if v, ok := d.GetOk("ai_ocr_input_index_list"); ok {
			aiOcrInputIndexListSet := v.(*schema.Set).List()
			for i := range aiOcrInputIndexListSet {
				aiOcrInputIndexList := aiOcrInputIndexListSet[i].(int)
				request.AiOcrInputIndexList = append(request.AiOcrInputIndexList, helper.IntUint64(aiOcrInputIndexList))
			}
		}
	}

	if d.HasChange("check_stream_broken") {
		if v, ok := d.GetOkExists("check_stream_broken"); ok {
			request.CheckStreamBroken = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("check_stream_low_frame_rate") {
		if v, ok := d.GetOkExists("check_stream_low_frame_rate"); ok {
			request.CheckStreamLowFrameRate = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("allow_monitor_report") {
		if v, ok := d.GetOkExists("allow_monitor_report"); ok {
			request.AllowMonitorReport = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("ai_format_diagnose") {
		if v, ok := d.GetOkExists("ai_format_diagnose"); ok {
			request.AiFormatDiagnose = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().ModifyLiveStreamMonitor(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update live streamMonitor failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLiveStreamMonitorRead(d, meta)
}

func resourceTencentCloudLiveStreamMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_stream_monitor.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}
	streamMonitorId := d.Id()

	if err := service.DeleteLiveStreamMonitorById(ctx, monitorId); err != nil {
		return err
	}

	return nil
}
