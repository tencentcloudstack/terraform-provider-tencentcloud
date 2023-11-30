package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCssTimeShiftStreamList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssTimeShiftStreamListRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The start time, which must be a Unix timestamp.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The end time, which must be a Unix timestamp.",
			},

			"stream_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The stream name.",
			},

			"domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The push domain.",
			},

			"domain_group": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The group the push domain belongs to.",
			},

			"total_size": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of records in the specified time period.",
			},

			"stream_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The information of the streams.Note: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The group the push domain belongs to.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The push domain.",
						},
						"app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The push path.",
						},
						"stream_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The stream name.",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The stream start time, which is a Unix timestamp.",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The stream end time (for streams that ended before the time of query), which is a Unix timestamp.",
						},
						"trans_code_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The transcoding template ID.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"stream_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The stream type. `0`: The original stream; `1`: The watermarked stream; `2`: The transcoded stream.",
						},
						"duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The storage duration (seconds) of the recording.Note: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudCssTimeShiftStreamListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_css_time_shift_stream_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("start_time"); ok {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("stream_name"); ok {
		paramMap["StreamName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain_group"); ok {
		paramMap["DomainGroup"] = helper.String(v.(string))
	}

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	var timeShiftStreamList []*css.TimeShiftStreamInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssTimeShiftStreamListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		timeShiftStreamList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(timeShiftStreamList))
	tmpList := make([]map[string]interface{}, 0, len(timeShiftStreamList))
	if timeShiftStreamList != nil {
		for _, timeShiftStreamInfo := range timeShiftStreamList {
			timeShiftStreamInfoMap := map[string]interface{}{}

			if timeShiftStreamInfo.DomainGroup != nil {
				timeShiftStreamInfoMap["domain_group"] = timeShiftStreamInfo.DomainGroup
			}

			if timeShiftStreamInfo.Domain != nil {
				timeShiftStreamInfoMap["domain"] = timeShiftStreamInfo.Domain
			}

			if timeShiftStreamInfo.AppName != nil {
				timeShiftStreamInfoMap["app_name"] = timeShiftStreamInfo.AppName
			}

			if timeShiftStreamInfo.StreamName != nil {
				timeShiftStreamInfoMap["stream_name"] = timeShiftStreamInfo.StreamName
			}

			if timeShiftStreamInfo.StartTime != nil {
				timeShiftStreamInfoMap["start_time"] = timeShiftStreamInfo.StartTime
			}

			if timeShiftStreamInfo.EndTime != nil {
				timeShiftStreamInfoMap["end_time"] = timeShiftStreamInfo.EndTime
			}

			if timeShiftStreamInfo.TransCodeId != nil {
				timeShiftStreamInfoMap["trans_code_id"] = timeShiftStreamInfo.TransCodeId
			}

			if timeShiftStreamInfo.StreamType != nil {
				timeShiftStreamInfoMap["stream_type"] = timeShiftStreamInfo.StreamType
			}

			if timeShiftStreamInfo.Duration != nil {
				timeShiftStreamInfoMap["duration"] = timeShiftStreamInfo.Duration
			}

			ids = append(ids, *timeShiftStreamInfo.Domain)
			tmpList = append(tmpList, timeShiftStreamInfoMap)
		}

		_ = d.Set("stream_list", tmpList)
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
