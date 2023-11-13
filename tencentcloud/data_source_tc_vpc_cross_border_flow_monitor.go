/*
Use this data source to query detailed information of vpc cross_border_flow_monitor

Example Usage

```hcl
data "tencentcloud_vpc_cross_border_flow_monitor" "cross_border_flow_monitor" {
  source_region = "ap-guangzhou"
  destination_region = "ap-singapore"
  ccn_id = "ccn-qd6z2ld1"
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcCrossBorderFlowMonitor() *schema.Resource {
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

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudVpcCrossBorderFlowMonitorRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_cross_border_flow_monitor.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("source_region"); ok {
		paramMap["SourceRegion"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("destination_region"); ok {
		paramMap["DestinationRegion"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ccn_id"); ok {
		paramMap["CcnId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ccn_uin"); ok {
		paramMap["CcnUin"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("period"); v != nil {
		paramMap["Period"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var crossBorderFlowMonitorData []*vpc.CrossBorderFlowMonitorData

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcCrossBorderFlowMonitorByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		crossBorderFlowMonitorData = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(crossBorderFlowMonitorData))
	tmpList := make([]map[string]interface{}, 0, len(crossBorderFlowMonitorData))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
