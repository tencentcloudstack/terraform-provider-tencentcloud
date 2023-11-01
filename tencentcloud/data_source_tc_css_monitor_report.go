/*
Use this data source to query detailed information of css monitor_report

Example Usage

```hcl
data "tencentcloud_css_monitor_report" "monitor_report" {
  monitor_id = "0e8a12b5-df2a-4a1b-aa98-97d5610aa142"
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

func dataSourceTencentCloudCssMonitorReport() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssMonitorReportRead,
		Schema: map[string]*schema.Schema{
			"monitor_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Monitor ID.",
			},

			"mps_result": {
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

func dataSourceTencentCloudCssMonitorReportRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_css_monitor_report.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var monitorId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("monitor_id"); ok {
		monitorId = v.(string)
		paramMap["MonitorId"] = helper.String(v.(string))
	}

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	var mPSResult *css.DescribeMonitorReportResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssMonitorReportByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		mPSResult = result
		return nil
	})
	if err != nil {
		return err
	}

	if mPSResult.MPSResult != nil {
		mPSResultMap := map[string]interface{}{}

		if mPSResult.MPSResult.AiAsrResults != nil {
			mPSResultMap["ai_asr_results"] = mPSResult.MPSResult.AiAsrResults
		}

		if mPSResult.MPSResult.AiOcrResults != nil {
			mPSResultMap["ai_ocr_results"] = mPSResult.MPSResult.AiOcrResults
		}

		_ = d.Set("mps_result", []interface{}{mPSResultMap})
	}

	if mPSResult.DiagnoseResult != nil {
		diagnoseResultMap := map[string]interface{}{}

		if mPSResult.DiagnoseResult.StreamBrokenResults != nil {
			diagnoseResultMap["stream_broken_results"] = mPSResult.DiagnoseResult.StreamBrokenResults
		}

		if mPSResult.DiagnoseResult.LowFrameRateResults != nil {
			diagnoseResultMap["low_frame_rate_results"] = mPSResult.DiagnoseResult.LowFrameRateResults
		}

		if mPSResult.DiagnoseResult.StreamFormatResults != nil {
			diagnoseResultMap["stream_format_results"] = mPSResult.DiagnoseResult.StreamFormatResults
		}

		_ = d.Set("diagnose_result", []interface{}{diagnoseResultMap})
	}

	d.SetId(helper.DataResourceIdsHash([]string{monitorId}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
