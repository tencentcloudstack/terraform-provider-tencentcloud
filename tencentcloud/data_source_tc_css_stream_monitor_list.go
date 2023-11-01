/*
Use this data source to query detailed information of css stream_monitor_list

Example Usage

```hcl
data "tencentcloud_css_stream_monitor_list" "stream_monitor_list" {
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCssStreamMonitorList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssStreamMonitorListRead,
		Schema: map[string]*schema.Schema{

			"live_stream_monitors": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The list of live stream monitoring tasks.Note: This field may return null, indicating that no valid value is available.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"monitor_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitoring task ID.Note: This field may return null, indicating that no valid value is available.",
						},
						"monitor_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitoring task name. Up to 128 bytes.Note: This field may return null, indicating that no valid value is available.",
						},
						"output_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Monitoring task output information.Note: This field may return null, indicating that no valid value is available.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"output_stream_width": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The width of the output stream in pixels for the monitoring task. The range is [1, 1920]. It is recommended to be at least 100 pixels.Note: This field may return null, indicating that no valid value is available.",
									},
									"output_stream_height": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The height of the output stream in pixels for the monitoring task. The range is [1, 1080]. It is recommended to be at least 100 pixels.Note: This field may return null, indicating that no valid value is available.",
									},
									"output_stream_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the output stream for the monitoring task.If not specified, the system will generate a name automatically.The name should be within 256 bytes and can only contain letters, numbers, `-`, `_`, and `.` characters.Note: This field may return null, indicating that no valid value is available.",
									},
									"output_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The playback domain for the monitoring task.It should be within 128 bytes and can only be filled with an enabled playback domain.Note: This field may return null, indicating that no valid value is available.",
									},
									"output_app": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The playback path for the monitoring task.It should be within 32 bytes and can only contain letters, numbers, `-`, `_`, and `.` characters.Note: This field may return null, indicating that no valid value is available.",
									},
								},
							},
						},
						"input_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The input stream information for the monitoring task.Note: This field may return null, indicating that no valid value is available.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"input_stream_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the input stream for the monitoring task.It should be within 256 bytes and can only contain letters, numbers, `-`, `_`, and `.` characters.Note: This field may return null, indicating that no valid value is available.",
									},
									"input_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The push domain for the input stream to be monitored.It should be within 128 bytes and can only be filled with an enabled push domain.Note: This field may return null, indicating that no valid value is available.",
									},
									"input_app": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The push path for the input stream to be monitored.It should be within 32 bytes and can only contain letters, numbers, `-`, `_`, and `.` characters.Note: This field may return null, indicating that no valid value is available.",
									},
									"input_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The push URL for the input stream to be monitored. In most cases, this parameter is not required.Note: This field may return null, indicating that no valid value is available.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the monitoring task.It should be within 256 bytes.Note: This field may return null, indicating that no valid value is available.",
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The status of the monitoring task.  0: Represents idle.  1: Represents monitoring in progress.Note: This field may return null, indicating that no valid value is available.",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The last start time of the monitoring task, in Unix timestamp format.Note: This field may return null, indicating that no valid value is available.",
						},
						"stop_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The last stop time of the monitoring task, in Unix timestamp format.Note: This field may return null, indicating that no valid value is available.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The creation time of the monitoring task, in Unix timestamp format.Note: This field may return null, indicating that no valid value is available.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The update time of the monitoring task, in Unix timestamp format.Note: This field may return null, indicating that no valid value is available.",
						},
						"notify_policy": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The notification policy for monitoring events.Note: This field may return null, indicating that no valid value is available.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notify_policy_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The type of notification policy: Range [0,1]  0: Represents no notification policy is used.  1: Represents the use of a global callback policy, where all events are notified to the CallbackUrl.Note: This field may return null, indicating that no valid value is available.",
									},
									"callback_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The callback URL for notifications. It should be of length [0,512] and only support URLs with the http and https types.Note: This field may return null, indicating that no valid value is available.",
									},
								},
							},
						},
						"audible_input_index_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "The list of input indices for the output audio.Note: This field may return null, indicating that no valid value is available.",
						},
						"ai_asr_input_index_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "The list of input indices for enabling intelligent speech recognition.Note: This field may return null, indicating that no valid value is available.",
						},
						"check_stream_broken": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to enable stream disconnection detection.Note: This field may return null, indicating that no valid value is available.",
						},
						"check_stream_low_frame_rate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to enable low frame rate detection.Note: This field may return null, indicating that no valid value is available.",
						},
						"asr_language": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The language for intelligent speech recognition:0: Disabled1: Chinese2: English3: Japanese4: KoreanNote: This field may return null, indicating that no valid value is available.",
						},
						"ocr_language": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The language for intelligent text recognition:0: Disabled1: Chinese and EnglishNote: This field may return null, indicating that no valid value is available.",
						},
						"ai_ocr_input_index_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "The list of input indices for enabling intelligent text recognition.Note: This field may return null, indicating that no valid value is available.",
						},
						"allow_monitor_report": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to store monitoring events in the monitoring report and allow querying of the monitoring report.Note: This field may return null, indicating that no valid value is available.",
						},
						"ai_format_diagnose": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to enable format diagnosis. Note: This field may return null, indicating that no valid value is available.",
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

func dataSourceTencentCloudCssStreamMonitorListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_css_stream_monitor_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	var streamMonitorList []*css.LiveStreamMonitorInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssStreamMonitorListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		streamMonitorList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(streamMonitorList))
	tmpList := make([]map[string]interface{}, 0, len(streamMonitorList))
	if streamMonitorList != nil {
		for _, liveStreamMonitorInfo := range streamMonitorList {
			liveStreamMonitorInfoMap := map[string]interface{}{}

			if liveStreamMonitorInfo.MonitorId != nil {
				liveStreamMonitorInfoMap["monitor_id"] = liveStreamMonitorInfo.MonitorId
			}

			if liveStreamMonitorInfo.MonitorName != nil {
				liveStreamMonitorInfoMap["monitor_name"] = liveStreamMonitorInfo.MonitorName
			}

			if liveStreamMonitorInfo.OutputInfo != nil {
				outputInfoMap := map[string]interface{}{}

				if liveStreamMonitorInfo.OutputInfo.OutputStreamWidth != nil {
					outputInfoMap["output_stream_width"] = liveStreamMonitorInfo.OutputInfo.OutputStreamWidth
				}

				if liveStreamMonitorInfo.OutputInfo.OutputStreamHeight != nil {
					outputInfoMap["output_stream_height"] = liveStreamMonitorInfo.OutputInfo.OutputStreamHeight
				}

				if liveStreamMonitorInfo.OutputInfo.OutputStreamName != nil {
					outputInfoMap["output_stream_name"] = liveStreamMonitorInfo.OutputInfo.OutputStreamName
				}

				if liveStreamMonitorInfo.OutputInfo.OutputDomain != nil {
					outputInfoMap["output_domain"] = liveStreamMonitorInfo.OutputInfo.OutputDomain
				}

				if liveStreamMonitorInfo.OutputInfo.OutputApp != nil {
					outputInfoMap["output_app"] = liveStreamMonitorInfo.OutputInfo.OutputApp
				}

				liveStreamMonitorInfoMap["output_info"] = []interface{}{outputInfoMap}
			}

			if liveStreamMonitorInfo.InputList != nil {
				inputListList := []interface{}{}
				for _, inputList := range liveStreamMonitorInfo.InputList {
					inputListMap := map[string]interface{}{}

					if inputList.InputStreamName != nil {
						inputListMap["input_stream_name"] = inputList.InputStreamName
					}

					if inputList.InputDomain != nil {
						inputListMap["input_domain"] = inputList.InputDomain
					}

					if inputList.InputApp != nil {
						inputListMap["input_app"] = inputList.InputApp
					}

					if inputList.InputUrl != nil {
						inputListMap["input_url"] = inputList.InputUrl
					}

					if inputList.Description != nil {
						inputListMap["description"] = inputList.Description
					}

					inputListList = append(inputListList, inputListMap)
				}

				liveStreamMonitorInfoMap["input_list"] = inputListList
			}

			if liveStreamMonitorInfo.Status != nil {
				liveStreamMonitorInfoMap["status"] = liveStreamMonitorInfo.Status
			}

			if liveStreamMonitorInfo.StartTime != nil {
				liveStreamMonitorInfoMap["start_time"] = liveStreamMonitorInfo.StartTime
			}

			if liveStreamMonitorInfo.StopTime != nil {
				liveStreamMonitorInfoMap["stop_time"] = liveStreamMonitorInfo.StopTime
			}

			if liveStreamMonitorInfo.CreateTime != nil {
				liveStreamMonitorInfoMap["create_time"] = liveStreamMonitorInfo.CreateTime
			}

			if liveStreamMonitorInfo.UpdateTime != nil {
				liveStreamMonitorInfoMap["update_time"] = liveStreamMonitorInfo.UpdateTime
			}

			if liveStreamMonitorInfo.NotifyPolicy != nil {
				notifyPolicyMap := map[string]interface{}{}

				if liveStreamMonitorInfo.NotifyPolicy.NotifyPolicyType != nil {
					notifyPolicyMap["notify_policy_type"] = liveStreamMonitorInfo.NotifyPolicy.NotifyPolicyType
				}

				if liveStreamMonitorInfo.NotifyPolicy.CallbackUrl != nil {
					notifyPolicyMap["callback_url"] = liveStreamMonitorInfo.NotifyPolicy.CallbackUrl
				}

				liveStreamMonitorInfoMap["notify_policy"] = []interface{}{notifyPolicyMap}
			}

			if liveStreamMonitorInfo.AudibleInputIndexList != nil {
				liveStreamMonitorInfoMap["audible_input_index_list"] = liveStreamMonitorInfo.AudibleInputIndexList
			}

			if liveStreamMonitorInfo.AiAsrInputIndexList != nil {
				liveStreamMonitorInfoMap["ai_asr_input_index_list"] = liveStreamMonitorInfo.AiAsrInputIndexList
			}

			if liveStreamMonitorInfo.CheckStreamBroken != nil {
				liveStreamMonitorInfoMap["check_stream_broken"] = liveStreamMonitorInfo.CheckStreamBroken
			}

			if liveStreamMonitorInfo.CheckStreamLowFrameRate != nil {
				liveStreamMonitorInfoMap["check_stream_low_frame_rate"] = liveStreamMonitorInfo.CheckStreamLowFrameRate
			}

			if liveStreamMonitorInfo.AsrLanguage != nil {
				liveStreamMonitorInfoMap["asr_language"] = liveStreamMonitorInfo.AsrLanguage
			}

			if liveStreamMonitorInfo.OcrLanguage != nil {
				liveStreamMonitorInfoMap["ocr_language"] = liveStreamMonitorInfo.OcrLanguage
			}

			if liveStreamMonitorInfo.AiOcrInputIndexList != nil {
				liveStreamMonitorInfoMap["ai_ocr_input_index_list"] = liveStreamMonitorInfo.AiOcrInputIndexList
			}

			if liveStreamMonitorInfo.AllowMonitorReport != nil {
				liveStreamMonitorInfoMap["allow_monitor_report"] = liveStreamMonitorInfo.AllowMonitorReport
			}

			if liveStreamMonitorInfo.AiFormatDiagnose != nil {
				liveStreamMonitorInfoMap["ai_format_diagnose"] = liveStreamMonitorInfo.AiFormatDiagnose
			}

			ids = append(ids, *liveStreamMonitorInfo.MonitorId)
			tmpList = append(tmpList, liveStreamMonitorInfoMap)
		}

		_ = d.Set("live_stream_monitors", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
