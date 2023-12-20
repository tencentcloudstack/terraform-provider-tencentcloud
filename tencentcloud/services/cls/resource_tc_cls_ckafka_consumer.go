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
					},
				},
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
							Required:    true,
							Description: "vip.",
						},
						"vport": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "vport.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "instance id.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "instance name.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "topic id of ckafka.",
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
		request.Content = &consumerContent
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

		_ = d.Set("content", []interface{}{contentMap})
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
	mutableArgs := []string{"need_content", "content", "ckafka", "compression"}

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
			request.Content = &consumerContent
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
