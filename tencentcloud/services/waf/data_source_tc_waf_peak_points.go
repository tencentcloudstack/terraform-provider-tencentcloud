package waf

import (
	"context"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWafPeakPoints() *schema.Resource {
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
				Description: "Only support sparta-waf and clb-waf. If not passed, there will be no filtering.",
			},
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "WAF instance ID, if not passed, there will be no filtering.",
			},
			"metric_name": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue(MetricNameList),
				Description:  "Thirteen values are available: access-Peak qps trend chart; botAccess- bot peak qps trend chart; down-Downstream peak bandwidth trend chart; up-Upstream peak bandwidth trend chart; attack-Trend chart of total number of web attacks; cc-Trend chart of total number of CC attacks; StatusServerError-Trend chart of the number of status codes returned by WAF to the server; StatusClientError-Trend chart of the number of status codes returned by WAF to the client; StatusRedirect-Trend chart of the number of status codes returned by WAF to the client; StatusOk-Trend chart of the number of status codes returned by WAF to the client; UpstreamServerError-Trend chart of the number of status codes returned to WAF by the origin site; UpstreamClientError-Trend chart of the number of status codes returned to WAF by the origin site; UpstreamRedirect-Trend chart of the number of status codes returned to WAF by the origin site.",
			},
			"points": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "point list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Second level timestamp.",
						},
						"access": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "qps.",
						},
						"bot_access": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bot qps.",
						},
						"attack": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of web attacks.",
						},
						"cc": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of cc attacks.",
						},
						"down": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Peak downlink bandwidth, unit B.",
						},
						"up": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Peak uplink bandwidth, unit B.",
						},
						"status_server_error": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart of the number of status codes returned by WAF to the server.",
						},
						"status_client_error": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart of the number of status codes returned by WAF to the client.",
						},
						"status_redirect": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart of the number of status codes returned by WAF to the client.",
						},
						"status_ok": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart of the number of status codes returned by WAF to the client.",
						},
						"upstream_server_error": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart of the number of status codes returned to WAF by the origin site.",
						},
						"upstream_client_error": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart of the number of status codes returned to WAF by the origin site.",
						},
						"upstream_redirect": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart of the number of status codes returned to WAF by the origin site.",
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
	defer tccommon.LogElapsed("data_source.tencentcloud_waf_peak_points.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		points  []*waf.PeakPointsItem
	)

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

	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceID"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("metric_name"); ok {
		paramMap["MetricName"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafPeakPointsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		points = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(points))

	if points != nil {
		for _, point := range points {
			dMap := map[string]interface{}{}

			if point.Time != nil {
				dMap["time"] = point.Time
			}

			if point.Access != nil {
				dMap["access"] = point.Access
			}

			if point.BotAccess != nil {
				dMap["bot_access"] = point.BotAccess
			}

			if point.Attack != nil {
				dMap["attack"] = point.Attack
			}

			if point.Cc != nil {
				dMap["cc"] = point.Cc
			}

			if point.Down != nil {
				dMap["down"] = point.Down
			}

			if point.Up != nil {
				dMap["up"] = point.Up
			}

			if point.StatusServerError != nil {
				dMap["status_server_error"] = point.StatusServerError
			}

			if point.StatusClientError != nil {
				dMap["status_client_error"] = point.StatusClientError
			}

			if point.StatusRedirect != nil {
				dMap["status_redirect"] = point.StatusRedirect
			}

			if point.StatusOk != nil {
				dMap["status_ok"] = point.StatusOk
			}

			if point.UpstreamServerError != nil {
				dMap["upstream_server_error"] = point.UpstreamServerError
			}

			if point.UpstreamClientError != nil {
				dMap["upstream_client_error"] = point.UpstreamClientError
			}

			if point.UpstreamRedirect != nil {
				dMap["upstream_redirect"] = point.UpstreamRedirect
			}

			tmpList = append(tmpList, dMap)
		}

		_ = d.Set("points", tmpList)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
