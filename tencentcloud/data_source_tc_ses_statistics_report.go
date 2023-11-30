package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSesStatisticsReport() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSesStatisticsReportRead,
		Schema: map[string]*schema.Schema{
			"start_date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start date.",
			},

			"end_date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End date.",
			},

			"domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sender domain.",
			},

			"receiving_mailbox_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Recipient address type, for example, gmail.com.",
			},

			"daily_volumes": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Daily email sending statistics.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"send_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"request_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of email requests.",
						},
						"accepted_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of email requests accepted by Tencent Cloud.",
						},
						"delivered_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of delivered emails.",
						},
						"opened_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of users (deduplicated) who opened emails.",
						},
						"clicked_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of recipients who clicked on links in emails.",
						},
						"bounce_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of bounced emails.",
						},
						"unsubscribe_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of users who canceled subscriptions. Note: this field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},

			"overall_volume": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Overall email sending statistics.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"send_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"request_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of email requests.",
						},
						"accepted_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of email requests accepted by Tencent Cloud.",
						},
						"delivered_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of delivered emails.",
						},
						"opened_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of users (deduplicated) who opened emails.",
						},
						"clicked_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of recipients who clicked on links in emails.",
						},
						"bounce_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of bounced emails.",
						},
						"unsubscribe_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of users who canceled subscriptions. Note: this field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudSesStatisticsReportRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ses_statistics_report.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("start_date"); ok {
		paramMap["StartDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_date"); ok {
		paramMap["EndDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("receiving_mailbox_type"); ok {
		paramMap["ReceivingMailboxType"] = helper.String(v.(string))
	}

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	var statisticsReport *ses.GetStatisticsReportResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSesStatisticsReportByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		statisticsReport = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(statisticsReport.DailyVolumes))
	tmpList := make([]map[string]interface{}, 0, len(statisticsReport.DailyVolumes))

	if statisticsReport.DailyVolumes != nil {
		for _, volume := range statisticsReport.DailyVolumes {
			volumeMap := map[string]interface{}{}

			if volume.SendDate != nil {
				volumeMap["send_date"] = volume.SendDate
			}

			if volume.RequestCount != nil {
				volumeMap["request_count"] = volume.RequestCount
			}

			if volume.AcceptedCount != nil {
				volumeMap["accepted_count"] = volume.AcceptedCount
			}

			if volume.DeliveredCount != nil {
				volumeMap["delivered_count"] = volume.DeliveredCount
			}

			if volume.OpenedCount != nil {
				volumeMap["opened_count"] = volume.OpenedCount
			}

			if volume.ClickedCount != nil {
				volumeMap["clicked_count"] = volume.ClickedCount
			}

			if volume.BounceCount != nil {
				volumeMap["bounce_count"] = volume.BounceCount
			}

			if volume.UnsubscribeCount != nil {
				volumeMap["unsubscribe_count"] = volume.UnsubscribeCount
			}

			ids = append(ids, *volume.SendDate)
			tmpList = append(tmpList, volumeMap)
		}

		_ = d.Set("daily_volumes", tmpList)
	}

	if statisticsReport.OverallVolume != nil {
		overallVolume := statisticsReport.OverallVolume
		volumeMap := map[string]interface{}{}

		if overallVolume.SendDate != nil {
			volumeMap["send_date"] = overallVolume.SendDate
		}

		if overallVolume.RequestCount != nil {
			volumeMap["request_count"] = overallVolume.RequestCount
		}

		if overallVolume.AcceptedCount != nil {
			volumeMap["accepted_count"] = overallVolume.AcceptedCount
		}

		if overallVolume.DeliveredCount != nil {
			volumeMap["delivered_count"] = overallVolume.DeliveredCount
		}

		if overallVolume.OpenedCount != nil {
			volumeMap["opened_count"] = overallVolume.OpenedCount
		}

		if overallVolume.ClickedCount != nil {
			volumeMap["clicked_count"] = overallVolume.ClickedCount
		}

		if overallVolume.BounceCount != nil {
			volumeMap["bounce_count"] = overallVolume.BounceCount
		}

		if overallVolume.UnsubscribeCount != nil {
			volumeMap["unsubscribe_count"] = overallVolume.UnsubscribeCount
		}

		_ = d.Set("overall_volume", []interface{}{volumeMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
