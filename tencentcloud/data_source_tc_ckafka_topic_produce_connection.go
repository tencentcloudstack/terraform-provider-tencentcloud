package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCkafkaTopicProduceConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaTopicProduceConnectionRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},

			"topic_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "TopicName.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "link information return result set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_addr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ip address.",
						},
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "connect time.",
						},
						"is_un_support_version": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is the supported version.",
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

func dataSourceTencentCloudCkafkaTopicProduceConnectionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ckafka_topic_produce_connection.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var (
		instanceId string
		topicName  string
	)
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["instance_id"] = helper.String(instanceId)
	}

	if v, ok := d.GetOk("topic_name"); ok {
		topicName = v.(string)
		paramMap["topic_name"] = helper.String(topicName)
	}

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*ckafka.DescribeConnectInfoResultDTO

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		topicProduceConnection, e := service.DescribeCkafkaTopicProduceConnectionByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = topicProduceConnection
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0)

	if result != nil {
		for _, describeConnectInfoResultDTO := range result {
			describeConnectInfoResultDTOMap := map[string]interface{}{}

			if describeConnectInfoResultDTO.IpAddr != nil {
				describeConnectInfoResultDTOMap["ip_addr"] = describeConnectInfoResultDTO.IpAddr
			}

			if describeConnectInfoResultDTO.Time != nil {
				describeConnectInfoResultDTOMap["time"] = describeConnectInfoResultDTO.Time
			}

			if describeConnectInfoResultDTO.IsUnSupportVersion != nil {
				describeConnectInfoResultDTOMap["is_un_support_version"] = describeConnectInfoResultDTO.IsUnSupportVersion
			}

			tmpList = append(tmpList, describeConnectInfoResultDTOMap)
		}

		_ = d.Set("result", tmpList)
	}

	d.SetId(instanceId + FILED_SP + topicName)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
