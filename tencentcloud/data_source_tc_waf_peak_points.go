/*
Use this data source to query detailed information of waf peak_points

Example Usage

```hcl
data "tencentcloud_waf_peak_points" "peak_points" {
  from_time = ""
  to_time = ""
  domain = ""
  edition = ""
  instance_i_d = ""
  metric_name = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudWafPeakPoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafPeakPointsRead,
		Schema: map[string]*schema.Schema{
			"from_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Begin time.",
			},

			"to_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},

			"domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The domain name to be queried. If all domain name data is queried, this parameter is not filled in.",
			},

			"edition": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Only two values ​​are valid, sparta-waf and clb-waf. If not passed, there will be no filtering.",
			},

			"instance_i_d": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "WAF instance ID, if not passed, there will be no filtering.",
			},

			"metric_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Seven values ​​are available:access-Peak qps trend chartbotAccess- bot peak qps trend chartdown-Downstream peak bandwidth trend chartup-Upstream peak bandwidth trend chartattack-Trend chart of total number of web attackscc-Trend chart of total number of CC attackshttp_status-Trend chart of each status code number.",
			},

			"points": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Data point.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timestamp.",
						},
						"access": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "QPS.",
						},
						"up": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Peak uplink bandwidth, unit B.",
						},
						"down": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Peak downlink bandwidth, unit B.",
						},
						"attack": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Web attack count.",
						},
						"cc": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CC attack count.",
						},
						"bot_access": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bot qps.",
						},
						"status_server_error": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart of the number of status codes 5xx returned by WAF to the clientNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"status_client_error": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart of the number of status codes 4xx returned by WAF to the clientNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"status_redirect": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart of the number of status codes 302 returned by WAF to the clientNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"status_ok": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart of the number of status codes 202 returned by WAF to the clientNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"upstream_server_error": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart of the number of status codes 5xx returned to WAF by the origin siteNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"upstream_client_error": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart of the number of status codes 4xx returned to WAF by the origin siteNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"upstream_redirect": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart of the number of status codes 302 returned to WAF by the origin siteNote: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudWafPeakPointsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_waf_peak_points.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("from_time"); ok {
		paramMap["FromTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("to_time"); ok {
		paramMap["ToTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("edition"); ok {
		paramMap["Edition"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_i_d"); ok {
		paramMap["InstanceID"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("metric_name"); ok {
		paramMap["MetricName"] = helper.String(v.(string))
	}

	service := WafService{client: meta.(*TencentCloudClient).apiV3Conn}

	var points []*waf.PeakPointsItem

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafPeakPointsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		points = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(points))
	tmpList := make([]map[string]interface{}, 0, len(points))

	if points != nil {
		for _, peakPointsItem := range points {
			peakPointsItemMap := map[string]interface{}{}

			if peakPointsItem.Time != nil {
				peakPointsItemMap["time"] = peakPointsItem.Time
			}

			if peakPointsItem.Access != nil {
				peakPointsItemMap["access"] = peakPointsItem.Access
			}

			if peakPointsItem.Up != nil {
				peakPointsItemMap["up"] = peakPointsItem.Up
			}

			if peakPointsItem.Down != nil {
				peakPointsItemMap["down"] = peakPointsItem.Down
			}

			if peakPointsItem.Attack != nil {
				peakPointsItemMap["attack"] = peakPointsItem.Attack
			}

			if peakPointsItem.Cc != nil {
				peakPointsItemMap["cc"] = peakPointsItem.Cc
			}

			if peakPointsItem.BotAccess != nil {
				peakPointsItemMap["bot_access"] = peakPointsItem.BotAccess
			}

			if peakPointsItem.StatusServerError != nil {
				peakPointsItemMap["status_server_error"] = peakPointsItem.StatusServerError
			}

			if peakPointsItem.StatusClientError != nil {
				peakPointsItemMap["status_client_error"] = peakPointsItem.StatusClientError
			}

			if peakPointsItem.StatusRedirect != nil {
				peakPointsItemMap["status_redirect"] = peakPointsItem.StatusRedirect
			}

			if peakPointsItem.StatusOk != nil {
				peakPointsItemMap["status_ok"] = peakPointsItem.StatusOk
			}

			if peakPointsItem.UpstreamServerError != nil {
				peakPointsItemMap["upstream_server_error"] = peakPointsItem.UpstreamServerError
			}

			if peakPointsItem.UpstreamClientError != nil {
				peakPointsItemMap["upstream_client_error"] = peakPointsItem.UpstreamClientError
			}

			if peakPointsItem.UpstreamRedirect != nil {
				peakPointsItemMap["upstream_redirect"] = peakPointsItem.UpstreamRedirect
			}

			ids = append(ids, *peakPointsItem.RequestId)
			tmpList = append(tmpList, peakPointsItemMap)
		}

		_ = d.Set("points", tmpList)
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
