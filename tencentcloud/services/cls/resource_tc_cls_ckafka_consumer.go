package cls

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsCkafkaConsumer() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsCkafkaConsumerCreate,
		Read:   resourceTencentCloudClsCkafkaConsumerRead,
		Update: resourceTencentCloudClsCkafkaConsumerUpdate,
		Delete: resourceTencentCloudClsCkafkaConsumerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "topic id.",
			},

			"need_content": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "whether to deliver the metadata information of the log.",
			},

			"content": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "metadata information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_tag": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "whether to deliver the TAG info.",
						},
						"meta_fields": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "metadata info list.",
						},
						"tag_json_not_tiled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "whether to tiling tag json.",
						},
						"timestamp_accuracy": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "delivery timestamp precision,1 for second, 2 for millisecond.",
						},
						"json_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "JSON format for delivery. Valid values:\n- `0`: consistent with the original log, no escaping.\n- `1`: escaped.",
						},
						"auto_convert_number": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "whether to automatically convert numeric strings (such as `123`) in the first-level values of the JSON structure to numeric types (int / float). Default value: `false`.",
						},
					},
				},
			},

			"role_arn": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "role access description name for cross-account delivery.",
			},

			"external_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "external id for cross-account delivery.",
			},

			"advanced_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "advanced configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"partition_hash_status": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "whether to deliver messages with the same hash value of fields to the same ckafka partition. Default value: `false`.",
						},
						"partition_fields": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "list of fields used to calculate hash. Supports up to 5 fields.",
						},
					},
				},
			},

			"dsl_filter": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "log pre-filtering - pre-filter the raw data written to ckafka.",
			},

			"ckafka": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "ckafka info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "vip.If the delivery task is created using a Role ARN, this field can be left empty.",
						},
						"vport": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "vport.If the delivery task is created using a Role ARN, this field can be left empty.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "instance id.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "instance name.If the delivery task is created using a Role ARN, this field can be left empty.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "topic id of ckafka. If the delivery task is created using a Role ARN, this field can be left empty.",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "topic name of ckafka.",
						},
					},
				},
			},

			"compression": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "compression method. 0 for NONE, 2 for SNAPPY, 3 for LZ4.",
			},
		},
	}
}

func resourceTencentCloudClsCkafkaConsumerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_ckafka_consumer.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = cls.NewCreateConsumerRequest()
		topicId string
	)
	if v, ok := d.GetOk("topic_id"); ok {
		topicId = v.(string)
		request.TopicId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("need_content"); ok {
		request.NeedContent = helper.Bool(v.(bool))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "content"); ok {
		consumerContent := cls.ConsumerContent{}
		if v, ok := dMap["enable_tag"]; ok {
			consumerContent.EnableTag = helper.Bool(v.(bool))
		}
		if v, ok := dMap["meta_fields"]; ok {
			metaFieldsSet := v.(*schema.Set).List()
			for i := range metaFieldsSet {
				metaFields := metaFieldsSet[i].(string)
				consumerContent.MetaFields = append(consumerContent.MetaFields, &metaFields)
			}
		}
		if v, ok := dMap["tag_json_not_tiled"]; ok {
			consumerContent.TagJsonNotTiled = helper.Bool(v.(bool))
		}
		if v, ok := dMap["timestamp_accuracy"]; ok {
			consumerContent.TimestampAccuracy = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["json_type"]; ok {
			consumerContent.JsonType = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["auto_convert_number"]; ok {
			consumerContent.AutoConvertNumber = helper.Bool(v.(bool))
		}
		request.Content = &consumerContent
	}

	if v, ok := d.GetOk("role_arn"); ok {
		request.RoleArn = helper.String(v.(string))
	}

	if v, ok := d.GetOk("external_id"); ok {
		request.ExternalId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "advanced_config"); ok {
		advancedConfig := cls.AdvancedConsumerConfiguration{}
		if v, ok := dMap["partition_hash_status"]; ok {
			advancedConfig.PartitionHashStatus = helper.Bool(v.(bool))
		}
		if v, ok := dMap["partition_fields"]; ok {
			partitionFieldsSet := v.(*schema.Set).List()
			for i := range partitionFieldsSet {
				partitionFields := partitionFieldsSet[i].(string)
				advancedConfig.PartitionFields = append(advancedConfig.PartitionFields, &partitionFields)
			}
		}
		request.AdvancedConfig = &advancedConfig
	}

	if v, ok := d.GetOk("dsl_filter"); ok {
		request.DSLFilter = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ckafka"); ok {
		ckafka := cls.Ckafka{}
		if v, ok := dMap["vip"]; ok {
			ckafka.Vip = helper.String(v.(string))
		}
		if v, ok := dMap["vport"]; ok {
			ckafka.Vport = helper.String(v.(string))
		}
		if v, ok := dMap["instance_id"]; ok {
			ckafka.InstanceId = helper.String(v.(string))
		}
		if v, ok := dMap["instance_name"]; ok {
			ckafka.InstanceName = helper.String(v.(string))
		}
		if v, ok := dMap["topic_id"]; ok {
			ckafka.TopicId = helper.String(v.(string))
		}
		if v, ok := dMap["topic_name"]; ok {
			ckafka.TopicName = helper.String(v.(string))
		}
		request.Ckafka = &ckafka
	}

	if v, ok := d.GetOkExists("compression"); ok {
		request.Compression = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateConsumer(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cls ckafkaConsumer failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(topicId)

	return resourceTencentCloudClsCkafkaConsumerRead(d, meta)
}

func resourceTencentCloudClsCkafkaConsumerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_ckafka_consumer.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	topicId := d.Id()

	ckafkaConsumer, err := service.DescribeClsCkafkaConsumerById(ctx, topicId)
	if err != nil {
		return err
	}

	if ckafkaConsumer == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsCkafkaConsumer` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("topic_id", topicId)

	if ckafkaConsumer.NeedContent != nil {
		_ = d.Set("need_content", ckafkaConsumer.NeedContent)
	}

	if ckafkaConsumer.Content != nil {
		contentMap := map[string]interface{}{}

		if ckafkaConsumer.Content.EnableTag != nil {
			contentMap["enable_tag"] = ckafkaConsumer.Content.EnableTag
		}

		if ckafkaConsumer.Content.MetaFields != nil {
			contentMap["meta_fields"] = ckafkaConsumer.Content.MetaFields
		}

		if ckafkaConsumer.Content.TagJsonNotTiled != nil {
			contentMap["tag_json_not_tiled"] = ckafkaConsumer.Content.TagJsonNotTiled
		}

		if ckafkaConsumer.Content.TimestampAccuracy != nil {
			contentMap["timestamp_accuracy"] = ckafkaConsumer.Content.TimestampAccuracy
		}

		if ckafkaConsumer.Content.JsonType != nil {
			contentMap["json_type"] = ckafkaConsumer.Content.JsonType
		}

		if ckafkaConsumer.Content.AutoConvertNumber != nil {
			contentMap["auto_convert_number"] = ckafkaConsumer.Content.AutoConvertNumber
		}

		_ = d.Set("content", []interface{}{contentMap})
	}

	if ckafkaConsumer.RoleArn != nil {
		_ = d.Set("role_arn", ckafkaConsumer.RoleArn)
	}

	if ckafkaConsumer.ExternalId != nil {
		_ = d.Set("external_id", ckafkaConsumer.ExternalId)
	}

	if ckafkaConsumer.AdvancedConfig != nil {
		advancedConfigMap := map[string]interface{}{}

		if ckafkaConsumer.AdvancedConfig.PartitionHashStatus != nil {
			advancedConfigMap["partition_hash_status"] = ckafkaConsumer.AdvancedConfig.PartitionHashStatus
		}

		if ckafkaConsumer.AdvancedConfig.PartitionFields != nil {
			advancedConfigMap["partition_fields"] = ckafkaConsumer.AdvancedConfig.PartitionFields
		}

		_ = d.Set("advanced_config", []interface{}{advancedConfigMap})
	}

	if ckafkaConsumer.DSLFilter != nil {
		_ = d.Set("dsl_filter", ckafkaConsumer.DSLFilter)
	}

	if ckafkaConsumer.Ckafka != nil {
		ckafkaMap := map[string]interface{}{}

		if ckafkaConsumer.Ckafka.Vip != nil {
			ckafkaMap["vip"] = ckafkaConsumer.Ckafka.Vip
		}

		if ckafkaConsumer.Ckafka.Vport != nil {
			ckafkaMap["vport"] = ckafkaConsumer.Ckafka.Vport
		}

		if ckafkaConsumer.Ckafka.InstanceId != nil {
			ckafkaMap["instance_id"] = ckafkaConsumer.Ckafka.InstanceId
		}

		if ckafkaConsumer.Ckafka.InstanceName != nil {
			ckafkaMap["instance_name"] = ckafkaConsumer.Ckafka.InstanceName
		}

		if ckafkaConsumer.Ckafka.TopicId != nil {
			ckafkaMap["topic_id"] = ckafkaConsumer.Ckafka.TopicId
		}

		if ckafkaConsumer.Ckafka.TopicName != nil {
			ckafkaMap["topic_name"] = ckafkaConsumer.Ckafka.TopicName
		}

		_ = d.Set("ckafka", []interface{}{ckafkaMap})
	}

	if ckafkaConsumer.Compression != nil {
		_ = d.Set("compression", ckafkaConsumer.Compression)
	}

	return nil
}

func resourceTencentCloudClsCkafkaConsumerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_ckafka_consumer.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cls.NewModifyConsumerRequest()

	topicId := d.Id()

	request.TopicId = &topicId

	needChange := false
	mutableArgs := []string{"need_content", "content", "ckafka", "compression", "role_arn", "external_id", "advanced_config", "dsl_filter"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOkExists("need_content"); ok {
			request.NeedContent = helper.Bool(v.(bool))
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "content"); ok {
			consumerContent := cls.ConsumerContent{}
			if v, ok := dMap["enable_tag"]; ok {
				consumerContent.EnableTag = helper.Bool(v.(bool))
			}
			if v, ok := dMap["meta_fields"]; ok {
				metaFieldsSet := v.(*schema.Set).List()
				for i := range metaFieldsSet {
					metaFields := metaFieldsSet[i].(string)
					consumerContent.MetaFields = append(consumerContent.MetaFields, &metaFields)
				}
			}
			if v, ok := dMap["tag_json_not_tiled"]; ok {
				consumerContent.TagJsonNotTiled = helper.Bool(v.(bool))
			}
			if v, ok := dMap["timestamp_accuracy"]; ok {
				consumerContent.TimestampAccuracy = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["json_type"]; ok {
				consumerContent.JsonType = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["auto_convert_number"]; ok {
				consumerContent.AutoConvertNumber = helper.Bool(v.(bool))
			}
			request.Content = &consumerContent
		}

		if v, ok := d.GetOk("role_arn"); ok {
			request.RoleArn = helper.String(v.(string))
		}

		if v, ok := d.GetOk("external_id"); ok {
			request.ExternalId = helper.String(v.(string))
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "advanced_config"); ok {
			advancedConfig := cls.AdvancedConsumerConfiguration{}
			if v, ok := dMap["partition_hash_status"]; ok {
				advancedConfig.PartitionHashStatus = helper.Bool(v.(bool))
			}
			if v, ok := dMap["partition_fields"]; ok {
				partitionFieldsSet := v.(*schema.Set).List()
				for i := range partitionFieldsSet {
					partitionFields := partitionFieldsSet[i].(string)
					advancedConfig.PartitionFields = append(advancedConfig.PartitionFields, &partitionFields)
				}
			}
			request.AdvancedConfig = &advancedConfig
		}

		if v, ok := d.GetOk("dsl_filter"); ok {
			request.DSLFilter = helper.String(v.(string))
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "ckafka"); ok {
			ckafka := cls.Ckafka{}
			if v, ok := dMap["vip"]; ok {
				ckafka.Vip = helper.String(v.(string))
			}
			if v, ok := dMap["vport"]; ok {
				ckafka.Vport = helper.String(v.(string))
			}
			if v, ok := dMap["instance_id"]; ok {
				ckafka.InstanceId = helper.String(v.(string))
			}
			if v, ok := dMap["instance_name"]; ok {
				ckafka.InstanceName = helper.String(v.(string))
			}
			if v, ok := dMap["topic_id"]; ok {
				ckafka.TopicId = helper.String(v.(string))
			}
			if v, ok := dMap["topic_name"]; ok {
				ckafka.TopicName = helper.String(v.(string))
			}
			request.Ckafka = &ckafka
		}

		if v, ok := d.GetOkExists("compression"); ok {
			request.Compression = helper.IntInt64(v.(int))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyConsumer(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cls ckafkaConsumer failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClsCkafkaConsumerRead(d, meta)
}

func resourceTencentCloudClsCkafkaConsumerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_ckafka_consumer.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	topicId := d.Id()

	if err := service.DeleteClsCkafkaConsumerById(ctx, topicId); err != nil {
		return err
	}

	return nil
}
