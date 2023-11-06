/*
Use this data source to query detailed information of monitor tmp_regions

Example Usage

```hcl
data "tencentcloud_monitor_tmp_regions" "tmp_regions" {
  pay_mode = 1
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMonitorTmpRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorTmpRegionsRead,
		Schema: map[string]*schema.Schema{
			"pay_mode": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Pay mode. `1`-Prepaid, `2`-Postpaid, `3`-All regions (default is all regions if not filled in).",
			},

			"region_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Region set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Region ID.",
						},
						"region_state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Region status (0-unavailable; 1-available).",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name.",
						},
						"region_short_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region short name.",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region area.",
						},
						"region_pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Region pay mode.",
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

func dataSourceTencentCloudMonitorTmpRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_monitor_tmp_regions.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("pay_mode"); ok {
		paramMap["PayMode"] = helper.IntInt64(v.(int))
	}

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	var regionSet []*monitor.PrometheusRegionItem
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorTmpRegionsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		regionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(regionSet))
	tmpList := make([]map[string]interface{}, 0, len(regionSet))

	if regionSet != nil {
		for _, prometheusRegionItem := range regionSet {
			prometheusRegionItemMap := map[string]interface{}{}

			if prometheusRegionItem.Region != nil {
				prometheusRegionItemMap["region"] = prometheusRegionItem.Region
			}

			if prometheusRegionItem.RegionId != nil {
				prometheusRegionItemMap["region_id"] = prometheusRegionItem.RegionId
			}

			if prometheusRegionItem.RegionState != nil {
				prometheusRegionItemMap["region_state"] = prometheusRegionItem.RegionState
			}

			if prometheusRegionItem.RegionName != nil {
				prometheusRegionItemMap["region_name"] = prometheusRegionItem.RegionName
			}

			if prometheusRegionItem.RegionShortName != nil {
				prometheusRegionItemMap["region_short_name"] = prometheusRegionItem.RegionShortName
			}

			if prometheusRegionItem.Area != nil {
				prometheusRegionItemMap["area"] = prometheusRegionItem.Area
			}

			if prometheusRegionItem.RegionPayMode != nil {
				prometheusRegionItemMap["region_pay_mode"] = prometheusRegionItem.RegionPayMode
			}

			ids = append(ids, *prometheusRegionItem.Region)
			tmpList = append(tmpList, prometheusRegionItemMap)
		}

		_ = d.Set("region_set", tmpList)
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
