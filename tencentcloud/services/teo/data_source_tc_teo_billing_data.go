package teo

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func DataSourceTencentCloudTeoBillingData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoBillingDataRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Start time.",
			},

			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "End time. The query time range (`EndTime` - `StartTime`) must be less than or equal to 31 days.",
			},

			"zone_ids": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    100,
				Description: "Site ID collection, up to 100 site IDs. Use `*` to query account-level data of all sites under the current TencentCloud main account.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"metric_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Billing metric name, such as `acc_flux`, `acc_bandwidth`, etc.",
			},

			"interval": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Query time granularity. Valid values: `5min`, `hour`, `day`.",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions. Each item contains `type` and `value`. Valid `type` values: `host`, `proxy-id`, `region-id`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter name. Valid values: `host`, `proxy-id`, `region-id`.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter value.",
						},
					},
				},
			},

			"group_by": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Grouping and aggregation dimensions, up to two dimensions. Valid values: `zone-id`, `host`, `proxy-id`, `region-id`.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Billing data point list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data timestamp.",
						},
						"value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Numeric value.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Site ID to which the data point belongs.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name to which the data point belongs.",
						},
						"proxy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Layer 4 proxy instance ID to which the data point belongs.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing region ID to which the data point belongs.",
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

func dataSourceTencentCloudTeoBillingDataRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_billing_data.read")()
	defer tccommon.InconsistentCheck(d, meta)

	logId := tccommon.GetLogId(nil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient()

	request := teo.NewDescribeBillingDataRequest()
	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}
	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}
	if v, ok := d.GetOk("zone_ids"); ok {
		zoneIdsSet := v.([]interface{})
		request.ZoneIds = helper.InterfacesStringsPoint(zoneIdsSet)
	}
	if v, ok := d.GetOk("metric_name"); ok {
		request.MetricName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("interval"); ok {
		request.Interval = helper.String(v.(string))
	}
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*teo.BillingDataFilter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filter := teo.BillingDataFilter{}
			filterMap := item.(map[string]interface{})
			if v, ok := filterMap["type"].(string); ok && v != "" {
				filter.Type = helper.String(v)
			}
			if v, ok := filterMap["value"].(string); ok && v != "" {
				filter.Value = helper.String(v)
			}
			tmpSet = append(tmpSet, &filter)
		}
		request.Filters = tmpSet
	}
	if v, ok := d.GetOk("group_by"); ok {
		groupBySet := v.([]interface{})
		request.GroupBy = helper.InterfacesStringsPoint(groupBySet)
	}

	ratelimit.Check(request.GetAction())

	var response *teo.DescribeBillingDataResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		resp, e := client.DescribeBillingDataWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = resp
		if response == nil || response.Response == nil {
			log.Printf("[DATASOURCE] read empty, skip SetId")
			return resource.NonRetryableError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read tencentcloud_teo_billing_data failed, reason:%s\n", logId, err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	dataList := make([]map[string]interface{}, 0)
	if response.Response.Data != nil {
		for _, item := range response.Response.Data {
			dataMap := map[string]interface{}{}
			if item.Time != nil {
				dataMap["time"] = item.Time
			}
			if item.Value != nil {
				dataMap["value"] = int(*item.Value)
			}
			if item.ZoneId != nil {
				dataMap["zone_id"] = item.ZoneId
			}
			if item.Host != nil {
				dataMap["host"] = item.Host
			}
			if item.ProxyId != nil {
				dataMap["proxy_id"] = item.ProxyId
			}
			if item.RegionId != nil {
				dataMap["region_id"] = item.RegionId
			}
			dataList = append(dataList, dataMap)
		}
	}
	_ = d.Set("data", dataList)

	metricName := d.Get("metric_name").(string)
	startTime := d.Get("start_time").(string)
	endTime := d.Get("end_time").(string)
	d.SetId(strings.Join([]string{metricName, startTime, endTime}, tccommon.FILED_SP))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataList); e != nil {
			return e
		}
	}

	return nil
}
