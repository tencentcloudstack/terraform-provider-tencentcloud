/*
Use this data source to query detailed information of antiddos pending risk info

Example Usage

```hcl
data "tencentcloud_antiddos_pending_risk_info" "pending_risk_info" {
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAntiddosPendingRiskInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAntiddosPendingRiskInfoRead,
		Schema: map[string]*schema.Schema{
			"is_paid_usr": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Is it a paid user? True: paid user, false: regular user.",
			},

			"attacking_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of resources in the attack.",
			},

			"blocking_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of resources in blockage.",
			},

			"expired_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of expired resources.",
			},

			"total": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total number of all pending risk events.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudAntiddosPendingRiskInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_antiddos_pending_risk_info.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}
	var pendingRiskInfoResponseParams *antiddos.DescribePendingRiskInfoResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAntiddosPendingRiskInfoByFilter(ctx)
		if e != nil {
			return retryError(e)
		}
		pendingRiskInfoResponseParams = result
		return nil
	})
	if err != nil {
		return err
	}
	resultMap := make(map[string]interface{})
	if pendingRiskInfoResponseParams.IsPaidUsr != nil {
		resultMap["is_paid_usr"] = pendingRiskInfoResponseParams.IsPaidUsr
		_ = d.Set("is_paid_usr", pendingRiskInfoResponseParams.IsPaidUsr)
	}

	if pendingRiskInfoResponseParams.AttackingCount != nil {
		resultMap["attacking_count"] = pendingRiskInfoResponseParams.AttackingCount
		_ = d.Set("attacking_count", pendingRiskInfoResponseParams.AttackingCount)
	}

	if pendingRiskInfoResponseParams.BlockingCount != nil {
		resultMap["blocking_count"] = pendingRiskInfoResponseParams.BlockingCount
		_ = d.Set("blocking_count", pendingRiskInfoResponseParams.BlockingCount)
	}

	if pendingRiskInfoResponseParams.ExpiredCount != nil {
		resultMap["expired_count"] = pendingRiskInfoResponseParams.ExpiredCount
		_ = d.Set("expired_count", pendingRiskInfoResponseParams.ExpiredCount)
	}

	if pendingRiskInfoResponseParams.Total != nil {
		resultMap["total"] = pendingRiskInfoResponseParams.Total
		_ = d.Set("total", pendingRiskInfoResponseParams.Total)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), resultMap); e != nil {
			return e
		}
	}
	return nil
}
