/*
Use this data source to query detailed information of live describe_monitor_report

Example Usage

```hcl
data "tencentcloud_live_describe_monitor_report" "describe_monitor_report" {
  monitor_id = ""
    }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLiveDescribeMonitorReport() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLiveDescribeMonitorReportRead,
		Schema: map[string]*schema.Schema{
			"monitor_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Monitor ID。.",
			},

			"m_p_s_result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The information about the media processing result.Note: This field may return null, indicating that no valid value was found.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ai_asr_results": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The result of intelligent speech recognition.Note: This field may return null, indicating that no valid value was found.",
						},
						"ai_ocr_results": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The result of intelligent text recognition.Note: This field may return null, indicating that no valid value was found.",
						},
					},
				},
			},

			"diagnose_result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The information about the media diagnostic result.Note: This field may return null, indicating that no valid value was found.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"stream_broken_results": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The information about the stream interruption.Note: This field may return null, indicating that no valid value was found.",
						},
						"low_frame_rate_results": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The information about low frame rate.Note: This field may return null, indicating that no valid value was found.",
						},
						"stream_format_results": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The information about the stream format diagnosis.Note: This field may return null, indicating that no valid value was found.",
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

func dataSourceTencentCloudLiveDescribeMonitorReportRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_live_describe_monitor_report.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("monitor_id"); ok {
		paramMap["MonitorId"] = helper.String(v.(string))
	}

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	var mPSResult []*live.MPSResult

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLiveDescribeMonitorReportByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		mPSResult = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(mPSResult))
	if mPSResult != nil {
		mPSResultMap := map[string]interface{}{}

		if mPSResult.AiAsrResults != nil {
			mPSResultMap["ai_asr_results"] = mPSResult.AiAsrResults
		}

		if mPSResult.AiOcrResults != nil {
			mPSResultMap["ai_ocr_results"] = mPSResult.AiOcrResults
		}

		ids = append(ids, *mPSResult.MonitorId)
		_ = d.Set("m_p_s_result", mPSResultMap)
	}

	if diagnoseResult != nil {
		diagnoseResultMap := map[string]interface{}{}

		if diagnoseResult.StreamBrokenResults != nil {
			diagnoseResultMap["stream_broken_results"] = diagnoseResult.StreamBrokenResults
		}

		if diagnoseResult.LowFrameRateResults != nil {
			diagnoseResultMap["low_frame_rate_results"] = diagnoseResult.LowFrameRateResults
		}

		if diagnoseResult.StreamFormatResults != nil {
			diagnoseResultMap["stream_format_results"] = diagnoseResult.StreamFormatResults
		}

		ids = append(ids, *diagnoseResult.MonitorId)
		_ = d.Set("diagnose_result", diagnoseResultMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), mPSResultMap); e != nil {
			return e
		}
	}
	return nil
}
