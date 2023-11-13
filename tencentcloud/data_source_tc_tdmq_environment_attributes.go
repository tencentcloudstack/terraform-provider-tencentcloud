/*
Use this data source to query detailed information of tdmq environment_attributes

Example Usage

```hcl
data "tencentcloud_tdmq_environment_attributes" "environment_attributes" {
    cluster_id = ""
              }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTdmqEnvironmentAttributes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqEnvironmentAttributesRead,
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Environment (namespace) name.",
			},

			"cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID of the Pulsar cluster.",
			},

			"msg_t_t_l": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Expiration time of unconsumed messages, unit second, maximum 1296000 (15 days).",
			},

			"rate_in_byte": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Consumption rate limit, unit byte/second, 0 unlimited rate.",
			},

			"rate_in_size": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Consumption rate limit, unit number/second, 0 is unlimited.",
			},

			"retention_hours": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Consumed message storage policy, unit hour, 0 will be deleted immediately after consumption.",
			},

			"retention_size": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Consumed message storage strategy, unit G, 0 Delete immediately after consumption.",
			},

			"replicas": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Duplicate number.",
			},

			"remark": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Remark.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTdmqEnvironmentAttributesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdmq_environment_attributes.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqEnvironmentAttributesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		msgTTL = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(msgTTL))
	if environmentId != nil {
		_ = d.Set("environment_id", environmentId)
	}

	if msgTTL != nil {
		_ = d.Set("msg_t_t_l", msgTTL)
	}

	if rateInByte != nil {
		_ = d.Set("rate_in_byte", rateInByte)
	}

	if rateInSize != nil {
		_ = d.Set("rate_in_size", rateInSize)
	}

	if retentionHours != nil {
		_ = d.Set("retention_hours", retentionHours)
	}

	if retentionSize != nil {
		_ = d.Set("retention_size", retentionSize)
	}

	if replicas != nil {
		_ = d.Set("replicas", replicas)
	}

	if remark != nil {
		_ = d.Set("remark", remark)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
