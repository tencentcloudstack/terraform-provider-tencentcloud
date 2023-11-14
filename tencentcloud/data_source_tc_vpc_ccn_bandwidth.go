/*
Use this data source to query detailed information of vpc ccn_bandwidth

Example Usage

```hcl
data "tencentcloud_vpc_ccn_bandwidth" "ccn_bandwidth" {
  filters {
		name = "source-region"
		values =

  }
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

func dataSourceTencentCloudVpcCcn_bandwidth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcCcn_bandwidthRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter condition.Currently, only one value is supported. The supported fields.1)source-region,the value is like ap-guangzhou.2)destination-region, the value is like ap-shanghai.3)ccn-ids,cloud network ID array, the value is like ccn-12345678.4)user-account-id,user account ID, the value is like 12345678.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attribute name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Attribute value.",
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

func dataSourceTencentCloudVpcCcn_bandwidthRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_ccn_bandwidth.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*vpc.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := vpc.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var ccnBandwidthSet []*vpc.CcnBandwidth

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcCcn_bandwidthByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		ccnBandwidthSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(ccnBandwidthSet))
	tmpList := make([]map[string]interface{}, 0, len(ccnBandwidthSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
