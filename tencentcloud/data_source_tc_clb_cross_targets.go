/*
Use this data source to query detailed information of clb cross_targets

Example Usage

```hcl
data "tencentcloud_clb_cross_targets" "cross_targets" {
  filters {
    name = "vpc-id"
    values = ["vpc-4owdpnwr"]
  }
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClbCrossTargets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbCrossTargetsRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions to query CVMs and ENIs: vpc-id - String - Required: No - (Filter condition) Filter by VPC ID, such as vpc-12345678. ip - String - Required: No - (Filter condition) Filter by real server IP, such as 192.168.0.1. listener-id - String - Required: No - (Filter condition) Filter by listener ID, such as lbl-12345678. location-id - String - Required: No - (Filter condition) Filter by forwarding rule ID of the layer-7 listener, such as loc-12345678.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter values.",
						},
					},
				},
			},

			"cross_target_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cross target set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"local_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID of the CLB instance.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID of the CVM or ENI instance.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP address of the CVM or ENI instance.",
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC name of the CVM or ENI instance.",
						},
						"eni_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ENI ID of the CVM instance.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the CVM instance.Note: This field may return null, indicating that no valid value was found.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the CVM instance. Note: This field may return null, indicating that no valid value was found.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region of the CVM or ENI instance.",
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

func dataSourceTencentCloudClbCrossTargetsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clb_cross_targets.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*clb.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := clb.Filter{}
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
		paramMap["Filters"] = tmpSet
	}

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var crossTargetSet []*clb.CrossTargets

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbCrossTargetsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		crossTargetSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(crossTargetSet))
	tmpList := make([]map[string]interface{}, 0, len(crossTargetSet))

	if crossTargetSet != nil {
		for _, crossTargets := range crossTargetSet {
			crossTargetsMap := map[string]interface{}{}

			if crossTargets.LocalVpcId != nil {
				crossTargetsMap["local_vpc_id"] = crossTargets.LocalVpcId
			}

			if crossTargets.VpcId != nil {
				crossTargetsMap["vpc_id"] = crossTargets.VpcId
			}

			if crossTargets.IP != nil {
				crossTargetsMap["ip"] = crossTargets.IP
			}

			if crossTargets.VpcName != nil {
				crossTargetsMap["vpc_name"] = crossTargets.VpcName
			}

			if crossTargets.EniId != nil {
				crossTargetsMap["eni_id"] = crossTargets.EniId
			}

			if crossTargets.InstanceId != nil {
				crossTargetsMap["instance_id"] = crossTargets.InstanceId
			}

			if crossTargets.InstanceName != nil {
				crossTargetsMap["instance_name"] = crossTargets.InstanceName
			}

			if crossTargets.Region != nil {
				crossTargetsMap["region"] = crossTargets.Region
			}

			ids = append(ids, *crossTargets.VpcId)
			tmpList = append(tmpList, crossTargetsMap)
		}

		_ = d.Set("cross_target_set", tmpList)
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
