package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTsfDeliveryConfigByGroupId() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfDeliveryConfigByGroupIdRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "groupId.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "configuration item for deliver to a Kafka.",
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

func dataSourceTencentCloudTsfDeliveryConfigByGroupIdRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_delivery_config_by_group_id.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var groupId string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		paramMap["GroupId"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var deliveryConfig *tsf.SimpleKafkaDeliveryConfig
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfDeliveryConfigByGroupIdByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		deliveryConfig = result
		return nil
	})
	if err != nil {
		return err
	}

	simpleKafkaDeliveryConfigMap := map[string]interface{}{}
	if deliveryConfig != nil {
		if deliveryConfig.ConfigId != nil {
			simpleKafkaDeliveryConfigMap["config_id"] = deliveryConfig.ConfigId
		}

		if deliveryConfig.ConfigName != nil {
			simpleKafkaDeliveryConfigMap["config_name"] = deliveryConfig.ConfigName
		}

		_ = d.Set("result", []interface{}{simpleKafkaDeliveryConfigMap})
	}

	d.SetId(groupId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), simpleKafkaDeliveryConfigMap); e != nil {
			return e
		}
	}
	return nil
}
