/*
Use this data source to query detailed information of as advices

Example Usage

```hcl
data "tencentcloud_as_advices" "advices" {
  auto_scaling_group_ids = ["asc-lo0b94oy"]
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAsAdvices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAsAdvicesRead,
		Schema: map[string]*schema.Schema{
			"auto_scaling_group_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of scaling groups to be queried. Upper limit: 100.",
			},

			"auto_scaling_advice_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "A collection of suggestions for scaling group configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_scaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Auto scaling group ID.",
						},
						"level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scaling group warning level. Valid values: NORMAL, WARNING, CRITICAL.",
						},
						"advices": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A collection of suggestions for scaling group configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"problem": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Problem Description.",
									},
									"detail": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Problem Details.",
									},
									"solution": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Recommended resolutions.",
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

func dataSourceTencentCloudAsAdvicesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_as_advices.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("auto_scaling_group_ids"); ok {
		autoScalingGroupIdsSet := v.(*schema.Set).List()
		paramMap["AutoScalingGroupIds"] = helper.InterfacesStringsPoint(autoScalingGroupIdsSet)
	}

	service := AsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var autoScalingAdviceSet []*as.AutoScalingAdvice

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAsAdvices(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		autoScalingAdviceSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(autoScalingAdviceSet))
	tmpList := make([]map[string]interface{}, 0, len(autoScalingAdviceSet))

	if autoScalingAdviceSet != nil {
		for _, autoScalingAdvice := range autoScalingAdviceSet {
			autoScalingAdviceMap := map[string]interface{}{}

			if autoScalingAdvice.AutoScalingGroupId != nil {
				autoScalingAdviceMap["auto_scaling_group_id"] = autoScalingAdvice.AutoScalingGroupId
			}

			if autoScalingAdvice.Level != nil {
				autoScalingAdviceMap["level"] = autoScalingAdvice.Level
			}

			if autoScalingAdvice.Advices != nil {
				advicesList := []interface{}{}
				for _, advices := range autoScalingAdvice.Advices {
					advicesMap := map[string]interface{}{}

					if advices.Problem != nil {
						advicesMap["problem"] = advices.Problem
					}

					if advices.Detail != nil {
						advicesMap["detail"] = advices.Detail
					}

					if advices.Solution != nil {
						advicesMap["solution"] = advices.Solution
					}

					advicesList = append(advicesList, advicesMap)
				}

				autoScalingAdviceMap["advices"] = []interface{}{advicesList}
			}

			ids = append(ids, *autoScalingAdvice.AutoScalingGroupId)
			tmpList = append(tmpList, autoScalingAdviceMap)
		}

		_ = d.Set("auto_scaling_advice_set", tmpList)
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
