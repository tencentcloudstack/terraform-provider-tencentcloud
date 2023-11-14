/*
Provides a resource to create a cls ckafka_consumer

Example Usage

```hcl
resource "tencentcloud_cls_ckafka_consumer" "ckafka_consumer" {
  topic_id = "5cd3a17e-fb0b-418c-afd7-77b365397426"
  need_content = true
  content {
		enable_tag = true
		meta_fields =
		tag_json_not_tiled = true
		timestamp_accuracy = 1

  }
  ckafka {
		vip = "1.1.1.1"
		vport = "8000"
		instance_id = "ckafka-xxxxx"
		instance_name = "test"
		topic_id = "topic-5cd3a17e-fb0b-418c-afd7-77b3653974xx"
		topic_name = "test"

  }
  compression = 0
}
```

Import

cls ckafka_consumer can be imported using the id, e.g.

```
terraform import tencentcloud_cls_ckafka_consumer.ckafka_consumer ckafka_consumer_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudClsCkafkaConsumer() *schema.Resource {
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
				Description: "Topic id.",
			},

			"need_content": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to deliver the metadata information of the log.",
			},

			"content": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Metadata information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_tag": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to deliver the TAG info.",
						},
						"meta_fields": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Metadata info list.",
						},
						"tag_json_not_tiled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to tiling tag json.",
						},
						"timestamp_accuracy": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Delivery timestamp precision,1 for second, 2 for millisecond.",
						},
					},
				},
			},

			"ckafka": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Ckafka info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Vip.",
						},
						"vport": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Vport.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance id.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance name.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Topic id of ckafka.",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Topic name of ckafka.",
						},
					},
				},
			},

			"compression": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Compression method. 0 for NONE, 2 for SNAPPY, 3 for LZ4.",
			},
		},
	}
}

func resourceTencentCloudClsCkafkaConsumerCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_ckafka_consumer.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cls.NewCreateConsumerRequest()
		response = cls.NewCreateConsumerResponse()
		topicId  string
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().CreateConsumer(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cls ckafkaConsumer failed, reason:%+v", logId, err)
		return err
	}

	topicId = *response.Response.TopicId
	d.SetId(topicId)

	return resourceTencentCloudClsCkafkaConsumerRead(d, meta)
}

func resourceTencentCloudClsCkafkaConsumerRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_ckafka_consumer.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}

	ckafkaConsumerId := d.Id()

	ckafkaConsumer, err := service.DescribeClsCkafkaConsumerById(ctx, topicId)
	if err != nil {
		return err
	}

	if ckafkaConsumer == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsCkafkaConsumer` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ckafkaConsumer.TopicId != nil {
		_ = d.Set("topic_id", ckafkaConsumer.TopicId)
	}

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
	defer logElapsed("resource.tencentcloud_cls_ckafka_consumer.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cls.NewModifyConsumerRequest()

	ckafkaConsumerId := d.Id()

	request.TopicId = &topicId

	immutableArgs := []string{"topic_id", "need_content", "content", "ckafka", "compression"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("topic_id") {
		if v, ok := d.GetOk("topic_id"); ok {
			request.TopicId = helper.String(v.(string))
		}
	}

	if d.HasChange("need_content") {
		if v, ok := d.GetOkExists("need_content"); ok {
			request.NeedContent = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("content") {
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
	}

	if d.HasChange("ckafka") {
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
	}

	if d.HasChange("compression") {
		if v, ok := d.GetOkExists("compression"); ok {
			request.Compression = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().ModifyConsumer(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cls ckafkaConsumer failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudClsCkafkaConsumerRead(d, meta)
}

func resourceTencentCloudClsCkafkaConsumerDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_ckafka_consumer.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	ckafkaConsumerId := d.Id()

	if err := service.DeleteClsCkafkaConsumerById(ctx, topicId); err != nil {
		return err
	}

	return nil
}
