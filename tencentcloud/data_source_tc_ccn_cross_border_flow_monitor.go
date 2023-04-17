/*
Use this data source to query detailed information of vpc cross_border_flow_monitor

Example Usage

```hcl
data "tencentcloud_ccn_cross_border_flow_monitor" "cross_border_flow_monitor" {
  source_region = "ap-guangzhou"
  destination_region = "ap-singapore"
  ccn_id = "ccn-39lqkygf"
  ccn_uin = "979137"
  period = 60
  start_time = "2023-01-01 00:00:00"
  end_time = "2023-01-01 01:00:00"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCcnCrossBorderFlowMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcCrossBorderFlowMonitorRead,
		Schema: map[string]*schema.Schema{
			"source_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "SourceRegion.",
			},

			"destination_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "DestinationRegion.",
			},

			"ccn_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "CcnId.",
			},

			"ccn_uin": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "CcnUin.",
			},

			"period": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "TimePeriod.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "StartTime.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "EndTime.",
			},

			"cross_border_flow_monitor_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "monitor data of cross border.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"in_bandwidth": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "in bandwidth, `bps`.",
						},
						"out_bandwidth": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "out bandwidth, `bps`.",
						},
						"in_pkg": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "in pkg, `pps`.",
						},
						"out_pkg": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "out pkg, `pps`.",
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

func dataSourceTencentCloudVpcCrossBorderFlowMonitorRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ccn_cross_border_flow_monitor.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var ccnId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("source_region"); ok {
		paramMap["source_region"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("destination_region"); ok {
		paramMap["destination_region"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ccn_id"); ok {
		ccnId = v.(string)
		paramMap["ccn_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ccn_uin"); ok {
		paramMap["ccn_uin"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("period"); v != nil {
		paramMap["period"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["start_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["end_time"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var crossBorderFlowMonitorData []*vpc.CrossBorderFlowMonitorData

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCcnCrossBorderFlowMonitorByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		crossBorderFlowMonitorData = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(crossBorderFlowMonitorData))

	if crossBorderFlowMonitorData != nil {
		for _, crossBorderFlowMonitorData := range crossBorderFlowMonitorData {
			crossBorderFlowMonitorDataMap := map[string]interface{}{}

			if crossBorderFlowMonitorData.InBandwidth != nil {
				crossBorderFlowMonitorDataMap["in_bandwidth"] = crossBorderFlowMonitorData.InBandwidth
			}

			if crossBorderFlowMonitorData.OutBandwidth != nil {
				crossBorderFlowMonitorDataMap["out_bandwidth"] = crossBorderFlowMonitorData.OutBandwidth
			}

			if crossBorderFlowMonitorData.InPkg != nil {
				crossBorderFlowMonitorDataMap["in_pkg"] = crossBorderFlowMonitorData.InPkg
			}

			if crossBorderFlowMonitorData.OutPkg != nil {
				crossBorderFlowMonitorDataMap["out_pkg"] = crossBorderFlowMonitorData.OutPkg
			}

			tmpList = append(tmpList, crossBorderFlowMonitorDataMap)
		}

		_ = d.Set("cross_border_flow_monitor_data", tmpList)
	}

	d.SetId(ccnId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
