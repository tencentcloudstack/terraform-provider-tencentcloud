/*
Use this data source to query detailed information of vpc subnet_resource_dashboard

Example Usage

```hcl
data "tencentcloud_vpc_subnet_resource_dashboard" "subnet_resource_dashboard" {
  subnet_ids = ["subnet-i9tpf6hq"]
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

func dataSourceTencentCloudVpcSubnetResourceDashboard() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcSubnetResourceDashboardRead,
		Schema: map[string]*schema.Schema{
			"subnet_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Subnet instance ID, such as `subnet-f1xjkw1b`.",
			},

			"resource_statistics_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information of resources returned.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC instance ID, such as vpc-f1xjkw1b.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet instance ID, such as `subnet-bthucmmy`.",
						},
						"ip": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of used IP addresses.",
						},
						"resource_statistics_item_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information of associated resources.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource type, such as CVM, ENI.",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource name.",
									},
									"resource_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of resources.",
									},
								},
							},
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

func dataSourceTencentCloudVpcSubnetResourceDashboardRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_subnet_resource_dashboard.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("subnet_ids"); ok {
		subnetIdsSet := v.(*schema.Set).List()
		paramMap["SubnetIds"] = helper.InterfacesStringsPoint(subnetIdsSet)
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var resourceStatisticsSet []*vpc.ResourceStatistics

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcSubnetResourceDashboardByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		resourceStatisticsSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(resourceStatisticsSet))
	tmpList := make([]map[string]interface{}, 0, len(resourceStatisticsSet))

	if resourceStatisticsSet != nil {
		for _, resourceStatistics := range resourceStatisticsSet {
			resourceStatisticsMap := map[string]interface{}{}

			if resourceStatistics.VpcId != nil {
				resourceStatisticsMap["vpc_id"] = resourceStatistics.VpcId
			}

			if resourceStatistics.SubnetId != nil {
				resourceStatisticsMap["subnet_id"] = resourceStatistics.SubnetId
			}

			if resourceStatistics.Ip != nil {
				resourceStatisticsMap["ip"] = resourceStatistics.Ip
			}

			if resourceStatistics.ResourceStatisticsItemSet != nil {
				resourceStatisticsItemSetList := []interface{}{}
				for _, resourceStatisticsItemSet := range resourceStatistics.ResourceStatisticsItemSet {
					resourceStatisticsItemSetMap := map[string]interface{}{}

					if resourceStatisticsItemSet.ResourceType != nil {
						resourceStatisticsItemSetMap["resource_type"] = resourceStatisticsItemSet.ResourceType
					}

					if resourceStatisticsItemSet.ResourceName != nil {
						resourceStatisticsItemSetMap["resource_name"] = resourceStatisticsItemSet.ResourceName
					}

					if resourceStatisticsItemSet.ResourceCount != nil {
						resourceStatisticsItemSetMap["resource_count"] = resourceStatisticsItemSet.ResourceCount
					}

					resourceStatisticsItemSetList = append(resourceStatisticsItemSetList, resourceStatisticsItemSetMap)
				}

				resourceStatisticsMap["resource_statistics_item_set"] = resourceStatisticsItemSetList
			}

			ids = append(ids, *resourceStatistics.SubnetId)
			tmpList = append(tmpList, resourceStatisticsMap)
		}

		_ = d.Set("resource_statistics_set", tmpList)
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
