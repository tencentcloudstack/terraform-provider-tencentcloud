package cls

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsKafkaConsumer() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsKafkaConsumerCreate,
		Read:   resourceTencentCloudClsKafkaConsumerRead,
		Update: resourceTencentCloudClsKafkaConsumerUpdate,
		Delete: resourceTencentCloudClsKafkaConsumerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"from_topic_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Log topic ID.",
			},
			"compression": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Compression method: 0-NONE, 2-SNAPPY, 3-LZ4.",
			},
			"consumer_content": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Kafka protocol consumption data format.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Content format: 0-original content, 1-JSON.",
						},
						"enable_tag": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether to deliver TAG information.",
						},
						"meta_fields": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Metadata field list.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"tag_transaction": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Tag data processing method: 1-not flattened, 2-flattened.",
						},
						"json_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Consumption data JSON format: 1-not escaped, 2-escaped.",
						},
					},
				},
			},
			"topic_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Topic parameter used when KafkaConsumer consumes.",
			},
		},
	}
}

func resourceTencentCloudClsKafkaConsumerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_kafka_consumer.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		request     = cls.NewOpenKafkaConsumerRequest()
		fromTopicId string
	)

	if v, ok := d.GetOk("from_topic_id"); ok {
		fromTopicId = v.(string)
		request.FromTopicId = helper.String(fromTopicId)
	}

	if v, ok := d.GetOkExists("compression"); ok {
		request.Compression = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("consumer_content"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			kafkaConsumerContent := cls.KafkaConsumerContent{}
			if v, ok := dMap["format"]; ok {
				kafkaConsumerContent.Format = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["enable_tag"]; ok {
				kafkaConsumerContent.EnableTag = helper.Bool(v.(bool))
			}
			if v, ok := dMap["meta_fields"]; ok {
				metaFieldsSet := v.([]interface{})
				for i := range metaFieldsSet {
					metaFields := metaFieldsSet[i].(string)
					kafkaConsumerContent.MetaFields = append(kafkaConsumerContent.MetaFields, helper.String(metaFields))
				}
			}
			if v, ok := dMap["tag_transaction"]; ok {
				kafkaConsumerContent.TagTransaction = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["json_type"]; ok {
				kafkaConsumerContent.JsonType = helper.IntInt64(v.(int))
			}
			request.ConsumerContent = &kafkaConsumerContent
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().OpenKafkaConsumerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cls kafka consumer failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(fromTopicId)

	return resourceTencentCloudClsKafkaConsumerRead(d, meta)
}

func resourceTencentCloudClsKafkaConsumerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_kafka_consumer.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		request     = cls.NewDescribeKafkaConsumerRequest()
		response    *cls.DescribeKafkaConsumerResponse
		fromTopicId = d.Id()
	)

	request.FromTopicId = helper.String(fromTopicId)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().DescribeKafkaConsumerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cls kafka consumer failed, reason:%+v", logId, err)
		return err
	}

	if response == nil || response.Response == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsKafkaConsumer` [%s] not found, please check if it has been deleted.\n", logId, fromTopicId)
		return nil
	}

	if response.Response.Status != nil && !*response.Response.Status {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsKafkaConsumer` [%s] status is false, resource has been closed.\n", logId, fromTopicId)
		return nil
	}

	_ = d.Set("from_topic_id", fromTopicId)

	if response.Response.Compression != nil {
		_ = d.Set("compression", response.Response.Compression)
	}

	if response.Response.TopicID != nil {
		_ = d.Set("topic_id", response.Response.TopicID)
	}

	if response.Response.ConsumerContent != nil {
		consumerContentMap := map[string]interface{}{}

		if response.Response.ConsumerContent.Format != nil {
			consumerContentMap["format"] = response.Response.ConsumerContent.Format
		}

		if response.Response.ConsumerContent.EnableTag != nil {
			consumerContentMap["enable_tag"] = response.Response.ConsumerContent.EnableTag
		}

		if response.Response.ConsumerContent.MetaFields != nil {
			consumerContentMap["meta_fields"] = response.Response.ConsumerContent.MetaFields
		}

		if response.Response.ConsumerContent.TagTransaction != nil {
			consumerContentMap["tag_transaction"] = response.Response.ConsumerContent.TagTransaction
		}

		if response.Response.ConsumerContent.JsonType != nil {
			consumerContentMap["json_type"] = response.Response.ConsumerContent.JsonType
		}

		_ = d.Set("consumer_content", []interface{}{consumerContentMap})
	}

	return nil
}

func resourceTencentCloudClsKafkaConsumerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_kafka_consumer.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		request     = cls.NewModifyKafkaConsumerRequest()
		fromTopicId = d.Id()
	)

	request.FromTopicId = helper.String(fromTopicId)

	if d.HasChange("compression") {
		if v, ok := d.GetOkExists("compression"); ok {
			request.Compression = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("consumer_content") {
		if v, ok := d.GetOk("consumer_content"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				kafkaConsumerContent := cls.KafkaConsumerContent{}
				if v, ok := dMap["format"]; ok {
					kafkaConsumerContent.Format = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["enable_tag"]; ok {
					kafkaConsumerContent.EnableTag = helper.Bool(v.(bool))
				}
				if v, ok := dMap["meta_fields"]; ok {
					metaFieldsSet := v.([]interface{})
					for i := range metaFieldsSet {
						metaFields := metaFieldsSet[i].(string)
						kafkaConsumerContent.MetaFields = append(kafkaConsumerContent.MetaFields, helper.String(metaFields))
					}
				}
				if v, ok := dMap["tag_transaction"]; ok {
					kafkaConsumerContent.TagTransaction = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["json_type"]; ok {
					kafkaConsumerContent.JsonType = helper.IntInt64(v.(int))
				}
				request.ConsumerContent = &kafkaConsumerContent
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyKafkaConsumerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cls kafka consumer failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudClsKafkaConsumerRead(d, meta)
}

func resourceTencentCloudClsKafkaConsumerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_kafka_consumer.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		request     = cls.NewCloseKafkaConsumerRequest()
		fromTopicId = d.Id()
	)

	request.FromTopicId = helper.String(fromTopicId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CloseKafkaConsumerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cls kafka consumer failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
