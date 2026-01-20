package ckafka

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func DataSourceTencentCloudCkafkaVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaVersionRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CKafka instance ID.",
			},

			"kafka_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current Kafka version.",
			},

			"cur_broker_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current broker version.",
			},

			"latest_broker_versions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of latest broker versions supported by the platform.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kafka_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kafka version.",
						},
						"broker_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Broker version.",
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

func dataSourceTencentCloudCkafkaVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_version.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(nil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Get("instance_id").(string)
	)

	var respData *ckafka.InstanceVersion
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCkafkaVersionByFilter(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	if respData != nil {
		if respData.KafkaVersion != nil {
			_ = d.Set("kafka_version", respData.KafkaVersion)
		}

		if respData.CurBrokerVersion != nil {
			_ = d.Set("cur_broker_version", respData.CurBrokerVersion)
		}

		latestBrokerVersionsList := make([]map[string]interface{}, 0, len(respData.LatestBrokerVersion))
		if respData.LatestBrokerVersion != nil {
			for _, latestBrokerVersions := range respData.LatestBrokerVersion {
				latestBrokerVersionsMap := map[string]interface{}{}
				if latestBrokerVersions.KafkaVersion != nil {
					latestBrokerVersionsMap["kafka_version"] = latestBrokerVersions.KafkaVersion
				}

				if latestBrokerVersions.BrokerVersion != nil {
					latestBrokerVersionsMap["broker_version"] = latestBrokerVersions.BrokerVersion
				}

				latestBrokerVersionsList = append(latestBrokerVersionsList, latestBrokerVersionsMap)
			}

			_ = d.Set("latest_broker_versions", latestBrokerVersionsList)
		}
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		outputData := map[string]interface{}{
			"kafka_version":          d.Get("kafka_version"),
			"cur_broker_version":     d.Get("cur_broker_version"),
			"latest_broker_versions": d.Get("latest_broker_versions"),
		}
		if e := tccommon.WriteToFile(output.(string), outputData); e != nil {
			return e
		}
	}

	return nil
}
