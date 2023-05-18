/*
Use this data source to query detailed information of vpc bandwidth_package_bill_usage

Example Usage

```hcl
data "tencentcloud_vpc_bandwidth_package_bill_usage" "bandwidth_package_bill_usage" {
  bandwidth_package_id = "bwp-234rfgt5"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcBandwidthPackageBillUsage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcBandwidthPackageBillUsageRead,
		Schema: map[string]*schema.Schema{
			"bandwidth_package_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the postpaid bandwidth package.",
			},

			"bandwidth_package_bill_bandwidth_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "current billing amount.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth_usage": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Current billing amount in Mbps.",
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

func dataSourceTencentCloudVpcBandwidthPackageBillUsageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_bandwidth_package_bill_usage.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		paramMap["BandwidthPackageId"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var bandwidthPackageBillBandwidthSet []*vpc.BandwidthPackageBillBandwidth

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcBandwidthPackageBillUsageByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		bandwidthPackageBillBandwidthSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(bandwidthPackageBillBandwidthSet))
	tmpList := make([]map[string]interface{}, 0, len(bandwidthPackageBillBandwidthSet))

	if bandwidthPackageBillBandwidthSet != nil {
		for _, bandwidthPackageBillBandwidth := range bandwidthPackageBillBandwidthSet {
			bandwidthPackageBillBandwidthMap := map[string]interface{}{}

			if bandwidthPackageBillBandwidth.BandwidthUsage != nil {
				bandwidthPackageBillBandwidthMap["bandwidth_usage"] = bandwidthPackageBillBandwidth.BandwidthUsage
			}

			ids = append(ids, fmt.Sprintf("%f", *bandwidthPackageBillBandwidth.BandwidthUsage))
			tmpList = append(tmpList, bandwidthPackageBillBandwidthMap)
		}

		_ = d.Set("bandwidth_package_bill_bandwidth_set", tmpList)
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
