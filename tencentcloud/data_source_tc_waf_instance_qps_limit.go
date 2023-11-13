/*
Use this data source to query detailed information of waf instance_qps_limit

Example Usage

```hcl
data "tencentcloud_waf_instance_qps_limit" "instance_qps_limit" {
  instance_id = ""
  type = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudWafInstanceQpsLimit() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafInstanceQpsLimitRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Unique ID of Instance.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance type.",
			},

			"qps_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Qps info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"elastic_billing_default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Elastic qps default valueNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"elastic_billing_min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum elastic qpsNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"elastic_billing_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum elastic qpsNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"q_p_s_extend_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum qps of extend packageNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"q_p_s_extend_intl_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum qps of extend package for overseasNote: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudWafInstanceQpsLimitRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_waf_instance_qps_limit.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	service := WafService{client: meta.(*TencentCloudClient).apiV3Conn}

	var qpsData []*waf.QpsData

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafInstanceQpsLimitByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		qpsData = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(qpsData))
	if qpsData != nil {
		qpsDataMap := map[string]interface{}{}

		if qpsData.ElasticBillingDefault != nil {
			qpsDataMap["elastic_billing_default"] = qpsData.ElasticBillingDefault
		}

		if qpsData.ElasticBillingMin != nil {
			qpsDataMap["elastic_billing_min"] = qpsData.ElasticBillingMin
		}

		if qpsData.ElasticBillingMax != nil {
			qpsDataMap["elastic_billing_max"] = qpsData.ElasticBillingMax
		}

		if qpsData.QPSExtendMax != nil {
			qpsDataMap["q_p_s_extend_max"] = qpsData.QPSExtendMax
		}

		if qpsData.QPSExtendIntlMax != nil {
			qpsDataMap["q_p_s_extend_intl_max"] = qpsData.QPSExtendIntlMax
		}

		ids = append(ids, *qpsData.InstanceId)
		_ = d.Set("qps_data", qpsDataMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), qpsDataMap); e != nil {
			return e
		}
	}
	return nil
}
