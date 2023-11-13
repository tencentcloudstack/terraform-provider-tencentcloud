/*
Use this data source to query detailed information of tsf describe_delivery_config_by_group_i_d

Example Usage

```hcl
data "tencentcloud_tsf_describe_delivery_config_by_group_i_d" "describe_delivery_config_by_group_i_d" {
  group_id = "group-yrjkln9v"
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTsfDescribeDeliveryConfigByGroupID() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfDescribeDeliveryConfigByGroupIDRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "GroupId.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Configuration item for deliver to a Kafka .",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Config ID. Note: This field may return null, which means that no valid value was obtained.",
						},
						"config_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Config Name. Note: This field may return null, which means that no valid value was obtained.",
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

func dataSourceTencentCloudTsfDescribeDeliveryConfigByGroupIDRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_describe_delivery_config_by_group_i_d.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*tsf.SimpleKafkaDeliveryConfig

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfDescribeDeliveryConfigByGroupIDByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		simpleKafkaDeliveryConfigMap := map[string]interface{}{}

		if result.ConfigId != nil {
			simpleKafkaDeliveryConfigMap["config_id"] = result.ConfigId
		}

		if result.ConfigName != nil {
			simpleKafkaDeliveryConfigMap["config_name"] = result.ConfigName
		}

		ids = append(ids, *result.GroupId)
		_ = d.Set("result", simpleKafkaDeliveryConfigMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), simpleKafkaDeliveryConfigMap); e != nil {
			return e
		}
	}
	return nil
}
